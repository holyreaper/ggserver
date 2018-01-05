package def

import "net"

// ServerType 服务器类型
type ServerType int32

const (
	//ServerTypeNormal lobby服务器
	ServerTypeNormal ServerType = iota + 1

	//ServerTypeDB db服务器
	ServerTypeDB

	//ServerTypeProxy 代理服务器
	ServerTypeProxy
)

//UID 玩家唯一标识
type UID int64

// SID 	服务器id
type SID int

// SERVERBASEVALUE 服务器id/SERVERBASEVALUE = 服务器类型 ServerType
const (
	SERVERBASEVALUE = 10000
)

//IsTimeOut check net err
func IsTimeOut(err error) bool {
	if NetErr, ok := err.(*net.OpError); ok && NetErr.Timeout() {
		return true
	}
	return false
}
