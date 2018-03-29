package packet

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/holyreaper/ggserver/util/convert"
)

//定义 client <-> server 的通信基本协议
/*
* 0000|0000|000...
* len |type|packet
 */
type Packet struct {
	Len  int32
	Type int32
	Data []byte
}

const (
	//PKUnKnown 初始
	PKUnKnown = iota
	//PKGHeartBeat 心跳
	PKGHeartBeat

	//PKGLogin 登录
	PKGLogin

	//PkLogOut 登出
	PkLogOut

	//PKGChat 聊天
	PKGChat
	//PKGChatPush 聊天推送
	PKGChatPush
)

const (
	//MAXPACKETLEN 包体最大长度
	MAXPACKETLEN = 4096 + 8
	//DEFAULTPACKLEN 默认包长度
	DEFAULTPACKLEN = 32
)

//Pack 打包
func (pk *Packet) Pack(tp int32, data proto.Message) error {

	dt, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("packet message error ", err)
		return err
	}
	pk.Len = 8 + int32(len(dt))
	pk.Type = tp
	pk.Data = dt
	return err
}

//UnPack 解析
func (pk *Packet) UnPack() (pb proto.Message, err error) {
	err = proto.Unmarshal(pk.Data, pb)
	if err != nil {
		fmt.Println("packet message error ", err)
		return
	}
	return
}

//FormatBuf 转化成buf
func (pk *Packet) FormatBuf() (buf []byte) {
	buf = append(buf, convert.Int32ToBytes(int32(pk.Len))...)
	buf = append(buf, convert.Int32ToBytes(int32(pk.Type))...)
	buf = append(buf, pk.Data...)
	return
}

//GetType 获取类型
func (pk *Packet) GetType() int32 {
	return pk.Type
}

//Clear clear data
func (pk *Packet) Clear() {
	pk.Type = PKUnKnown
	pk.Data = make([]byte, 0)
	pk.Len = 0
}
