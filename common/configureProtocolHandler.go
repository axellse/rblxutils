package common

import "golang.org/x/sys/windows/registry"

func RegisterAsProtocolHandler() {
	handler := `"` + BinPath + `" %1`

	//protocol handlers
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\roblox-player\shell\open\command`, registry.ALL_ACCESS)
	if err != nil {
		FatalError(err)
	}
	key.SetStringValue("", handler)

	key, _, err = registry.CreateKey(registry.CURRENT_USER, `Software\Classes\roblox\shell\open\command`, registry.ALL_ACCESS)
	if err != nil {
		FatalError(err)
	}
	key.SetStringValue("", handler)

	//icons
	key, _, err = registry.CreateKey(registry.CURRENT_USER, `Software\Classes\roblox-player\DefaultIcon`, registry.ALL_ACCESS)
	if err != nil {
		FatalError(err)
	}
	key.SetStringValue("", BinPath)

	key, _, err = registry.CreateKey(registry.CURRENT_USER, `Software\Classes\roblox\DefaultIcon`, registry.ALL_ACCESS)
	if err != nil {
		FatalError(err)
	}
	key.SetStringValue("", BinPath)
}