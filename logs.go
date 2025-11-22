package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func FindAndOpenLog() {
	logDir := filepath.Join(LocalAppData, "Roblox", "logs")

	foundLog := ""
	for foundLog == "" {
		fmt.Println("now searching for log files...")
		files, err := os.ReadDir(logDir)
		if err != nil {
			FatalError(err)
		}

		for _, file := range files {
			if file.IsDir() {continue}
			info, err := file.Info()
			if err != nil {
				Error(err)
				continue
			}

			if time.Since(info.ModTime()).Seconds() <= 15 {
				fmt.Println("found valid log file from " + strconv.Itoa(int(time.Since(info.ModTime()).Seconds())) + "s ago!")
				foundLog = file.Name()
				break
			}
		}
		if foundLog == "" {
			fmt.Println("couldnt find any new enough log files, waiting a little while...")
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Print("alright, now opening log file... ")
	file, err := os.Open(filepath.Join(logDir, foundLog))
	if err != nil {
		fmt.Println("failiure!")
		FatalError(err)
	}

	fmt.Println("success")
	for {
		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			if err == io.EOF {
				time.Sleep(2 * time.Second)
				continue
			} else {
				FatalError(err)
			}
		}
	}

}