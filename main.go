package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"axell.me/rblxutils/bootstrapper"
	"axell.me/rblxutils/common"
	"axell.me/rblxutils/configurator"
	"axell.me/rblxutils/proxy"
	"axell.me/rblxutils/resources"
	"axell.me/rblxutils/uninstaller"
	"github.com/gen2brain/beeep"
)

var hide_helper string
var keep_helper_alive string


func main() {
	fmt.Println("rblxutils started")
	fmt.Println(resources.CatAscii)
	time.Sleep(300 * time.Millisecond) //give the fs changes some time to marinate

	beeep.AppName = "Rblxutils"
	common.DefineEnvs()

	common.LoadConfiguration()
	common.LoadState()
	if len(os.Args) > 1 && os.Args[1] == "-helper" {
		proxy.StartProxy()
		return
	}

	fmt.Println("first up, registering as protocol handler.")
	common.RegisterProtocolHandler()
	fmt.Println("okay, now making sure proxy helper is set up.")
	InstallFlow()
	fmt.Println("config and envs loaded, now determining what to do.")

	if len(os.Args) == 1 {
		fmt.Println("no args, starting configurator...")
		configurator.LaunchConfigurator(false, nil)
	} else if len(os.Args) > 1 && os.Args[1] == "uninstall" {
		uninstaller.LaunchUninstaller()
	} else {
		fmt.Println("launching bootstrapper (as new instance)")
		bootstrapper.LaunchBootstrapper(true, strings.Join(os.Args[1:], " "))
	}
}

func init() {
	runtime.LockOSThread()
}