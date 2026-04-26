package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

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
	args, _ := syscall.UTF16PtrFromString(`/create /tn "rblxutils-proxy-helper" /tr "` + common.BinPath + ` -helper" /sc once /st 00:00 /sd 2000/01/01 /rl highest`)
	null := uint16(0)
	err := windows.ShellExecute(0, verb, program, args, &null, 1)
	if err != nil {
		common.FatalErrorStr("Could not setup rblxutils proxy helper: " + err.Error())
	}
	common.Notification("Everything was sucessfully installed!")
}

func ModifyHostsFile() {
	ba, err := os.ReadFile("C:\\Windows\\System32\\drivers\\etc\\hosts")
	if err != nil {
		common.FatalError(err)
	}
	
	hosts := strings.ReplaceAll(string(ba), "\r", "")
	hosts = strings.ReplaceAll(hosts, "\n", "\r\n") //make sure the host file is clean

	lines := []string{}
	for line := range strings.SplitSeq(hosts, "\r\n") {
		if !strings.Contains(line, "fts.rbxcdn.com") && !strings.Contains(line, "assetdelivery.roblox.com") && !strings.Contains(line, "rblxutils") {
			lines = append(lines, line)
		}
	}

	if common.State.HelperAction == "hosts-add" {
		lines = append(lines, "# The following two lines were inserted by rblxutils. They should be automatically removed when rblxutils exits.")
		lines = append(lines, "  127.0.0.1     fts.rbxcdn.com")
		lines = append(lines, "  127.0.0.1     assetdelivery.roblox.com")
	}

	finalBa := strings.Join(lines, "\r\n")
	err = os.WriteFile("C:\\Windows\\System32\\drivers\\etc\\hosts", []byte(finalBa), 0666)
	if err != nil {
		common.FatalError(err)
	}

	time.Sleep(100 * time.Millisecond)
	err = exec.Command("ipconfig", "/flushdns").Run()
	if err != nil {
		common.FatalError(err)
	}
	time.Sleep(100 * time.Millisecond)
}