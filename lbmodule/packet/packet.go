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
	//PKGHeartBeat 心跳
	PKGHeartBeat = iota + 1
	//PKGLogin 登录
	PKGLogin
	//PKGChat 聊天
	PKGChat
)

const (
	//MAXPACKETLEN 包体最大长度
	MAXPACKETLEN = 4096 + 8
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
