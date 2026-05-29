package configurator

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/axellse/rblxutils/common"
)

func RenderDiscordRPC(win *nucular.Window) {
	win.Row(25).Dynamic(1)
	win.Label("This is essentially a recreation of Bloxstrap's RPC.", label.Align("LC"))
	win.Row(25).Dynamic(1)
	win.CheckboxText("Discord Rich Presence", &common.Config.DiscordRPC.Enabled)
	if common.Config.DiscordRPC.Enabled {
		win.Row(25).Dynamic(1)
		win.CheckboxText("Join server button", &common.Config.DiscordRPC.AllowJoin)
		win.Row(25).Dynamic(1)
		win.CheckboxText("Show account on profile", &common.Config.DiscordRPC.ShowUserProfile)
	}

}