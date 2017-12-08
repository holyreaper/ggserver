package chat

import (
	"github.com/holyreaper/ggserver/lbmodule/funcall"
	. "github.com/holyreaper/ggserver/lbmodule/packet"

	"github.com/holyreaper/ggserver/lbmodule/pb/message"
)

func init() {
	funcall.BindFunc(PKGChat, Chat)
}

//Chat chat
func Chat(req message.Message) (rsp message.Message) {

	//fmt.Printf("user %d send to user %d message %s", req.Fuid, req.Tuid, req.Msg)
	//rsp.Result = 1
	return
}
