package common

import (

	"github.com/axellse/rblxutils/resources"
	"github.com/gen2brain/beeep"
	"github.com/sqweek/dialog"
)

func FatalError(err error) {
	FatalErrorStr(err.Error())
}

func Error(err error) {
	ErrorStr(err.Error())
}

func ErrorStr(err string) {
	if Config.UI.ErrorStyle == 0 {
		dialog.Message("%s", err).Title("Rblxutils error").Error()
	} else {
		beeep.Alert("rblxutils errored:", err, resources.ProgramLogo)
	}
}

func Notification(text string) {
	if Config.UI.ErrorStyle == 0 {
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
	panic(err)
}
