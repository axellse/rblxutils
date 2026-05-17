package configurator

import (
	"bytes"
	"image"
	"image/draw"
	"strconv"

	"axell.me/rblxutils/common"
	"axell.me/rblxutils/resources"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
)

func RenderLivePanel(win *nucular.Window, instance *common.Instance, i int) {
	win.Row(10).Dynamic(1)
	win.Label("Welcome to the live panel!", label.Align("LC"))
	win.Row(160).Dynamic(1)
	Oimg, _, err := image.Decode(bytes.NewReader(resources.ApartmentsLjmsImage))
	if err != nil {
		common.Error(err)
	}

	bounds := Oimg.Bounds()
	Nimg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(Nimg, Nimg.Bounds(), Oimg, bounds.Min, draw.Src)

	win.Image(Nimg)
	if win.Input().Mouse.HoveringRect(win.LastWidgetBounds) {
		win.Tooltip("Apartments, Late July Midsummer")
	}

	action := "Idling"
	if instance.ServerData.GameData.Name != "" {
		action = "Playing " + instance.ServerData.GameData.Name
	}

	win.Row(10).Dynamic(1)
	win.Label(action, label.Align("LC"))

	if action != "Idling" {
		win.Row(20).Dynamic(1)
		if win.TreePushNamed(nucular.TreeTab, "ServerInfo-" + strconv.Itoa(i), "Server Info", true) {
			RenderServerInfo(win, instance.ServerData, i)
			win.TreePop()
		}
	}

	win.Row(20).Dynamic(1)
	if win.TreePushNamed(nucular.TreeNode, "LivePanelNetwork-"+strconv.Itoa(i), "Network Stats", false) {
		win.Row(10).Dynamic(1)
		win.Label("Average AssetDelivery ModifyResponse Delay:"+strconv.FormatFloat(float64(UIStates.CurrentProxyStats.AvgModifyResponseAssetDeliveryDelayNs)/1_000_000, 'g', -1, 64)+"ms", label.Align("LC"))
		win.Row(10).Dynamic(1)
		win.Label("Average Cdn ModifyResponse Delay:"+strconv.FormatFloat(float64(UIStates.CurrentProxyStats.AvgModifyResponseCdnDelayNs)/1_000_000, 'g', -1, 64)+"ms", label.Align("LC"))
		win.Row(10).Dynamic(1)
		win.Label("Average Rewrite Delay:"+strconv.FormatFloat(float64(UIStates.CurrentProxyStats.AvgRewriteDelayNs)/1_000_000, 'g', -1, 64)+"ms", label.Align("LC"))
		win.TreePop()
	}
}