package lbservice

import (
	"fmt"

	"github.com/holyreaper/ggserver/lbmodule/net"
)

//LBService ...
type LBService struct {
	lbnet *lbnet.LBNet
}

//NewLBService new
func NewLBService() *LBService {
	return &LBService{
		lbnet: lbnet.NewLBNet(),
	}
}

//init 初始化
func (lb *LBService) init() error {
	bindip := "127.0.0.1:8091"
	return lb.lbnet.Init("tcp", bindip)
}

//Start start service
func (lb *LBService) Start(exitCh <-chan bool) (err error) {
	err = lb.init()
	if err != nil {
		return
	}
	defer func() {
		lb.Stop()
	}()
	go lb.lbnet.Start()
	<-exitCh
	return
}

//Stop stop server
func (lb *LBService) Stop() {
	lb.lbnet.Stop()
	fmt.Println("lbserver stop ...")
}
