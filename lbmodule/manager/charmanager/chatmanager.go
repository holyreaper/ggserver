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
