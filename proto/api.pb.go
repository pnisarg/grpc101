// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	PingMessage
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type PingMessage struct {
	Greeting string `protobuf:"bytes,1,opt,name=greeting" json:"greeting,omitempty"`
}

func (m *PingMessage) Reset()                    { *m = PingMessage{} }
func (m *PingMessage) String() string            { return proto1.CompactTextString(m) }
func (*PingMessage) ProtoMessage()               {}
func (*PingMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PingMessage) GetGreeting() string {
	if m != nil {
		return m.Greeting
	}
	return ""
}

func init() {
	proto1.RegisterType((*PingMessage)(nil), "proto.PingMessage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Ping service

type PingClient interface {
	SayHello(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*PingMessage, error)
}

type pingClient struct {
	cc *grpc.ClientConn
}

func NewPingClient(cc *grpc.ClientConn) PingClient {
	return &pingClient{cc}
}

func (c *pingClient) SayHello(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*PingMessage, error) {
	out := new(PingMessage)
	err := grpc.Invoke(ctx, "/proto.Ping/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Ping service

type PingServer interface {
	SayHello(context.Context, *PingMessage) (*PingMessage, error)
}

func RegisterPingServer(s *grpc.Server, srv PingServer) {
	s.RegisterService(&_Ping_serviceDesc, srv)
}

func _Ping_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Ping/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingServer).SayHello(ctx, req.(*PingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ping_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Ping",
	HandlerType: (*PingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Ping_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto1.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0x32, 0xe9, 0xf9, 0xf9, 0xe9, 0x39,
	0xa9, 0xfa, 0x89, 0x05, 0x99, 0xfa, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99, 0xf9, 0x79,
	0xc5, 0x10, 0x45, 0x4a, 0x9a, 0x5c, 0xdc, 0x01, 0x99, 0x79, 0xe9, 0xbe, 0xa9, 0xc5, 0xc5, 0x89,
	0xe9, 0xa9, 0x42, 0x52, 0x5c, 0x1c, 0xe9, 0x45, 0xa9, 0xa9, 0x25, 0x99, 0x79, 0xe9, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x70, 0xbe, 0x91, 0x3f, 0x17, 0x0b, 0x48, 0xa9, 0x90, 0x3b, 0x17,
	0x47, 0x70, 0x62, 0xa5, 0x47, 0x6a, 0x4e, 0x4e, 0xbe, 0x90, 0x10, 0xc4, 0x18, 0x3d, 0x24, 0x33,
	0xa4, 0xb0, 0x88, 0x29, 0x09, 0x37, 0x5d, 0x7e, 0x32, 0x99, 0x89, 0x57, 0x89, 0x43, 0xbf, 0xcc,
	0x50, 0xbf, 0x20, 0x33, 0x2f, 0xdd, 0x8a, 0x51, 0x2b, 0x89, 0x0d, 0xac, 0xce, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0xcc, 0xfe, 0x04, 0x6b, 0xb4, 0x00, 0x00, 0x00,
}
