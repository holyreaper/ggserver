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
	lbnet lbnet.Server
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
	err = gLBService.init("0.0.0.0:" + strconv.Itoa(info.Port))
	if err != nil {
		return
	}
	return nil
}

//Start start service
func Start() (err error) {
	go gLBService.lbnet.Start()
	return
}

//Stop service
func Stop() {
	if gLBService != nil {
		gLBService.Stop()
	}
}

//init 初始化
func (lb *LBService) init(addr string) error {
	return lb.lbnet.Init(addr)
}

//Stop stop server
func (lb *LBService) Stop() {
	lb.lbnet.Stop()
	fmt.Println("lbserver stop ...")
}
