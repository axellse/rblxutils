package common

import (
	"encoding/json"
	"os"
)

type StateT struct {
	RequiresModApplication bool //this reports if the last action was a change to mods in config rather than a mod application
	HelperAction string //what the helper should do when it starts up. this should be set to "" by the helper when it's done
}

var State StateT

func LoadState() {
	ba, err := os.ReadFile(LPath("./state.json"))
	if err != nil {
		FatalError(err)
	}

	err = json.Unmarshal(ba, &State)
	if err != nil {
		FatalError(err)
	}
}

func WriteState() error {
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
	LoadState()
	return nil
}