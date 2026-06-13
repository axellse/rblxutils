package installer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/common/shortcut"
	"golang.org/x/sys/windows"
)

func RemoveHelper() {
	verb, _ := syscall.UTF16PtrFromString("runas")
	program, _ := syscall.UTF16PtrFromString("schtasks")
	args, _ := syscall.UTF16PtrFromString(`/delete /f /tn "rblxutils-proxy-helper"`)
	null := uint16(0)
	err := windows.ShellExecute(0, verb, program, args, &null, 1)
	if err != nil {
		common.FatalErrorStr("Could not uninstall rblxutils proxy helper: " + err.Error())
	}
}

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
	RemoveHelper()
	fmt.Println("now removing protocol handler")
	common.RemoveAsProtocolHandler()
	fmt.Println("removing old versions")
	err := os.RemoveAll(common.LPath("./versions"))
	if err != nil {
		common.FatalError(err)
	}
	fmt.Println("removing desktop shortcut")
	err = common.DeleteShortcut(shortcut.Desktop)
	if err != nil {
		common.FatalError(err)
	}
	fmt.Println("removing start menu shortcut")
	err = common.DeleteShortcut(shortcut.StartMenu)
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

		cmd := exec.Command(common.LPath("RobloxPlayerInstaller.exe"))
		err = cmd.Start()
		if err != nil {
			common.FatalError(err)
		}

		err = cmd.Process.Release()
		if err != nil {
			common.FatalError(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	os.Exit(0)
}
