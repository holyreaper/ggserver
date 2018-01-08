package rpclog

import (
	"github.com/holyreaper/ggserver/glog"
)

var grpcLog *glog.GLog

func init() {
	grpcLog = glog.NewGLog("grpc_")
}

//GetLogger  GetLogger
func GetLogger() *glog.GLog {
	return grpcLog
}
