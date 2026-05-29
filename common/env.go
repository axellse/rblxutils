package common

import (
	"os"
	"path/filepath"
)

var RobloxAppData string
var LocalAppData string
var DotSlash string
var BinPath string

func LPath(p string) string {
	return filepath.Join(DotSlash, p)
}


func DefineEnvs() {
	binPath, err := os.Executable()
	BinPath = binPath
	if err != nil {
		FatalError(err)
	}
	DotSlash = filepath.Dir(binPath)

	LocalAppData, err = os.UserCacheDir()
	if err != nil {
		FatalError(err)
	}

	RobloxAppData = filepath.Join(LocalAppData, "Roblox")
}
