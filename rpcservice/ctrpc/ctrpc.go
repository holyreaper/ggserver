package ctrpc

import (
	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"
	"golang.org/x/net/context"
)

func init() {
}

//CTRPC
type CTRPC struct {
	//Srv interface{}
}

//KeepAlive  hello ...
func (chat *CTRPC) KeepAlive(cont context.Context, chatmsg *ctrpcpt.KeepAliveRequest) (*ctrpcpt.KeepAliveReply, error) {
	return &ctrpcpt.KeepAliveReply{Result: 2}, nil
}

//Login  hello ...
func (chat *CTRPC) Login(cont context.Context, chatmsg *ctrpcpt.LoginRequest) (*ctrpcpt.LoginReply, error) {
	return &ctrpcpt.LoginReply{Result: "HELO"}, nil
}
