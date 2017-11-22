package lbservice

import (
	"net"
	"fmt"
	"sync"
	"github.com/holyreaper/ggserver/lbmodule/net"
)

//LBService ...
type LBService struct {
	exitCh    chan struct{}
	lbnet     *lbnet.LBNet
}

//NewLBService new
func NewLBService() *LBService {
	return &LBService{
		exitCh:    make(chan struct{}),
		lbnet :&lbnet.LBNet{
			exitCh: make(chan struct{})
			listen: &net.TCPListener{}
			waitGroup: &sync.WaitGroup{},
			readTimeOut :  time.Duration(30),
			writeTimeOut:  time.Duration(30),

		}
	}
}

//init 初始化
func (lb *LBService) init() bool {
	bindip := "192.168.1.177:8091"
	lb.lbnet.Init("ipv4", "tcp", bindip)
	return true
}

//Start start service
func (lb *LBService) Start() {
	if !lb.init() {
		return
	}
	defer func() {
		lb.lbnet.Stop()
	}()
	lb.lbnet.Start()

}

//Stop stop server
func (lb *LBService) Stop() {
	close(lb.exitCh)
	lb.lbnet.Stop();
	fmt.Println("lbserver stop ...")
}
