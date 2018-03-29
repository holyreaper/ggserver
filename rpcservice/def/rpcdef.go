package rpcdef

//RPCPushType rpc推送协议的枚举类型
type RPCPushType int32

const (
	//RPCPTypeUserLogin user login
	RPCPTypeUserLogin RPCPushType = iota + 1
	//RPCPTypeUserLogout user logout
	RPCPTypeUserLogout
	//RPCPTypeUserChat user chat
	RPCPTypeUserChat
)
