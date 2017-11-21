package charmanager

import (
	. "github.com/holyreaper/ggserver/def"
)

//UserMng single user mng
type UserMng struct {
	IManager
	uid UID
}

//Login user
func (cm *UserMng) Login(uid UID) bool {
	cm.uid = uid
	return true
}

//LogOut user logout
func (cm *UserMng) LogOut(uid UID) bool {
	return true
}
