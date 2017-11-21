package chat

import (
	"github.com/holyreaper/ggserver/lbmodule/funcall"
	. "github.com/holyreaper/ggserver/lbmodule/package"
)

func init() {
	funcall.BindFunc(PKGChat, Chat)
}

//Chat chat
func Chat(uid uint32) bool {
	return true
}
