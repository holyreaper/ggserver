package mysqldao

import (
	"errors"

	"github.com/holyreaper/ggserver/database/mysql"
	. "github.com/holyreaper/ggserver/glog"
)

//登录相关的数据库操作
//userField .
var userField []string

const (
	//UserTableName .
	UserTableName = "t_user"
)
const (
	//UserUID .
	UserUID = "uid"
	//UserUname .
	UserUname = "uname"
	//UserCreateTime .
	UserCreateTime = "create_time"
	//UserLastLGITime .
	UserLastLGITime = "last_login_time"
	//UserLastLGOTime .
	UserLastLGOTime = "last_logout_time"
	//UserExp .
	UserExp = "exp"
	//UserLevel .
	UserLevel = "level"
	//UserFigure .
	UserFigure = "figure"
	//UserExtData .
	UserExtData = "et_data"
)

func init() {
	userField = []string{
		UserUID,
		UserUname,
		UserCreateTime,
		UserLastLGITime,
		UserLastLGOTime,
		UserExp,
		UserLevel,
		UserFigure,
		UserExtData}
}

//UserSelectSingle .
func UserSelectSingle(field []string, uid int32) (rst map[string]interface{}, err error) {
	if len(field) <= 0 {
		err = errors.New("empty field")
		return
	}
	command := mysql.Command{}
	command.Select(UserTableName, field).Where("uid")
	tmpRst, err := command.Query(uid)
	if err != nil {
		LogFatal("mysql SelectSingleUser error %v ", err)
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

//UserSelectMulity .
func UserSelectMulity(field []string, where []string, value []interface{}) (rst map[int32]map[string]interface{}, err error) {
	if len(field) <= 0 {
		err = errors.New("empty field")
		return
	}
	command := mysql.Command{}
	command.Select(UserTableName, field)
	for _, v := range where {
		command.Where(v)
	}
	rst, err = command.Query(value...)
	if err != nil {
		LogFatal("mysql UserSelectMulity error %v ", err)
		return
	}
	return
}

//UserUpdateSingle
func UserUpdateSingle(field []string, uid int32, value []interface{}) (err error) {
	if len(field) <= 0 {
		err = errors.New("empty field")
		return
	}

	command := mysql.Command{}
	command.Update(UserTableName).Set(field).Where("uid")

	tmpValue := make([]interface{}, len(value)+1)

	tmpValue[0] = uid

	tmpValue = append(tmpValue[0:], value[:])

	tmpRst, err := command.Exec(tmpValue...)
	if err != nil {
		LogFatal("mysql UserUpdateSingle Exec error %v ", err)
		return
	}
	cnt, err := tmpRst.RowsAffected()
	if err != nil {
		LogFatal("mysql UserUpdateSingle RowsAffected error %v ", err)
	}
	if cnt <= 0 {
		LogFatal(" mysql UserUpdateSingle RowsAffected rows 0  ")
		err = errors.New(" effect 0 rows ")
	}
	return
}

//UserInsertSingle .
func UserInsertSingle(value []interface{}) (err error) {

	command := mysql.Command{}
	command.InsertInto(UserTableName).Set(userField)

	tmpRst, err := command.Exec(value...)
	if err != nil {
		LogFatal("mysql UserUpdateSingle Exec error %v ", err)
		return
	}
	cnt, err := tmpRst.RowsAffected()
	if err != nil {
		LogFatal("mysql UserUpdateSingle UserInsertSingle error %v ", err)
	}
	if cnt <= 0 {
		LogFatal(" mysql UserUpdateSingle UserInsertSingle rows 0  ")
		err = errors.New(" effect 0 rows ")
	}
	return
}
