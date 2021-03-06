package user

import (
	"fmt"

	"github.com/holyreaper/ggserver/lbmodule/packet"

	. "github.com/holyreaper/ggserver/def"
	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/lbmodule/funcall"
	"github.com/holyreaper/ggserver/lbmodule/manager/charmanager"
	. "github.com/holyreaper/ggserver/lbmodule/packet"
	"github.com/holyreaper/ggserver/lbmodule/pb/message"
)

func init() {
	fmt.Println("bind login message  ")
	funcall.BindFunc(PKGLogin, Login)
}

//Login Login
func Login(rpack chan<- packet.Packet, spack chan<- packet.Packet, exitCh chan bool, req message.Message) (rsp message.Message) {

	fmt.Printf("user %d Login success !!", req.LoginRequest.GetUid())
	LogInfo("user %d Login success !!", req.LoginRequest.GetUid())
	//rsp.Pack(packet.PKGLogin, &message.LoginMsgReply{Result: 2018})
	rsp.LoginReply = &message.LoginMsgReply{}
	rsp.LoginReply.Result = 2018
	charmanager.Login(rpack, spack, exitCh, UID(req.LoginRequest.GetUid()))
	return
}
