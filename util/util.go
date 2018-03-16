package util

import (
	"sync"
	"time"

	"github.com/holyreaper/ggserver/common"
	. "github.com/holyreaper/ggserver/glog"
)

//MaxOffset .
const MaxOffset = 16777215

var offset int32
var lastTime int32
var mutex sync.Mutex

//GenerateID 根据当前id生成一个全服唯一的id
func GenerateID() (id int64) {
	defer func() {
		mutex.Unlock()
	}()
	mutex.Lock()
	now := int32(time.Now().Unix())
	if now != lastTime {
		offset = 0
	}
	if MaxOffset <= offset {
		LogFatal("util::GenerateID %d have got max offset ", MaxOffset)
		time.Sleep(time.Second)
	}

	serverid := int32(common.GetServerID())
	//先搞时间
	id += int64(now >> 24)
	id <<= 8
	id += int64(now >> 16)
	id <<= 8
	id += int64(now >> 8)
	id <<= 8
	id += int64(byte(now))
	//服务器id
	id <<= 8
	id += int64(serverid >> 8)
	id <<= 8
	id += int64(byte(serverid))
	//顺序
	id <<= 8
	id = id + int64(offset) + 1

	return
}
