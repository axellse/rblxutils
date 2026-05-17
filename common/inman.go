package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/coder/websocket"
)

//inman - instance manager

type ServerData struct {
	JobId string
	PlaceId int
	UDMUXAddress string
	RCCAddress string
	ServerAddress string //actual server ip, if unprotected it's RCCAddress, if protected its UDMUXAddress
	UniverseId int
	UserId int
	Players []string
	GameData GameData
}

type Instance struct {
	parentInman *Inman
	ServerData ServerData
	Process *os.Process
}

type Inman struct {
	instanceRecord []*Instance
	Conn *websocket.Conn
}

func (i *Inman) GetInstances() []*Instance {
	return i.instanceRecord
}

func (i *Inman) AllocateInstance() *Instance {
	if i.instanceRecord == nil {
		i.instanceRecord = []*Instance{}
	}
	
	inPtr := &Instance{
		parentInman: i,
		ServerData: ServerData{},
	}
	i.instanceRecord = append(i.instanceRecord, inPtr)

	return inPtr
}

//updates inman's record to declare this instance as closed, and exits rblxutils if no more instances are alive
func (i *Instance) MarkAsClosed() {
	ii := slices.Index(i.parentInman.instanceRecord, i)
	i.parentInman.instanceRecord[ii] = i.parentInman.instanceRecord[len(i.parentInman.instanceRecord)-1]
	i.parentInman.instanceRecord = i.parentInman.instanceRecord[:len(i.parentInman.instanceRecord)-1]

	if len(i.parentInman.instanceRecord) == 0 {
		fmt.Println("inman: no more instances alive, cleaning up then exiting.")
		i.parentInman.Conn.Close(websocket.StatusNormalClosure, "close")
		os.Exit(0)
	}
}

//Close shuts down the instance process, then it calls MarkAsClosed
func (i *Instance) Close() {
	t1 := time.Now()
	err := i.Process.Kill()
	if err != nil {
		FatalError(err)
	}
	i.Process.Wait()
	fmt.Println("took", time.Since(t1).Milliseconds(), "ms to close process")
	
	i.MarkAsClosed()
}

type v1GamesResponse struct {
	Data []GameData `json:"data"`
}

type GameData struct {
	Name string`json:"name"`
	Description string `json:"description"`
}

func (i *Instance) QueryPlaceInfo() {
	if i.ServerData.UniverseId == 0 {
		fmt.Println("skipping quering place info: unpopulated universeid.")
		return
	}

	resp, err := http.Get("https://games.roblox.com/v1/games?universeIds=" + strconv.Itoa(i.ServerData.UniverseId))
	if err != nil {
		FatalError(err)
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		FatalError(err)
	}

	var v1GameResp v1GamesResponse
	err = json.Unmarshal(ba, &v1GameResp)
	if err != nil {
		FatalError(err)
	}

	if len(v1GameResp.Data) == 0 {
		FatalErrorStr("games response is nil or empty.")
	}

	i.ServerData.GameData = v1GameResp.Data[0]
}