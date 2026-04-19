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

func PerformKVMod(fileType string, key string, value string) {
	if fileType == "json" {

	}
}

func ApplyFileMod(installDir string, modBa []byte) {
	var mod rcmfFile
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

			}

			


			

			
		}
	}

}

type rcmfFile struct {
	Spec  string     `json:"spec"`
	Rules []rcmfRule `json:"rules"`
}

type Sources struct {
	Expressions []string `json:"expressions"`
	Ids []string `json:"ids"`
	Types []string `json:"types"`
	Files []string `json:"files"`
}

type rcmfRule struct {
	Sources  Sources `json:"sources"`
	Data rcmfData `json:"data"`
}

type rcmfData struct {
	Blob  []byte `json:"blob"`
	Key   string `json:"key"` //dots can be used for nesting (eg. Settings.ContentFolder)
	Value string `json:"value"`
}