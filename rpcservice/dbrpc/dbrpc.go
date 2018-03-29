package dbrpc

import (
	"errors"

	"github.com/holyreaper/ggserver/database/mysqldao"
	. "github.com/holyreaper/ggserver/glog"
	"github.com/holyreaper/ggserver/rpcservice/pb/dbrpc"
	"github.com/holyreaper/ggserver/rpcservice/pb/rpcutil"
	"golang.org/x/net/context"
)

func init() {
}

//DBRPC ..
type DBRPC struct {
	//Srv interface{}
}

//KeepAlive  .
func (chat *DBRPC) KeepAlive(cont context.Context, req *dbrpcpt.KeepAliveRequest) (*dbrpcpt.KeepAliveReply, error) {
	return &dbrpcpt.KeepAliveReply{Result: 0}, nil
}

//Login  get user info
func (chat *DBRPC) Login(cont context.Context, req *dbrpcpt.LoginRequest) (*dbrpcpt.LoginReply, error) {
	if req.GetUid() <= 0 {
		return nil, errors.New("user uid is nil")
	}
	ret, err := mysqldao.UserSelectSingle(mysqldao.UserField, req.GetUid())
	if err != nil {
		LogFatal("user login fail %d %d %v", req.GetUid(), req.GetServerId(), err)
		return nil, err
	}
	rsp := &dbrpcpt.LoginReply{}
	rsp.User = rpcutil.FillUserField(ret)
	return rsp, nil
}

//GetOffLineMsg  .
func (chat *DBRPC) GetOffLineMsg(cont context.Context, req *dbrpcpt.GetOffLineMsgRequest) (*dbrpcpt.GetOffLineMsgReply, error) {
	if req.GetUid() <= 0 {
		return nil, errors.New("user uid is nil")
	}
	where := []string{"uid"}
	var value []interface{}
	value = append(value, req.GetUid())
	ret, err := mysqldao.OffLineMsgSelectMulity(mysqldao.OffLineMsgField, where, value)
	if err != nil {
		LogFatal("user GetOffLineMsg fail %d %v", req.GetUid(), err)
		return nil, err
	}
	rsp := &dbrpcpt.GetOffLineMsgReply{}
	for _, v := range ret {
		field := rpcutil.FillOffLineMsgField(v)
		rsp.Info[field.Id] = field
	}
	return rsp, nil
}

//SaveOffLineMsg  ...
func (chat *DBRPC) SaveOffLineMsg(cont context.Context, req *dbrpcpt.SaveOffLineMsgRequest) (*dbrpcpt.SaveOffLineMsgReply, error) {
	for _, v := range req.GetInfo() {
		field := rpcutil.ConvertOffLineMsg(v)
		err := mysqldao.OffLineMsgInsert(field)
		if err != nil {
			LogFatal("user SaveOffLineMsg fail %v ", err)
			return nil, err
		}
	}
	return &dbrpcpt.SaveOffLineMsgReply{Result: 0}, nil
}

//DelOffLineMsg  ...
func (chat *DBRPC) DelOffLineMsg(cont context.Context, req *dbrpcpt.DelOffLineMsgRequest) (*dbrpcpt.DelOffLineMsgReply, error) {
	if req.GetUid() <= 0 {
		return nil, errors.New("user uid is nil")
	}
	value := []interface{}{req.GetUid()}
	where := []string{"uid"}
	if req.GetId() > 0 {
		where = append(where, "id")
		value = append(value, req.GetId())
	}
	mysqldao.OffLineMsgDelete(where, value)

	return &dbrpcpt.DelOffLineMsgReply{}, nil
}
