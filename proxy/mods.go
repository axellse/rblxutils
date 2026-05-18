package proxy

import (
	"encoding/json"
	"fmt"

	"github.com/axellse/rblxutils/common"
)

func ConsolidateMods() []common.RcmfRule {
	fmt.Println("reading/consolidating mods...")
	rules := []common.RcmfRule{}
	for _, modF := range common.Config.Mods {
		if !modF.Enabled {
			continue
		}

		var mod common.RcmfFile
		err := json.Unmarshal(modF.Binary, &mod)
		if err != nil {
			common.FatalError(err)
		}

		if mod.Spec != "https://rcmf.axell.me/v1" {
			common.FatalErrorStr("Unknown spec '" + mod.Spec + "'. Rblxutils is compatible with https://rcmf.axell.me/v1")
			continue
		}

		rules = append(rules, mod.Rules...)
	}

	return rules
}
