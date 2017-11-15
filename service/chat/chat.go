package chat

import (
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
func (chat *Chat) Chat(context.Context, *chatrpc.ChatMsgRequest) (*chatrpc.ChatMsgReply, error) {
	return &chatrpc.ChatMsgReply{Result: 1}, nil
}
