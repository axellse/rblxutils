package bootstrapper

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"axell.me/rblxutils/common"
)

func PerformKVMod(fileType string, file []byte, key string, value any) []byte {
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
		return result
	case ".xml":
		var baseXml any
		err := xml.Unmarshal(file, &baseXml)
		if err != nil {
			common.FatalErrorStr("failed Unmarshal xml for KV mod, is your install corrupted?")
		}

		object := baseXml
		path := strings.Split(key, ".")
		for ti, target := range path {
			nextXml, ok := object.(map[string]any)
			if !ok {
				common.FatalErrorStr("cant apply kv mod: xml has unexpected structure.")
			}

			if ti == len(path)-1 {
				nextXml[target] = value
			} else if nextXml[target] == nil {
				nextXml[target] = map[string]any{}
			}
			object = nextXml[target]
		}

		result, err := xml.Marshal(baseXml)
		if err != nil {
			common.FatalError(err)
		}
		return result
	}

	return []byte{}
}