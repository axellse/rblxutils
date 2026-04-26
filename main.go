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
	beeep.AppName = resources.ProgramName
	common.DefineEnvs()
	if len(os.Args) > 1 && os.Args[1] == "-helper" {
		fmt.Println("Just rblxutils's helper program!")
		fmt.Println(resources.CatAscii)
		time.Sleep(400 * time.Millisecond) //give the fs changes some time to marinate
		common.LoadState()
		if common.State.HelperAction == "hosts-add" || common.State.HelperAction == "hosts-remove" {
			ModifyHostsFile()
		}

		common.State.HelperAction = ""
		err := common.WriteState()
		if err != nil {
			common.FatalError(err)
		}
		return
	}

	common.LoadConfiguration()
	common.LoadState()
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