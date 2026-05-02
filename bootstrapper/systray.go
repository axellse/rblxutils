//go:build windows
// +build windows

package bootstrapper

import (
	"axell.me/rblxutils/common"
	"axell.me/rblxutils/configurator"
	"axell.me/rblxutils/lib/winsystray"
	"axell.me/rblxutils/resources"
)

func TieChanToCallback[T any](ch chan T, cb func(T)) {
	go func() {
		val := <-ch
		cb(val)
	}()
}

func StartSystray() {
	hwnd, err := winsystray.CreateTrayWindow("RblxutilsTray")
	if err != nil {
		common.FatalError(err)
	}

	for {
		ti, err := winsystray.NewTrayIcon(hwnd)
		if err != nil {
			common.FatalError(err)
		}

		ti.SetIconFromBytes(resources.ProgramLogoIco)
		ti.ProcessEvents()

		ti.Dispose()
		configurator.LaunchConfigurator(true)
	}


	


}
