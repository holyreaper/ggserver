package main

import (
	//"context"
	"fmt"

	"github.com/holyreaper/ggserver/conf"

	"runtime"

	"net"

	"github.com/holyreaper/ggserver/pb"
	"github.com/holyreaper/ggserver/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":8090"
	msg  = "helo lhm"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	service.Hello()
	conf := conf.GetConf()
	if data, ok := conf["info"]; ok {
		sliceArray := data.([]interface{})
		for k, v := range sliceArray {
			fmt.Printf("key %d, value %s", k, v)
		}
	}

	defer func() {
		println("finish ...")
	}()
	listen, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("listen fail ....")
		panic(-1)
	}

	LoginServer := grpc.NewServer()

	login.RegisterLoginServer(LoginServer, &LoginServerExt{})
	LoginServer.Serve(listen)
}

//LoginServerExt struct
type LoginServerExt struct {
}

/*
SayHello ...
helo world
*/
func (s *LoginServerExt) SayHello(context.Context, *login.HelloRequest) (*login.HelloReply, error) {
	return &login.HelloReply{Message: msg}, nil
}
