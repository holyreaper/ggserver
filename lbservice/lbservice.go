package lbservice

import (
	"fmt"

	"github.com/holyreaper/ggserver/lbmodule/net"
)

//LBService ...
type LBService struct {
	exitCh chan struct{}
	lbnet  *lbnet.LBNet
}

//NewLBService new
func NewLBService() *LBService {
	return &LBService{
		exitCh: make(chan struct{}),
		lbnet:  lbnet.NewLBNet(),
	}
}

//init 初始化
func (lb *LBService) init() bool {
	bindip := "127.0.0.1:8091"
	lb.lbnet.Init("tcp", "tcp", bindip)
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
	lb.lbnet.Stop()
	fmt.Println("lbserver stop ...")
}
