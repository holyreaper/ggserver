package callbackmanager

import (
	"fmt"
	"sync"

	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbmodule/manager/charmanager"
	"github.com/holyreaper/ggserver/lbmodule/packet"

	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/lbmodule/pb/message"
	"github.com/holyreaper/ggserver/rpcservice/def"
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
				switch rpcdef.RPCPushType(data.Type) {
				case rpcdef.RPCPTypeUserChat:
					DealWithChat(data)

				}
			}
		}
	}

}

//DealWithChat .
func DealWithChat(data ctrpcpt.PushMessageReply) {
	msg := message.Message{}
	msg.ChatmsgPush.Msg.Fuid = data.ChatMsg.Fuid
	msg.ChatmsgPush.Msg.Funame = data.ChatMsg.Funame
	msg.ChatmsgPush.Msg.Ffigure = data.ChatMsg.Ffigure
	msg.ChatmsgPush.Msg.Flevel = data.ChatMsg.Flevel
	msg.ChatmsgPush.Msg.Fvip = data.ChatMsg.Fvip
	msg.ChatmsgPush.Msg.Tuid = data.ChatMsg.Tuid
	msg.ChatmsgPush.Msg.Msg = data.ChatMsg.Msg
	//pack := packet.Packet{}
	//pack.Pack(packet.PKGChatPush, &msg)
	//req := message.Message{LoginRequest: , Command: packet.PKGChatPush}
	err := charmanager.SendMessageToUser(def.UID(data.ChatMsg.Tuid), packet.PKGChatPush, &msg)
	if err != nil {
		//save to db
		LogFatal("user %d is not online save data to db ", data.ChatMsg.Tuid)
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
