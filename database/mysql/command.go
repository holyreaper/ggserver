package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

//简单的mysql查询存储字段封装...

//Command mysql command
type Command struct {
	command string
	table   string
	what    string
	from    string
	where   string
	//params  []interface{}
}

//Select select from db
func (cm *Command) Select(table string, columns []string) *Command {
	cm.command = "select"
	cm.table = table
	cm.from = "from"
	count := len(columns)
	if count == 0 {
		cm.what = " * "
	} else {
		for i, v := range columns {
			cm.what += v
			if count > 1 && i != count-1 {
				cm.what += ","
			}
		}
	}
	return cm
}

//Update ...
func (cm *Command) Update(table string) *Command {
	cm.command = "update"
	cm.table = table
	return cm
}

//InsertInto ...
func (cm *Command) InsertInto(table string) *Command {
	cm.command = "insert"
	cm.table = table
	return cm
}

//Delete ...
func (cm *Command) Delete(table string) *Command {
	cm.command = "delete"
	cm.from = "from"
	cm.table = table
	return cm
}

//Set set from db
func (cm *Command) Set(columns []string) *Command {
	cm.what = "set "
	count := len(columns)
	if count == 0 {
		fmt.Println("Fatal Set error column is null  ...")
	}
	for i, v := range columns {
		cm.what += v
		cm.what += " = ? "
		if count > 1 && i != count-1 {
			cm.what += ","
		}
	}
	return cm
}

//Where select from db where
func (cm *Command) Where(where string) *Command {
	if len(cm.from) != 0 {
		cm.where += " and "
		cm.where += where
		cm.where += " = ? "
	}
	cm.where = " where " + where + " = ?"

	return cm
}

//Params ...
/*
func (cm *Command) Params(params ...interface{}) *Command {
	cm.params = make([]interface{}, len(params))
	for i, v := range params {
		cm.params[i] = v
	}
	return cm
}
*/

//Query ...
func (cm *Command) Query(params ...interface{}) (map[int32]map[string]interface{}, error) {
	switch cm.command {
	case "select":
		return Query(cm.command+" "+cm.what+" "+cm.from+" "+cm.table+" "+cm.where, params...)
	default:
		fmt.Println("Mysql Query unsupport command ", cm.command)
		return nil, errors.New(" mysql Query not support command ")
	}
}

//Exec ...
func (cm *Command) Exec(params ...interface{}) (sql.Result, error) {
	switch cm.command {
	case "update":
		return Exec(cm.command+" "+cm.table+" "+cm.what+" "+cm.where, params...)
	case "delete":
		return Exec(cm.command+" "+cm.from+" "+cm.table+" "+cm.where, params...)
	case "insert":
		return Exec(cm.command+" "+cm.table+" "+cm.what+" "+cm.where, params...)
	default:
		fmt.Println("Mysql Exec unsupport command ", cm.command)
		return nil, errors.New("unsupport command ")
	}
}

//Clear ...
func (cm *Command) Clear() {
	cm.command = ""
	cm.table = ""
	cm.what = ""
	cm.from = ""
	cm.where = ""
}
