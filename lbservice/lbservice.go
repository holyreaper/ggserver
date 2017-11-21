package lbservice

import (
	"fmt"
	"net"
	"sync"
)

//LBService ...
type LBService struct {
	waitGroup *sync.WaitGroup
	exitCh    chan struct{}
	listen    *net.TCPListener
}

//NewLBService new
func NewLBService() *LBService {
	return &LBService{
		exitCh:    make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

//init 初始化
func (lb *LBService) init() bool {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.1.177:8091")
	lb.listen, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("LbServer Listen error ", err)
		return false
	}
	return true
}

//Start start service
func (lb *LBService) Start() {
	if !lb.init() {
		return
	}
	defer func() {
		lb.listen.Close()
	}()

	for {

		select {
		case <-lb.exitCh:
			break
		}
		cnn, err := lb.listen.Accept()
		if err != nil {
			continue
			//if NetErr, ok := err.(*net.OpError); ok && NetErr.Timeout() {
		}

	}
	//cnn.Close()

}

//Stop stop server
func (lb *LBService) Stop() {
	close(lb.exitCh)
	lb.waitGroup.Wait()
	fmt.Println("lbserver stop ...")
}
