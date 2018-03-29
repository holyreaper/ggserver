package rpcutil

//提供rpc操作的一些公用接口
import (
	"github.com/holyreaper/ggserver/database/mysqldao"
	pbmysql "github.com/holyreaper/ggserver/rpcservice/pb/mysql"
)

// FillUserField 填充user信息
func FillUserField(data map[string]interface{}) (ret *pbmysql.UserField) {
	for k, v := range data {
		if k == mysqldao.UserUID {
			ret.Uid = v.(int64)
		} else if k == mysqldao.UserUname {
			ret.Uname = v.(string)
		} else if k == mysqldao.UserCreateTime {
			ret.CreateTime = v.(int32)
		} else if k == mysqldao.UserLastLGITime {
			ret.LastLoginTime = v.(int32)
		} else if k == mysqldao.UserLastLGOTime {
			ret.LastLogoutTime = v.(int32)
		} else if k == mysqldao.UserExp {
			ret.Exp = v.(int32)
		} else if k == mysqldao.UserLevel {
			ret.Level = v.(int32)
		} else if k == mysqldao.UserFigure {
			ret.Figure = v.(int32)
		} else if k == mysqldao.UserExtData {
			ret.EtData = v.(string)
		}

	}
	return
}

//离线相关
//FillOffLineMsgField .
func FillOffLineMsgField(data map[string]interface{}) (ret *pbmysql.OffLineMsgField) {
	for k, v := range data {
		if k == mysqldao.OfflineMsgID {
			ret.Id = v.(int64)
		} else if k == mysqldao.OfflineMsgUID {
			ret.Uid = v.(int64)
		} else if k == mysqldao.OfflineMsgSendTime {
			ret.SendTime = v.(int32)
		} else if k == mysqldao.OfflineMsgFUID {
			ret.Fuid = v.(int64)
		} else if k == mysqldao.OfflineMsgFUname {
			ret.Funame = v.(string)
		} else if k == mysqldao.OfflineMsgFfigure {
			ret.Ffigure = v.(int32)
		} else if k == mysqldao.OfflineMsgETData {
			ret.EtData = v.(string)
		}

	}
	return
}

//ConvertOffLineMsg .
func ConvertOffLineMsg(data *pbmysql.OffLineMsgField) []interface{} {

	ret := make([]interface{}, 7)
	ret[0] = data.GetFuid()
	ret[1] = data.GetFfigure()
	ret[2] = data.GetFuname()
	ret[3] = data.GetSendTime()
	ret[4] = data.GetId()
	ret[5] = data.GetUid()
	ret[6] = data.GetEtData()
	return ret
}
