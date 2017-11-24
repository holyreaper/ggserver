package login

import (
	"github.com/holyreaper/ggserver/rpcservice/pb/login"
	"golang.org/x/net/context"
)

func init() {
	// init
}

//Login 登录实现
type Login struct {
	Srv interface{}
}

//Login hello ...
func (loginServer *Login) Login(context.Context, *loginrpc.LoginRequest) (*loginrpc.LoginReply, error) {
	return &loginrpc.LoginReply{Message: "helo"}, nil
}
