package def

import (
	"net"
	"syscall"
)

// ServerType 服务器类型
type ServerType int32

const (
	//ServerTypeNormal lobby服务器
	ServerTypeNormal ServerType = iota + 1

	//ServerTypeDB db服务器
	ServerTypeDB

	//ServerTypeCenter 中心服务器
	ServerTypeCenter
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
	if operr, ok := err.(*net.OpError); ok && operr.Timeout() {
		return true
	}
	return false
}

//IsFdBlock check net block err
func IsFdBlock(err error) bool {
	if operr, ok := err.(*net.OpError); ok {
		if numerr, ok2 := operr.Err.(syscall.Errno); ok2 {
			if numerr == syscall.EAGAIN || numerr == syscall.EWOULDBLOCK || numerr == syscall.EINTR {
				return true
			}
		}
	}
	return false
}

//GetServerType ..
func GetServerType(id SID) ServerType {
	s := ServerType(id / SERVERBASEVALUE)
	return s
}
