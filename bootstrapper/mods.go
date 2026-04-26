package bootstrapper

import (
	"encoding/json"
	"os"
	"path/filepath"

	"axell.me/rblxutils/common"
)

func EvalSpecialFile(file string) string {
	switch file {
	case "GlobalBasicSettings_13.xml":
		return filepath.Join(common.RobloxAppData, "GlobalBasicSettings_13.xml")
	case "GlobalSettings_13.xml":
		return filepath.Join(common.RobloxAppData, "GlobalSettings_13.xml")
	}
	return ""
}

func ApplyFileMod(installDir string, modBa []byte) {
	var mod common.RcmfFile
	err := json.Unmarshal(modBa, &mod)
	if err != nil {
		common.FatalError(err)
	}

	if mod.Spec != "https://rcmf.axell.me/v1" {
		common.FatalErrorStr("Unknown spec '" + mod.Spec + "'. Rblxutils is compatible with https://rcmf.axell.me/v1")
		return
	}

	for _, rule := range mod.Rules {
		for _, path := range rule.Sources.Files {
			if nPath := EvalSpecialFile(path); nPath != "" {
				path = nPath
			} else if !filepath.IsLocal(path) {
				common.FatalErrorStr("mod has invalid path: tried modifying files outside roblox directory.")
			} else {
				path = filepath.Join(installDir, path)
			}

			if len(rule.Data.Blob) > 0 {
				err := os.WriteFile(path, rule.Data.Blob, 0666)
				if err != nil {
					common.FatalErrorStr("could not apply mod, file write error: " + err.Error())
				}
			} else if rule.Data.Key != "" {
				ba, err := os.ReadFile(path)
				if err != nil {
					common.FatalErrorStr("could not apply mod, file write error: " + err.Error())
				}
				PerformKVMod(filepath.Ext(path), ba, rule.Data.Key, rule.Data.Value)
			}
		}
	}
}

func ApplyFileMods(installDir string) {
	for _, mod := range common.Config.Mods {
		if !mod.Enabled {
			Println("skipping disabled mod", mod.Name)
			continue
		}

		Println("now applying mod", mod.Name)
		ApplyFileMod(installDir, mod.Binary)
	}
}


