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
		idling := rpc.instance.ServerData.GameData.Name == "" || rpc.instance.ServerData.PlaceId == 0
		activity := client.Activity{}

		activity.SmallImage = "rblxutils_hires"
		activity.SmallText = "using rblxutils"
		if common.Config.DiscordRPC.ShowUserProfile && !idling {
			activity.SmallImage = rpc.instance.ServerData.HeadshotURL
			activity.SmallText = "Playing on " + rpc.instance.ServerData.User.DisplayName + " (@" + rpc.instance.ServerData.User.Name + ")"
		}

		if common.Config.DiscordRPC.AllowJoin && !idling {
			activity.Buttons = append(activity.Buttons, &client.Button{
				Label: "Join",
				Url: "https://api.axell.me/rblxutils/join/?p=" + strconv.Itoa(rpc.instance.ServerData.PlaceId) + "&j=" + rpc.instance.ServerData.JobId,
			})
		}

		activity.State = "using rblxutils"
		if !idling {
			switch stateI {
			case 1:
				activity.State = "server in " + rpc.instance.ServerData.Location.City + ", " + common.GetCountry(rpc.instance.ServerData.Location.Country)
			case 2:
				activity.State = "by " + rpc.instance.ServerData.GameData.Creator.Name
				if rpc.instance.ServerData.GameData.Creator.Verified {
					activity.State += "✅"
				}
			}
		}

		activity.Details = "idling"
		if !idling {
			activity.Details = rpc.instance.ServerData.GameData.Name
		}

		activity.LargeImage = "https://apis.axell.me/termusic/v1/idling-images/from-style/nature-1" //todo: move to api.axell.me/rblxutils
		activity.LargeText = "<3" //image credit could go here
		if !idling {
			activity.LargeImage = rpc.instance.ServerData.GameData.IconURL
			activity.LargeText = rpc.instance.ServerData.GameData.Name
		}

		activity.Timestamps = &client.Timestamps{
			Start: &rpc.instance.AllocationTime,
		}
		if !idling {
			activity.Timestamps.Start = &rpc.instance.ServerData.JoinTime
		}

		err = client.SetActivity(activity)
		if err != nil {
			fmt.Println("setActvity err, closing client.")
			client.Logout()
		}

		stateI++
		if stateI == 3 {
			stateI = 0
		}
	}
}