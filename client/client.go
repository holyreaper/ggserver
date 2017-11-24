package client

import (
	"fmt"
	"io"
	"net"

	"github.com/holyreaper/ggserver/lbmodule/packet"

	"github.com/holyreaper/ggserver/util/convert"

	"github.com/golang/protobuf/proto"
	"github.com/holyreaper/ggserver/lbmodule/pb/user"
	"github.com/holyreaper/ggserver/rpcservice/pb/chat"
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
	ex := make(chan bool)
	go UserClient(ex)

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

//UserClient client
func UserClient(ex <-chan bool) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8091")
	if err != nil {
		fmt.Println("resolve tcp addr error ", err)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("client return  ", err)
		}
	}()
	var (
		bLen  = make([]byte, 4)
		bType = make([]byte, 4)
	)
	cnn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("cnn server failed ", err)
		return
	}
	fmt.Println("client cnn server succ !")
	//send data
	req := ptuser.LoginMsgRequest{Uid: 1001}
	data, err := proto.Marshal(&req)
	if err != nil {
		fmt.Println("mashal login request fail !")
		return
	}
	var reqData []byte
	lenSlice := convert.Int32ToBytes(int32(len(data) + 8))

	fmt.Println("client start send lenSlice !", len(data), data)
	reqData = append(reqData, lenSlice...)
	typeSlice := convert.Int32ToBytes(int32(packet.PKGLogin))
	reqData = append(reqData, typeSlice...)
	fmt.Println("client start send data !", reqData)
	reqData = append(reqData, data...)
	ln, err := cnn.Write(reqData)
	if err != nil || ln != len(reqData) {
		fmt.Println("send login request fail !")
		return
	}
	fmt.Println("client send data succ  !")

	//recv data
	fmt.Println("client read dataing ... !")
	ln, err = io.ReadFull(cnn, bLen)
	if err != nil || ln != 4 {
		fmt.Println("read len  fail !")
		return
	}
	fmt.Println("client read len succ ... !")
	ln, err = io.ReadFull(cnn, bType)
	if err != nil || ln != 4 {
		fmt.Println("read len  fail !")
		return
	}
	fmt.Println("client read type succ ... !")
	ulen := convert.BytesToInt32(bLen)

	bData := make([]byte, ulen)
	ln, err = io.ReadFull(cnn, bData)
	if err != nil || ln != int(ulen) {
		fmt.Println("read bData  fail !")
		return
	}
	if convert.BytesToInt32(bType) == packet.PKGLogin {
		rsp := ptuser.LoginMsgReply{}
		err := proto.Unmarshal(bData, &rsp)
		if err != nil {
			fmt.Println("unmashal data fail ")
		}
		//unmashal succ
		fmt.Printf("unmashal data succ result :%d", rsp.GetResult())
	} else {
		fmt.Println("unknown message type ")
	}
	cnn.Close()
	<-ex
	return
}
