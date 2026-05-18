package bootstrapper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"axell.me/rblxutils/common"
)

type LogProcessor struct{
	instance *common.Instance
	lineBuf []byte
}

//these were taken from bloxstrap's activity watcher
var GameJoinPattern = regexp.MustCompile(`! Joining game '([0-9a-f\-]{36})' place ([0-9]+) at ([0-9\.]+)`)
var GameJoinLoadTime = regexp.MustCompile(`universeid:([0-9]+).*userid:([0-9]+)`)
var GameJoinUdmux = regexp.MustCompile(`UDMUX Address = ([0-9\.]+), Port = [0-9]+ \| RCC Server Address = ([0-9\.]+), Port = [0-9]+`)
var PlayerStateChanged = regexp.MustCompile(`Warning: (?:added|removed)  ([^\n\r ]*)`)
var LightningTechnology = regexp.MustCompile(`\[FLog::Graphics\] (.*) shadows`)


func (lp *LogProcessor) Write(p []byte) (int, error) {
	lp.lineBuf = append(lp.lineBuf, p...)

	if !bytes.ContainsRune(lp.lineBuf, '\n') {
		return len(p), nil
	}

	line := string(lp.lineBuf)
	lp.lineBuf = []byte{}

	if strings.Contains(line, "App, internal browser session end") {
		lp.instance.MarkAsClosed()
	} else if strings.Contains(line, "[FLog::Output] ! Joining game") {
		matches := GameJoinPattern.FindStringSubmatch(line)
		if len(matches) != 4 {
			return len(p), errors.New("invalid matches: " + line)
		}

		lp.instance.ServerData.JobId = matches[1]
		placeId, err := strconv.Atoi(matches[2])
		if err != nil {
			common.FatalError(err)
		}

		lp.instance.ServerData.PlaceId = placeId

		lp.instance.ServerData.RCCAddress = matches[3]
		lp.instance.ServerData.ServerAddress = matches[3]
		lp.instance.QueryPlaceInfo()
	} else if strings.Contains(line, "[FLog::GameJoinLoadTime] Report game_join_loadtime:") {
		fmt.Println("game join load time found")
		matches := GameJoinLoadTime.FindStringSubmatch(line)
		if len(matches) != 3 {
			return len(p), errors.New("invalid matches: " + line)
		}

		universeId, err := strconv.Atoi(matches[1])
		if err != nil {
			common.FatalError(err)
		}
		lp.instance.ServerData.UniverseId = universeId
		lp.instance.QueryPlaceInfo()

		userId, err := strconv.Atoi(matches[2])
		if err != nil {
			common.FatalError(err)
		}
		lp.instance.ServerData.UserId = userId
	} else if strings.Contains(line, "[FLog::Network] UDMUX Address = ") {
		matches := GameJoinUdmux.FindStringSubmatch(line)
		if len(matches) != 3 {
			return len(p), errors.New("invalid matches: " + line)
		}

		if lp.instance.ServerData.RCCAddress != matches[2] {
			common.FatalErrorStr("RCC Address mismatch!")
		}

		lp.instance.ServerData.UDMUXAddress = matches[1]
		lp.instance.ServerData.ServerAddress = matches[1]
	} else if strings.Contains(line, "[FLog::SingleSurfaceApp] leaveUGCGameInternal") {
		fmt.Println("leaving game, clearing game/server data.")
		if common.Config.ServerHistoryEnabled {
			common.State.ServerHistory = append(common.State.ServerHistory, lp.instance.ServerData)
			common.WriteState()
		}
		lp.instance.ServerData = common.ServerData{}
	} else if strings.Contains(line, "[FLog::Warning] Warning: added") {
		matches := PlayerStateChanged.FindStringSubmatch(line)
		if len(matches) != 2 {
			return len(p), errors.New("invalid matches: " + line)
		}

		fmt.Println("player added", matches[1])
		lp.instance.ServerData.Players = append(lp.instance.ServerData.Players, matches[1])
	} else if strings.Contains(line, "[FLog::Warning] Warning: removed") {
		matches := PlayerStateChanged.FindStringSubmatch(line)
		fmt.Println("player removed", matches[1])
		if len(matches) != 2 {
			return len(p), errors.New("invalid matches: " + line)
		}

		i := slices.Index(lp.instance.ServerData.Players, matches[1])
		if i == - 1{
			fmt.Println("player not in list?")
		} else {
			fmt.Println("total", len(lp.instance.ServerData.Players), "before")
			lp.instance.ServerData.Players = append(lp.instance.ServerData.Players[:i], lp.instance.ServerData.Players[i+1:]...)
			fmt.Println("total", len(lp.instance.ServerData.Players), "after")
		}
	}
	return len(p), nil
}

func FindAndOpenLog(instance *common.Instance) {
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

	fmt.Println("LOG FILE IS", foundLog)

	for _, in := range GlobalInman.GetInstances() {
		if in.LogFileName == foundLog {
			fmt.Println("conflicting logfilename!")
		}
	}

	lp := LogProcessor{
		instance: instance,
	}
	for {
		_, err = io.Copy(&lp, file)
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Second)
				continue
			} else {
				common.FatalError(err)
			}
		}
	}
}
