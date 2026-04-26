package bootstrapper

import (
	"fmt"
	"os/exec"
	"time"

	"axell.me/rblxutils/common"
)

func StartProxy() {
	common.LoadState()
	common.State.HelperAction = "start-proxy"
	err := common.WriteState()
	if err != nil {
		common.FatalError(err)
	}

	cmd := exec.Command("schtasks", `/run`, `/tn`, `rblxutils-proxy-helper`)
	ba, err := cmd.CombinedOutput()
	fmt.Println(string(ba))
	if err != nil {
		common.FatalError(err)
	}

	for {
		fmt.Println("waiting for proxy to start up...")
		time.Sleep(50 * time.Millisecond)
		common.LoadState()
		if common.State.HelperAction == "" {
			break
		}
	}
}