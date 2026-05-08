package bootstrapper

import (
	"fmt"
	"os"
	"slices"

	"github.com/coder/websocket"
)

//inman - instance manager

type Instance struct {
	parentInman *Inman
}

type Inman struct {
	instanceRecord []*Instance
	conn *websocket.Conn
}

//only inman should be poking around in here 

func (i *Inman) AllocateInstance() *Instance {
	if i.instanceRecord == nil {
		i.instanceRecord = []*Instance{}
	}
	
	inPtr := &Instance{
		parentInman: i,
	}
	i.instanceRecord = append(i.instanceRecord, inPtr)

	return inPtr
}

//updates inman's record to declare this instance as closed, and exits rblxutils if no more instances are alive
func (i *Instance) Close() {
	ii := slices.Index(i.parentInman.instanceRecord, i)
	i.parentInman.instanceRecord[ii] = i.parentInman.instanceRecord[len(i.parentInman.instanceRecord)-1]
	i.parentInman.instanceRecord = i.parentInman.instanceRecord[:len(i.parentInman.instanceRecord)-1]

	if len(i.parentInman.instanceRecord) == 0 {
		fmt.Println("inman: no more instances alive, cleaning up then exiting.")
		i.parentInman.conn.Close(websocket.StatusNormalClosure, "close")
		os.Exit(0)
	}
}