package charmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/holyreaper/ggserver/lbmodule/packet"

	"github.com/golang/protobuf/proto"

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
		manager: make(map[UID]*CharManager),
		wrLock:  &sync.RWMutex{},
	}
)

//OnLineMng user  manager
type OnLineMng struct {
	manager map[UID]*CharManager
	wrLock  *sync.RWMutex
}

//AddUser 增加user
func (cm *OnLineMng) AddUser(rpack chan<- packet.Packet, spack chan<- packet.Packet, uid UID) bool {
	cm.wrLock.Lock()
	if mng, ok := cm.manager[uid]; ok {
		//TODO release module
		mng.LogOut(uid)
	}
	cmg := NewCharMng(rpack, spack, uid)
	cm.manager[uid] = cmg
	cm.wrLock.Unlock()
	if !cmg.Login(uid) {
		return false
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
	for k, v := range cm.manager {
		//TODO release module
		//		fmt.Println(v)
		v.LogOut(v.uid)
		delete(cm.manager, k)
	}
	return true
}

//SendMessageToUser SendMessageToUser
func (cm *OnLineMng) SendMessageToUser(uid UID, tp int32, data proto.Message) (err error) {
	fmt.Println("find user   ", uid)
	mng := cm.GetUser(uid)
	if mng != nil {
		fmt.Println("find user succ  ", uid)
		charmng := mng.(*CharManager)
		err = charmng.SendMessage(tp, data)
		return
	}
	err = errors.New("user not online ")
	fmt.Println("can not find user ", uid)
	return
}

//AddMessageToUser AddMessageToUser
func (cm *OnLineMng) AddMessageToUser(uid UID, tp int32, data proto.Message) (err error) {
	fmt.Println("find user   ", uid)
	mng := cm.GetUser(uid)
	if mng != nil {
		fmt.Println("find user succ  ", uid)
		charmng := mng.(*CharManager)
		err = charmng.AddMessage(tp, data)
		return
	}
	err = errors.New("user not online ")
	fmt.Println("can not find user ", uid)
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
func AddUser(rpack chan<- packet.Packet, spack chan<- packet.Packet, uid UID) bool {
	return onlineMng.AddUser(rpack, spack, uid)
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
func Login(rpack chan<- packet.Packet, spack chan<- packet.Packet, uid UID) interface{} {
	//各种login
	AddUser(rpack, spack, uid)
	return nil
}

//SendMessageToUser SendMessageToUser
func SendMessageToUser(uid UID, tp int32, data proto.Message) (err error) {
	err = onlineMng.SendMessageToUser(uid, tp, data)
	return
}

//玩家具体管理类的定义
//CharManager
type CharManager struct {
	Manager
	rpack     chan<- packet.Packet
	spack     chan<- packet.Packet
	rwlock    *sync.RWMutex
	keepAlive time.Duration
	//userMng
	userMng *UserMng
	chatMng *ChatMng
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

//SendMessage send data to user
func (cm *CharManager) SendMessage(tp int32, data proto.Message) (err error) {
	//SendMessage ..
	fmt.Println("SendMessage to user start ....")

	var pack packet.Packet
	err = pack.Pack(tp, data)
	if err != nil {
		fmt.Printf("sendMessage fail err %s", err)
		return
	}
	cm.spack <- pack
	return
}

//AddMessageToUser server  add  a Message to user's message channel  , pretent client message
func AddMessageToUser(uid UID, tp int32, data proto.Message) (err error) {
	err = onlineMng.AddMessageToUser(uid, tp, data)
	return
}

//AddMessage server  add  a Message to user's message channel  , pretent client message
func (cm *CharManager) AddMessage(tp int32, data proto.Message) (err error) {
	fmt.Println("addMessage to user start ...")
	var pack packet.Packet
	err = pack.Pack(tp, data)
	if err != nil {
		fmt.Printf("addMessage fail err %s", err)
		return
	}
	cm.rpack <- pack
	return
}

//NewCharMng new char mng
func NewCharMng(rpack chan<- packet.Packet, spack chan<- packet.Packet, uid UID) *CharManager {
	return &CharManager{
		rpack:   rpack,
		spack:   spack,
		userMng: NewUserMng(uid),
		chatMng: NewChatMng(uid),
		rwlock:  &sync.RWMutex{},
	}
}
