package user

import (
	. "github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbmodule/funcall"
	"github.com/holyreaper/ggserver/lbmodule/manager/charmanager"
	. "github.com/holyreaper/ggserver/lbmodule/package"
)

func init() {
	funcall.BindFunc(PKGLogin, Login)
}

//Login Login
func Login(uid UID) bool {
	charmanager.AddUser(uid)
	return true
}
