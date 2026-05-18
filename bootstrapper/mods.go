package bootstrapper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/resources"
	"github.com/beevik/etree"
)

// removed globalbasicsettings due to possible misuse
func EvalSpecialFile(file string) string {
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

				ba, ok := PerformKVMod(filepath.Ext(path), ba, rule.Data.Key, rule.Data.Value)
				if ok {
					os.WriteFile(path, ba, 0666)
				} else {
					fmt.Println("skipping applying mod rule, PerformKVMod reports not okay status")
				}
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

	certFile := filepath.Join(installDir, "ssl", "cacert.pem")
	ba, err := os.ReadFile(certFile)
	if err != nil {
		common.FatalError(err)
	}

	ba = append(ba, []byte("\n"+string(resources.CACert))...)
	err = os.WriteFile(certFile, ba, 0666)
	if err != nil {
		common.FatalError(err)
	}
}

func PerformKVMod(fileType string, file []byte, key string, value string) (out []byte, ok bool) {
	switch fileType {
	case ".json":
		var baseJson any
		err := json.Unmarshal(file, &baseJson)
		if err != nil {
			common.FatalErrorStr("failed Unmarshal json for KV mod, is your install corrupted?")
		}

		object := baseJson
		path := strings.Split(key, ".")
		for ti, target := range path {
			nextJson, ok := object.(map[string]any)
			if !ok {
				common.FatalErrorStr("cant apply kv mod: json has unexpected structure.")
			}

			if ti == len(path)-1 {
				nextJson[target] = value
			} else if nextJson[target] == nil {
				nextJson[target] = map[string]any{}
			}
			object = nextJson[target]
		}

		result, err := json.Marshal(baseJson)
		if err != nil {
			common.FatalError(err)
		}

		return result, true
	case ".xml":
		doc := etree.NewDocument()
		err := doc.ReadFromBytes(file)
		if err != nil {
			common.FatalErrorStr("cant apply kv mod: xml has unexpected structure.")
		}

		for e := range doc.FindElementsSeq(key) {
			for i := range e.Child {
				e.RemoveChildAt(i)
			}
			e.CreateText(value)
		}

		doc.IndentTabs()

		result, err := doc.WriteToBytes()
		return result, err == nil
	}

	return []byte{}, false
}
