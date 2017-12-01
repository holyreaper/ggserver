package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "gg:123456@tcp(192.168.1.177:3306)/gg?charset=utf8")
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.Ping()
}

//Query ...
func Query(req string, params ...interface{}) (ret map[int32]map[string]interface{}, err error) {
	stmt, err := db.Prepare(req)
	if err != nil {
		fmt.Println("mysql query fail ", err)
		return
	}
	rows, err := stmt.Query(params...)
	if err != nil {
		return
	}
	stmt.Close()
	return FetchResult(rows), err
}

//Exec ...
func Exec(req string, params ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(req)
	if err != nil {
		fmt.Println("mysql query fail ", err)
		return nil, errors.New("params is nil")
	}

	result, err := stmt.Exec(params...)
	if err != nil {
		return nil, err
	}
	stmt.Close()
	return result, err
}

//FetchResult ...
func FetchResult(rows *sql.Rows) (ret map[int32]map[string]interface{}) {
	columns, _ := rows.Columns()
	args := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	ret = make(map[int32]map[string]interface{})
	for i := range values {
		args[i] = &values[i]
	}
	var k int32
	k = 0
	for rows.Next() {
		columnMap := make(map[string]interface{})
		err := rows.Scan(args...)
		if err != nil {
			fmt.Println("mysql FetchResult err  ", err)
			return
		}
		for i, v := range values {
			columnMap[columns[i]] = v
		}
		ret[k] = columnMap
		k++
	}
	//rows.Close()
	return
}

//Close ...
func Close() {
	if db != nil {
		db.Close()
	}
}
