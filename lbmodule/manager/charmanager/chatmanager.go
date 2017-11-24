package charmanager

import (
	. "github.com/holyreaper/ggserver/def"
)

//ChatMng single chat mng
type ChatMng struct {
	Manager
}

//Login user
func (cm *ChatMng) Login(uid UID) bool {
	cm.uid = uid
	return true
}

//SendMessageToUser SendMessageToUser
func (cm *ChatMng) SendMessageToUser(uid UID, data []byte) (err error) {
	err = SendMessageToUser(uid, data)
	if err != nil {
		//save to off line message
	}
	return
}
