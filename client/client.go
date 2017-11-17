package client

import (
	"fmt"

	"github.com/holyreaper/ggserver/pb/chat"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() {

}

//Start client
func Start() {
	fmt.Println("client start ...")
	client, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure())
	if err != nil {
		fmt.Println("client exit error msg ", err)
		return
	}

	defer func() {
		client.Close()

	}()

	//chat rpc
	cnn := chatrpc.NewChatRpcClient(client)

	msg, err := cnn.Chat(context.Background(), &chatrpc.ChatMsgRequest{Name: "xiaodian"})
	if err != nil {
		fmt.Println("client get msg error ", err)
		return
	}
	fmt.Println("client get msg :", *msg)

	//stream chat list

	fmt.Println("start chatlist stream... ")

	streamClient, err := cnn.ChatList(context.Background())
	if err != nil {
		fmt.Println("call chatlist error ", err)
		return
	}
	err = streamClient.Send(&chatrpc.ChatMsgRequest{Name: "helo"})
	if err != nil {
		fmt.Println("call chatlist send  error ", err)
		return
	}

	for {
		msg, err := streamClient.Recv()
		if err != nil {
			fmt.Println("call chatlist recv error ", err)
			return
		}
		fmt.Println("recv msg result :", msg.Result)
	}

}
