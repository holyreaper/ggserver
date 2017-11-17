package main

import (
	"time"

	"github.com/holyreaper/ggserver/client"

	"github.com/holyreaper/ggserver/def"
	//"context"
	"fmt"

	"runtime"

	"github.com/holyreaper/ggserver/conf"
	"github.com/holyreaper/ggserver/service"

	"flag"
)

const (
	port = ":8090"
	msg  = "helo lhm"
)

var mode = flag.Int("mode", 0, "server mode ")
var exit = make(chan int)

func init() {
	flag.Parse()
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("server mode : %d\n", *mode)
	conf := conf.GetConf()
	if data, ok := conf["info"]; ok {
		sliceArray := data.([]interface{})
		for k, v := range sliceArray {
			fmt.Printf("key %d, value %s", k, v)
		}
	}
	defer func() {
		println("finish ...")
	}()
	fmt.Println("start .service ")
	ggserver := service.NewGGService(def.ServerTypeNormal)
	ggserver.Start()
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
