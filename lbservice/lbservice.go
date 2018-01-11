package lbservice

import (
	"fmt"
	"strconv"

	"github.com/holyreaper/ggserver/consul"
	. "github.com/holyreaper/ggserver/def"
	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/lbmodule/net"
)

//LBService ...
type LBService struct {
	lbnet *lbnet.LBNet
	id    SID
}

var gLBService *LBService

//Init ...
func Init(sid SID) (err error) {
	gLBService = &LBService{
		lbnet: lbnet.NewLBNet(),
		id:    sid,
	}
	info, err := consul.GetSingleServerInfo(gLBService.id)
	if err != nil {
		LogFatal("lbserver Init fail %s ", err)
		return
	}
	err = gLBService.init(info.Address + ":" + strconv.Itoa(info.Port))
	if err != nil {
		return
	}
	return nil
}

//init 初始化
func (lb *LBService) init(addr string) error {
	return lb.lbnet.Init("tcp", addr)
}

//Start start service
func Start(exitCh <-chan bool) (err error) {
	go gLBService.lbnet.Start()

	defer func() {
		gLBService.Stop()
	}()

	<-exitCh
	return
}

//Stop stop server
func (lb *LBService) Stop() {
	lb.lbnet.Stop()
	fmt.Println("lbserver stop ...")
}
