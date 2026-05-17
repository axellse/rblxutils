package bootstrapper

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"axell.me/rblxutils/common"
	"axell.me/rblxutils/configurator"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TryLock() (conn *websocket.Conn, state string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "wss://localhost/lock/", &websocket.DialOptions{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	})

	if err != nil {
		return nil, "not_online"
	}

	err = wsjson.Write(ctx, conn, common.SocketMessage{
		Type: "lock",
		DataB: []byte(strings.Join(os.Args[1:], " ")),
	})
	if err != nil {
		return nil, "not_online"
	}

	var resp common.SocketMessage
	err = wsjson.Read(ctx, conn, &resp)
	if err != nil {
		return nil, "not_online"
	}

	if resp.Error != "" {
		return nil, "lock_rejected"
	}

	return conn, "ok"
}

func RunBootstrapperReadLoop(conn *websocket.Conn) {
	ctx := context.Background()

	for {
		var msg common.SocketMessage
		err := wsjson.Read(ctx, conn, &msg)
		if err != nil {
			common.FatalError(err)
		}

		switch msg.Type {
		case "proxy_stats":
			configurator.UiStateMutex.Lock()
			configurator.UIStates.CurrentProxyStats = msg.Stats
			configurator.UiStateMutex.Unlock()
			if configurator.UIStates.Update != nil {
				configurator.UIStates.Update()
			}
		case "new_instance":
			fmt.Println("helper wants us to launch a new instance.")
			LaunchBootstrapper(false, string(msg.DataB))
		}
	}
}

func StartProxy(inman *common.Inman) (*websocket.Conn) {
	conn, state := TryLock()
	switch state {
	case "lock_rejected":
		os.Exit(0)
	case "ok":
		fmt.Println("startproxy return ok??? very weird but we'll go with it.")
		go RunBootstrapperReadLoop(conn)
		return conn
	}

	cmd := exec.Command("schtasks", `/run`, `/tn`, `rblxutils-proxy-helper`)
	ba, err := cmd.CombinedOutput()
	fmt.Println(string(ba))
	if err != nil {
		common.FatalError(err)
	}

	for {
		fmt.Println("waiting for proxy to start up...")
		time.Sleep(100 * time.Millisecond)
		conn, state = TryLock()
		switch state {
		case "ok":
			go RunBootstrapperReadLoop(conn)
			return conn
		case "lock_rejected":
			common.FatalErrorStr("someone stole OUR helper and locked onto it before us, something very cursed is happening on this machine.")
		}
	}
}