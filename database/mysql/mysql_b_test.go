package mysql

import (
	"testing"
)

//Benchmark_Select ...
func Benchmark_Select(b *testing.B) {
	bb := []string{"uid", "name", "level"}
	command := Command{}
	command.Select("t_user", bb).Where("uid")
	for i := 1; i < b.N; i++ {
		//fmt.Println("helo....", b.N)
		command.Query("21002")

		//fmt.Println("return ....", b.N)
	}
}
