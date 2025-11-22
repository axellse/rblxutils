package main

import (
	_ "embed"
	"os"

	"github.com/gen2brain/beeep"
)

//go:embed placeholder.png
var ProgramLogo []byte
var ProgramName = "rblxutils" //you never know if you might want to change it

var LocalAppData = ""

func main() {
	var err error
	LocalAppData, err = os.UserCacheDir()
	if err != nil {
		FatalError(err)
	}
	beeep.AppName = ProgramName

	LoadConfiguration()
	go LaunchUI()
	go FindAndOpenLog()
	go OpenCacheDatabase()
	select {}
}