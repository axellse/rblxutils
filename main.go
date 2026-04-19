package main

import (
	"fmt"
	"os"

	"axell.me/rblxutils/bootstrapper"
	"axell.me/rblxutils/common"
	"axell.me/rblxutils/configurator"
	"axell.me/rblxutils/resources"
	"github.com/gen2brain/beeep"
)
func main() {
	beeep.AppName = resources.ProgramName
	common.LoadConfiguration()
	common.DefineEnvs()
	fmt.Println("first up, registering as protocol handler.")
	common.RegisterAsProtocolHandler()
	fmt.Println("config and envs loaded, now determining what to do.")
	if len(os.Args) == 1 {
		fmt.Println("no args, starting configurator...")
		configurator.LaunchConfigurator()
	} else {
		fmt.Println("launching bootstrapper")
		bootstrapper.LaunchBootstrapper()
	}
}