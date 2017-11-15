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
