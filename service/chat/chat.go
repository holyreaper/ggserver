package chat

import (
	"fmt"

	"github.com/holyreaper/ggserver/pb/chat"
	"golang.org/x/net/context"
)

func init() {
	// init
}

//Chat 登录实现
type Chat struct {
}

//Chat  hello ...
func (chat *Chat) Chat(cont context.Context, chatmsg *chatrpc.ChatMsgRequest) (*chatrpc.ChatMsgReply, error) {
	fmt.Println("get client msg ", *chatmsg)
	return &chatrpc.ChatMsgReply{Result: 2}, nil
}
