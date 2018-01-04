package main

import (
	"time"

	"github.com/holyreaper/ggserver/client"
	"github.com/holyreaper/ggserver/lbservice"

	"github.com/holyreaper/ggserver/def"

	"fmt"

	"runtime"

	"github.com/holyreaper/ggserver/conf"
	"github.com/holyreaper/ggserver/rpcservice"

	"flag"
)

const (
	port = ":8090"
	msg  = "helo lhm"
)

var mode = flag.Int("mode", 0, "server mode ")
var exit = make(chan int)
var exit2 = make(chan bool, 1)

func init() {
	flag.Parse()
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("server mode : %d\n", *mode)
	exit2 <- true
	switch {
	case <-exit2:

	}
	switch {
	case <-exit2:

	}
	conf := conf.GetConf()
	for _, _ = range conf {

	}
	defer func() {
		println("finish ...")
	}()
	fmt.Println("start .service ")
	ggserver := rpcservice.NewGGService(def.ServerTypeNormal)
	ggserver.Start()
	lbserver := lbservice.NewLBService()
	go lbserver.Start()
	go client.Start()
	go Tick()
	<-exit

}

//Tick ...
func Tick() {
	tick := time.NewTicker(10 * time.Second)
	for _ = range tick.C {
		fmt.Println("tick ", tick.C)
	}
	fmt.Println("tick exit")
}
