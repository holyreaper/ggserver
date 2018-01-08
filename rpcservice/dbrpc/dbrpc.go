package dbrpc

import (
	"github.com/holyreaper/ggserver/rpcservice/pb/dbrpc"
	"golang.org/x/net/context"
)

func init() {
}

//DBRPC ..
type DBRPC struct {
	//Srv interface{}
}

//KeepAlive  hello ...
func (chat *DBRPC) KeepAlive(cont context.Context, chatmsg *dbrpcpt.KeepAliveRequest) (*dbrpcpt.KeepAliveReply, error) {
	return &dbrpcpt.KeepAliveReply{Result: 2}, nil
}

//Login  hello ...
func (chat *DBRPC) Login(cont context.Context, chatmsg *dbrpcpt.LoginRequest) (*dbrpcpt.LoginReply, error) {
	return &dbrpcpt.LoginReply{}, nil
}
