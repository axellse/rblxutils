package bootstrapper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/bootstrapper/rich-go/client"
)

type DiscordRPC struct {
	instance *common.Instance
}

func (rpc *DiscordRPC) RunRPC() {
	if !common.Config.DiscordRPC.Enabled {
		return
	}

	fmt.Println("now running discord rpc")
	stateI := 0
	for range time.Tick(2 * time.Second) {
		err := client.Login("1509961963613065379")
		if err != nil && strings.Contains(err.Error(), "Timed out") {
			continue
		} else if err != nil {
			common.FatalError(err)
		}

		//TODO: do something when idling
		if rpc.instance.ServerData.GameData.Name == "" || rpc.instance.ServerData.PlaceId == 0 {
			continue
		}

		smallImage := "rblxutils_hires"
		smallImageText := "using rblxutils"
		if common.Config.DiscordRPC.ShowUserProfile {
			smallImage = rpc.instance.ServerData.HeadshotURL
			smallImageText = "Playing on " + rpc.instance.ServerData.User.DisplayName + " (@" + rpc.instance.ServerData.User.Name + ")"
		}

		buttons := []*client.Button{}
		if common.Config.DiscordRPC.AllowJoin {
			buttons = append(buttons, &client.Button{
				Label: "Join",
				Url: "https://api.axell.me/rblxutils/join/?p=" + strconv.Itoa(rpc.instance.ServerData.PlaceId) + "&j=" + rpc.instance.ServerData.JobId,
			})
		}

		state := "using rblxutils"
		switch stateI {
		case 1:
			state = "Server in " + rpc.instance.ServerData.Location.City + ", " + common.GetCountry(rpc.instance.ServerData.Location.Country)
		case 2:
			state = "by " + rpc.instance.ServerData.GameData.Creator.Name
			if rpc.instance.ServerData.GameData.Creator.Verified {
				state += "✅"
			}
		case 3:
			state = strconv.Itoa(len(rpc.instance.ServerData.Players)) + "/" + strconv.Itoa(rpc.instance.ServerData.GameData.MaxPlayers) + " Players in server"
		}

		err = client.SetActivity(client.Activity{
			Details: rpc.instance.ServerData.GameData.Name,
			State: state,
			LargeImage: rpc.instance.ServerData.GameData.IconURL,
			LargeText: rpc.instance.ServerData.GameData.Name,
			SmallImage: smallImage,
			SmallText: smallImageText,
			Timestamps: &client.Timestamps{
				Start: &rpc.instance.ServerData.JoinTime,
			},
			Buttons: buttons,
		})

		if err != nil {
			fmt.Println("setActvity err, closing client.")
			client.Logout()
		}

		stateI++
		if stateI == 4 {
			stateI = 0
		}
	}
}