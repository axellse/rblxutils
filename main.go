package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/axellse/rblxutils/bootstrapper"
	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/configurator"
	"github.com/axellse/rblxutils/proxy"
	"github.com/axellse/rblxutils/uninstaller"
	"github.com/gen2brain/beeep"
)

var hide_helper string
var keep_helper_alive string

func main() {
	fmt.Println("rblxutils started")
	time.Sleep(300 * time.Millisecond) //give the fs changes some time to marinate

	beeep.AppName = "Rblxutils"
	common.DefineEnvs()
	common.InitCountryCodeMap()

	common.LoadConfiguration()
	common.LoadState()
	fmt.Println("config envs, and state loaded")
	if len(os.Args) > 1 && os.Args[1] == "-helper" {
		proxy.StartProxy()
		return
	}

	fmt.Println("not running as proxy")
	fmt.Println("first up, registering as protocol handler.")
	common.RegisterProtocolHandler()
	fmt.Println("okay, now making sure proxy helper is set up.")
	InstallFlow()
	fmt.Println("okay, now asserting desktop shortcut status")
	common.AssertDesktopShortcut()
	fmt.Println("cleaning up enviroument...")
	os.Remove(common.LPath("./update.bat"))
	fmt.Println("everything ready, now determining what to do")
	fmt.Println("--------------------------------------------------------------")

	if len(os.Args) == 1 {
		fmt.Println("no args, starting configurator...")
		configurator.LaunchConfigurator(nil, nil)
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
