package configurator

import (
	"slices"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/clipboard"
	"github.com/aarzilli/nucular/label"
	"github.com/axellse/rblxutils/common"
)

func RenderServerHistory(win *nucular.Window) {
	win.Row(10).Dynamic(1)
	win.Label("Server History allows you to view servers you have been in.", label.Align("LC"))
	win.Row(10).Dynamic(1)
	win.Label("From there, you can rejoin, view the player list, etc.", label.Align("LC"))

	win.Row(20).Dynamic(1)
	win.CheckboxText("Server History", &common.Config.ServerHistoryEnabled)

	if common.Config.ServerHistoryEnabled {
		servers := make([]common.ServerData, len(common.State.ServerHistory))
		copy(servers, common.State.ServerHistory)
		slices.Reverse(servers)
		for i, server := range servers {
			if win.TreePushNamed(nucular.TreeTab, "ActualServerHistory-"+strconv.Itoa(i), server.GameData.Name, false) {
				RenderServerInfo(win, server, i)
				win.TreePop()
			}
		}

		win.Row(20).Static(100)
		if win.ButtonText("Clear History") && common.YesNo("Are you sure you want to clear your server history?") {
			common.State.ServerHistory = []common.ServerData{}
			common.WriteState()
		}
	}
}

func RenderServerInfo(win *nucular.Window, server common.ServerData, i int) {
	win.Row(207).Dynamic(1)
	win.Image(server.GameData.Thumbnail)
	if win.Input().Mouse.HoveringRect(win.LastWidgetBounds) {
		win.Tooltip(server.GameData.Name)
	}

	win.Row(20).Static(70, 190, 80)
	win.Label("Link Type: ", label.Align("LC"))
	UIStates.LinkTypeCombo = win.ComboSimple([]string{
		"Link to specific server",
		"Link to game",
	}, UIStates.LinkTypeCombo, 15)

	if win.ButtonText("Copy Link") {
		link := "https://api.axell.me/rblxutils/join/?p=" + strconv.Itoa(server.PlaceId)
		if UIStates.LinkTypeCombo == 0 {
			link += "&j=" + server.JobId
		}
		clipboard.Set(link)
		common.Notification("Copied to clipboard!")
	}

	//does not seem to work reliably
	/*if win.TreePushNamed(nucular.TreeTab, "Players-" + strconv.Itoa(i), "Players", false) {
		for _, v := range server.Players {
			win.Label(v, label.Align("LC"))
		}
		win.TreePop()
	}*/

	win.Row(10).Dynamic(1)
	win.Label("Server address: "+server.ServerAddress, label.Align("LC"))
	if win.Input().Mouse.HoveringRect(win.LastWidgetBounds) {
		win.Tooltip("This is the actual address the client talks to.")
	}
	win.Row(10).Dynamic(1)
	win.Label("UDMUX Server address: "+server.UDMUXAddress, label.Align("LC"))
	if win.Input().Mouse.HoveringRect(win.LastWidgetBounds) {
		win.Tooltip("UDMUX is some kind of ratelimiting system Roblox uses to prevent DOS attacks.")
	}
	win.Row(10).Dynamic(1)
	win.Label("RCC Server address: "+server.RCCAddress, label.Align("LC"))
	if win.Input().Mouse.HoveringRect(win.LastWidgetBounds) {
		win.Tooltip("This is the address of the RCC server that runs the game server.")
	}
}
