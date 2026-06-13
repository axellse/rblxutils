package installer

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/axellse/rblxutils/common"
	hidehelper "github.com/axellse/rblxutils/installer/hideHelper"
	keepalive "github.com/axellse/rblxutils/installer/keepAlive"
	"golang.org/x/sys/windows"
)

type SCHTASKSTaskExec struct {
	Command string
	Arguments string
}
type SCHTASKSTask struct {
	Actions []SCHTASKSTaskExec `xml:"Actions>Exec"`
}

func IsHelperInstalled() bool {
	cmd := exec.Command("schtasks", `/query`, `/tn`, `rblxutils-proxy-helper`, `/xml`)
	ba, err := cmd.CombinedOutput()
	fmt.Println(string(ba))
	if err != nil {
		return false
	}

	if strings.Contains(string(ba), "ERROR") {
		return false
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(ba))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset != "UTF-16" {
			return nil, errors.New("what the actual FUCK is microsoft doing")
		}

		return input, nil //its fine..
	}

	var task SCHTASKSTask
	err = decoder.Decode(&task)
	if err != nil {
		common.FatalError(err)
	}

	if len(task.Actions) != 1 {
		return false
	}

	if task.Actions[0] != GetHelperExec() {
		return false
	}

	return true
}

func HelperInstallFlow() {
	if !IsHelperInstalled() {
		if common.YesNo("Rblxutils needs to install its helper which requires adminstrator rights. Would you like to continue?") {
			CreateHelperTask()
		} else {
			os.Exit(0)
		}
	}
}

func GetHelperExec() (exec SCHTASKSTaskExec) {
	exec.Command = common.BinPath
	exec.Arguments = "-helper"
	if keepalive.KeepHelperAlive {
		exec.Command = "C:\\Windows\\System32\\cmd.exe"
		exec.Arguments = "/k " + common.BinPath + " -helper"
	} else if hidehelper.HideHelper {
		exec.Command = "C:\\Windows\\System32\\conhost.exe"
		exec.Arguments = "--headless " + common.BinPath + " -helper"
	}

	return
}

func CreateHelperTask() {
	verb, _ := syscall.UTF16PtrFromString("runas")
	program, _ := syscall.UTF16PtrFromString("schtasks")

	exc := GetHelperExec()
	args, _ := syscall.UTF16PtrFromString(`/create /f /tn "rblxutils-proxy-helper" /tr "` + exc.Command + " " + exc.Arguments + `" /sc once /st 00:00 /sd 2000/01/01 /rl highest`)
	null := uint16(0)
	err := windows.ShellExecute(0, verb, program, args, &null, 1)
	if err != nil {
		common.FatalErrorStr("Could not setup rblxutils proxy helper: " + err.Error())
	}
	common.Notification("Everything was sucessfully installed!")
}
