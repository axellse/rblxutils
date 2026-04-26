package configurator

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"time"

	"axell.me/rblxutils/common"
	"axell.me/rblxutils/resources"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/style"
)

var UIStates UIState //its giving yandere simulator or whatever its called

func LaunchUI() {
	wnd := nucular.NewMasterWindowSize(nucular.WindowHelp, "Rblxutils", image.Point{400, 500}, renderWindow)
	wnd.OnClose(func() {
		if common.ChangesMade() && common.YesNo("You have unsaved changes. Do you want to save them before closing?") {
			fmt.Println("saving changes...")
			err := common.WriteConfiguration()
			if err != nil {
				common.FatalError(err)
			}
		}
	})
	wnd.SetStyle(style.FromTheme(style.Theme(common.Config.Misc.Theme), 1.0))
	wnd.Main()
}

func renderWindow(win *nucular.Window) {
	windowWidth := win.Bounds.W - 4

	win.Row(5).Dynamic(1)
	welcome := win.TreePushNamed(nucular.TreeTab, "Welcome", "Welcome! (Update v1.2.1)", !common.Config.Misc.DisableWelcomeScreen)
	if welcome {
		win.Row(10).Dynamic(1)
		win.Label("Welcome to rblxutils!", label.Align("LC"))
		win.Row(120).Dynamic(1)

		Oimg, _, err := image.Decode(bytes.NewReader(resources.WelcomeCatImage))
		if err != nil {
			common.Error(err)
		}

		bounds := Oimg.Bounds()
		Nimg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(Nimg, Nimg.Bounds(), Oimg, bounds.Min, draw.Src)

		win.Image(Nimg)
		win.Row(10).Dynamic(1)
		win.Label("This is a program that provides a collection of ways", label.Align("LC"))
		win.Row(10).Dynamic(1)
		win.Label("to extend Roblox, many of which are experimental.", label.Align("LC"))

		win.Row(25).Dynamic(1)
		win.CheckboxText("Hide this welcome page next time", &common.Config.Misc.DisableWelcomeScreen)
		win.TreePop()
	}

	/*win.Row(5).Dynamic(1)
	ptt := win.TreePushNamed(nucular.TreeTab, "PTT", "Push-to-talk (PTT)", false)
	if ptt {
		win.Row(20).Dynamic(1)
		win.CheckboxText("Enabled", &Config.PTTEnabled)

		if (Config.PTTEnabled) {
			win.Row(20).Static(windowWidth-22-150, 150)
			win.Label("Activation key: ", "LL")
			win.ButtonText("(Click to record)")
			win.Row(20).Static(windowWidth-22-150, 150)
			win.Label("Microphone:", "LC")
			win.ComboSimple([]string{
				"dummy mic device 1",
				"shitty soundboard trademark",
			}, 0, 15)
		}
		win.TreePop()
	}*/

	mods := win.TreePushNamed(nucular.TreeTab, "Mods", "Client Modifications", false)
	if mods {
		win.Row(10).Dynamic(1)
		win.Label("Mods can change local files, assets or both.", label.Align("LC")) //Rblxutils can apply and manage traditional mods that modify local files (eg. Bloxstap mods), but also mods that modify assets downloaded by the client from Roblox's CDN. Rlbxutils uses the Roblox Modding Interchange Format (rcmf) but can also automatically convert from a couple of other common formats.
		win.Row(10).Dynamic(1)
		win.Label("The following formats are supported: rcmf, zip", label.Align("LC"))

		win.Row(20).Static(150)
		if win.ButtonText("Import from file") {
			ModImportWizard()
		}

		win.Row(5).Dynamic(1)
		win.Spacing(1)
		win.Row(15).Dynamic(1)
		win.Label("Manage mods:", label.Align("LC"))

		for i := range common.Config.Mods {
			if len(common.Config.Mods) <= i {
				continue
			}
			manage := win.TreePush(nucular.TreeTab, common.Config.Mods[i].Name, false)
			if manage {
				win.Row(20).Dynamic(1)
				win.CheckboxText("Enabled", &common.Config.Mods[i].Enabled)
				win.Row(20).Dynamic(3)
				if win.ButtonText("Remove") {
					common.Config.Mods = append(common.Config.Mods[:i], common.Config.Mods[i+1:]...)
				}

				if win.ButtonText("Download rcmf") {
					SaveModWizard(common.Config.Mods[i])
				}
				win.TreePop()
			}
		}

		win.TreePop()
	}

	misc := win.TreePushNamed(nucular.TreeTab, "Misc", "Miscellaneous Options", false)
	if misc {
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
		win.TreePop()
	}

	win.Row(20).Static(windowWidth-88, 70)
	win.Label(UIStates.SaveLabel, label.Align("LC"))
	if common.ChangesMade() {
		UIStates.SaveLabel = "You have unsaved changes!"
	} else if UIStates.SaveLabelInauguration.Add(3 * time.Second).Before(time.Now()) {
		UIStates.SaveLabel = ""
	}
	
	if win.ButtonText("Save") {
		if common.ChangesMade() {
			if !reflect.DeepEqual(common.Config.Mods, common.ConfigFileState.Mods) {
				common.LoadState()
				common.State.RequiresModApplication = true
				err := common.WriteState()
				if err != nil {
					common.Error(err)
				}
			}
			err := common.WriteConfiguration()
			if err != nil {
				UIStates.SaveLabel = err.Error()
			} else {
				UIStates.SaveLabel = "Saved changes"
			}
		} else {
			UIStates.SaveLabel = "No changes made"
		}
		UIStates.SaveLabelInauguration = time.Now()
	}
}

type UIState struct {
	SaveLabel             string
	SaveLabelInauguration time.Time
}
