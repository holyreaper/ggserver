package rpcservice

import (
	"fmt"
	"net"

	"google.golang.org/grpc/keepalive"

	"github.com/holyreaper/ggserver/def"
	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/rpcservice/ctrpc"
	"github.com/holyreaper/ggserver/rpcservice/dbrpc"
	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"
	"github.com/holyreaper/ggserver/rpcservice/pb/dbrpc"
	"google.golang.org/grpc"
)

func init() {

}

// GGService 服务
type GGService struct {
	st def.ServerType
}

var gRPCService *GGService

//Init Init Rpcservice
func Init(serverType def.ServerType) {
	gRPCService = &GGService{st: serverType}
}

// Start start service
func Start(exitCh <-chan bool) {

	listen, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("listen fail ....")
		return
	}

	keepaliveParam := grpc.KeepaliveParams(keepalive.ServerParameters{})

	server := grpc.NewServer(keepaliveParam)

	gRPCService.RegisterModule(server)

	LogInfo("rpc server type %d start !!!", gRPCService.st)
	go server.Serve(listen)

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" GGServer Start Error : %s ", err)
		}
		server.Stop()
		LogInfo("rpc server type %d stop !!!", gRPCService.st)
	}()
	<-exitCh

}

//RegisterModule 注册服务
func (s *GGService) RegisterModule(rpcServer *grpc.Server) {
	if s.st == def.ServerTypeNormal {

	} else if s.st == def.ServerTypeDB {
		dbrpcpt.RegisterDBRPCServer(rpcServer, &dbrpc.DBRPC{})
	} else if s.st == def.ServerTypeCenter {
		ctrpcpt.RegisterCTRPCServer(rpcServer, &ctrpc.CTRPC{})
	}

}
