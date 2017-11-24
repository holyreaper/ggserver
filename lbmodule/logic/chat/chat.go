package chat

import (
	"fmt"

	"github.com/holyreaper/ggserver/lbmodule/funcall"
	. "github.com/holyreaper/ggserver/lbmodule/packet"

	"github.com/holyreaper/ggserver/lbmodule/pb/chat"
)

func init() {
	funcall.BindFunc(PKGChat, Chat)
}

//Chat chat
func Chat(req ptchat.ChatMsgRequest) (rsp ptchat.ChatMsgReply) {

	fmt.Printf("user %d send to user %d message %s", req.Fuid, req.Tuid, req.Msg)
	rsp.Result = 1
	return
}
