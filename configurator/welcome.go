package configurator

import (
	"image/color"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/richtext"
	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/resources"
)

func FetchWelcome() {
	uResp := common.CheckForUpdates()

	resp, err := http.Get("https://api.axell.me/rblxutils/v1/updates/image")
	if err != nil {
		common.FatalError(err)
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		common.FatalError(err)
	}
	img := common.LoadImageUI(ba, 0, 0)

	UiStateMutex.Lock()
	UIStates.UpdateText = uResp.Text
	UIStates.UpdateTitle = uResp.Title
	UIStates.OriginalUpdateImage = img
	UIStates.UpdateImage = img
	UIStates.UpdateTextCols = uResp.TextCols
	UIStates.OpenUpdatePanel = true
	UiStateMutex.Unlock()
}

var updateText = richtext.New(richtext.AutoWrap)
var linkPattern = regexp.MustCompile(`\[(.*)\]\[(.*)\]`)

func RenderWelcome(win *nucular.Window) {
	if UIStates.UpdateText == "" {
		win.Row(10).Dynamic(1)
		win.Label("Fetching content...", label.Align("LC"))
		return
	}
	win.Row(UIStates.UpdateImageHeight).Dynamic(1)
	if UIStates.UpdateImageWidth != common.CalcWidth(win) {
		img, size := common.AutoResize(UIStates.OriginalUpdateImage, win)
		UiStateMutex.Lock()
		UIStates.UpdateImage = img
		UIStates.UpdateImageHeight = size
		UiStateMutex.Unlock()
	}
	win.Image(UIStates.UpdateImage)

	win.Row(UIStates.UpdateTextCols * 20).Dynamic(1)
	cstr := updateText.Widget(win, false)
	if cstr != nil {
		links := linkPattern.FindAllStringSubmatch(UIStates.UpdateText, -1)
		for i, text := range linkPattern.Split(UIStates.UpdateText, -1) {
			cstr.SetStyle(richtext.TextStyle{Color: win.Master().Style().Text.Color})
			cstr.Text(text)

			if i >= len(links) {
				continue
			}

			cstr.SetStyle(richtext.TextStyle{Color: color.RGBA{23, 96, 255, 255}})
			if cstr.Link(links[i][1], color.RGBA{43, 109, 252, 255}, nil) {
				err := exec.Command("cmd", "/c", "start", strings.ReplaceAll(links[i][2], "&", "^&")).Run()
				if err != nil {
					common.FatalError(err)
				}
			}
		}

		/*for _, sel := range linkSels {
			cstr.SetStyleForSel(sel, richtext.TextStyle{Color: color.RGBA{23, 96, 255, 255}})
		}*/
		cstr.End()
	}

	win.Row(25).Dynamic(1)
	win.CheckboxText("Hide this welcome page next time", &common.Config.UI.DisableWelcomeScreen)
}

func RenderQuickLaunch(win *nucular.Window) {
	win.Row(50).Dynamic(2)
	if win.Button(label.IT(common.LoadImageUI(resources.RobloxRLogo, 0, 0), "Roblox App", label.Align("CC")), false) {
		LaunchRoblox(0, "", win.Close)
	}
	if win.Button(label.IT(common.LoadImageUI(resources.BuilderClubLogo, 0, 0), "Mod Studio", label.Align("CC")), false) {
		common.Notification("🚧🚧UNDER CONSTRUCTION OKAY IM WORKING ON IT🚧🚧")
	}
}
