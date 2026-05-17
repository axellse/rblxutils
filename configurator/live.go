package configurator

import (
	"strconv"

	"axell.me/rblxutils/common"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
)

func RenderLivePanel(win *nucular.Window, instance *common.Instance, i int) {
	win.Row(10).Dynamic(1)
	win.Label("Welcome to the live panel!", label.Align("LC"))


	action := "Idling"
	if instance.ServerData.GameData.Name != "" {
		action = "Playing " + instance.ServerData.GameData.Name
	}

	win.Row(10).Dynamic(1)
	win.Label(action, label.Align("LC"))

	if action != "Idling" {
		RenderServerInfo(win, instance.ServerData, i)
	}

	win.Row(20).Dynamic(1)
	if win.TreePushNamed(nucular.TreeTab, "LivePanelNetwork-"+strconv.Itoa(i), "Network Stats", false) {
		win.Row(10).Dynamic(1)
		win.Label("Average AssetDelivery ModifyResponse Delay:"+strconv.FormatFloat(float64(UIStates.CurrentProxyStats.AvgModifyResponseAssetDeliveryDelayNs)/1_000_000, 'g', -1, 64)+"ms", label.Align("LC"))
		win.Row(10).Dynamic(1)
		win.Label("Average Cdn ModifyResponse Delay:"+strconv.FormatFloat(float64(UIStates.CurrentProxyStats.AvgModifyResponseCdnDelayNs)/1_000_000, 'g', -1, 64)+"ms", label.Align("LC"))
		win.Row(10).Dynamic(1)
		win.Label("Average Rewrite Delay:"+strconv.FormatFloat(float64(UIStates.CurrentProxyStats.AvgRewriteDelayNs)/1_000_000, 'g', -1, 64)+"ms", label.Align("LC"))
		win.TreePop()
	}
}