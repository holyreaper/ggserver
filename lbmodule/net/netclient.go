package lbnet

import (
	"context"
	"net"
	"sync"
)

/* ----------------------------------------------------------------
	客户端连接服务器
---------------------------------------------------------------- */
type ServerCnn struct {
	addr      string
	port      int32
	cnn       *net.TCPConn
	ctx       context.Context
	waitGroup sync.WaitGroup
}
