package configurator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"axell.me/rblxutils/common"
	"github.com/sqweek/dialog"
)

func SaveModWizard(mod common.Mod) {
	filename, err := dialog.File().SetStartFile(mod.Name+".rcmf").Title("Save/Export mod as RCMF").Filter(".rcmf", "rcmf").Save()
	if err == dialog.ErrCancelled {
		return
	}
	if err != nil {
		common.Error(err)
		return
	}

	err = os.WriteFile(filename, mod.Binary, 0755)
	if err != nil {
		common.Error(err)
	}
}

func ModImportWizard() {
	filename, err := dialog.File().Title("Import mod from file").Filter(".rcmf", "rcmf").Load()
	if err == dialog.ErrCancelled {
		return
	}
	if err != nil {
		common.Error(err)
		return
	}

	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	for _, mod := range common.Config.Mods {
		if mod.Name == name {
			common.ErrorStr("Mod with that name already exists! If you really want to import this mod, rename the file.")
			return
		}
	}

	ba, err := os.ReadFile(filename)
	common.Config.Mods = append(common.Config.Mods, common.Mod{
		Name:    name,
		Enabled: true,
		Binary:  ba,
	})
	fmt.Println("import", filepath.Base(filename), "ok")
}
