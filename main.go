package main

import (
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
	ggserver := service.NewGGService(def.ServerTypeNormal)
	ggserver.Start()
	defer func() {
		println("finish ...")
	}()

}
