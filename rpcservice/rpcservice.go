package rpcservice

import (
	"fmt"
	"net"

	"github.com/holyreaper/ggserver/def"
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

// Start start service
func (s *GGService) Start() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" GGServer Start Error : %s ", err)
		}
	}()
	listen, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("listen fail ....")
		panic(-1)
	}
	LoginServer := grpc.NewServer()
	s.RegisterModule(LoginServer)
	fmt.Println("server start ...")
	go LoginServer.Serve(listen)

}

//RegisterModule 注册服务
func (s *GGService) RegisterModule(rpcServer *grpc.Server) {
	if s.st == def.ServerTypeNormal {

	} else if s.st == def.ServerTypeDB {
		dbrpcpt.RegisterDBRPCServerServer(rpcServer, &dbrpc.DBRPC{})
	} else if s.st == def.ServerTypeCenter {
		ctrpcpt.RegisterCTRPCServerServer(rpcServer, &ctrpc.CTRPC{})
	}

}

//NewGGService new 服务
func NewGGService(serverType def.ServerType) *GGService {
	return &GGService{st: serverType}
}
