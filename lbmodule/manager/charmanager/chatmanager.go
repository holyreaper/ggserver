package charmanager

import (
	"github.com/golang/protobuf/proto"
	. "github.com/holyreaper/ggserver/def"
)

//ChatMng single chat mng
type ChatMng struct {
	Manager
	uid UID
}

//Login user
func (cm *ChatMng) Login(uid UID) bool {
	cm.uid = uid
	return true
}

//SendMessageToUser SendMessageToUser
func (cm *ChatMng) SendMessageToUser(uid UID, tp int32, data proto.Message) (err error) {
	err = SendMessageToUser(uid, tp, data)
	if err != nil {
		//save to off line message
	}
	return
}

//NewChatMng .
func NewChatMng(id UID) *ChatMng {
	return &ChatMng{
		uid: id,
	}
}
