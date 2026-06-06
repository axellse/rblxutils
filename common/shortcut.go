package common

import (
	"fmt"
	"os"

	"github.com/axellse/rblxutils/common/shortcut"
)

func DeleteShortcut(sType shortcut.ShortcutType) error {
	err := shortcut.DeleteShortcut("rblxutils", sType)
	if os.IsNotExist(err) {
		fmt.Println("no action required for assertion")
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func AssertShortcuts() {
	var err error
	if Config.Misc.DesktopShortcutEnabled {
		err = shortcut.CreateShortcut("rblxutils", BinPath, BinPath+",0", shortcut.Desktop)
	} else {
		err = DeleteShortcut(shortcut.Desktop)
	}

	if err != nil {
		FatalError(err)
	}

	if !Config.Misc.DisableStartmenuShortcut {
		err = shortcut.CreateShortcut("rblxutils", BinPath, BinPath+",0", shortcut.StartMenu)
	} else {
		err = DeleteShortcut(shortcut.StartMenu)
	}

	if err != nil {
		FatalError(err)
	}
}
