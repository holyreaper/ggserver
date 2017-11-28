package mysql

import (
	"database/sql"
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
func Query(req string, params ...interface{}) (ret map[int32]map[string]interface{}) {
	rows, err := db.Query(req, params)
	if err != nil {
		fmt.Println("mysql query fail ", err)
		return
	}
	return FetchResult(rows)
}

//Exec ...
func Exec(req string) bool {
	return true
}

//FetchResult ...
func FetchResult(rows *sql.Rows) (ret map[int32]map[string]interface{}) {
	columns, _ := rows.Columns()
	args := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		args[i] = &values[i]
	}
	var k int32
	k = 0
	for rows.Next() {
		err := rows.Scan(args)
		if err != nil {
			fmt.Println("mysql FetchResult err  ", err)
			return
		}
		for i, v := range values {
			ret[k][columns[i]] = v
		}
		k++
	}
	return
}
