package configurator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/axellse/rblxutils/common"
	"github.com/sqweek/dialog"
)

func RenderMods(win *nucular.Window) {
	win.Row(10).Dynamic(1)
	win.Label("Mods can change local files, assets or both.", label.Align("LC")) //Rblxutils can apply and manage traditional mods that modify local files (eg. Bloxstap mods), but also mods that modify assets downloaded by the client from Roblox's CDN. Rlbxutils uses the Roblox Community Modding Format (rcmf) but can also automatically convert from a couple of other common formats.
	win.Row(10).Dynamic(1)
	win.Label("The following formats are supported: rcmf, zip", label.Align("LC"))

	win.Row(20).Static(150)
	if win.ButtonText("Import from file") {
		ModImportWizard()
	}

	win.Row(5).Dynamic(1)
	win.Spacing(1)
	win.Row(15).Dynamic(1)
	win.Label("Manage mods:", label.Align("LC"))

	for i := range common.Config.Mods {
		if len(common.Config.Mods) <= i {
			continue
		}
		manage := win.TreePush(nucular.TreeTab, common.Config.Mods[i].Name, false)
		if manage {
			win.Row(20).Dynamic(1)
			win.CheckboxText("Enabled", &common.Config.Mods[i].Enabled)
			win.Row(20).Dynamic(3)
			if win.ButtonText("Remove") {
				common.Config.Mods = append(common.Config.Mods[:i], common.Config.Mods[i+1:]...)
			}

			if win.ButtonText("Download rcmf") {
				SaveModWizard(common.Config.Mods[i])
			}
			win.TreePop()
		}
	}
}

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
