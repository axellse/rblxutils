//go:build windows

package bootstrapper

import (
	"time"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/configurator"
	"github.com/axellse/rblxutils/winsystray"
	"github.com/axellse/rblxutils/resources"
	"github.com/gen2brain/beeep"
)

func TieChanToCallback[T any](ch chan T, cb func(T)) {
	go func() {
		val := <-ch
		cb(val)
	}()
}

func StartSystray() {
	if !common.Config.Misc.DisableLaunchNotification {
		beeep.Alert("rblxutils", "you may now access rblxutils from the systray!", resources.ProgramLogo)
	}

	hwnd, err := winsystray.CreateTrayWindow("RblxutilsTray")
	if err != nil {
		common.FatalError(err)
	}

	windowCloseCh := make(chan struct{})
	for {
		ti, err := winsystray.NewTrayIcon(hwnd)
		if err != nil {
			common.FatalError(err)
		}
		ti.SetIconFromBytes(resources.ProgramLogoIco)

		ti.ProcessEvents()
		err = ti.Dispose()
		if err != nil {
			common.FatalError(err)
		}
		go configurator.LaunchConfigurator(windowCloseCh, GlobalInman)
		<-windowCloseCh
		time.Sleep(100 * time.Millisecond)
	}
}
