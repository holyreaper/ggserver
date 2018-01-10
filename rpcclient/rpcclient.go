package rpcclient

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/util"

	"github.com/holyreaper/ggserver/rpcservice/pb/dbrpc"

	"google.golang.org/grpc/keepalive"

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
func (cl *RPCClient) StartClient() {
	cl.cnn()
}
func (cl *RPCClient) cnn() {
	keepaliveParam := grpc.WithKeepaliveParams(keepalive.ClientParameters{PermitWithoutStream: true})
	client, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure(), keepaliveParam)
	if err != nil {
		fmt.Println("client exit error msg ", err)
		return
	}
	cl.client = client
}

//GetCTRPC .
func (cl *RPCClient) GetCTRPC() (cnn ctrpcpt.CTRPCClient, err error) {
	if cl.tp == def.ServerTypeNormal {
		cnn = ctrpcpt.NewCTRPCClient(cl.client)
	} else if cl.tp == def.ServerTypeCenter {
		err = errors.New("have no ctrpc function")
	} else if cl.tp == def.ServerTypeDB {
		err = errors.New("have no ctrpc function")
	}

	return cnn, err
}

//GetDBRPC .
func (cl *RPCClient) GetDBRPC() (cnn dbrpcpt.DBRPCClient, err error) {
	if cl.tp == def.ServerTypeNormal {
		err = errors.New("have no dbrpc function")
	} else if cl.tp == def.ServerTypeCenter {
		err = errors.New("have no dbrpc function")
	} else if cl.tp == def.ServerTypeDB {
		cnn = dbrpcpt.NewDBRPCClient(cl.client)
	}
	return cnn, err
}

func (cl *RPCClient) exit() {
	if cl.client != nil {
		cl.client.Close()
	}
}

//RPCClientMng .
type RPCClientMng struct {
	client  map[SID]*RPCClient
	rwMutex sync.RWMutex
}

var grpcmng *RPCClientMng

var gserverID SID
var gexitCh chan bool

func init() {
	grpcmng = &RPCClientMng{
		client:  make(map[SID]*RPCClient, 10),
		rwMutex: sync.RWMutex{},
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
			tp:   util.GetServerType(SID(id)),
		}
		if gserverID != cl.id {
			if _, ok := grpcmng.client[cl.id]; ok {
				//have got yet
			} else {
				//new service
				if cl.tp != util.GetServerType(gserverID) {
					cl.StartClient()
					grpcmng.client[cl.id] = &cl
				}

			}
		}

	}
}

func (mng *RPCClientMng) exit() {
	rpclog.GetLogger().LogInfo("rpcservice all exit !!!")
	return
}
