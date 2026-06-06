//go:build windows
// +build windows

package shortcut

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// Shortcut the shortcut (.lnk file) property struct
type Shortcut struct {
	// Shortcut (.lnk file) path
	ShortcutPath string
	// Shortcut target: a file path or a website
	Target string
	// Shortcut icon path, default: "%SystemRoot%\\System32\\SHELL32.dll,0"
	IconLocation string
	// Arguments of shortcut
	Arguments string
	// Description of shortcut
	Description string
	// Hotkey of shortcut
	Hotkey string
	// WindowStyle, "1"(default) for default size and location; "3" for maximized window; "7" for minimized window
	WindowStyle string
	// Working directory of shortcut
	WorkingDirectory string
}

// CreateShortcut create a desktop shortcut with name, target and shortcut type
// target is a file or website
// if iconPath is empty string, icon would be "%SystemRoot%\\System32\\SHELL32.dll,0"
func CreateShortcut(name, target, iconPath string, sType ShortcutType) error {
	dir, err := GetShortcutTypeDirectory(sType)
	if err != nil {
		return err
	}
	shortcutPath := filepath.Join(dir, name+".lnk")
	shortcut := Shortcut{
		ShortcutPath:     shortcutPath,
		Target:           target,
		IconLocation:     iconPath,
		Arguments:        "",
		Description:      "",
		Hotkey:           "",
		WindowStyle:      "1",
		WorkingDirectory: "",
	}
	return Create(shortcut)
}

type ShortcutType int
const (
	Desktop ShortcutType = 0
	StartMenu ShortcutType = 1
)

func GetShortcutTypeDirectory(shortcutType ShortcutType) (string, error) {
	switch shortcutType {
	case Desktop:
		u, err := user.Current()
		if err != nil {
			return "", err
		}
		return u.HomeDir, nil
	case StartMenu:
		appdata, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(appdata, "./Microsoft/Windows/Start Menu/Programs"), nil
	}
	return "", errors.New("invalid shortcut type")
}

func DeleteShortcut(name string, shortcutType ShortcutType) error {
	dir, err := GetShortcutTypeDirectory(shortcutType)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(dir, name+".lnk"))
}

// CreateShortcut create with a shortcut object
func Create(shortcut Shortcut) error {
	if shortcut.IconLocation == "" {
		shortcut.IconLocation = "%SystemRoot%\\System32\\SHELL32.dll,0"
	}
	if shortcut.WindowStyle == "" {
		shortcut.WindowStyle = "1"
	}
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcut.ShortcutPath)
	if err != nil {
		return err
	}

	idispatch := cs.ToIDispatch()
	_, err = oleutil.PutProperty(idispatch, "IconLocation", shortcut.IconLocation)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "TargetPath", shortcut.Target)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "Arguments", shortcut.Arguments)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "Description", shortcut.Description)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "Hotkey", shortcut.Hotkey)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "WindowStyle", shortcut.WindowStyle)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "WorkingDirectory", shortcut.WorkingDirectory)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(idispatch, "Save")
	if err != nil {
		return err
	}
	return nil
}