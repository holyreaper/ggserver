package lbnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	. "github.com/holyreaper/ggserver/glog"

	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbmodule/funcall"
	_ "github.com/holyreaper/ggserver/lbmodule/logic/chat"
	_ "github.com/holyreaper/ggserver/lbmodule/logic/user"
	"github.com/holyreaper/ggserver/lbmodule/packet"
	"github.com/holyreaper/ggserver/lbmodule/pb/message"
	"github.com/holyreaper/ggserver/util/convert"
)

//LBNet net
type LBNet struct {
	listen       *net.TCPListener
	exitCh       chan bool
	waitGroup    *sync.WaitGroup
	readTimeOut  time.Duration
	writeTimeOut time.Duration
	oc           sync.Once
}

//Sender 发送者
type Sender interface {
	Send(packet.Packet) error
}

//Receiver 接收
type Receiver interface {
	Receive() (packet.Packet, error)
}

//ClientCnn 客户端的连接
type ClientCnn struct {
	//接收缓冲区
	receivePacket chan packet.Packet
	//待发送缓冲
	sendPacket chan packet.Packet
	//标记
	exitCh chan bool
	//地址
	addr string
	//cnn
	cnn *net.TCPConn
	//one
	one sync.Once
}

//Send send data
func (client *ClientCnn) Send(pack packet.Packet) error {
	client.sendPacket <- pack
	defer func() {
		if err := recover(); err != nil {
			LogFatal("client send %s fail err %s", client.addr, err)
		}
	}()
	return nil
}

//Receive receive data
func (client *ClientCnn) Receive() (packet.Packet, error) {
	pack, ok := <-client.receivePacket
	if !ok {
		LogFatal("client Receive %s fail chan has been closed", client.addr)
		return packet.Packet{}, errors.New("chan has been closed")
	}
	return pack, nil
}

//HandlePacket .
func (client *ClientCnn) HandlePacket(exitCh chan bool) {
	var uid int64
	fmt.Println("HandlePacket ing ...")
	for {
		select {
		case <-client.exitCh:
			fmt.Printf("Stop HandlePacket \r\n")
			return
		case <-exitCh:
			fmt.Println("Stop Handle Send Packet exit flag ")
			return
		case p := <-client.receivePacket:
			{
				var err error
				var ret []reflect.Value
				var lmsg message.Message
				err = proto.Unmarshal(p.Data, &lmsg)
				if err != nil {
					fmt.Println("message unmarshal fail ", err)
					continue
				}
				if p.Type == packet.PKGLogin {
					ret, err = funcall.Call(p.Type, lmsg)
					uid = lmsg.Uid
				} else {
					if uid <= 0 {
						//illegal method
						LogFatal("net have an illegal request disconnect %s", client.GetAddr())
						p.Clear()
						lmsg.LogoutReply.Result = "illegal request disconnect !!"
						p.Pack(packet.PkLogOut, &lmsg)
						//client.sendPacket <- p
						client.Send(p)
						time.Sleep(1)
						exitCh <- true
						return
					}
					lmsg.Uid = uid
					ret, err = funcall.Call(p.Type, lmsg)
				}
				if err != nil {
					fmt.Printf("func call type  %d error %s \r\n", p.Type, err)
					continue
				}
				for _, value := range ret {
					v := value.Interface()
					//tmppack[index] = v.(packet.Packet)
					//spack <- tmppack[index]
					tmp := v.(message.Message)
					tmppack := packet.Packet{}
					tmppack.Pack(tmp.Command, &tmp)
				}
			}
		}

	}
}

//HandleSend .
func (client *ClientCnn) HandleSend(exitCh chan bool) {
	for {
		select {
		case <-client.exitCh:
			fmt.Printf("Stop Handle Send Packet exit ch \r\n")
			return
		case <-exitCh:
			fmt.Println("Stop Handle Send Packet exit flag ")
			return
		case p := <-client.sendPacket:
			fmt.Printf("send packeting ... %d len %d ", p.Type, p.Len)
			buff := p.FormatBuf()
			currSend := 0
			for {
				ln, err := client.cnn.Write(buff[currSend:])
				if err != nil && ln != len(buff[currSend:]) {
					if def.IsTimeOut(err) {
						currSend += ln
						continue
					} else {
						fmt.Println("send byte to client  err ", err)
						return
					}
				}
				if currSend+ln >= len(buff) {
					break
				}
			}
		}

	}
}

//HandleReceive .
func (client *ClientCnn) HandleReceive(exitCh chan bool) {
	var (
		pLen    = make([]byte, 4)
		pType   = make([]byte, 4)
		pPacLen int32
		pacData = make([]byte, packet.MAXPACKETLEN)

		cLen    = 0   //curr read pLen
		cType   = 0   //curr read pType
		cPacLen int32 //curr read pPackLen

	)
	for {
		if client.IsClosed() {
			LogInfo("client handlereceive %s exit", client.addr)
			return
		}
		client.SetReadDeadline(time.Now().Add(30 * time.Second))
		if cLen < 4 {
			if n, err := io.ReadFull(client.cnn, pLen[cLen:]); err != nil && n != 4 {
				if def.IsTimeOut(err) {
					cLen += n
					fmt.Printf("Read pLen time out  \r\n")
					continue
				}
				fmt.Printf("Read pLen failed: %v\r\n", err)
				return //exit
			}
			cLen += 4
			if pPacLen = convert.BytesToInt32(pLen); pPacLen > packet.MAXPACKETLEN {
				fmt.Printf("pacLen larger than pPacLen\r\n")
				return
			}
		}
		if cType < 4 {
			if n, err := io.ReadFull(client.cnn, pType); err != nil && n != 4 {
				if err == io.EOF {
					LogInfo("net ioread fail cnn %s hase been closed!! ", client.cnn.RemoteAddr().String())
				} else if def.IsTimeOut(err) {
					cType += n
					fmt.Printf("Read pType time out  \r\n")
					continue
				}
				fmt.Printf("Read pType failed: %v\r\n", err)
				return
			}
			cType += 4
		}
		if cPacLen < pPacLen {
			if n, err := io.ReadFull(client.cnn, pacData[cPacLen:pPacLen-8]); err != nil && n != int(pPacLen-8-cPacLen) {
				if err == io.EOF {
					LogInfo("net ioread fail cnn %s hase been closed!! ", client.cnn.RemoteAddr().String())
				} else if def.IsTimeOut(err) {
					cPacLen += int32(n)
					fmt.Printf("Read pacData time out  \r\n")
					continue
				}
				fmt.Printf("Read pacData failed: %v\r\n", err)
				return
			}
			cPacLen += pPacLen
		}

		fmt.Println("HandleCnn get full packet data   ...")
		client.receivePacket <- packet.Packet{
			Len:  pPacLen,
			Type: convert.BytesToInt32(pType),
			Data: pacData[:pPacLen-8],
		}
		cLen = 0
		cType = 0
		cPacLen = 0
	}
}

//Close close
func (client *ClientCnn) Close() {
	client.one.Do(func() {
		close(client.exitCh)
	})
}

//GetAddr addr
func (client *ClientCnn) GetAddr() string {
	return client.addr
}

//SetReadDeadline 阻塞时间
func (client *ClientCnn) SetReadDeadline(t time.Time) error {
	return client.cnn.SetReadDeadline(t)
}

//IsClosed addr
func (client *ClientCnn) IsClosed() bool {
	select {
	case <-client.exitCh:
		return false
	default:
	}
	return true
}

const (
	//MaxSend .
	MaxSend int32 = 100
	//MaxRead .
	MaxRead int32 = 100
	//MaxReadWriteTime 超时时间
	MaxReadWriteTime int32 = 10
)

//NewLBNet new lbnet
func NewLBNet() *LBNet {
	return &LBNet{
		exitCh:       make(chan bool),
		listen:       &net.TCPListener{},
		waitGroup:    &sync.WaitGroup{},
		readTimeOut:  time.Duration(30),
		writeTimeOut: time.Duration(30),
	}
}

// Start deal cnn
func (lbnet *LBNet) Start() {
	LogInfo("lbnet start ...")
	defer func() {
		if err := recover(); err != nil {
			LogFatal("lbnet  have panic ", err)
		}
		LogInfo("lbnet end ...")
		lbnet.Stop()
	}()
	for {
		select {
		case <-lbnet.exitCh:
			break
		default:
		}
		lbnet.listen.SetDeadline(time.Now().Add(lbnet.readTimeOut * time.Second))
		cnn, err := lbnet.listen.AcceptTCP()
		if err != nil {
			if def.IsTimeOut(err) {
				continue
			} else {
				lbnet.Stop()
			}
		}
		LogInfo("have a Cnn accept start to serve it  ")
		go lbnet.HandleCnn(cnn)
	}

}

// Stop deal cnn
func (lbnet *LBNet) Stop() {
	lbnet.oc.Do(func() {
		close(lbnet.exitCh)
		lbnet.waitGroup.Wait()
	})
}

// Init deal cnn
func (lbnet *LBNet) Init(nettype string /*tcp*/, ipaddr string) (err error) {
	tcpAddr, err := net.ResolveTCPAddr(nettype, ipaddr)
	if err != nil {
		LogFatal("lbserver ResolveTCPAddr  addr %s error %s   ", ipaddr, err)
		return
	}
	lbnet.listen, err = net.ListenTCP(nettype, tcpAddr)
	if err != nil {
		LogFatal("lbserver Listen error %s", err)
		return
	}
	return nil
}

//HandleCnn handle cnn
func (lbnet *LBNet) HandleCnn(tpcCnn *net.TCPConn) {
	lbnet.waitGroup.Add(1)
	client := &ClientCnn{
		receivePacket: make(chan packet.Packet, MaxRead),
		sendPacket:    make(chan packet.Packet, MaxSend),
		exitCh:        make(chan bool),
		addr:          tpcCnn.RemoteAddr().String(),
		cnn:           tpcCnn,
		one:           sync.Once{},
	}
	defer func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("cnn panic  from addr ", tpcCnn.RemoteAddr().String())
			}
		}()
		client.Close()
		tpcCnn.Close()
		lbnet.waitGroup.Done()
	}()
	go lbnet.HandleSend(client)
	go lbnet.HandleReceive(client)
	go lbnet.HandlePacket(client)
}

//HandleReceive handle receive
func (lbnet *LBNet) HandleReceive(client *ClientCnn) {
	LogInfo("start service for  the new client  %s", client.GetAddr())

}

//HandlePacket handle packet
func (lbnet *LBNet) HandlePacket(client *ClientCnn) {
	lbnet.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle deal packet panic error %v \r\n", err)
		}
		lbnet.waitGroup.Done()
	}()
	client.HandlePacket(lbnet.exitCh)
}

//HandleSend handle send packet...
func (lbnet *LBNet) HandleSend(client *ClientCnn) {
	lbnet.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle send packet panic error %v \r\n", err)
		}
		lbnet.waitGroup.Done()
	}()
}
