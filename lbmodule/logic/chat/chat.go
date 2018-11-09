package chat

import (
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

	return nil
}
