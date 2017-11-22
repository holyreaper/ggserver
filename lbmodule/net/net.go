package lbnet

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/holyreaper/ggserver/lbmodule/funcall"
	"github.com/holyreaper/ggserver/lbmodule/packet"
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

// Start deal cnn
func (lbnet *LBNet) Start() {

	for {

		select {
		case <-lbnet.exitCh:
			break
		}
		cnn, err := lbnet.listen.AcceptTCP()
		if err != nil {
			continue
			//if NetErr, ok := err.(*net.OpError); ok && NetErr.Timeout() {
		}
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
	lbnet.listen, err = net.ListenTCP(netproto, tcpAddr)
	if err != nil {
		fmt.Println("LbServer Listen error ", err)
		os.Exit(-2)
	}
}

//HandleCnn handle cnn
func (lbnet *LBNet) HandleCnn(cnn *net.TCPConn) {
	lbnet.waitGroup.Add(1)
	recvPacket := make(chan *packet.Packet, 10)

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
		pPacLen uint32
	)
	for {
		select {
		case <-lbnet.exitCh:
			fmt.Printf("Stop HandleCnn\r\n")
			return
		default:
		}
		cnn.SetReadDeadline(time.Now().Add(lbnet.readTimeOut))

		if n, err := io.ReadFull(cnn, pLen); err != nil && n != 4 {
			fmt.Printf("Read pLen failed: %v\r\n", err)
			return
		}
		if n, err := io.ReadFull(cnn, pType); err != nil && n != 4 {
			fmt.Printf("Read pType failed: %v\r\n", err)
			return
		}
		if pPacLen = uint32(convert.BytesToInt32(pLen)); pPacLen > packet.MAXPACKETLEN {
			fmt.Printf("pacLen larger than pPacLen\r\n")
			return
		}

		pacData := make([]byte, pPacLen-8)
		if n, err := io.ReadFull(cnn, pacData); err != nil && n != int(pPacLen) {
			fmt.Printf("Read pacData failed: %v\r\n", err)
			return
		}
		recvPacket <- &packet.Packet{
			Len:  pPacLen,
			Type: uint32(convert.BytesToInt32(pType)),
			Data: pacData,
		}
	}
}

//HandlePacket handle packet
func (lbnet *LBNet) HandlePacket(cnn *net.TCPConn, pack <-chan *packet.Packet) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handle packet panic error %v \r\n", err)
		}
	}()

	for {
		select {
		case <-lbnet.exitCh:
			fmt.Printf("Stop HandlePacket \r\n")
			return
		case p := <-pack:
			{
				ret, err := funcall.Call(p.Type, p.Data)
				if err != nil {
					fmt.Printf("func call type  %d error %s \r\n", p.Type, err)
					continue
				}
				//deal with return value
				for _, value := range ret {
					fmt.Println(value.Bytes())
					//todo send return value to client
				}
			}
		}

	}

}
