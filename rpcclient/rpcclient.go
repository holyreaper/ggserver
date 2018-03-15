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
	. "github.com/holyreaper/ggserver/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/holyreaper/ggserver/lbmodule/manager/callbackmanager"
)

// Stat 服务器状态
type Stat int32

const (
	//Ready 正常
	Ready = iota
	//NoReady 断开连接
	NoReady
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
	stat   Stat
}

//StartClient start rpc client
func (cl *RPCClient) StartClient() {
	cl.cnn()
}
func (cl *RPCClient) cnn() {
	keepaliveParam := grpc.WithKeepaliveParams(keepalive.ClientParameters{PermitWithoutStream: true})
	addr := fmt.Sprintf("%s:%d", cl.addr, cl.port)
	client, err := grpc.Dial(addr, grpc.WithInsecure(), keepaliveParam)
	if err != nil {
		fmt.Println("client exit error msg ", err)
		return
	}
	if cl.tp == def.ServerTypeCenter {
		LogInfo("rpcclient cnn ct server %s:%d ", cl.addr, cl.port)
		go cl.rpcStreamClient()
	}
	cl.client = client
	fmt.Printf(" connect to server %s %v ", cl.addr, cl.client)

}

//GetCTRPC .
func (cl *RPCClient) GetCTRPC() (cnn ctrpcpt.CTRPCClient, err error) {
	if cl.tp == def.ServerTypeNormal {
		err = errors.New("have no ctrpc function")
	} else if cl.tp == def.ServerTypeCenter {
		if cl.client == nil {
			err = errors.New("have no ctrpc function nill client ")
		} else {
			cnn = ctrpcpt.NewCTRPCClient(cl.client)
		}

	} else if cl.tp == def.ServerTypeDB {
		err = errors.New("have no ctrpc function")
	}
	fmt.Printf(" GetCTRPC to server %s %v %v", cl.addr, cl.client, cnn)
	return cnn, err
}

//GetDBRPC .
func (cl *RPCClient) GetDBRPC() (cnn dbrpcpt.DBRPCClient, err error) {
	if cl.tp == def.ServerTypeNormal {
		err = errors.New("have no dbrpc function")
	} else if cl.tp == def.ServerTypeCenter {
		err = errors.New("have no dbrpc function")
	} else if cl.tp == def.ServerTypeDB {
		if cl.client == nil {
			err = errors.New("have no dbrpc function nill client")
		} else {
			cnn = dbrpcpt.NewDBRPCClient(cl.client)
		}

	}

	return cnn, err
}

func (cl *RPCClient) rpcStreamClient() {
	defer func() {
		if err := recover(); err != nil {
			LogFatal("rpcStreamClient fail exit err %v ", err)
		}
	}()
	req := ctrpcpt.PushMessageRequest{
		ServerId: int32(cl.id),
	}
	for {
		select {
		case <-grpcmng.exitCh:
			return
		default:

		}
		LogInfo("stat %d ", cl.client.GetState())
		rpc, err := cl.GetCTRPC()
		if err != nil {
			LogFatal("rpcStreamClient GetCTRPC cnn fail server id  ", cl.id)
			return
		}

		//ctx, _ := context.WithTimeout(context.Background(), 30000*time.Microsecond)
		LogInfo("start pushStream %v ", cl.addr)
		client, err := rpc.PushStream(context.Background(), &req)
		if err != nil {
			LogFatal("rpcStreamClient pushStream fail err %v ", err)
			time.Sleep(2 * time.Second)
			continue
		}
		var rsp *ctrpcpt.PushMessageReply
		for {
			rsp, err = client.Recv()
			if err != nil {
				LogFatal("rpcStreamClient recv data fail err %v ", err)
				time.Sleep(2 * time.Second)
				break
			}
			callbackmanager.Put(*rsp)

			LogInfo("rpcStreamClient from %d get message type %d  %v", cl.id, rsp.GetType(), cl.client)
		}

	}

}
func (cl *RPCClient) sayHello() {
	if cl.tp == def.ServerTypeCenter {
		rpc, err := cl.GetCTRPC()
		if err != nil {
			LogFatal("sayHello GetCTRPC cnn fail server id  %d error %v", cl.id, err)
			cl.stat = NoReady
			return
		}
		req := ctrpcpt.KeepAliveRequest{
			Time: int64(time.Now().Unix()),
		}
		_, err = rpc.KeepAlive(context.Background(), &req)
		if err != nil {
			LogFatal(" sayHello ctservice %v discnnect ", cl.addr)
			cl.stat = NoReady
		} else {
			cl.stat = Ready
		}

	} else if cl.tp == def.ServerTypeDB {
		rpc, err := cl.GetDBRPC()
		if err != nil {
			LogFatal("sayHello GetDBRPC cnn fail server id  ", cl.id)
			cl.stat = NoReady
			return
		}
		req := dbrpcpt.KeepAliveRequest{
			Time: int64(time.Microsecond),
		}
		_, err = rpc.KeepAlive(context.Background(), &req)
		if err != nil {
			LogFatal(" sayHello dbservice %v discnnect ", cl.addr)
			cl.stat = NoReady
		} else {
			cl.stat = Ready
		}
	}

}

func (cl *RPCClient) exit() {
	if cl.client != nil {
		cl.client.Close()
	}
}

//RPCClientMng d
type RPCClientMng struct {
	client  map[SID]*RPCClient
	rwMutex sync.RWMutex
	exitCh  chan bool
}

var grpcmng *RPCClientMng

var gserverID SID

func init() {
	grpcmng = &RPCClientMng{
		client:  make(map[SID]*RPCClient, 10),
		rwMutex: sync.RWMutex{},
		exitCh:  make(chan bool),
	}
	//gexitCh = make(chan bool)
}

//Start .
func Start(id SID) {
	gserverID = id
	go grpcmng.checkSvr()
}

//Stop .
func Stop() {
	if grpcmng != nil {
		close(grpcmng.exitCh)
	}
}

//GetCTRPC .
func GetCTRPC() (cnn ctrpcpt.CTRPCClient, err error) {
	cl := grpcmng.GetRPCClientFromType(ServerTypeCenter)
	if cl == nil {
		return nil, err
	}
	cnn, err = cl.GetCTRPC()
	return
}

//GetDBRPC .
func GetDBRPC() (cnn dbrpcpt.DBRPCClient, err error) {
	cl := grpcmng.GetRPCClientFromType(ServerTypeDB)
	if cl == nil {
		return nil, err
	}
	cnn, err = cl.GetDBRPC()
	return
}

//checkSvr check new server
func (mng *RPCClientMng) checkSvr() {
	mng.rfreshSvr()
	tick := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-tick.C:
			mng.rfreshSvr()
			//check
		case <-grpcmng.exitCh:
			mng.exit()
			return
			//exit
		}
	}
}

//GetRPCClientFromID  .
func (mng *RPCClientMng) GetRPCClientFromID(id SID) (cl *RPCClient) {
	if cl, ok := mng.client[id]; ok {
		return cl
	}
	return nil
}

//GetRPCClientFromType .
func (mng *RPCClientMng) GetRPCClientFromType(tp ServerType) (cl *RPCClient) {
	for _, v := range mng.client {
		if v.tp == tp {
			return v
		}
	}
	return nil
}

func (mng *RPCClientMng) rfreshSvr() {
	mng.rwMutex.Lock()
	defer func() {
		mng.rwMutex.Unlock()
	}()

	svr, err := consul.GetServices()
	if err != nil {
		LogFatal("rpcclient getservice fail %s ", err)
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
			if value, ok := grpcmng.client[cl.id]; ok {
				//have got yet
				go value.sayHello()
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
	LogInfo("rpcclient  all exit !!!")
	return
}
