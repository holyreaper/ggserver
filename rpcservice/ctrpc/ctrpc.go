package ctrpc

import (
	"time"

	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"
	"golang.org/x/net/context"

	. "github.com/holyreaper/ggserver/glog"
)

func init() {
}

//CTRPC .
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

//PushStream push stream
func (chat *CTRPC) PushStream(req *ctrpcpt.PushMessageRequest, stream ctrpcpt.CTRPC_PushStreamServer) error {
	var i int32
	LogInfo("ctrpc recv data  %v", *req)
	for {
		rsp := &ctrpcpt.PushMessageReply{
			Type: i,
		}
		err := stream.Send(rsp)
		if err != nil {
			LogFatal(" rpcservice  push stream got err %v", err)
			return err
		}
		time.Sleep(2 * time.Second)
		i++
	}
	return nil
}
