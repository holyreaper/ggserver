package lbnet

import (
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

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
	exitCh       chan struct{}
	waitGroup    *sync.WaitGroup
	readTimeOut  time.Duration
	writeTimeOut time.Duration
}

//MAX_SEND .
const MAX_SEND int32 = 10
const MAX_RECV int32 = 10

//NewLBNet new lbnet
func NewLBNet() *LBNet {
	return &LBNet{
		exitCh:       make(chan struct{}),
		listen:       &net.TCPListener{},
		waitGroup:    &sync.WaitGroup{},
		readTimeOut:  time.Duration(30),
		writeTimeOut: time.Duration(30),
	}
}

// Start deal cnn
func (lbnet *LBNet) Start() {
	fmt.Println("lbnet start ...")
	defer func() {

		if err := recover(); err != nil {
			fmt.Println("lbnet  have panic ", err)
		}
		fmt.Println("lbnet end ...")
	}()
	for {

		select {
		case <-lbnet.exitCh:
			break
		default:
		}
		cnn, err := lbnet.listen.AcceptTCP()
		if err != nil {
			fmt.Println("have an invalide connect ", err)
			continue
			//if NetErr, ok := err.(*net.OpError); ok && NetErr.Timeout() {
		}
		fmt.Println("have a Cnn accept start to serve her ")
		go lbnet.HandleCnn(cnn)
	}

}

// Stop deal cnn
func (lbnet *LBNet) Stop() {
	close(lbnet.exitCh)
	lbnet.waitGroup.Wait()
}

// Init deal cnn
func (lbnet *LBNet) Init(netproto string /*proto4 */, nettype string /*tcp*/, ipaddr string) {
	tcpAddr, err := net.ResolveTCPAddr(nettype, ipaddr)
	if err != nil {
		fmt.Println("err")
		os.Exit(-1)
	}
	lbnet.listen, err = net.ListenTCP(nettype, tcpAddr)
	if err != nil {
		fmt.Println("LbServer Listen error ", err)
		os.Exit(-2)
	}
}

//HandleCnn handle cnn
func (lbnet *LBNet) HandleCnn(cnn *net.TCPConn) {
	lbnet.waitGroup.Add(1)
	recvPacket := make(chan packet.Packet, MAX_RECV)
	sendPacket := make(chan packet.Packet, MAX_SEND)
	addr := cnn.RemoteAddr().String()

	defer func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("cnn panic  from addr ", addr)
			}
		}()
		lbnet.waitGroup.Done()
		cnn.Close()
	}()

	var (
		pLen    = make([]byte, 4)
		pType   = make([]byte, 4)
		pPacLen int32
		pacData = make([]byte, packet.MAXPACKETLEN)

		cLen    = 0   //curr read pLen
		cType   = 0   //curr read pType
		cPacLen int32 //curr read pPackLen

	)
	go lbnet.HandleSend(cnn, sendPacket)
	go lbnet.HandlePacket(cnn, recvPacket, sendPacket)
	fmt.Println("start serve the new client  ", cnn.RemoteAddr().String())
	for {
		select {
		case <-lbnet.exitCh:
			fmt.Printf("Stop HandleCnn\r\n")
			return
		default:
		}
		cnn.SetReadDeadline(time.Now().Add(10e9))
		if cLen < 4 {
			if n, err := io.ReadFull(cnn, pLen[cLen:]); err != nil && n != 4 {
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
			if n, err := io.ReadFull(cnn, pType); err != nil && n != 4 {
				if def.IsTimeOut(err) {
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
			if n, err := io.ReadFull(cnn, pacData[cPacLen:pPacLen-8]); err != nil && n != int(pPacLen-8-cPacLen) {
				if def.IsTimeOut(err) {
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
		recvPacket <- packet.Packet{
			Len:  pPacLen,
			Type: convert.BytesToInt32(pType),
			Data: pacData[:pPacLen-8],
		}
		cLen = 0
		cType = 0
		cPacLen = 0
	}
}

//HandlePacket handle packet
func (lbnet *LBNet) HandlePacket(cnn *net.TCPConn, rpack chan packet.Packet, spack chan<- packet.Packet) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle deal packet panic error %v \r\n", err)
		}
	}()
	fmt.Println("HandlePacket ing ...")
	for {
		select {
		case <-lbnet.exitCh:
			fmt.Printf("Stop HandlePacket \r\n")
			return
		case p := <-rpack:
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
					ret, err = funcall.Call(p.Type, rpack, spack, lmsg)
				} else {
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
					spack <- tmppack
				}
			}
		}

	}

}

//HandleSend handle send packet...
func (lbnet *LBNet) HandleSend(net *net.TCPConn, spack <-chan packet.Packet) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle send packet panic error %v \r\n", err)
		}
	}()
	for {
		select {
		case <-lbnet.exitCh:
			fmt.Printf("Stop Handle Send Packet \r\n")
			return
		case p := <-spack:
			fmt.Printf("send packeting ... %d len %d ", p.Type, p.Len)
			buff := p.FormatBuf()
			currSend := 0
			for {

				ln, err := net.Write(buff[currSend:])
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
