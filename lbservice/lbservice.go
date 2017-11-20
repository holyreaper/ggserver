package lbservice

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/holyreaper/ggserver/lbmodule/funcall"
)

//LBService ...
type LBService struct {
	waitGroup *sync.WaitGroup
	exitCh    chan struct{}
	funcCall  *funcall.FunCall
	listen    *net.TCPListener
}

//NewLBService new
func NewLBService() *LBService {
	return &LBService{
		exitCh:    make(chan struct{}),
		funcCall:  funcall.NewFuncCall(),
		waitGroup: &sync.WaitGroup{},
	}
}

//init 初始化
func (lb *LBService) init() bool {
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.1.177:8091")
	lb.listen, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("LbServer Listen error ", err)
		return false
	}
	return true
}

//Start start service
func (lb *LBService) Start() {
	if !lb.init() {
		return
	}
	for {
		lb.listen.SetDeadline(time.Now().Add(100))
		cnn, err := lb.listen.Accept()
		if err != nil {

		}
		cnn.Close()
	}

}
