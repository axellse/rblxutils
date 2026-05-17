package configurator

import (
	"axell.me/rblxutils/common"
	"axell.me/rblxutils/uninstaller"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/style"
)

func RenderMisc(win *nucular.Window, windowWidth int) {
	win.Row(20).Static(100, 150, windowWidth-100-150-50-30, 50)
	win.Label("Error style: ", label.Align("LC"))

	common.Config.Misc.ErrorStyle = win.ComboSimple([]string{
		"Dialog box",
		"Notification",
	}, common.Config.Misc.ErrorStyle, 15)
	win.Spacing(1)

	if win.ButtonText("Test") {
		common.ErrorStr("This is a test of the error dialog/notification!")
	}

	win.Row(20).Static(100, 150)
	win.Label("UI Theme: ", label.Align("LC"))

	newTheme := win.ComboSimple([]string{
		"Default",
		"White",
		"Red",
		"Dark",
	}, common.Config.Misc.Theme, 15)

	if common.Config.Misc.Theme != newTheme {
		win.Master().SetStyle(style.FromTheme(style.Theme(newTheme), 1.0))
	}
	common.Config.Misc.Theme = newTheme

	win.Row(25).Dynamic(1)
	win.CheckboxText("Disable launch notification", &common.Config.Misc.DisableLaunchNotification)

	win.Row(25).Dynamic(1)
	win.CheckboxText("Enable Desktop Shortcut", &common.Config.Misc.DesktopShortcutEnabled)

	win.Row(20).Static(200, 150)
	if !common.State.RequiresModApplication && win.ButtonText("Force Mod Reapplication") {
		common.State.RequiresModApplication = true
		err := common.WriteState()
		if err != nil {
			common.FatalError(err)
		}
	}

	if !UIStates.LiveMode && win.ButtonText("Uninstall rblxutils") {
		uninstaller.LaunchUninstaller()
	}
}
