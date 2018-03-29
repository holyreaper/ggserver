package ctrpc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/holyreaper/ggserver/rpcservice/def"
	"github.com/holyreaper/ggserver/rpcservice/pb/ctrpc"
	"golang.org/x/net/context"

	. "github.com/holyreaper/ggserver/glog"
)

//InitChanSize .
const InitChanSize = 1000

//InitServerSize .
const InitServerSize = 10

func init() {
}

//UserData .
type UserData struct {
	uid      int64
	uname    string
	figure   int32
	level    int32
	vip      int32
	serverid int32
}

//InitUserSize .
const InitUserSize = 10000

//CTRPC .
type CTRPC struct {
	//Srv interface{}
	PushMsgMap map[int32]chan ctrpcpt.PushMessageReply
	PushLock   sync.RWMutex

	UserMap  map[int64]UserData
	UserLock sync.RWMutex
}

//Init .
func (rpc *CTRPC) Init() {
	rpc.PushMsgMap = make(map[int32]chan ctrpcpt.PushMessageReply, InitServerSize)
	rpc.UserMap = make(map[int64]UserData, InitUserSize)
}

//KeepAlive   ...
func (rpc *CTRPC) KeepAlive(cont context.Context, req *ctrpcpt.KeepAliveRequest) (*ctrpcpt.KeepAliveReply, error) {
	return &ctrpcpt.KeepAliveReply{Result: 0}, nil
}

//Login   ...
func (rpc *CTRPC) Login(cont context.Context, req *ctrpcpt.LoginRequest) (*ctrpcpt.LoginReply, error) {
	old := rpc.AddUser(*req)
	if old.serverid != 0 {
		rpc.AddPushMessage(ctrpcpt.PushMessageReply{Type: int32(rpcdef.RPCPTypeUserLogout)}, old.serverid)
	}
	return &ctrpcpt.LoginReply{Result: 0}, nil
}

//Logout   ...
func (rpc *CTRPC) Logout(cont context.Context, req *ctrpcpt.LogoutRequest) (*ctrpcpt.LogoutReply, error) {
	rpc.DelUser(req.GetUid())
	return &ctrpcpt.LogoutReply{Result: 0}, nil
}

//Chat   ...
func (rpc *CTRPC) Chat(cont context.Context, req *ctrpcpt.ChatRequest) (*ctrpcpt.ChatReply, error) {
	if req.GetChatMsg().GetTuid() > 0 {
		userData, err := rpc.GetUser(req.GetChatMsg().GetTuid())
		if err != nil {
			return &ctrpcpt.ChatReply{Result: 1}, err
		}
		err = rpc.AddPushMessage(ctrpcpt.PushMessageReply{Type: int32(rpcdef.RPCPTypeUserChat), ChatMsg: req.GetChatMsg()}, userData.serverid)
		if err != nil {
			return &ctrpcpt.ChatReply{Result: 2}, err
		}
	}

	return &ctrpcpt.ChatReply{Result: 0}, nil
}

//PushStream push stream
func (rpc *CTRPC) PushStream(req *ctrpcpt.PushMessageRequest, stream ctrpcpt.CTRPC_PushStreamServer) error {
	//var i int32
	defer func(serverid int32) {
		if err := recover(); err != nil {
			LogFatal("ctrpc svr push stream have error %v ", err)
		}
		LogInfo("ctrpc svr push stream stop ")
		rpc.DelPushMessage(serverid)
	}(req.GetServerId())

	LogInfo("ctrpc PushStream recv  data  %v", *req)
	rpc.InitChanServer(req.GetServerId())
	timer := time.NewTicker(2 * time.Second)
	for {
		ch, err := rpc.GetPushMessage(req.GetServerId())
		if err != nil {
			LogFatal(" server %d have no chan ", req.GetServerId())
			//time.Sleep(2 * time.Second)
			return err
		}
		select {
		case rsp, close := <-ch:
			if close {
				return nil
			}
			err = stream.Send(&rsp)
			if err != nil {
				LogFatal(" rpcservice  push stream got err %v", err)
				return err
			}
		case <-timer.C:
		}

	}
	//return nil
}

//InitChanServer .
func (rpc *CTRPC) InitChanServer(serverid int32) {
	defer func() {
		rpc.PushLock.Unlock()
	}()
	rpc.PushLock.Lock()
	if _, ok := rpc.PushMsgMap[serverid]; !ok {
		rpc.PushMsgMap[serverid] = make(chan ctrpcpt.PushMessageReply, InitChanSize)
	}
}

//AddPushMessage .
func (rpc *CTRPC) AddPushMessage(msg ctrpcpt.PushMessageReply, serverid int32) error {
	defer func() {
		rpc.PushLock.RLocker()
	}()
	rpc.PushLock.RUnlock()
	if _, ok := rpc.PushMsgMap[serverid]; ok {
		rpc.PushMsgMap[serverid] <- msg
	} else {
		errstr := fmt.Sprintf("no register server %d", serverid)
		return errors.New(errstr)
	}
	return nil
}

//GetPushMessage .
func (rpc *CTRPC) GetPushMessage(serverid int32) (msg chan ctrpcpt.PushMessageReply, err error) {
	defer func() {
		rpc.PushLock.RUnlock()
	}()
	err = errors.New("no register server")
	rpc.PushLock.RLock()
	if v, ok := rpc.PushMsgMap[serverid]; ok {

		return v, nil
	}
	return msg, err
}

//DelPushMessage .
func (rpc *CTRPC) DelPushMessage(serverid int32) {
	defer func() {
		rpc.PushLock.Unlock()
	}()
	rpc.PushLock.Lock()
	if v, ok := rpc.PushMsgMap[serverid]; ok {
		close(v)
	}
	delete(rpc.PushMsgMap, serverid)
}

//AddUser .
func (rpc *CTRPC) AddUser(user ctrpcpt.LoginRequest) (old UserData) {
	defer func() {
		rpc.UserLock.Unlock()
	}()
	userData := UserData{
		uid:      user.GetUid(),
		uname:    user.GetUname(),
		figure:   user.GetFigure(),
		level:    user.GetLevel(),
		vip:      user.GetVip(),
		serverid: user.GetServerId()}
	rpc.UserLock.Lock()
	if v, ok := rpc.UserMap[user.GetUid()]; ok {
		old = v
		rpc.UserMap[user.GetUid()] = userData
	}
	return
}

//DelUser .
func (rpc *CTRPC) DelUser(uid int64) {
	defer func() {
		rpc.UserLock.Unlock()
	}()
	rpc.UserLock.Lock()
	delete(rpc.UserMap, uid)
	return
}

//GetUser .
func (rpc *CTRPC) GetUser(uid int64) (UserData, error) {
	defer func() {
		rpc.UserLock.RUnlock()
	}()
	rpc.UserLock.RLock()
	if data, ok := rpc.UserMap[uid]; ok {
		return data, nil
	}
	return UserData{}, errors.New("have no user ")
}
