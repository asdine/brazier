// Code generated by protoc-gen-go.
// source: saver.proto
// DO NOT EDIT!

/*
Package internal is a generated protocol buffer package.

It is generated from these files:
	saver.proto
	getter.proto
	deleter.proto
	lister.proto

It has these top-level messages:
	SaveRequest
	SaveReply
	GetRequest
	GetReply
	DeleteRequest
	DeleteReply
	ListRequest
	ListReply
*/
package internal

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

// The request message containing the user's data.
type SaveRequest struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Data   []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *SaveRequest) Reset()                    { *m = SaveRequest{} }
func (m *SaveRequest) String() string            { return proto.CompactTextString(m) }
func (*SaveRequest) ProtoMessage()               {}
func (*SaveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The response message containing the status
type SaveReply struct {
	Status int32 `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
}

func (m *SaveReply) Reset()                    { *m = SaveReply{} }
func (m *SaveReply) String() string            { return proto.CompactTextString(m) }
func (*SaveReply) ProtoMessage()               {}
func (*SaveReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*SaveRequest)(nil), "internal.SaveRequest")
	proto.RegisterType((*SaveReply)(nil), "internal.SaveReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Saver service

type SaverClient interface {
	// Saves user data
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error)
}

type saverClient struct {
	cc *grpc.ClientConn
}

func NewSaverClient(cc *grpc.ClientConn) SaverClient {
	return &saverClient{cc}
}

func (c *saverClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error) {
	out := new(SaveReply)
	err := grpc.Invoke(ctx, "/internal.Saver/Save", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Saver service

type SaverServer interface {
	// Saves user data
	Save(context.Context, *SaveRequest) (*SaveReply, error)
}

func RegisterSaverServer(s *grpc.Server, srv SaverServer) {
	s.RegisterService(&_Saver_serviceDesc, srv)
}

func _Saver_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SaverServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/internal.Saver/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SaverServer).Save(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Saver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "internal.Saver",
	HandlerType: (*SaverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _Saver_Save_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("saver.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 169 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x4e, 0x2c, 0x4b,
	0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xc8, 0xcc, 0x2b, 0x49, 0x2d, 0xca, 0x4b,
	0xcc, 0x51, 0xf2, 0xe6, 0xe2, 0x0e, 0x4e, 0x2c, 0x4b, 0x0d, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e,
	0x11, 0x12, 0xe3, 0x62, 0x4b, 0x2a, 0x4d, 0xce, 0x4e, 0x2d, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0,
	0x0c, 0x82, 0xf2, 0x84, 0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x98, 0xc0, 0x82, 0x20, 0xa6,
	0x90, 0x10, 0x17, 0x4b, 0x4a, 0x62, 0x49, 0xa2, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x98,
	0xad, 0xa4, 0xcc, 0xc5, 0x09, 0x31, 0xac, 0x20, 0xa7, 0x12, 0x64, 0x54, 0x71, 0x49, 0x62, 0x49,
	0x69, 0x31, 0xd8, 0x28, 0xd6, 0x20, 0x28, 0xcf, 0xc8, 0x96, 0x8b, 0x15, 0xa4, 0xa8, 0x48, 0xc8,
	0x84, 0x8b, 0x05, 0xc4, 0x10, 0x12, 0xd5, 0x83, 0xb9, 0x46, 0x0f, 0xc9, 0x29, 0x52, 0xc2, 0xe8,
	0xc2, 0x05, 0x39, 0x95, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0x1f, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x95, 0x61, 0x37, 0x11, 0xd0, 0x00, 0x00, 0x00,
}