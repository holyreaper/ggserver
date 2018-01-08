package rpcclient

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"

	"github.com/holyreaper/ggserver/consul"
	. "github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/rpcservice/rpclog"
	"google.golang.org/grpc"
)

//RPCClient .
type RPCClient struct {
	name   string
	addr   string
	port   int32
	tp     ServerType
	id     SID
	data   chan bool
	client *grpc.ClientConn
}

//StartClient start rpc client
func (cl *RPCClient) StartClient(wg *sync.WaitGroup) {
	cl.cnn()
}
func (cl *RPCClient) cnn() {
	client, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure())
	if err != nil {
		fmt.Println("client exit error msg ", err)
		return
	}
	cl.client = client
}

//GetCTRPC .
func (cl *RPCClient) GetCTRPC() ctrpcpt.CTRPCServerClient {
	cnn := ctrpcpt.NewCTRPCServerClient(cl.client)
	cl.client.GetState()
	return cnn
}

func (cl *RPCClient) execute(data interface{}) (ret interface{}) {

	return
}
func (cl *RPCClient) exit() {
	if cl.client != nil {
		cl.client.Close()
	}
}

//RPCClientMng .
type RPCClientMng struct {
	client    map[SID]*RPCClient
	rwMutex   sync.RWMutex
	waitGroup *sync.WaitGroup
}

var grpcmng *RPCClientMng

var gserverID SID
var gexitCh chan bool

func init() {
	grpcmng = &RPCClientMng{
		client:    make(map[SID]*RPCClient, 10),
		rwMutex:   sync.RWMutex{},
		waitGroup: &sync.WaitGroup{},
	}
	//gexitCh = make(chan bool)
}

//Start .
func Start(id SID, exitCh chan bool) {
	gserverID = id
	gexitCh = exitCh
	go checkSvr()
}

//checkSvr check new server
func checkSvr() {
	tick := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-tick.C:
			grpcmng.rfreshSvr()
			//check
		case <-gexitCh:
			grpcmng.exit()
			//exit
		}
	}

}

func (mng *RPCClientMng) rfreshSvr() {
	mng.rwMutex.Lock()
	defer func() {
		mng.rwMutex.Unlock()
	}()

	svr, err := consul.GetServices()
	if err != nil {
		rpclog.GetLogger().LogFatal("rpcclient getservice fail %s ", err)
		return
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
		if gserverID != cl.id {
			if _, ok := grpcmng.client[cl.id]; ok {
				//have got yet
			} else {
				//new service
				cl.StartClient(mng.waitGroup)
				grpcmng.client[cl.id] = &cl
			}
		}

	}
}

func (mng *RPCClientMng) exit() {
	mng.waitGroup.Wait()
	rpclog.GetLogger().LogInfo("rpcservice all exit !!!")
	return
}
