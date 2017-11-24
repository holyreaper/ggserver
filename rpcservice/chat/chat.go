package chat

import (
	"fmt"
	"runtime"
	"time"

	"github.com/holyreaper/ggserver/rpcservice/pb/chat"
	"golang.org/x/net/context"
)

func init() {
	// init
	fmt.Println("sdfsdafasd")
}

//Chat 登录实现
type Chat struct {
	Srv interface{}
}

//Chat  hello ...
func (chat *Chat) Chat(cont context.Context, chatmsg *chatrpc.ChatMsgRequest) (*chatrpc.ChatMsgReply, error) {
	fmt.Println("get client msg ", *chatmsg)
	return &chatrpc.ChatMsgReply{Result: 2}, nil
}

//ChatList stream
func (chat *Chat) ChatList(chatstream chatrpc.ChatRpc_ChatListServer) error {

	fmt.Println("current goroutin num :", runtime.NumGoroutine())
	for {

		//err := chatstream.Send(&reply)
		time.Sleep(time.Second * 10)
	}
}
