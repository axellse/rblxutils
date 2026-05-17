package bootstrapper

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"math/rand/v2"

	"axell.me/rblxutils/common"
	"axell.me/rblxutils/resources"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/aarzilli/nucular/style"
)

var DockSplitWindowFlags = nucular.WindowBorder|nucular.WindowMovable|nucular.WindowScalable|nucular.WindowTitle|nucular.WindowNonmodal|nucular.WindowNoScrollbar

func LaunchUI() {
	wnd := nucular.NewMasterWindowSize(nucular.WindowHelp, "rblxutils bootstrapper", image.Point{700, 400}, func(w *nucular.Window) {})
	wnd.OnClose(func() {
		fmt.Println("youre not terminating the program just because the window was closed lil bro")
	})
	UiState.Update = wnd.Changed
	UiState.CloseWindow = wnd.Close

	go func() {
		dockSplit := wnd.ResetWindows()
		l, r := dockSplit.Split(false, 0)

		l.Open("Bootstrapper Log", DockSplitWindowFlags, rect.Rect{W: 300, H: 300}, true, renderWindowLog)
		r.Open("Bootstrapper", DockSplitWindowFlags, rect.Rect{W: 300, H: 300}, true, renderWindowProgress)
		wnd.SetStyle(style.FromTheme(style.Theme(common.Config.UI.Theme), 1.0))
		//go common.SetWindowStyle()
		wnd.Main()
	}()
}

type BootstrapperUIState struct {
	LogOutput *nucular.TextEditor
	Progress int
	CurrentOperation string
	Update func()
	CloseWindow func()
}

var UiState = BootstrapperUIState{
	LogOutput: &nucular.TextEditor{},
	Progress: 0,
	CurrentOperation: "Waiting for bootstrapper...",
}

func renderWindowLog(win *nucular.Window) {
	UiState.LogOutput.Buffer = []rune(BootstrapperLog.String())
	UiState.LogOutput.SingleLine = false
	UiState.LogOutput.Active = true
	UiState.LogOutput.Flags = nucular.EditReadOnly | nucular.EditMultiline | nucular.EditSelectable | nucular.EditClipboard
	win.Row(0).Dynamic(1)
	UiState.LogOutput.Edit(win)
}

var randPicInt = rand.IntN(6)

func GetBootstrapperImage() []byte {
	if common.Config.UI.BootstrapperImage == 1 {
		return resources.ProgramLogo
	}
	
	switch randPicInt {
	case 0:
		return resources.CatPic1
	case 1:
		return resources.CatPic2
	case 2:
		return resources.CatPic3
	case 3:
		return resources.CatPic4
	case 4:
		return resources.CatPic5
	case 5:
		return resources.CatPic6
	case 6:
		return resources.CatPic7
	}

	return resources.CatPic1 //never happens
}

func renderWindowProgress(win *nucular.Window) {
	width := win.Bounds.W
	height := win.Bounds.H
	win.Row(height/2 - ((15 + 20 + 100)/2) - 30).Dynamic(1)
	win.Spacing(1)
	win.Row(100).Static((width -100 -29) /2, 100, (width -100 -29) /2)
	
	Oimg, _, err := image.Decode(bytes.NewReader(GetBootstrapperImage()))
	if err != nil {
		common.FatalError(err)
	}
	bounds := Oimg.Bounds()
	Nimg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(Nimg, Nimg.Bounds(), Oimg, bounds.Min, draw.Src)
	win.Spacing(1)
	win.Image(Nimg)
	win.Spacing(1)

	win.Row(20).Static((width -150 -29) /2, 150, (width -150 -29) /2)
	win.Spacing(1)
	win.Progress(&UiState.Progress, 100, false)
	win.Spacing(1)
	win.Row(15).Dynamic(1)
	win.Label(UiState.CurrentOperation, "CC")
}