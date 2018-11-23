// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat.proto

/*
Package message is a generated protocol buffer package.

It is generated from these files:
	chat.proto
	message.proto
	user.proto

It has these top-level messages:
	ChatMessage
	ChatRequest
	ChatReply
	Message
	LoginRequest
	LoginReply
	LogOutRequest
	LogOutReply
	RegisterRequest
	RegisterReply
	KeepAliveRequest
	KeepAliveReply
*/
package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ChatMessage struct {
	Fuid    int64  `protobuf:"varint,1,opt,name=fuid" json:"fuid,omitempty"`
	Funame  string `protobuf:"bytes,2,opt,name=funame" json:"funame,omitempty"`
	Ffigure int32  `protobuf:"varint,3,opt,name=ffigure" json:"ffigure,omitempty"`
	Flevel  int32  `protobuf:"varint,4,opt,name=flevel" json:"flevel,omitempty"`
	Fvip    int32  `protobuf:"varint,5,opt,name=fvip" json:"fvip,omitempty"`
	Tuid    int64  `protobuf:"varint,6,opt,name=tuid" json:"tuid,omitempty"`
	Msg     string `protobuf:"bytes,7,opt,name=msg" json:"msg,omitempty"`
}

func (m *ChatMessage) Reset()                    { *m = ChatMessage{} }
func (m *ChatMessage) String() string            { return proto.CompactTextString(m) }
func (*ChatMessage) ProtoMessage()               {}
func (*ChatMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ChatMessage) GetFuid() int64 {
	if m != nil {
		return m.Fuid
	}
	return 0
}

func (m *ChatMessage) GetFuname() string {
	if m != nil {
		return m.Funame
	}
	return ""
}

func (m *ChatMessage) GetFfigure() int32 {
	if m != nil {
		return m.Ffigure
	}
	return 0
}

func (m *ChatMessage) GetFlevel() int32 {
	if m != nil {
		return m.Flevel
	}
	return 0
}

func (m *ChatMessage) GetFvip() int32 {
	if m != nil {
		return m.Fvip
	}
	return 0
}

func (m *ChatMessage) GetTuid() int64 {
	if m != nil {
		return m.Tuid
	}
	return 0
}

func (m *ChatMessage) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type ChatRequest struct {
	Msg *ChatMessage `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *ChatRequest) Reset()                    { *m = ChatRequest{} }
func (m *ChatRequest) String() string            { return proto.CompactTextString(m) }
func (*ChatRequest) ProtoMessage()               {}
func (*ChatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ChatRequest) GetMsg() *ChatMessage {
	if m != nil {
		return m.Msg
	}
	return nil
}

// The response message containing the greetings
type ChatReply struct {
	Result int32 `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *ChatReply) Reset()                    { *m = ChatReply{} }
func (m *ChatReply) String() string            { return proto.CompactTextString(m) }
func (*ChatReply) ProtoMessage()               {}
func (*ChatReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ChatReply) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

func init() {
	proto.RegisterType((*ChatMessage)(nil), "message.ChatMessage")
	proto.RegisterType((*ChatRequest)(nil), "message.ChatRequest")
	proto.RegisterType((*ChatReply)(nil), "message.ChatReply")
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xcd, 0x6a, 0xc4, 0x20,
	0x14, 0x85, 0xb1, 0xf9, 0x23, 0x37, 0x9b, 0x22, 0xa5, 0xb8, 0x0c, 0x29, 0x94, 0xac, 0xb2, 0x68,
	0xe9, 0x13, 0x74, 0xdd, 0x8d, 0x6f, 0x60, 0xdb, 0x6b, 0x12, 0x30, 0x93, 0x4c, 0xd4, 0xc0, 0x3c,
	0xd0, 0xbc, 0xe7, 0xe0, 0xd5, 0x81, 0xd9, 0x9d, 0xef, 0xea, 0xf1, 0x9c, 0x2b, 0xc0, 0xdf, 0xa4,
	0xdc, 0xb0, 0xed, 0xab, 0x5b, 0x79, 0xb5, 0xa0, 0xb5, 0x6a, 0xc4, 0xee, 0xca, 0xa0, 0xf9, 0x9e,
	0x94, 0xfb, 0x89, 0xcc, 0x39, 0xe4, 0xda, 0xcf, 0xff, 0x82, 0xb5, 0xac, 0xcf, 0x24, 0x69, 0xfe,
	0x0a, 0xa5, 0xf6, 0x27, 0xb5, 0xa0, 0x78, 0x6a, 0x59, 0x5f, 0xcb, 0x44, 0x5c, 0x40, 0xa5, 0xf5,
	0x3c, 0xfa, 0x1d, 0x45, 0xd6, 0xb2, 0xbe, 0x90, 0x77, 0x24, 0x87, 0xc1, 0x03, 0x8d, 0xc8, 0xe9,
	0x20, 0x11, 0xbd, 0x7e, 0xcc, 0x9b, 0x28, 0x68, 0x4a, 0x3a, 0xcc, 0x5c, 0x48, 0x2c, 0x63, 0x62,
	0xd0, 0xfc, 0x19, 0xb2, 0xc5, 0x8e, 0xa2, 0xa2, 0xb8, 0x20, 0xbb, 0xaf, 0x58, 0x53, 0xe2, 0xd9,
	0xa3, 0x75, 0xfc, 0x3d, 0x5e, 0x08, 0x2d, 0x9b, 0x8f, 0x97, 0x21, 0x6d, 0x33, 0x3c, 0x6c, 0x12,
	0x6d, 0x6f, 0x50, 0x47, 0xdb, 0x66, 0x2e, 0xa1, 0xd5, 0x8e, 0xd6, 0x1b, 0x47, 0xbe, 0x42, 0x26,
	0xfa, 0x2d, 0xe9, 0x4f, 0x3e, 0x6f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xce, 0xb7, 0xb8, 0xdc, 0x21,
	0x01, 0x00, 0x00,
}
