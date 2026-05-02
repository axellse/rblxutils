package bootstrapper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"axell.me/rblxutils/common"
)

type LogProcessor struct{}

func (*LogProcessor) Write(p []byte) (int, error) {
	if strings.Contains(string(p), "App, internal browser session end") {
		common.KillHelper()
		os.Exit(0)
	}
	return len(p), nil
}

func FindAndOpenLog() {
	logDir := filepath.Join(common.LocalAppData, "Roblox", "logs")

	foundLog := ""
	for foundLog == "" {
		fmt.Println("now searching for log files...")
		files, err := os.ReadDir(logDir)
		if err != nil {
			common.FatalError(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			info, err := file.Info()
			if err != nil {
				common.Error(err)
				continue
			}

			if time.Since(info.ModTime()).Seconds() <= 3 {
				fmt.Println("found valid log file from " + strconv.Itoa(int(time.Since(info.ModTime()).Seconds())) + "s ago!")
				foundLog = file.Name()
				break
			}
		}
		if foundLog == "" {
			fmt.Println("couldnt find any new enough log files, waiting a little while...")
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Print("alright, now opening log file... ")
	file, err := os.Open(filepath.Join(logDir, foundLog))
	if err != nil {
		fmt.Println("failiure!")
		common.FatalError(err)
	}
	fmt.Println("success")

	lp := LogProcessor{}
	for {
		_, err = io.Copy(&lp, file)
		if err != nil {
			if err == io.EOF {
				time.Sleep(1 * time.Second)
				continue
			} else {
				common.FatalError(err)
			}
		}
	}
}
