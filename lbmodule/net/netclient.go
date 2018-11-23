package lbnet

import (
	"errors"
	"fmt"
	"net"
	"time"

	. "github.com/holyreaper/ggserver/glog"
)

/* ----------------------------------------------------------------
	客户端连接服务器
---------------------------------------------------------------- */
//Client .
type Client interface {
	Cnn() error
	Close()
	Status() int32
	Init(addr string) error
	KeepAlive(bool) bool
}

//TCPClient .
type TCPClient struct {
	Connection
	keepAlive bool
	exitCh    *BoolChan
}

//Cnn .
func (client *TCPClient) Cnn() error {
	ipAddr, err := net.ResolveTCPAddr("tcp", client.addr)
	if err != nil {
		return err
	}
	client.cnn, err = net.DialTCP("tcp", nil, ipAddr)
	if err != nil {
		return err
	}
	return nil
}

//GoKeepAlive keepalive
func (client *TCPClient) GoKeepAlive() {
	var i = 0
	for {
		time.AfterFunc(1*time.Second, func() {
			i++
			if i >= 30 {
				//send keepAlive 包
			}
		})
		select {
		case <-client.exitCh.GetChan():
			LogInfo("net client keepalive exit")
			return
		default:
		}
	}
}

//Close .
func (client *TCPClient) Close() {
	if client.cnn != nil {
		client.cnn.Close()
		client.cnn = nil
	}
	client.exitCh.Close()
}

//KeepAlive keep
func (client *TCPClient) KeepAlive(keepAlive bool) bool {
	client.keepAlive = keepAlive
	return true
}

//Init .
func (client *TCPClient) Init(addr string, keepAlive bool) error {
	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	client.addr = addr
	client.exitCh = NewBoolChan(0)
	return nil
}

//ProxyClient 代理服务器的处理
type ProxyClient struct {
	TCPClient
}

//HandlePacket 代理服务器特化处理
func (client *ProxyClient) HandlePacket(chch CheckChan) (err error) {
	fmt.Println("HandlePacket ing ...")
	for {
		if !chch.CheckChan() {
			err = errors.New("exit packeting ")
			goto END
		}
		if !client.CheckChan() {
			err = errors.New("exit packeting ")
			goto END
		}
		select {
		case p := <-client.receivePacket:
			{
				//找到 这个连接 绑定的service然后转发
				LogDebug("proxy client receive package %d", p.GetType())
			}
		}
	}
END:
	return err
}
