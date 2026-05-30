package configurator

import (
	"github.com/aarzilli/nucular"
	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/uninstaller"
)

func RenderMisc(win *nucular.Window, windowWidth int) {
	win.Row(25).Dynamic(1)
	win.CheckboxText("Disable launch notification", &common.Config.Misc.DisableLaunchNotification)

	win.Row(25).Dynamic(1)
	win.CheckboxText("Enable Desktop Shortcut", &common.Config.Misc.DesktopShortcutEnabled)
	win.Row(25).Dynamic(1)
	win.CheckboxText("Keep running when all Roblox instances close", &common.Config.Misc.InmanStayAlive)

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
