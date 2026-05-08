package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"axell.me/rblxutils/common"
	"golang.org/x/sys/windows"
)

func IsHelperInstalled() bool {
	cmd := exec.Command("schtasks", `/query`, `/tn`, `rblxutils-proxy-helper`)
	ba, err := cmd.CombinedOutput()
	fmt.Println(string(ba))
	if err != nil {
		return false
	}

	if strings.Contains(string(ba), "ERROR") {
		return false
	}
	return true
}

func InstallFlow() {
	if !IsHelperInstalled() {
		if common.YesNo("Rblxutils needs to install its helper which requires adminstrator rights. Would you like to continue?") {
			CreateHelperTask()
		} else {
			os.Exit(0)
		}
	}
}

func CreateHelperTask() {
	verb, _ := syscall.UTF16PtrFromString("runas")
	program, _ := syscall.UTF16PtrFromString("schtasks")
	trPrefix := ""
	if keep_helper_alive == "true" {
		trPrefix = "cmd.exe /k "
	} else if hide_helper == "true" {
		trPrefix = "conhost.exe --headless "
	}

	args, _ := syscall.UTF16PtrFromString(`/create /tn "rblxutils-proxy-helper" /tr "` + trPrefix + common.BinPath + ` -helper" /sc once /st 00:00 /sd 2000/01/01 /rl highest`) 
	null := uint16(0)
	err := windows.ShellExecute(0, verb, program, args, &null, 1)
	if err != nil {
		common.FatalErrorStr("Could not setup rblxutils proxy helper: " + err.Error())
	}
	common.Notification("Everything was sucessfully installed!")
}