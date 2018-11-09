package pool

import (
	"container/list"
	"sync"
	"github.com/holyreaper/ggserver/glog"
)
/** 
	主要提供连接池的操作
*/
//IPoll interface
type IPool interface {
	//初始化
	Init(addr string, port uint, maxSize uint, minSize uint, maxIdleSize uint, param ...interface{})
	//Get 获取一个实例
	Get() interface{}
	//Back 还回去一个实例
	Back(interface{})
	//Release 释放
	Release()
}

//Param 参数
type Param struct {
	//最大数量
	maxSize uint
	//最小数量
	minSize uint
	//最大空闲数量
	maxIdleSize uint
	//地址
	addr string
	//端口
	port uint
	//连接参数
	param []interface{}
	//互斥量
	mu sync.Mutex

	//实例
	instance *list.List
}

//Pool 模板
type Pool struct {
	Param
}

//Init 初始化
func (t *Pool) Init(addr string, port uint, maxSize uint, minSize uint, maxIdleSize uint, param ...interface{}) {
	t.addr = addr
	t.port = port
	t.maxSize = maxSize
	t.maxIdleSize = maxIdleSize
	t.minSize = minSize
	t.param = make([]interface{}, len(param))
	for p, v := range param {
		t.param[p] = v
	}
	t.instance = list.New()
	t.instance.Init()
	t.mu = sync.Mutex{}

	//初始化连接池
	t.resize()
}

//连接池初始化
func (t *Pool) resize() {

}

//Get 获取一个实例
func (t *Pool) Get() interface{} {
	t.mu.Lock()
	defer func() {
		t.mu.Unlock()
	}()
	e := t.instance.Front()
	if e != nil {
		return t.instance.Remove(e)
	}
	return nil
}

//Back 还回去
func (t *Pool) Back(value interface{}) {
	t.mu.Lock()
	defer func() {
		t.mu.Unlock()
	}()
	t.instance.PushBack(value)
}

//Release 连接池中的对象 必须要实现Release接口
type Release interface {
	release()
}

//Release 释放全部
func (t *Pool) Release() {
	for k := t.instance.Front(); k != nil; k = k.Next() {
		v, ok := k.Value.(Release)
		if ok {
			v.release()
		} else {
			glog.LogFatal("connect pool err no release func ")
		}
	}
	t.instance.Init()
}
