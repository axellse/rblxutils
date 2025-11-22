package main

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
)

//Ui state: (quite horrible)
var UIPTTEnabled bool
var HideWelcome bool

func LaunchUI() {
	wnd := nucular.NewMasterWindowSize(nucular.WindowHelp, ProgramName, image.Point{400, 500},renderWindow)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 1.0))
	wnd.Main()
}

func renderWindow(win *nucular.Window) {
	windowWidth := win.Bounds.W - 4
	
	win.Row(5).Dynamic(1)
	welcome := win.TreePushNamed(nucular.TreeTab, "Welcome", "Welcome! (Update v1.2.1)", true)
	if welcome {
		win.Label("Welcome to rblxutils!", "LL")
		win.CheckboxText("hide this stupid welcome thing next time", &HideWelcome)
		win.TreePop()
	}

	win.Row(5).Dynamic(1)
	ptt := win.TreePushNamed(nucular.TreeTab, "PTT", "Push-to-talk (PTT)", false)
	if ptt {
		win.Row(20).Dynamic(1)
		win.CheckboxText("Enabled", &UIPTTEnabled)

		if (UIPTTEnabled) {
			win.Row(20).Dynamic(1)
			win.Label("Coming soon!", "LL")
			/*
			win.Row(20).Static(windowWidth-22-150, 150)
			win.Label("Activation key: ", "LL")
			win.ButtonText("(Click to record)")
			win.Row(20).Static(windowWidth-22-150, 150)
			win.Label("Microphone:", "LC")
			win.ComboSimple([]string{
				"dummy mic device 1",
				"shitty soundboard trademark",
			}, 0, 15)*/
		}
		win.TreePop()
	}

	mods := win.TreePushNamed(nucular.TreeTab, "Mods", "Client Modifications", false)
	if mods {
		manage := win.TreePush(nucular.TreeNode, "Manage Installed Modifications", false)
		if manage {
			example := win.TreePush(nucular.TreeNode, "Loud Footsteps", false)
			if example {
				win.Row(20).Dynamic(1)
				win.CheckboxText("Enabled", &UIPTTEnabled)
				win.Row(20).Dynamic(3)
				win.ButtonText("Uninstall")
				win.ButtonText("Download manifest")
				win.TreePop()
			}
			win.TreePop()
		}
		win.TreePop()
	}

	env := win.TreePushNamed(nucular.TreeTab, "Environment", "Configure Environment", false)
	if env {
		win.Row(20).Static(windowWidth-22-150, 150)
		win.Label("Use Preset:", "LC")
		win.ComboSimple([]string{
			"Stock Bootstrapper",
			"Bloxstrap",
			"Fishtrap",
		}, 0, 15)
		win.TreePop()

		env := win.TreePushNamed(nucular.TreeTab, "adv_env", "Advanced", false)
		if env {
			win.Row(20).Static(windowWidth-22-150, 150)
			win.Label("Use Preset:", "LC")
			win.ComboSimple([]string{
				"Stock Bootstrapper",
				"Bloxstrap",
				"Fishtrap",
			}, 0, 15)
			win.TreePop()
		}
	}

	win.Row(20).Static(windowWidth - 88, 70)
	win.Spacing(1)
	win.ButtonText("Save")
}