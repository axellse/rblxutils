package common

import (
	"fmt"
	"image"
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
	JobId         string
	PlaceId       int
	UDMUXAddress  string
	RCCAddress    string
	ServerAddress string //actual server ip, if unprotected it's RCCAddress, if protected its UDMUXAddress
	UniverseId    int
	UserId        int
	Players []string
	GameData      GameData
	JoinTime      time.Time
	LeaveTime     time.Time
	Location      ServerLocationInfo
	HeadshotURL   string
	User          User
}

type Instance struct {
	parentInman *Inman
	LogFileName string
	ServerData  ServerData
	process     *os.Process
}

type Inman struct {
	instanceRecord      []*Instance
	Conn                *websocket.Conn
	LaunchBootstrapperF func(newProcess bool, robloxArgs string)
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
		ServerData:  ServerData{},
	}
	i.instanceRecord = append(i.instanceRecord, inPtr)

	return inPtr
}

// updates inman's record to declare this instance as closed, and exits rblxutils if no more instances are alive
func (i *Instance) MarkAsClosed() {
	ii := slices.Index(i.parentInman.instanceRecord, i)
	if ii == -1 {
		return
	}
	i.parentInman.instanceRecord[ii] = i.parentInman.instanceRecord[len(i.parentInman.instanceRecord)-1]
	i.parentInman.instanceRecord = i.parentInman.instanceRecord[:len(i.parentInman.instanceRecord)-1]

	if len(i.parentInman.instanceRecord) == 0 {
		fmt.Println("inman: no more instances alive, cleaning up then exiting.")
		i.parentInman.Conn.Close(websocket.StatusNormalClosure, "close")
		os.Exit(0)
	}
}

// WaitForInstance waits for the instance process to exit, then calls MarkAsClosed.
func (i *Instance) WaitForInstance() {
	if i.process == nil {
		fmt.Println("nil process?")
		return
	}

	_, err := i.process.Wait()
	if err != nil {
		FatalError(err)
	}

	i.MarkAsClosed()
}

// SetProcess set Instance.process, then calls WaitForInstance in a new goroutine
func (i *Instance) SetProcess(p *os.Process) {
	i.process = p
	go i.WaitForInstance()
}

// Close shuts down the instance process, then it calls MarkAsClosed
func (i *Instance) Close() {
	t1 := time.Now()
	err := i.process.Kill()
	if err != nil {
		FatalError(err)
	}
	i.process.Wait()
	fmt.Println("took", time.Since(t1).Milliseconds(), "ms to close process")

	i.MarkAsClosed()
}

type v1GamesResponse struct {
	Data []GameData `json:"data"`
}

type GameData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxPlayers  int    `json:"maxPlayers"`
	RootPlaceId int    `json:"rootPlaceId"`
	Creator     struct {
		Name     string `json:"name"`
		Verified bool   `json:"hasVerifiedBadge"`
	} `json:"creator"`
	Thumbnail *image.RGBA
	IconURL   string
}

type v1GamesThumbnailResponse struct {
	Data []struct {
		Thumbnails []Thumbnail `json:"thumbnails"`
	} `json:"data"`
}

type v1GameIconsResponse struct {
	Data []Thumbnail `json:"data"`
}

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type Thumbnail struct {
	ImageUrl string `json:"imageUrl"`
}

func (i *Instance) QueryPlaceInfo() {
	if i.ServerData.UniverseId == 0 || i.ServerData.PlaceId == 0 || i.ServerData.UserId == 0 {
		fmt.Println("skipping quering place info: unpopulated universeid, placeid or userid.")
		fmt.Println(i.ServerData)
		return
	}

	fmt.Println("quering place info...")

	games := &v1GamesResponse{}
	QueryAndUnmarshal("https://games.roblox.com/v1/games?universeIds="+strconv.Itoa(i.ServerData.UniverseId), games)
	if len(games.Data) == 0 {
		FatalErrorStr("games response is nil or empty.")
	}
	i.ServerData.GameData = games.Data[0]

	thumbnails := v1GamesThumbnailResponse{}
	QueryAndUnmarshal("https://thumbnails.roblox.com/v1/games/multiget/thumbnails?format=png&size=384x216&countPerUniverse=1&universeIds="+strconv.Itoa(i.ServerData.UniverseId), &thumbnails)
	if len(thumbnails.Data) == 0 || len(thumbnails.Data[0].Thumbnails) == 0 {
		FatalErrorStr("games thumb response is nil or empty.")
	}

	resp, err := http.Get(thumbnails.Data[0].Thumbnails[0].ImageUrl)
	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		FatalError(err)
	}

	i.ServerData.GameData.Thumbnail = LoadImageUI(ba, 368, 0)

	icons := v1GameIconsResponse{}
	QueryAndUnmarshal("https://thumbnails.roblox.com/v1/places/gameicons?format=png&size=512x512&placeIds="+strconv.Itoa(i.ServerData.PlaceId), &icons)
	if len(icons.Data) == 0 {
		FatalErrorStr("games icon response is nil or empty.")
	}
	i.ServerData.GameData.IconURL = icons.Data[0].ImageUrl

	headshots := v1GameIconsResponse{}
	QueryAndUnmarshal("https://thumbnails.roblox.com/v1/users/avatar-headshot?size=720x720&format=Png&isCircular=false&userIds="+strconv.Itoa(i.ServerData.UserId), &headshots)
	if len(icons.Data) == 0 {
		FatalErrorStr("user headshots response is nil or empty.")
	}
	i.ServerData.HeadshotURL = headshots.Data[0].ImageUrl

	user := User{}
	QueryAndUnmarshal("https://users.roblox.com/v1/users/"+strconv.Itoa(i.ServerData.UserId), &user)
	i.ServerData.User = user
	fmt.Println("query ok")
}

type ServerLocationInfo struct {
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Loc     string `json:"loc"`
	Bogon   bool   `json:"bogon"`
}

func (i *Instance) QueryServerLocation() {
	if i.ServerData.ServerAddress == "" {
		fmt.Println("skipping quering server location: unpopulated server address.")
		return
	}
	fmt.Println("quering server location")

	location := ServerLocationInfo{}
	QueryAndUnmarshal("https://ipinfo.io/" + i.ServerData.ServerAddress + "/json", &location)
	if location.Bogon {
		fmt.Println("bogon address, probably we're connecting via udmux and been given a bogon rcc address.")
		return
	} else if location.Country == "" {
		FatalErrorStr("ip response is nil or empty.")
	}

	i.ServerData.Location = location
}
