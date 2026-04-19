package common

import (
	"fmt"
	"syscall"
	"time"

	"github.com/TheTitanrain/w32"
	"golang.org/x/sys/windows"
)

//work in progress
func SetWindowStyle() {
	var foundMatch uintptr
	for {
		time.Sleep(20 * time.Millisecond)
		fmt.Println("trying to find window...")
		ourProcessId := windows.GetCurrentProcessId()
		cb := syscall.NewCallback(func(hwnd windows.Handle, lparam uintptr) uintptr {
			var winProcessId uint32
			_, err := windows.GetWindowThreadProcessId(windows.HWND(hwnd), &winProcessId)
			if err != nil {
				Error(err)
			}

			if ourProcessId == winProcessId {
				fmt.Println("match!", winProcessId, "is us,", hwnd, "is our window!")
				foundMatch = uintptr(hwnd)
				return 0
			}

			return 1
		})

		err := windows.EnumWindows(cb, nil)
		if err != nil && foundMatch == 0 {
			Error(err)
		}

		if foundMatch != 0 {
			break
		} else {
			fmt.Println("no match :(")
		}
		time.Sleep(100 * time.Millisecond)
	}

	style := w32.GetWindowLongPtr(w32.HWND(foundMatch), w32.GWL_STYLE)
	style = style &^ uintptr(w32.WS_CAPTION)
	w32.SetWindowLongPtr(w32.HWND(foundMatch), w32.GWL_STYLE, style)

	w32.SetWindowPos(
		w32.HWND(foundMatch),
		0,
		0,
		0,
		0,
		0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_FRAMECHANGED,
	)

	w32.SetWindowPos(w32.HWND(foundMatch), w32.HWND_TOPMOST, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE)
}
