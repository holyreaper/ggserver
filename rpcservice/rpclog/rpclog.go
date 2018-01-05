package rpclog

import (
	"github.com/holyreaper/ggserver/glog"
)

var grpcLog *glog.GLog

func init() {
	grpcLog = glog.NewGLog("grpc_")
}

//GetRPCLogger getrpcLogger
func GetRPCLogger() *glog.GLog {
	return grpcLog
}
