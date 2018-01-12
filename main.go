package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/holyreaper/ggserver/client"
	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbservice"

	"github.com/holyreaper/ggserver/util"

	"fmt"

	"runtime"

	"github.com/holyreaper/ggserver/rpcservice"

	"flag"

	"net/http"
	_ "net/http/pprof"

	. "github.com/holyreaper/ggserver/glog"
)

const (
	port = ":8090"
	msg  = "helo lhm"
)

var gserverID *int

var gexitCh = make(chan bool)

var gserverType def.ServerType

func init() {

}
func main() {

	gserverID = flag.Int("serverid", 0, "serverid ")

	flag.Parse()

	gserverType = util.GetServerType(def.SID(*gserverID))

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("server mode : %d\n", *gserverID)

	InitLog(gserverType, def.SID(*gserverID))

	LogInfo("server start server mode  %d  ", *gserverID)
	//	conf := conf.GetConf()
	//for _, _ = range conf {

	//}
	//debug list
	//go Pprof()
	defer func() {
		println("finish ...")
	}()
	fmt.Println("start .service ")
	if gserverType != def.ServerTypeNormal {
		rpcservice.Init(gserverType)
		go rpcservice.Start(gexitCh)
	} else {
		err := lbservice.Init(def.SID(*gserverID))
		if err != nil {
			LogFatal("start lbservice fail err %s", err)
		}
		go lbservice.Start(gexitCh)
	}
	go client.Start()
	//go Tick()
	go Signal()
	<-gexitCh
}

//Tick ...
func Tick() {
	tick := time.NewTicker(10 * time.Second)
	for _ = range tick.C {
		fmt.Println("tick ", tick.C)
	}
	fmt.Println("tick exit")
}

//Signal signal
func Signal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	msg := <-sig
	LogInfo("server get signal %s", msg)

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
