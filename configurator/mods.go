package configurator

import (
	"encoding/json"
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

func TranslateBloxstrapZipMod(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		common.FatalError(err)
	}

	defer f.Close()
	stats, err := f.Stat()
	if err != nil {
		common.FatalError(err)
	}

	rules, err := common.GetFiles(f, stats.Size())
	if err != nil {
		common.FatalError(err)
	}

	outputRcmf := common.RcmfFile{
		Spec: "https://rcmf.axell.me/v1",
	}

	for path, data := range rules {
		outputRcmf.Rules = append(outputRcmf.Rules, common.RcmfRule{
			Sources: common.Sources{
				Files: []string{path},
			},
			Data: common.RcmfData{
				Blob: data,
			},
		})
	}

	ba, err := json.Marshal(outputRcmf)
	if err != nil {
		common.FatalError(err)
	}

	return ba
}

func ModImportWizard() {
	filename, err := dialog.File().Title("Import mod from file").Filter(".rcmf (Roblox Community Modding Format)", "rcmf").Filter(".zip (Bloxstrap-style standard file modding zip)", "zip").Load()
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
			common.ErrorStr("Mod with that name already exists! If you really want to import this mod, rename the file or remove the existing mod.")
			return
		}
	}

	var ba []byte
	if filepath.Ext(filename) == ".zip" {
		fmt.Println("importing as zip")
		ba = TranslateBloxstrapZipMod(filename)
	} else {
		fmt.Println("importing as rcmf")
		ba, err = os.ReadFile(filename)
		if err != nil {
			common.Error(err)
			return
		}
	}

	common.Config.Mods = append(common.Config.Mods, common.Mod{
		Name:    name,
		Enabled: true,
		Binary:  ba,
	})
	fmt.Println("import", filepath.Base(filename), "ok")
}
