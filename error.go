package main

import (
	"fmt"
	"os"

	"github.com/gen2brain/beeep"
)

func FatalError(err error) {
	FatalErrorStr(err.Error())
}

func Error(err error) {
	ErrorStr(err.Error())
}

func ErrorStr(err string) {
	beeep.Notify("rblxutils errored:", err, ProgramLogo)
}

func FatalErrorStr(err string) {
	ErrorStr(err)
	fmt.Println(err)
	os.Exit(1)
}