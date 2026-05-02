package uninstaller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"axell.me/rblxutils/common"
	"golang.org/x/sys/windows"
)

func LaunchUninstaller() {
	if !common.YesNo("Are you sure you would like to uninstall rblxutils?") {
		return
	}

	if !common.YesNo("Would you like to keep your config (your mods, settings, etc)?") {
		err := os.Remove(common.LPath("./config.json"))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("now removing helper")
	verb, _ := syscall.UTF16PtrFromString("runas")
	program, _ := syscall.UTF16PtrFromString("schtasks")
	args, _ := syscall.UTF16PtrFromString(`/delete /f /tn "rblxutils-proxy-helper"`)
	null := uint16(0)
	err := windows.ShellExecute(0, verb, program, args, &null, 1)
	if err != nil {
		common.FatalErrorStr("Could not uninstall rblxutils proxy helper: " + err.Error())
	}

	fmt.Println("now removing protocol handler")
	common.RemoveAsProtocolHandler()
	fmt.Println("removing old versions")
	err = os.RemoveAll(common.LPath("./versions"))
	if err != nil {
		common.FatalError(err)
	}
	fmt.Println("base uninstallation complete")
	common.Notification("Rblxutils has been uninstalled")

	if common.YesNo("Would you like to install Roblox with the regular bootstrapper again?") {
		resp, err := http.Get("http://setup.rbxcdn.com/RobloxPlayerInstaller.exe")
		if err != nil {
			common.FatalError(err)
		}

		f, err := os.Create(common.LPath("./RobloxPlayerInstaller.exe"))
		if err != nil {
			common.FatalError(err)
		}
		defer f.Close()

		wi, err := io.Copy(f, resp.Body)
		if err != nil {
			common.FatalError(err)
		}
		f.Close()

		fmt.Println("downloaded roblox installer, total", wi, "bytes")
		fmt.Println("now running installer...")

		p, err := os.StartProcess(common.LPath("RobloxPlayerInstaller.exe"), []string{common.LPath("RobloxPlayerInstaller.exe")}, &os.ProcAttr{
			Files: []*os.File{
				nil, nil, nil,
			},
		})
		if err != nil {
			common.FatalError(err)
		}

		err = p.Release()
		if err != nil {
			common.FatalError(err)
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}
}