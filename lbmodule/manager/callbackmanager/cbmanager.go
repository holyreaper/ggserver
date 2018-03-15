package callbackmanager

import (
	"fmt"
	"sync"

	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"
)

//CBManager .
type CBManager struct {
	data   chan ctrpcpt.PushMessageReply
	exitCh chan bool
}

var cbmanager *CBManager
var onece sync.Once

const (
	initNum = 1000
)

func init() {
	cbmanager = &CBManager{
		data:   make(chan ctrpcpt.PushMessageReply, initNum),
		exitCh: make(chan bool),
	}
}

//DealCBData deal call back data
func (cbm *CBManager) DealCBData() {
	defer func() {
		if err := recover(); err != nil {
			LogFatal("callbackmanager exit err %v ", err)
		}
		LogInfo("cbmanager exit ...")
	}()
	LogInfo("cbmanager DealCBData start")
	for {
		select {
		case <-cbm.exitCh:
			break
		case data := <-cbm.data:
			{
				LogInfo(" lbserver recv data type %d ", data.GetType())

				//TODO deal with center req
			}
		}
	}

}

//Start .
func Start() {
	go cbmanager.DealCBData()
}

//Stop .
func Stop() {
	fmt.Printf("cbmanager stop")
	onece.Do(func() {
		close(cbmanager.exitCh)
	})
}

//Put .
func Put(data ctrpcpt.PushMessageReply) {
	if cbmanager == nil {
		LogFatal("callbackmanager is not init ...")
	}
	cbmanager.data <- data
}
