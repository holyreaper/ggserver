package funcall

// FuncCallEnum .
type FuncCallEnum uint32

const (
	//FCHeartBeat  心跳包
	FCHeartBeat FuncCallEnum = 100001

	//FCUserModule user 模块
	FCUserModule = FCHeartBeat + 100000
	//FCChatModule chat 模块
	FCChatModule = FCHeartBeat + 200000
)
