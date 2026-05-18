package configurator

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/style"
	"github.com/axellse/rblxutils/common"
)

func RenderUiAndStyling(win *nucular.Window, windowWidth int) {
	win.Row(20).Static(150, 150, windowWidth-150-150-50-30, 50)
	win.Label("Error style: ", label.Align("LC"))

	common.Config.UI.ErrorStyle = win.ComboSimple([]string{
		"Dialog box",
		"Notification",
	}, common.Config.UI.ErrorStyle, 15)
	win.Spacing(1)

	if win.ButtonText("Test") {
		common.ErrorStr("This is a test of the error dialog/notification!")
	}

	win.Row(20).Static(150, 150)
	win.Label("UI Theme: ", label.Align("LC"))

	newTheme := win.ComboSimple([]string{
		"Default",
		"White",
		"Red",
		"Dark",
	}, common.Config.UI.Theme, 15)

	win.Row(20).Static(150, 150)
	win.Label("Bootstrapper image: ", label.Align("LC"))

	common.Config.UI.BootstrapperImage = win.ComboSimple([]string{
		"Random cat",
		"Rblxutils",
	}, common.Config.UI.BootstrapperImage, 15)

	if common.Config.UI.Theme != newTheme {
		win.Master().SetStyle(style.FromTheme(style.Theme(newTheme), 1.0))
	}
	common.Config.UI.Theme = newTheme
}
