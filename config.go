package main

import (
	"encoding/json"
	"os"
)

type PTT struct {
	Enabled bool
	Device string
	Modifier string
	Key string
}
type Configuration struct {
	PushToTalk PTT
}

var Config Configuration

func LoadConfiguration() {
	ba, err := os.ReadFile("./config.json")
	if err != nil {
		FatalError(err)
	}

	err = json.Unmarshal(ba, &Config)
	if err != nil {
		FatalError(err)
	}
}

func WriteConfiguration() {
	ba, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		Error(err)
		return
	}

	err = os.WriteFile("./config.json", ba, 0666)
	if err != nil {
		Error(err)
	}
}