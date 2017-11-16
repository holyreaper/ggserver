package service

import (
	"fmt"
	"net"

	"github.com/holyreaper/ggserver/pb/chat"
	"github.com/holyreaper/ggserver/pb/login"
	"github.com/holyreaper/ggserver/service/chat"
	"github.com/holyreaper/ggserver/service/login"

	"github.com/holyreaper/ggserver/def"
	"golang.org/x/net/context"
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

	//client
	go Client()
}

//Client ...
func Client() {
	fmt.Println("client start ...")
	client, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure())
	if err != nil {
		fmt.Println("client exit error msg ", err)
		return
	}
	defer client.Close()
	cnn := chatrpc.NewChatRpcClient(client)

	msg, err := cnn.Chat(context.Background(), &chatrpc.ChatMsgRequest{Name: "xiaodian"})
	if err != nil {
		fmt.Println("client get msg error ", err)
	} else {
		fmt.Println("client get msg :", *msg)
	}
}

//RegisterModule 注册服务
func (s *GGService) RegisterModule(rpcServer *grpc.Server) {
	if s.st == def.ServerTypeNormal {
		loginrpc.RegisterLoginRpcServer(rpcServer, &login.Login{})
		chatrpc.RegisterChatRpcServer(rpcServer, &chat.Chat{})

	} else if s.st == def.ServerTypeDB {

	} else {

	}

}

//NewGGService new 服务
func NewGGService(serverType def.ServerType) *GGService {
	return &GGService{st: serverType}
}
