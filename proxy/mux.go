package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"axell.me/rblxutils/common"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func CalcAvg(slice []int) int {
	if len(slice) == 0 {return 0}
	cum := 0
	for _, n := range slice {
		cum += n
	}

	return cum / len(slice)
}

var LockedConn *websocket.Conn

func RunProxyWriteLoop(ctx context.Context) {
	for {
		wsjson.Write(ctx, LockedConn, common.SocketMessage{
			Type: "proxy_stats",
			Stats: common.ProxyStats{
				AvgRewriteDelayNs: CalcAvg(RewriteDelaysNs),
				AvgModifyResponseAssetDeliveryDelayNs: CalcAvg(ModifyResponseAssetDeliveryDelaysNs),
				AvgModifyResponseCdnDelayNs: CalcAvg(ModifyResponseCdnDelaysNs),
			},
		})

		time.Sleep(1 * time.Second)
	}
}

func GetProxyServemux(proxy *httputil.ReverseProxy, closeF *func() error) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/lock/", func (w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			common.FatalError(err)
		}
		defer conn.CloseNow()


		ctx := context.Background()

		for {
			var msg common.SocketMessage
			err := wsjson.Read(ctx, conn, &msg)
			if err != nil {
				if !strings.Contains(err.Error(), "close") {
					common.Error(err)
				}

				fmt.Println("now shutting down.")
				break
			}

			switch msg.Type {
			case "lock":
				if LockedConn == nil {
					err = wsjson.Write(ctx, conn, common.SocketMessage{
						Type: "lock_resp",
						Error: "",
					})
					LockedConn = conn
					go RunProxyWriteLoop(ctx)

					if err != nil {
						common.FatalError(err)
					}
					fmt.Println("acquired lock!")
				} else {
					err := wsjson.Write(ctx, conn, common.SocketMessage{
						Type: "lock_resp",
						Error: "already_locked",
					})
					if err != nil {
						common.FatalError(err)
					}

					wsjson.Write(ctx, LockedConn, common.SocketMessage{
						Type: "new_instance",
						DataB: msg.DataB,
					})
					return 
				}
			}
		}

		err = (*closeF)()
		if err != nil {
			common.FatalErrorStr("closeF err:" + err.Error())
		}
	})
	mux.Handle("/", proxy)

	return mux
}

