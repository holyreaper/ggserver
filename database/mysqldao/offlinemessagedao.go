package mysqldao

import (
	"errors"

	"github.com/holyreaper/ggserver/database/mysql"
	. "github.com/holyreaper/ggserver/glog"
)

var OffLineMsgField []string

const (
	//OffTableName  .
	OffTableName = "t_offline_message"
)

const (
	//OfflineMsgID .
	OfflineMsgID = "id"
	//OfflineMsgUID  .
	OfflineMsgUID = "uid"
	//OfflineMsgSendTime  .
	OfflineMsgSendTime = "send_time"
	//OfflineMsgFUID .
	OfflineMsgFUID = "fuid"
	//OfflineMsgFUname .
	OfflineMsgFUname = "funame"
	//OfflineMsgFfigure  .
	OfflineMsgFfigure = "ffigure"
	// OfflineMsgETData .
	OfflineMsgETData = "et_data"
)

func init() {
	OffLineMsgField = []string{
		OfflineMsgID,
		OfflineMsgUID,
		OfflineMsgSendTime,
		OfflineMsgFUID,
		OfflineMsgFUname,
		OfflineMsgFfigure,
		OfflineMsgETData}
}

//OffLineMsgSelectSingle .
func OffLineMsgSelectSingle(field []string, uid int32) (rst map[string]interface{}, err error) {
	if len(field) <= 0 {
		err = errors.New("empty field")
		return
	}
	command := mysql.Command{}
	command.Select(UserTableName, field).Where("uid")
	tmpRst, err := command.Query(uid)
	if err != nil {
		LogFatal("mysql OffLineMsgSelectSingle error %v ", err)
		return
	}
	value, ok := tmpRst[0]
	if ok {
		rst = value
	} else {
		err = errors.New("no data ")
	}

	return
}

//OffLineMsgSelectMulity .
func OffLineMsgSelectMulity(field []string, where []string, value []interface{}) (rst map[int32]map[string]interface{}, err error) {
	if len(field) <= 0 {
		err = errors.New("empty field")
		return
	}
	command := mysql.Command{}
	command.Select(OffTableName, field)
	for _, v := range where {
		command.Where(v)
	}
	rst, err = command.Query(value...)
	if err != nil {
		LogFatal("mysql OffLineMsgSelectMulity error %v ", err)
		return
	}
	return
}

//OffLineMsgDelete .
func OffLineMsgDelete(where []string, value []interface{}) (err error) {

	command := mysql.Command{}
	command.Delete(UserTableName)
	for _, v := range where {
		command.Where(v)
	}
	tmpRst, err := command.Exec(value...)
	if err != nil {
		LogFatal("mysql OffLineMsgDelete Exec error %v ", err)
		return
	}
	cnt, err := tmpRst.RowsAffected()
	if err != nil {
		LogFatal("mysql OffLineMsgDelete RowsAffected error %v ", err)
	}
	if cnt <= 0 {
		LogFatal(" mysql OffLineMsgDelete RowsAffected rows 0  ")
		err = errors.New(" effect 0 rows ")
	}
	return
}

//OffLineMsgInsert .
func OffLineMsgInsert(value []interface{}) (err error) {
	command := mysql.Command{}
	command.InsertInto(UserTableName)
	command.Set(OffLineMsgField)
	tmpRst, err := command.Exec(value...)
	if err != nil {
		LogFatal("mysql OffLineMsgDelete Exec error %v ", err)
		return
	}
	cnt, err := tmpRst.RowsAffected()
	if err != nil {
		LogFatal("mysql OffLineMsgInsert RowsAffected error %v ", err)
	}
	if cnt <= 0 {
		LogFatal(" mysql OffLineMsgInsert RowsAffected rows 0  ")
		err = errors.New(" effect 0 rows ")
	}
	return
}
