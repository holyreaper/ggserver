package chat

import (
	"context"

	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbmodule/manager/charmanager"

	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"

	"github.com/holyreaper/ggserver/rpcclient"

	"github.com/holyreaper/ggserver/lbmodule/funcall"
	. "github.com/holyreaper/ggserver/lbmodule/packet"

	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/lbmodule/pb/message"
)

func init() {
	funcall.BindFunc(PKGChat, Chat)
}

//Chat chat
func Chat(req message.Message) (rsp message.Message) {

	//fmt.Printf("user %d send to user %d message %s", req.ChatRequest.Fuid, req.ChatRequest.Tuid, req.ChatRequest.Msg)
	err := SendMsg(req)
	if err != nil {
		LogFatal(" user %d chat fail error %v ", req.GetUid(), err)
		rsp.ChatReply.Result = 1
	}
	rsp.Command = PKGChat

	return
}

//SendMsg ,
func SendMsg(req message.Message) error {

	if charmanager.GetUser(def.UID(req.ChatRequest.Msg.Tuid)) != nil {
		msg := message.Message{}
		msg.ChatmsgPush.Msg.Fuid = req.ChatRequest.Msg.Fuid
		msg.ChatmsgPush.Msg.Funame = req.ChatRequest.Msg.Funame
		msg.ChatmsgPush.Msg.Ffigure = req.ChatRequest.Msg.Ffigure
		msg.ChatmsgPush.Msg.Flevel = req.ChatRequest.Msg.Flevel
		msg.ChatmsgPush.Msg.Fvip = req.ChatRequest.Msg.Fvip
		msg.ChatmsgPush.Msg.Tuid = req.ChatRequest.Msg.Tuid
		msg.ChatmsgPush.Msg.Msg = req.ChatRequest.Msg.Msg
		err := charmanager.SendMessageToUser(def.UID(req.ChatRequest.Msg.Tuid), PKGChatPush, &msg)
		if err != nil {
			return err
		}
	} else {
		client, err := rpcclient.GetCTRPC()
		if err != nil {
			return err
		}
		msg := &ctrpcpt.ChatRequest{}
		msg.ChatMsg.Fuid = req.ChatRequest.Msg.Fuid
		msg.ChatMsg.Funame = req.ChatRequest.Msg.Funame
		msg.ChatMsg.Ffigure = req.ChatRequest.Msg.Ffigure
		msg.ChatMsg.Flevel = req.ChatRequest.Msg.Flevel
		msg.ChatMsg.Fvip = req.ChatRequest.Msg.Fvip
		msg.ChatMsg.Tuid = req.ChatRequest.Msg.Tuid
		msg.ChatMsg.Msg = req.ChatRequest.Msg.Msg
		_, err = client.Chat(context.Background(), msg)
		if err != nil {
			return err
		}

	}
	return nil
}
