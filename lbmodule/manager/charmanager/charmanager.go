package charmanager

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"

	"sync"

	. "github.com/holyreaper/ggserver/def"
	"github.com/holyreaper/ggserver/util/convert"
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
func (cm *OnLineMng) AddUser(cnn *net.TCPConn, uid UID) bool {
	cm.wrLock.Lock()
	if mng, ok := cm.manager[uid]; ok {
		//TODO release module
		mng.LogOut(uid)
	}
	cmg := NewCharMng(cnn, uid)
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
func AddUser(cnn *net.TCPConn, uid UID) bool {
	return onlineMng.AddUser(cnn, uid)
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
	onlineMng.AddUser(cnn, uid)
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
	cnn       *net.TCPConn
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

	cm.rwlock.Lock()
	defer cm.rwlock.Unlock()

	currSend := 0
	var sendbuf []byte
	rsp, err := proto.Marshal(data)
	if err != nil {
		fmt.Println(" SendMessage to client  proto message mashal err \r\n", err)
		return
	}
	sendbuf = append(sendbuf, convert.Int32ToBytes(int32(len(rsp)+8))...)
	sendbuf = append(sendbuf, convert.Int32ToBytes(int32(tp))...)
	sendbuf = append(sendbuf, rsp...)

	cm.cnn.SetWriteDeadline(time.Now().Add(time.Duration(10e9)))
	for {
		ln, err2 := cm.cnn.Write(sendbuf[currSend:])
		if err2 != nil && ln != len(sendbuf[currSend:]) {
			if IsTimeOut(err2) {
				currSend += ln
				continue
			} else {
				fmt.Println("send byte to client  err ", err)
				return err2
			}
		}
		if currSend+ln >= len(sendbuf) {
			break
		}
	}
	fmt.Println("send data to client succ !!!")
	return
}

//NewCharMng new char mng
func NewCharMng(cn *net.TCPConn, uid UID) *CharManager {
	return &CharManager{
		keepAlive: 1e9,
		cnn:       cn,
		userMng:   NewUserMng(uid),
		chatMng:   NewChatMng(uid),
		rwlock:    &sync.RWMutex{},
	}
}
