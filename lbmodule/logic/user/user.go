package user

import (
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"

	. "github.com/holyreaper/ggserver/def"
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
func Login(cnn *net.TCPConn, reqData []byte) (rsp message.LoginMsgReply) {
	//user login ...
	var req message.LoginMsgRequest
	if err := proto.Unmarshal(reqData, &req); err != nil {
		rsp.Result = -1
		return
	}

	fmt.Printf("user %d Login success !!", req.Uid)
	//charmanager.AddUser(UID(req.GetUid()))
	charmanager.Login(cnn, UID(req.GetUid()))
	//rsp.Result = 189
	return
}
