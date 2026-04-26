package bootstrapper

import (
	"fmt"
	"os/exec"
	"time"

	"axell.me/rblxutils/common"
)

func ModifyHostsFile(addEntries bool) {
	common.State.HelperAction = "hosts-add"
	if !addEntries {
		common.State.HelperAction = "hosts-remove"
	}
	
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
		fmt.Println("waiting for helper to finish...")
		time.Sleep(50 * time.Millisecond)
		common.LoadState()
		if common.State.HelperAction == "" {
			break
		}
	}
}