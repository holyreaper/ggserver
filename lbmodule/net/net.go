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
	"github.com/holyreaper/ggserver/def"
	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/lbmodule/funcall"
	"github.com/holyreaper/ggserver/lbmodule/packet"
	"github.com/holyreaper/ggserver/lbmodule/pb/message"
	"github.com/holyreaper/ggserver/util"
	"github.com/holyreaper/ggserver/util/convert"
)

/* ----------------------------------------------------------------
	服务器处理客户端的连接
---------------------------------------------------------------- */
//BoolChan 简单的chan包装
type BoolChan struct {
	exitChan chan bool
	one      sync.Once
}

//NewBoolChan new chan
func NewBoolChan(len int32) *BoolChan {
	return &BoolChan{
		exitChan: make(chan bool, len),
		one:      sync.Once{},
	}
}

//Close 退出
func (bc *BoolChan) Close() {
	bc.one.Do(func() {
		close(bc.exitChan)
	})
}

//GetChan .
func (bc *BoolChan) GetChan() chan bool {
	return bc.exitChan
}

//CheckChan 是否关闭
type CheckChan interface {
	CheckChan() bool
}

//Address addr
type Address interface {
	Addr() string
}

//HandleSender 发送者
type HandleSender interface {
	HandleSend(CheckChan) error
}

//HandleReceiver 接收
type HandleReceiver interface {
	HandleReceive(CheckChan) error
}

//HandlePackger 打包
type HandlePackger interface {
	HandlePacket(CheckChan) error
}

//Connection 客户端连接服务器
type Connection struct {
	//接收缓冲区
	receivePacket chan packet.Packet
	//待发送缓冲
	sendPacket chan packet.Packet
	//标记
	exitCh *BoolChan
	//地址
	addr string
	//cnn
	cnn *net.TCPConn
	//连接id
	id uint64
}

//包处理 逻辑
func (client *Connection) dealPacket(p packet.Packet) (err error) {
	var ret []reflect.Value
	var lmsg message.Message
	err = proto.Unmarshal(p.Data, &lmsg)
	if err != nil {
		LogFatal("message unmarshal fail ", err)
		return
	}
	if lmsg.GetUid() <= 0 && lmsg.GetCommand() != uint32(funcall.FCUserModule+1) { //注册
		LogFatal("net client connection illegal user ")
		err = errors.New("net client connection illegal user")
		return
	}
	ret, err = funcall.Call(p.Type, lmsg)
	for _, value := range ret {
		_ = value.Interface()
		//tmppack[index] = v.(packet.Packet)
		//spack <- tmppack[index]
		//tmp := v.(message.Message)
		//tmppack := packet.Packet{}
		//tmppack.Pack(tmp.Command, &tmp)
	}
	return nil
}

//HandlePacket .
func (client *Connection) HandlePacket(chch CheckChan) (err error) {
	fmt.Println("HandlePacket ing ...")
	for {
		if !client.CheckChan() {
			err = errors.New("handlePacket exit")
			goto END
		}
		if !chch.CheckChan() {
			err = errors.New("handlePacket exit")
			goto END
		}
		select {
		case p := <-client.receivePacket:
			{
				err = client.dealPacket(p)
				if err != nil {
					goto END
				}
			}
		}
	}
END:
	return err
}

//HandleSend .
func (client *Connection) HandleSend(chch CheckChan) (err error) {

	for {
		if !client.CheckChan() {
			err = errors.New("handle send exit")
			goto END
		}
		if !chch.CheckChan() {
			err = errors.New("handle send exit")
			goto END
		}
		select {
		case p := <-client.sendPacket:
			fmt.Printf("send packeting ... %d len %d ", p.Type, p.Len)
			buff := p.FormatBuf()
			currSend := 0
			for {
				client.SetWriteDeadline(time.Now().Add(MaxReadWriteTime * time.Second))
				ln, err := client.cnn.Write(buff[currSend:])
				if err != nil {
					if !def.IsTimeOut(err) && !def.IsFdBlock(err) {
						fmt.Println("send byte to client  err  ", err.Error())
						return errors.New("handleSendt exit")
					}
				}
				currSend += ln
				if currSend >= len(buff) {
					break
				}
			}
		}
	}
END:
	return err
}

//HandleReceive .
func (client *Connection) HandleReceive(chch CheckChan) (err error) {
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
		if !chch.CheckChan() {
			fmt.Printf("Stop Handle Send Packet exit ch \r\n")
			err = errors.New("exit chan ")
			goto END
		}
		if !client.CheckChan() {
			fmt.Printf("Stop Handle Send Packet exit ch \r\n")
			goto END
		}

		if cLen < 4 {
			client.SetReadDeadline(time.Now().Add(MaxReadWriteTime * time.Second))
			if n, err := io.ReadFull(client.cnn, pLen[cLen:]); n < 4 {
				if def.IsTimeOut(err) || def.IsFdBlock(err) {
					cLen += n
					fmt.Printf("Read pLen time out or block \r\n")
					continue
				}
				fmt.Printf("Read pLen failed: %s\r\n", err.Error())
				err = errors.New("exit chan ")
				goto END
			}
			cLen += 4
			if pPacLen = convert.BytesToInt32(pLen); pPacLen > packet.MAXPACKETLEN {
				fmt.Printf("pacLen larger than pPacLen\r\n")
				err = errors.New("exit chan ")
				goto END
			}
		}
		if cType < 4 {
			client.SetReadDeadline(time.Now().Add(MaxReadWriteTime * time.Second))
			if n, err := io.ReadFull(client.cnn, pType); n < 4 {
				if def.IsTimeOut(err) || def.IsFdBlock(err) {
					cType += n
					fmt.Printf("Read pType time out or block \r\n")
					continue
				}
				err = errors.New("exit chan ")
				goto END
			}
			cType += 4
		}
		if cPacLen < pPacLen {
			client.SetReadDeadline(time.Now().Add(MaxReadWriteTime * time.Second))
			if n, err := io.ReadFull(client.cnn, pacData[cPacLen:pPacLen-8]); err != nil && n != int(pPacLen-8-cPacLen) {
				if def.IsTimeOut(err) || def.IsFdBlock(err) {
					cType += n
					fmt.Printf("Read pacData time out or block \r\n")
					continue
				}
				err = errors.New("exit chan ")
				goto END
			}
			cPacLen += pPacLen
		}

		fmt.Println("HandleCnn get full packet data   ...")
		client.receivePacket <- packet.Packet{
			Len:  pPacLen,
			Type: funcall.FuncCallEnum(convert.BytesToInt32(pType)),
			Data: pacData[:pPacLen-8],
		}
		cLen = 0
		cType = 0
		cPacLen = 0
	}
END:
	return err
}
func (client *Connection) getID() uint64 {
	return client.id
}
func (client *Connection) setID(id uint64) {
	client.id = id
}

//Close close
func (client *Connection) Close() {
	client.exitCh.Close()
}

//Addr addr
func (client *Connection) Addr() string {
	return client.addr
}

//SetReadDeadline 阻塞时间
func (client *Connection) SetReadDeadline(t time.Time) error {
	return client.cnn.SetReadDeadline(t)
}

//SetWriteDeadline 阻塞时间
func (client *Connection) SetWriteDeadline(t time.Time) error {
	return client.cnn.SetWriteDeadline(t)
}

//CheckChan 检查chan是否可用
func (client *Connection) CheckChan() bool {
	select {
	case <-client.exitCh.GetChan():
		return true
	default:
	}
	return false
}

const (
	//MaxSend .
	MaxSend int32 = 100
	//MaxRead .
	MaxRead int32 = 100
	//MaxReadWriteTime 超时时间
	MaxReadWriteTime time.Duration = 10
)

//Server ,
type Server interface {
	Start()
	Stop()
	Init(ipaddr string) (err error)
}

//TCPServer net
type TCPServer struct {
	listen       *net.TCPListener
	exitCh       *BoolChan
	waitGroup    *sync.WaitGroup
	readTimeOut  time.Duration
	writeTimeOut time.Duration
	oc           sync.Once
	cnnMap       map[uint64]*Connection
	idGenerator  *util.IDGenerator
}

//NewLBNet new lbnet
func NewLBNet() *TCPServer {
	return &TCPServer{
		exitCh:       NewBoolChan(0),
		listen:       &net.TCPListener{},
		waitGroup:    &sync.WaitGroup{},
		readTimeOut:  time.Duration(30),
		writeTimeOut: time.Duration(30),
		cnnMap:       make(map[uint64]*Connection, 10000),
		idGenerator:  util.NewIDGenerator(),
	}
}

// Start deal cnn
func (lbnet *TCPServer) Start() {
	LogInfo("lbnet start ...")
	defer func() {
		if err := recover(); err != nil {
			LogFatal("lbnet  have panic ", err)
		}
		LogInfo("lbnet end ...")
		lbnet.Stop()
	}()
	for {
		if !lbnet.CheckChan() {
			goto END
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
END:
}

// Stop deal cnn
func (lbnet *TCPServer) Stop() {
	lbnet.exitCh.Close()
}

//CheckChan 检查chan是否可用
func (lbnet *TCPServer) CheckChan() bool {
	select {
	case <-lbnet.exitCh.GetChan():
		return false
	default:
	}
	return true
}

// Init deal cnn
func (lbnet *TCPServer) Init(ipaddr string) (err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ipaddr)
	if err != nil {
		LogFatal("lbserver ResolveTCPAddr  addr %s error %s   ", ipaddr, err)
		return
	}
	lbnet.listen, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		LogFatal("lbserver Listen error %s", err)
		return
	}
	return nil
}

func (lbnet *TCPServer) addCnn(client *Connection) {
	lbnet.cnnMap[client.getID()] = client
}

//HandleCnn handle cnn
func (lbnet *TCPServer) HandleCnn(tpcCnn *net.TCPConn) {
	lbnet.waitGroup.Add(1)
	client := &Connection{
		receivePacket: make(chan packet.Packet, MaxRead),
		sendPacket:    make(chan packet.Packet, MaxSend),
		exitCh:        NewBoolChan(0),
		addr:          tpcCnn.RemoteAddr().String(),
		cnn:           tpcCnn,
		id:            0,
	}
	defer lbnet.disConnect(client)
	go lbnet.HandleSend(client)
	go lbnet.HandleReceive(client)
	go lbnet.HandlePacket(client)
	select {
	case <-client.exitCh.GetChan():
	}
	lbnet.waitGroup.Done()
}

//disConnect .
func (lbnet *TCPServer) disConnect(client *Connection) {
	if client == nil {
		return
	}
	if e := recover(); e != nil {
		fmt.Println("cnn panic  from addr ", client.Addr())
	}
	if client.getID() <= 0 {
		return
	}
	if _, ok := lbnet.cnnMap[client.getID()]; ok {
		delete(lbnet.cnnMap, client.getID())
	} else {
		LogFatal("cannot find index %d", client.getID())
	}
	client.Close()
	client = nil
}

//HandleReceive handle receive
func (lbnet *TCPServer) HandleReceive(receiver HandleReceiver) {
	lbnet.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle deal packet panic error %v \r\n", err)
		}
		lbnet.waitGroup.Done()
	}()
	receiver.HandleReceive(lbnet)
}

//HandlePacket handle packet
func (lbnet *TCPServer) HandlePacket(packger HandlePackger) {
	lbnet.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle deal packet panic error %v \r\n", err)
		}
		lbnet.waitGroup.Done()
	}()
	packger.HandlePacket(lbnet)
}

//HandleSend handle send packet...
func (lbnet *TCPServer) HandleSend(sender HandleSender) {
	lbnet.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle send packet panic error %v \r\n", err)
		}
		lbnet.waitGroup.Done()
	}()
	sender.HandleSend(lbnet)
}

//ProxyServer .
type ProxyServer struct {
	TCPServer
}

//加入连接
func (lbnet *ProxyServer) addCnn(client *Connection) {
	client.setID(lbnet.idGenerator.GenerateID())
	lbnet.cnnMap[client.getID()] = client
}
