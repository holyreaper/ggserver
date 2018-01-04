package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//confMap := make(map[string]interface)
//fileName 服务器配置
var fileName = "service.json"

//GetConf 获取服务器配置
func GetConf() []map[string]string {
	var sliceConfMap []map[string]string
	conf, err := readFile()
	if err == nil {
		jserr := json.Unmarshal(conf, &sliceConfMap)
		if jserr == nil {
			return sliceConfMap
		}
		fmt.Printf("unMarshal data error  :%s", jserr.Error())
	}
	return nil

}
func init() {
	fmt.Println("init conf...")
}

func readFile() ([]byte, error) {
	conf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("cannot readfile%s error msg :%s", fileName, err.Error())
		return nil, err
	}
	return conf, nil
}
