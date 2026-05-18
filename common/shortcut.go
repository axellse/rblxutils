package common

import (
	"fmt"
	"os"

	"github.com/axellse/rblxutils/lib/shortcut"
)

func DeleteDesktopShortcut() error {
	err := shortcut.DeleteDesktopShortcut("rblxutils")
	if os.IsNotExist(err) {
		fmt.Println("no action required for assertion")
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func AssertDesktopShortcut() {
	var err error
	if Config.Misc.DesktopShortcutEnabled {
		err = shortcut.CreateDesktopShortcut("rblxutils", BinPath, BinPath+",0")
	} else {
		err = DeleteDesktopShortcut()
	}

	if err != nil {
		FatalError(err)
	}
}
