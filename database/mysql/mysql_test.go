package mysql

import (
	"fmt"
	"testing"
)

//test select
func Test_Select(t *testing.T) {
	bb := []string{"uid", "name", "level"}
	command := Command{}
	command.Select("t_user", bb).Where("uid")
	for i := 1; i < 10; i++ {
		result, err := command.Query("21001")
		if err != nil {
			fmt.Println("err ..", err)
		} else {
			fmt.Println("have no err ..")
		}
		for _, v := range result {
			fmt.Println("get result :")
			for m, n := range v {
				fmt.Printf("key:%v,value:%v", m, n)
			}
		}
		
	}
}
