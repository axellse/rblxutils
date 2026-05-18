package winsystray

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"crypto/rand"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

const TrayIconMsg = WM_APP + 1

// This function taken from https://github.com/getlantern/systray/blob/01dc414284987aa070498fafbcdac794657bf2e1/systray_windows.go#L772
func iconBytesToFilePath(iconBytes []byte) (string, error) {
	bh := md5.Sum(iconBytes)
	dataHash := hex.EncodeToString(bh[:])
	iconFilePath := filepath.Join(os.TempDir(), "systray_temp_icon_"+dataHash)

	if _, err := os.Stat(iconFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(iconFilePath, iconBytes, 0644); err != nil {
			return "", err
		}
	}
	return iconFilePath, nil
}

func getWndProc() func(hWnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	return func(hWnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
		switch msg {
		case TrayIconMsg:
			switch nmsg := LOWORD(uint32(lParam)); nmsg {
			case WM_LBUTTONDOWN:
				PostQuitMessage(0)
			}
		case WM_DESTROY:
			PostQuitMessage(0)
		default:
			r, _ := DefWindowProc(hWnd, msg, wParam, lParam)
			return r
		}
		return 0
	}
}

func CreateTrayWindow(class string) (uintptr, error) {
	hInstance, err := GetModuleHandle(nil)
	if err != nil {
		return 0, err
	}
	wndClass := windows.StringToUTF16Ptr(class)

	var wcex WNDCLASSEX

	wcex.CbSize = uint32(unsafe.Sizeof(wcex))
	wcex.LpfnWndProc = windows.NewCallback(getWndProc())
	wcex.HInstance = hInstance
	wcex.LpszClassName = wndClass
	if _, err := RegisterClassEx(&wcex); err != nil {
		return 0, err
	}

	hwnd, err := CreateWindowEx(
		0,
		wndClass,
		windows.StringToUTF16Ptr("Winsystray system tray"),
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		0,
		0,
		hInstance,
		nil)
	if err != nil {
		return 0, err
	}

	return hwnd, nil
}

func newGUID() GUID {
	var buf [16]byte
	rand.Read(buf[:])
	return *(*GUID)(unsafe.Pointer(&buf[0]))
}

type TrayIcon struct {
	hwnd      uintptr
	guid      GUID
	eventChan chan string
}

func (ti *TrayIcon) ProcessEvents() {
	var msg MSG
	for {
		res, err := GetMessage(&msg, 0, 0, 0)
		if err != nil || res == 0 {
			return
		}
		_, _ = TranslateMessage(&msg)
		_, _ = DispatchMessage(&msg)
	}
}

func NewTrayIcon(trayWindowHwnd uintptr) (*TrayIcon, error) {
	ti := &TrayIcon{guid: newGUID()}

	ti.hwnd = trayWindowHwnd
	data := ti.initData()
	data.UFlags |= NIF_MESSAGE
	data.UCallbackMessage = TrayIconMsg
	if _, err := Shell_NotifyIcon(NIM_ADD, data); err != nil {
		return nil, err
	}
	return ti, nil
}

func (ti *TrayIcon) initData() *NOTIFYICONDATA {
	var data NOTIFYICONDATA
	data.CbSize = uint32(unsafe.Sizeof(data))
	data.UFlags = NIF_GUID
	data.HWnd = ti.hwnd
	data.GUIDItem = ti.guid
	return &data
}

func (ti *TrayIcon) Dispose() error {
	_, err := Shell_NotifyIcon(NIM_DELETE, ti.initData())
	return err
}

func (ti *TrayIcon) SetIconFromFile(iconFilename string) error {
	icon, err := LoadImage(
		0,
		windows.StringToUTF16Ptr(iconFilename),
		IMAGE_ICON,
		0,
		0,
		LR_DEFAULTSIZE|LR_LOADFROMFILE)

	if err != nil {
		panic(err)
	}
	data := ti.initData()
	data.UFlags |= NIF_ICON
	data.HIcon = icon
	_, err = Shell_NotifyIcon(NIM_MODIFY, data)
	return err
}

func (ti *TrayIcon) SetIconFromBytes(iconBytes []byte) error {

	iconFilePath, err := iconBytesToFilePath(iconBytes)
	if err != nil {
		panic(fmt.Sprintf("Unable to write icon data to temp file: %v", err))
	}
	err = ti.SetIconFromFile(iconFilePath)
	return err
}

func (ti *TrayIcon) SetTooltip(tooltip string) error {
	data := ti.initData()
	data.UFlags |= NIF_TIP
	copy(data.SzTip[:], windows.StringToUTF16(tooltip))
	_, err := Shell_NotifyIcon(NIM_MODIFY, data)
	return err
}

func (ti *TrayIcon) ShowBalloonNotification(title, text string) error {
	data := ti.initData()
	data.UFlags |= NIF_INFO
	if title != "" {
		copy(data.SzInfoTitle[:], windows.StringToUTF16(title))
	}
	copy(data.SzInfo[:], windows.StringToUTF16(text))
	_, err := Shell_NotifyIcon(NIM_MODIFY, data)
	return err
}