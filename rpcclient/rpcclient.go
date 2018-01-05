package rpcclient

import (
	"strconv"
	"time"

	"github.com/holyreaper/ggserver/consul"
	. "github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/rpcservice/rpclog"
)

//RPCClient .
type RPCClient struct {
	name string
	addr string
	port int32
	tp   ServerType
	id   SID
	new  bool
}

//RPCClientMng .
type RPCClientMng struct {
	client map[SID]RPCClient
}

var grpcmng *RPCClientMng

var serverID SID

func init() {
	grpcmng = &RPCClientMng{
		client: make(map[SID]RPCClient, 10),
	}
}

//Start .
func Start(id SID, exitCh <-chan bool) {
	go checkSvr(id, exitCh)
}
func rfreshSvr(sid SID) {
	svr, err := consul.GetServices()
	if err != nil {
		rpclog.GetRPCLogger().LogFatal("rpcclient getservice fail %s ", err)
	}
	for _, v := range svr {
		id, _ := strconv.Atoi(v.ID)
		cl := RPCClient{
			name: v.Service,
			addr: v.Address,
			port: int32(v.Port),
			id:   SID(id),
			tp:   ServerType(id / SERVERBASEVALUE),
		}
		//have no myself ...
		if sid != cl.id {
			if _, ok := grpcmng.client[cl.id]; ok {
				//have got yet
			} else {
				//new service
				cl.new = true
				grpcmng.client[cl.id] = cl
			}
		}

	}
}

//checkSvr go checkSvr
func checkSvr(id SID, exitCh <-chan bool) {
	tick := time.NewTicker(1 * time.Minute)
	for {
		select {

		case <-tick.C:
			//check
		case <-exitCh:
			//exit
		}
	}

}
