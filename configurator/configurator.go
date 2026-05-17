package configurator

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"axell.me/rblxutils/common"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/style"
)

var UiStateMutex sync.Mutex
var UIStates UIState //its giving yandere simulator or whatever its called

func LaunchConfigurator(ch chan struct{}, inman *common.Inman) {
	if UIStates.Active {
		fmt.Println("ui active, sorry boss")
		return
	}
	LaunchUI(ch, inman)
}

func LaunchUI(ch chan struct{}, inman *common.Inman) {
	UIStates = UIState{}
	UIStates.LiveMode = inman != nil
	UIStates.Active = true
	UIStates.UpdateTitle = "Loading updates..."
	UIStates.Inman = inman
	wnd := nucular.NewMasterWindowSize(nucular.WindowHelp, "Rblxutils", image.Point{400, 500}, renderWindow)

	UIStates.Update = wnd.Changed
	wnd.OnClose(func() {
		UiStateMutex.Lock()
		UIStates.Active = false
		UiStateMutex.Unlock()
		if common.ChangesMade() && common.YesNo("You have unsaved changes. Do you want to save them before closing?") {
			fmt.Println("saving changes...")
			err := common.WriteConfiguration()
			if err != nil {
				common.FatalError(err)
			}
		}

		if !UIStates.LiveMode {
			common.AssertDesktopShortcut()
			os.Exit(0)
		} else {
			ch <- struct{}{}
		}
	})
	go FetchWelcome()

	wnd.SetStyle(style.FromTheme(style.Theme(common.Config.UI.Theme), 1.0))
	wnd.Main()
}

func LoadImageUI(ba []byte) *image.RGBA {
	img, _, err := image.Decode(bytes.NewReader(ba))
	if err != nil {
		common.Error(err)
	}

	bounds := img.Bounds()
	rgbImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgbImg, rgbImg.Bounds(), img, bounds.Min, draw.Src)
	return rgbImg
}

func renderWindow(win *nucular.Window) {
	windowWidth := int(float64(win.Bounds.W-4) * 1 / 1)

	win.Row(5).Dynamic(1)
	if !UIStates.LiveMode {
		if win.TreePushNamed(nucular.TreeTab, "QuickLaunch", "Quick Launch Options", true) {
			RenderQuickLaunch(win)
			win.TreePop()
		}
		if win.TreePushNamed(nucular.TreeTab, "Welcome", UIStates.UpdateTitle, false) {
			RenderWelcome(win)
			win.TreePop()
		}
		if UIStates.OpenUpdatePanel {
			win.TreeOpen("Welcome")
			UiStateMutex.Lock()
			UIStates.OpenUpdatePanel = false
			UiStateMutex.Unlock()
		}
	} else {
		for i, instance := range UIStates.Inman.GetInstances() {
			panelName := "• LIVE Panel"
			if instance.ServerData.GameData.Name != "" {
				panelName += " - " + instance.ServerData.GameData.Name
			}

			if win.TreePushNamed(nucular.TreeTab, "LivePanel-"+strconv.Itoa(i), panelName, true) {
				RenderLivePanel(win, instance, i)
				win.TreePop()
			}
		}

	}

	if win.TreePushNamed(nucular.TreeTab, "Mods", "Mods", false) {
		RenderMods(win)
		win.TreePop()
	}

	if win.TreePushNamed(nucular.TreeTab, "ServerHistory", "Server History", false) {
		RenderServerHistory(win)
		win.TreePop()
	}

	if win.TreePushNamed(nucular.TreeTab, "UI", "UI and Styling", false) {
		RenderUiAndStyling(win, windowWidth)
		win.TreePop()
	}

	if win.TreePushNamed(nucular.TreeTab, "Misc", "Miscellaneous Options", false) {
		RenderMisc(win, windowWidth)
		win.TreePop()
	}

	win.Row(20).Static(windowWidth-88, 70)
	win.Label(UIStates.SaveLabel, label.Align("LC"))
	if common.ChangesMade() {
		UiStateMutex.Lock()
		UIStates.SaveLabel = "You have unsaved changes!"
		UiStateMutex.Unlock()
	} else if UIStates.SaveLabelInauguration.Add(3 * time.Second).Before(time.Now()) {
		UiStateMutex.Lock()
		UIStates.SaveLabel = ""
		UiStateMutex.Unlock()
	}

	if win.ButtonText("Save") {
		UiStateMutex.Lock()
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
		UiStateMutex.Unlock()
	}
}

type UIState struct {
	SaveLabel             string
	SaveLabelInauguration time.Time
	LiveMode              bool
	Active                bool
	CurrentProxyStats     common.ProxyStats
	Update                func()
	Inman                 *common.Inman
	LinkTypeCombo         int
	OpenUpdatePanel   bool
	UpdateText            string
	UpdateTitle           string
	UpdateImage           *image.RGBA
	UpdateTextCols int
}
