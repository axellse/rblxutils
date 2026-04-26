package main

import (
	"fmt"
	"os"
	"time"

	"axell.me/rblxutils/bootstrapper"
	"axell.me/rblxutils/common"
	"axell.me/rblxutils/configurator"
	"axell.me/rblxutils/resources"
	"github.com/gen2brain/beeep"
)
func main() {
	time.Sleep(300 * time.Millisecond) //give the fs changes some time to marinate
	fmt.Println("rblxutils started")
	fmt.Println(resources.CatAscii)

	beeep.AppName = "Rblxutils"
	common.DefineEnvs()

	common.LoadConfiguration()
	common.LoadState()
	if len(os.Args) > 1 && os.Args[1] == "-helper" {
		fmt.Println("launched rblxutils helper!")
		switch common.State.HelperAction {
		case "start-proxy":
			StartProxy()
		}

		return
	}

	fmt.Println("first up, registering as protocol handler.")
	common.RegisterAsProtocolHandler()
	fmt.Println("okay, now making sure proxy helper is set up.")
	InstallFlow()
	fmt.Println("config and envs loaded, now determining what to do.")

	if len(os.Args) == 1 {
		fmt.Println("no args, starting configurator...")
		configurator.LaunchConfigurator()
	} else {
		fmt.Println("launching bootstrapper")
		bootstrapper.LaunchBootstrapper()
	}
}