package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type StateT struct {
	RequiresModApplication bool //this reports if the last action was a change to mods in config rather than a mod application
	ServerHistory []ServerData
}

var State StateT
var stateMutex sync.Mutex

func LoadState() {
	stateMutex.Lock()
	ba, err := os.ReadFile(LPath("./state.json"))
	if err != nil {
		if DotSlash == filepath.Join(LocalAppData, "rblxutils") {
			ba = []byte("{}")
		} else {
			FatalErrorStr("Rblxutils does not appear to be installed. Move the executable to %LOCALAPPDATA%\\rblxutils and run it again.")
		}
	}

	err = json.Unmarshal(ba, &State)
	if err != nil {
		FatalError(err)
	}
	stateMutex.Unlock()
}

func WriteState() error {
	stateMutex.Lock()
	ba, err := json.MarshalIndent(State, "", "    ")
	if err != nil {
		Error(err)
		return err
	}

	err = os.WriteFile(LPath("./state.json"), ba, 0666)
	if err != nil {
		Error(err)
		return err
	}
	stateMutex.Unlock()
	LoadState()
	return nil
}