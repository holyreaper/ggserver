// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dbrpc.proto

/*
Package dbrpcpt is a generated protocol buffer package.

It is generated from these files:
	dbrpc.proto

It has these top-level messages:
	KeepAliveRequest
	KeepAliveReply
	LoginRequest
	LoginReply
*/
package dbrpcpt

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// KeepAlive
type KeepAliveRequest struct {
	Time int64 `protobuf:"varint,1,opt,name=time" json:"time,omitempty"`
}

func (m *KeepAliveRequest) Reset()                    { *m = KeepAliveRequest{} }
func (m *KeepAliveRequest) String() string            { return proto.CompactTextString(m) }
func (*KeepAliveRequest) ProtoMessage()               {}
func (*KeepAliveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *KeepAliveRequest) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type KeepAliveReply struct {
	Result int32 `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *KeepAliveReply) Reset()                    { *m = KeepAliveReply{} }
func (m *KeepAliveReply) String() string            { return proto.CompactTextString(m) }
func (*KeepAliveReply) ProtoMessage()               {}
func (*KeepAliveReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *KeepAliveReply) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

// login
type LoginRequest struct {
	Uid int64 `protobuf:"varint,1,opt,name=uid" json:"uid,omitempty"`
}

func (m *LoginRequest) Reset()                    { *m = LoginRequest{} }
func (m *LoginRequest) String() string            { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()               {}
func (*LoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *LoginRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type LoginReply struct {
	Uid        int64  `protobuf:"varint,1,opt,name=uid" json:"uid,omitempty"`
	Createtime int64  `protobuf:"varint,2,opt,name=createtime" json:"createtime,omitempty"`
	Name       string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Level      int32  `protobuf:"varint,4,opt,name=level" json:"level,omitempty"`
	Exp        int32  `protobuf:"varint,5,opt,name=exp" json:"exp,omitempty"`
	Data       string `protobuf:"bytes,6,opt,name=data" json:"data,omitempty"`
}

func (m *LoginReply) Reset()                    { *m = LoginReply{} }
func (m *LoginReply) String() string            { return proto.CompactTextString(m) }
func (*LoginReply) ProtoMessage()               {}
func (*LoginReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LoginReply) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *LoginReply) GetCreatetime() int64 {
	if m != nil {
		return m.Createtime
	}
	return 0
}

func (m *LoginReply) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LoginReply) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *LoginReply) GetExp() int32 {
	if m != nil {
		return m.Exp
	}
	return 0
}

func (m *LoginReply) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*KeepAliveRequest)(nil), "dbrpcpt.KeepAliveRequest")
	proto.RegisterType((*KeepAliveReply)(nil), "dbrpcpt.KeepAliveReply")
	proto.RegisterType((*LoginRequest)(nil), "dbrpcpt.LoginRequest")
	proto.RegisterType((*LoginReply)(nil), "dbrpcpt.LoginReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DBRPCServer service

type DBRPCServerClient interface {
	// KeepAlive rpc
	KeepAlive(ctx context.Context, in *KeepAliveRequest, opts ...grpc.CallOption) (*KeepAliveReply, error)
	// Login
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
}

type dBRPCServerClient struct {
	cc *grpc.ClientConn
}

func NewDBRPCServerClient(cc *grpc.ClientConn) DBRPCServerClient {
	return &dBRPCServerClient{cc}
}

func (c *dBRPCServerClient) KeepAlive(ctx context.Context, in *KeepAliveRequest, opts ...grpc.CallOption) (*KeepAliveReply, error) {
	out := new(KeepAliveReply)
	err := grpc.Invoke(ctx, "/dbrpcpt.DBRPCServer/KeepAlive", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBRPCServerClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := grpc.Invoke(ctx, "/dbrpcpt.DBRPCServer/Login", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DBRPCServer service

type DBRPCServerServer interface {
	// KeepAlive rpc
	KeepAlive(context.Context, *KeepAliveRequest) (*KeepAliveReply, error)
	// Login
	Login(context.Context, *LoginRequest) (*LoginReply, error)
}

func RegisterDBRPCServerServer(s *grpc.Server, srv DBRPCServerServer) {
	s.RegisterService(&_DBRPCServer_serviceDesc, srv)
}

func _DBRPCServer_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeepAliveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBRPCServerServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dbrpcpt.DBRPCServer/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBRPCServerServer).KeepAlive(ctx, req.(*KeepAliveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBRPCServer_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBRPCServerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dbrpcpt.DBRPCServer/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBRPCServerServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DBRPCServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dbrpcpt.DBRPCServer",
	HandlerType: (*DBRPCServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "KeepAlive",
			Handler:    _DBRPCServer_KeepAlive_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _DBRPCServer_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dbrpc.proto",
}

func init() { proto.RegisterFile("dbrpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 264 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0xcb, 0x4e, 0xc3, 0x30,
	0x10, 0xc4, 0xa4, 0x09, 0xea, 0x16, 0xa1, 0x6a, 0x79, 0x99, 0x1e, 0x50, 0xe4, 0x03, 0xca, 0x29,
	0x07, 0x10, 0x1f, 0x50, 0xe0, 0x06, 0x07, 0x64, 0xbe, 0x20, 0x6d, 0x56, 0x28, 0x92, 0xdb, 0x18,
	0xd7, 0x89, 0xe8, 0x17, 0x70, 0xe3, 0x9b, 0x91, 0xb7, 0xa1, 0x8a, 0x02, 0xb7, 0xd9, 0xf5, 0xcc,
	0x68, 0x76, 0x0c, 0x93, 0x72, 0xe1, 0xec, 0x32, 0xb7, 0xae, 0xf6, 0x35, 0x1e, 0xf1, 0x60, 0xbd,
	0xba, 0x81, 0xe9, 0x33, 0x91, 0x9d, 0x9b, 0xaa, 0x25, 0x4d, 0x1f, 0x0d, 0x6d, 0x3c, 0x22, 0x8c,
	0x7c, 0xb5, 0x22, 0x29, 0x52, 0x91, 0x45, 0x9a, 0xb1, 0xca, 0xe0, 0xa4, 0xc7, 0xb3, 0x66, 0x8b,
	0x17, 0x90, 0x38, 0xda, 0x34, 0xc6, 0x33, 0x2f, 0xd6, 0xdd, 0xa4, 0x52, 0x38, 0x7e, 0xa9, 0xdf,
	0xab, 0xf5, 0xaf, 0xdb, 0x14, 0xa2, 0xa6, 0x2a, 0x3b, 0xb3, 0x00, 0xd5, 0xb7, 0x00, 0xe8, 0x28,
	0xc1, 0xe8, 0x0f, 0x01, 0xaf, 0x01, 0x96, 0x8e, 0x0a, 0x4f, 0x1c, 0xe3, 0x90, 0x1f, 0x7a, 0x9b,
	0x10, 0x70, 0x5d, 0xac, 0x48, 0x46, 0xa9, 0xc8, 0xc6, 0x9a, 0x31, 0x9e, 0x41, 0x6c, 0xa8, 0x25,
	0x23, 0x47, 0x9c, 0x66, 0x37, 0x04, 0x6f, 0xfa, 0xb4, 0x32, 0xe6, 0x5d, 0x80, 0x41, 0x5b, 0x16,
	0xbe, 0x90, 0xc9, 0x4e, 0x1b, 0xf0, 0xed, 0x97, 0x80, 0xc9, 0xd3, 0x83, 0x7e, 0x7d, 0x7c, 0x23,
	0xd7, 0x92, 0xc3, 0x39, 0x8c, 0xf7, 0xc7, 0xe2, 0x55, 0xde, 0x75, 0x95, 0x0f, 0x8b, 0x9a, 0x5d,
	0xfe, 0xf7, 0x64, 0xcd, 0x56, 0x1d, 0xe0, 0x3d, 0xc4, 0x7c, 0x22, 0x9e, 0xef, 0x39, 0xfd, 0x56,
	0x66, 0xa7, 0xc3, 0x35, 0xcb, 0x16, 0x09, 0x7f, 0xcf, 0xdd, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x96, 0x49, 0x29, 0xc7, 0xad, 0x01, 0x00, 0x00,
}