package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/holyreaper/ggserver/common"

	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/lbservice"

	"fmt"

	"runtime"

	"flag"

	"net/http"
	_ "net/http/pprof"

	. "github.com/holyreaper/ggserver/glog"
)

var gexitCh = make(chan bool)

func init() {

}
func main() {

	serverid := flag.Int("serverid", 0, "serverid ")
	flag.Parse()
	common.SetServerID(def.SID(*serverid))

	//serverType = util.GetServerType(def.SID(*serverid))

	if common.GetServerType() < def.ServerTypeNormal || common.GetServerType() > def.ServerTypeCenter {
		fmt.Printf("invalid gserverid %d ", serverid)
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("server mode : %d\n", common.GetServerID())

	InitLog(common.GetServerType(), common.GetServerID())

	LogInfo("server start server mode  %d  ", common.GetServerID())
	//	conf := conf.GetConf()
	//for _, _ = range conf {

	//}
	//debug list
	//go Pprof()
	defer func() {
		println("finish ...")

	}()
	fmt.Println("start .service ")
	if common.GetServerType() == def.ServerTypeDB || common.GetServerType() == def.ServerTypeCenter {
	} else {
		err := lbservice.Init(common.GetServerID())
		if err != nil {
			LogFatal("start lbservice fail err %s", err)
			return
		}
		lbservice.Start()
	}
	//go client.UserClient(gexitCh)
	//go Tick()
	//Stop()
	go Signal()
	<-gexitCh
}

//Stop ...
func Stop() {
	if common.GetServerType() == def.ServerTypeDB || common.GetServerType() == def.ServerTypeCenter {

	} else {
		lbservice.Stop()
	}
	close(gexitCh)

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
	Stop()
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
