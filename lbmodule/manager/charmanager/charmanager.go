package charmanager

import (
	"errors"
	"fmt"
	"net"
	"time"

	//def
	"sync"

	. "github.com/holyreaper/ggserver/def"
)

const (
	//UserMngEnum ...
	UserMngEnum = iota
	//ChatMngENum  ...
	ChatMngENum
)

//Manager manager interface
type Manager struct {
	uid UID
}

//Login .
func (*Manager) Login(UID) bool {
	return false
}

//LogOut .
func (*Manager) LogOut(UID) bool {
	return false
}

var (
	onlineMng = &OnLineMng{
		manager: make(map[UID]map[int32]*Manager),
		wrLock:  sync.RWMutex{},
	}
)

//OnLineMng user  manager
type OnLineMng struct {
	manager map[UID]map[int32]*Manager
	wrLock  sync.RWMutex
}

//AddUser 增加user
func (cm *OnLineMng) AddUser(uid UID) bool {
	cm.wrLock.Lock()
	defer cm.wrLock.Unlock()
	if _, ok := cm.manager[uid]; ok {
		//TODO release module
	}

	return true
}

//DelUser ...
func (cm *OnLineMng) DelUser(uid UID) bool {
	cm.wrLock.Lock()
	defer cm.wrLock.Unlock()
	if _, ok := cm.manager[uid]; !ok {
		return false
	}
	//TODO release module
	delete(cm.manager, uid)
	return true
}

//DeleteAll delete all
func (cm *OnLineMng) DeleteAll() bool {
	cm.wrLock.Lock()
	defer cm.wrLock.Unlock()
	for k, _ := range cm.manager {
		//TODO release module
		//		fmt.Println(v)
		delete(cm.manager, k)
	}
	return true
}

//SendMessageToUser SendMessageToUser
func (cm *OnLineMng) SendMessageToUser(uid UID, data []byte) (err error) {
	mng := cm.GetUser(uid)
	if mng != nil {
		charmng := mng.(*CharManager)
		err = charmng.SendMessage(data)
		return
	}
	err = errors.New("user not online ")
	return
}

//GetUser get user
func (cm *OnLineMng) GetUser(uid UID) interface{} {
	cm.wrLock.RLock()
	defer cm.wrLock.RUnlock()
	if value, ok := cm.manager[uid]; ok {
		return value
	}

	return nil
}

//AddUser add user
func AddUser(uid UID) bool {
	return onlineMng.AddUser(uid)
}

//DelUser delete user
func DelUser(uid UID) bool {
	return onlineMng.DelUser(uid)
}

//DeleteAll delete all  user
func DeleteAll() bool {
	return onlineMng.DeleteAll()
}

//GetUser get user info
func GetUser(uid UID) interface{} {
	return onlineMng.GetUser(uid)
}

//Login user login
func Login(cnn *net.TCPConn, uid UID) interface{} {
	//各种login

	return nil
}

//SendMessageToUser SendMessageToUser
func SendMessageToUser(uid UID, data []byte) (err error) {
	err = onlineMng.SendMessageToUser(uid, data)
	return
}

//玩家具体管理类的定义
//CharManager
type CharManager struct {
	Manager
	cnn       *net.TCPConn
	keepAlive time.Duration
	//userMng
	userMng *Manager
	chatMng *Manager
}

//Login .
func (cm *CharManager) Login(uid UID) bool {
	cm.userMng.Login(uid)
	cm.chatMng.Login(uid)
	return false
}

//LogOut .
func (cm *CharManager) LogOut(uid UID) bool {
	cm.userMng.LogOut(uid)
	cm.chatMng.LogOut(uid)
	return false
}

//SendMessage ..
func (cm *CharManager) SendMessage(data []byte) (err error) {
	cm.cnn.SetWriteDeadline(time.Now().Add(time.Duration(1000)))
	ln, err := cm.cnn.Write(data)
	if err != nil || ln != len(data) {
		fmt.Println("send byte to client  err ", err)
	}
	return
}
