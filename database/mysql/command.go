package mysql

import (
	"fmt"
)

//简单的mysql查询存储字段封装...

var command = [...]string{"select", "update", "insert", "delete"}

//Command mysql command
type Command struct {
	command string
	from    string
	where   string
	set     string
	tp      string
	//params  []interface{}
}

//Select select from db
func (cm *Command) Select(columns []string) *Command {
	if len(cm.command) != 0 {
		fmt.Println("Fatal Select error have call func yet ...")
	}
	cm.tp = "select"
	cm.command = "select "
	count := len(columns)
	for i, v := range columns {
		cm.command += v
		if count > 1 && i != count-1 {
			cm.command += ","
		}
	}
	return cm
}

//Update ...
func (cm *Command) Update(table string) *Command {
	if len(cm.command) != 0 {
		fmt.Println("Fatal Select error have call func yet ...")
	}
	cm.tp = "update"
	cm.command = "update " + table
	return cm
}

//Set set from db
func (cm *Command) Set(columns []string) *Command {
	if len(cm.set) != 0 {
		fmt.Println("Fatal Set error have call func yet ...")
	}
	cm.tp = "select"
	cm.command = "select "
	count := len(columns)
	for i, v := range columns {
		cm.command += v
		if count > 1 && i != count-1 {
			cm.command += ","
		}
	}
	return cm
}

//From select from db
func (cm *Command) From(fm string) *Command {
	if len(cm.from) != 0 {
		fmt.Println("Fatal From error have call func yet ...")
	} else {
		cm.from = "from " + fm
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
	cm.from = " where " + where + " = ?"

	return cm
}

//InsertInto ..
func (cm *Command) InsertInto(table string) *Command {

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
func (cm *Command) Query(params ...interface{}) (ret map[int32]map[string]interface{}) {
	return Query(cm.command+cm.from+cm.where, params)
}
