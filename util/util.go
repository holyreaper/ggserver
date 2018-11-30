package util

import (
	"sync"
	"time"

	"github.com/holyreaper/ggserver/common"
	. "github.com/holyreaper/ggserver/glog"
)

//IDGenerator .
type IDGenerator struct {
	maxOffset uint32
	offset    uint32
	lastTime  uint32
	mutex     sync.Mutex
}

//GenerateID 生成一个全服唯一的id
func (idg *IDGenerator) GenerateID() (id uint64) {
	defer func() {
		idg.mutex.Unlock()
	}()
	idg.mutex.Lock()
	now := uint32(time.Now().Unix())
	if now != idg.lastTime {
		idg.offset = 0
	}
	if idg.maxOffset <= idg.offset {
		LogInfo("util::GenerateID %d have got max offset ", idg.maxOffset)
		time.Sleep(time.Second)
	}
	idg.offset++
	serverid := int32(common.GetServerID())
	//先搞时间
	id += uint64(now)
	id <<= 32
	//服务器id
	id += uint64(serverid & 0x0000ffff)
	id <<= 16
	//顺序
	id = id + uint64(idg.offset) + 1
	return
}

//NewIDGenerator .
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		maxOffset: 65535,
		offset:    0,
		lastTime:  0,
	}
}
