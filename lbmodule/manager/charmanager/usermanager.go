package charmanager

import (
	. "github.com/holyreaper/ggserver/def"
)

//UserMng single user mng
type UserMng struct {
	Manager
	uid UID
}

//Login user
func (cm *UserMng) Login(uid UID) bool {
	cm.uid = uid

	return true
}

//LogOut user logout
func (cm *UserMng) LogOut() bool {
	return true
}

//NewUserMng user mng
func NewUserMng(id UID) *UserMng {
	return &UserMng{
		uid: id,
	}
}
