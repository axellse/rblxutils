package common

import (
	"encoding/json"
	"os"
	"reflect"
)

type Configuration struct {
	UI struct {
		DisableWelcomeScreen bool
		ErrorStyle           int //0=dialog, 1=notification
		Theme                int //0=default, 1=white, 2=red, 3=dark
		BootstrapperImage    int//0=random cat, 1=rblxutils logo
	}
	Misc struct {
		DisableLaunchNotification bool
		DesktopShortcutEnabled    bool
		DebugOptions              []string
	}

	ServerHistoryEnabled bool
	Mods                 []Mod
}
type Mod struct {
	Name    string //Extracted from filename on import
	Enabled bool
	Binary  []byte
}

var Config Configuration
var ConfigFileState Configuration

func LoadConfiguration() {
	ba, err := os.ReadFile(LPath("./config.json"))
	if err != nil {
		ba = []byte("{}")
	}

	err = json.Unmarshal(ba, &Config)
	if err != nil {
		FatalError(err)
	}

	err = json.Unmarshal(ba, &ConfigFileState)
	if err != nil {
		FatalError(err)
	}
}

func ChangesMade() bool {
	return !reflect.DeepEqual(Config, ConfigFileState)
}

func WriteConfiguration() error {
	ba, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		Error(err)
		return err
	}

	err = os.WriteFile(LPath("./config.json"), ba, 0666)
	if err != nil {
		Error(err)
		return err
	}
	LoadConfiguration() //update ConfigFileState
	return nil
}

/*func GetConfigApplicationRequirement() string {
	modified := map[string]any{}
	err := mapstructure.Decode(Config, &modified)
	if err != nil {
		FatalError(err)
	}

	inital := map[string]any{}
	err = mapstructure.Decode(InitalConfig, &inital)
	if err != nil {
		FatalError(err)
	}

	restartType := ""
	for key, val := range modified {
		if inital[key] == val {
			continue
		}

		fmt.Println("Changed", key)
		if slices.Contains(RequiresRobloxRestart, key) {
			restartType = "roblox"
		} else if slices.Contains(RequiresSoftRestart, key) && restartType == "" {
			restartType = "soft"
		} else {
			restartType = "rblxutils"
		}
	}

	return restartType
}*/
