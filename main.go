package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/holyreaper/ggserver/client"
	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbservice"

	"fmt"

	"runtime"

	"github.com/holyreaper/ggserver/rpcservice"

	"flag"

	"net/http"
	_ "net/http/pprof"

	"github.com/holyreaper/ggserver/lbmodule/lblog"
)

const (
	port = ":8090"
	msg  = "helo lhm"
)

var mode = flag.Int("mode", 0, "server mode ")
var exit = make(chan int)

var serverID = flag.Int("id", 0, "serverid ")

func init() {
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("server mode : %d\n", *mode)
	lblog.GetLogger().LogInfo("server start server mode  %d  ", *mode)
	//	conf := conf.GetConf()
	//for _, _ = range conf {

	//}
	//debug list
	//go Pprof()
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

//Pprof 检查
func Pprof() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
	f, err := os.OpenFile("cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
}
