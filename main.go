package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/axellse/rblxutils/bootstrapper"
	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/configurator"
	"github.com/axellse/rblxutils/installer"
	"github.com/axellse/rblxutils/proxy"
	"github.com/gen2brain/beeep"
)

func main() {
	fmt.Println("rblxutils started")
	time.Sleep(300 * time.Millisecond) //give the fs changes some time to marinate

	beeep.AppName = "Rblxutils"
	common.DefineEnvs()
	common.InitCountryCodeMap()

	foundConfig := common.LoadConfiguration()
	common.LoadState()
	fmt.Println("base stuff has been setup (config, state, envs, etc...)")
	if len(os.Args) > 1 && os.Args[1] == "-helper" {
		proxy.StartProxy()
		return
	}

	if !foundConfig && common.DotSlash != filepath.Join(common.LocalAppData, "rblxutils") {
		fmt.Println("config not found, asking user what to do.")
		installer.LaunchInstaller()
		return
	}

	fmt.Println("first up, registering as protocol handler.")
	common.RegisterProtocolHandler()
	fmt.Println("okay, now making sure proxy helper is set up.")
	installer.HelperInstallFlow()
	fmt.Println("okay, now asserting shortcut status")
	common.AssertShortcuts()
	fmt.Println("everything ready, now determining what to do")
	fmt.Println("--------------------------------------------------------------")

	if len(os.Args) == 1 {
		fmt.Println("no args, starting configurator...")
		configurator.LaunchConfigurator(nil, nil)
	} else if len(os.Args) > 1 && os.Args[1] == "uninstall" {
		installer.LaunchUninstaller()
	} else {
		fmt.Println("launching bootstrapper (as new instance)")
		bootstrapper.LaunchBootstrapper(true, strings.Join(os.Args[1:], " "))
	}
}

func init() {
	runtime.LockOSThread()
}
