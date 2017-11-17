package def

// ServerType 服务器类型
type ServerType int

const (
	//ServerTypeNormal lobby服务器
	ServerTypeNormal ServerType = iota

	//ServerTypeDB db服务器
	ServerTypeDB

	//ServerTypeProxy 代理服务器
	ServerTypeProxy
)

//UID 玩家唯一标识
type UID int64

// SID 	服务器id
type SID int
