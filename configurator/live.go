package configurator

import (
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/axellse/rblxutils/common"
)

func RenderLivePanel(win *nucular.Window, instance *common.Instance, i int) {
	win.Row(10).Dynamic(1)
	win.Label("Welcome to the live panel!", label.Align("LC"))

	if instance.ServerData.GameData.Name != "" {
		RenderServerInfo(win, instance.ServerData, i, true)
	} else {
		win.Row(10).Dynamic(1)
		win.Label("Idling", label.Align("LC"))
	}

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
