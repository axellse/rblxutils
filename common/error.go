package common

import (
	"fmt"
	"os"

	"axell.me/rblxutils/resources"
	"github.com/gen2brain/beeep"
	"github.com/sqweek/dialog"
)

//ik its giving JavaAbstractDefaultFactoryBuilderFactory but this is just to be able to pass to both configurator and bootstrapper
type ErrorUtils struct {}

func FatalError(err error) {
	FatalErrorStr(err.Error())
}

func Error(err error) {
	ErrorStr(err.Error())
}

func ErrorStr(err string) {
	if Config.Misc.ErrorStyle == 0 {
		dialog.Message("%s", err).Title("Rblxutils error").Error()
	} else {
		beeep.Alert("rblxutils errored:", err, resources.ProgramLogo)
	}
}

func Notification(text string) {
	if Config.Misc.ErrorStyle == 0 {
		dialog.Message("%s", text).Title("Rblxutils").Info()
	} else {
		beeep.Alert("Rblxutils", text, resources.ProgramLogo)
	}
}

func YesNo(text string) bool {
	return dialog.Message("%s", text).Title("Rblxutils").YesNo()
}

func FatalErrorStr(err string) {
	ErrorStr(err)
	fmt.Println(err)
	os.Exit(1)
}
