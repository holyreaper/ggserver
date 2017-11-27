package charmanager

import (
	"fmt"

	. "github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbmodule/packet"
	"github.com/holyreaper/ggserver/lbmodule/pb/user"
)

//UserMng single user mng
type UserMng struct {
	Manager
	uid UID
}

//Login user
func (cm *UserMng) Login(uid UID) bool {
	cm.uid = uid
	// load data from d
	data := ptuser.LoginMsgReply{
		Result: 1314,
	}
	fmt.Println("userManager send packet to client ")
	SendMessageToUser(uid, packet.PKGLogin, &data)
	return true
}

//LogOut user logout
func (cm *UserMng) LogOut(uid UID) bool {
	return true
}

//NewUserMng user mng
func NewUserMng(id UID) *UserMng {
	return &UserMng{
		uid: id,
	}
}
