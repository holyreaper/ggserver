package rpcservice

import (
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc/keepalive"

	"github.com/holyreaper/ggserver/consul"
	"github.com/holyreaper/ggserver/def"
	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/rpcservice/ctrpc"
	"github.com/holyreaper/ggserver/rpcservice/dbrpc"
	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"
	"github.com/holyreaper/ggserver/rpcservice/pb/dbrpc"
	"github.com/holyreaper/ggserver/util"
	"google.golang.org/grpc"
)

func init() {

}

// GGService 服务
type GGService struct {
	st     def.ServerType
	id     def.SID
	server *grpc.Server
}

var gRPCService *GGService

//Init Init Rpcservice
func Init(id def.SID) {
	gRPCService = &GGService{id: id,
		st: util.GetServerType(id),
	}

}

// Start start service
func Start() (err error) {

	info, err := consul.GetSingleServerInfo(gRPCService.id)
	if err != nil {
		LogFatal("gRPCService GetSingleServerInfo fail %s ", err)
		return
	}
	LogInfo("start rpc service addr %s:%d ", info.Address, info.Port)
	listen, err := net.Listen("tcp", info.Address+":"+strconv.Itoa(info.Port))
	if err != nil {
		fmt.Println("listen fail ....")
		return
	}
	keepaliveParam := grpc.KeepaliveParams(keepalive.ServerParameters{})

	gRPCService.server = grpc.NewServer(keepaliveParam)

	gRPCService.RegisterModule()

	LogInfo("rpc server type %d start !!!", gRPCService.st)

	go gRPCService.server.Serve(listen)

	return
}

//Stop stop the server
func Stop() {
	if gRPCService.server != nil {
		gRPCService.server.Stop()
		LogInfo("rpc service %d stop !!!", gRPCService.id)
	}
}

//RegisterModule 注册服务
func (s *GGService) RegisterModule() {
	if s.st == def.ServerTypeNormal {

	} else if s.st == def.ServerTypeDB {
		dbrpcpt.RegisterDBRPCServer(s.server, &dbrpc.DBRPC{})
	} else if s.st == def.ServerTypeCenter {
		ctrpcpt.RegisterCTRPCServer(s.server, &ctrpc.CTRPC{})
	}

}
