// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kv/kvserver/storage_services.proto

package kvserver

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MultiRaftClient is the client API for MultiRaft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MultiRaftClient interface {
	RaftMessageBatch(ctx context.Context, opts ...grpc.CallOption) (MultiRaft_RaftMessageBatchClient, error)
	RaftSnapshot(ctx context.Context, opts ...grpc.CallOption) (MultiRaft_RaftSnapshotClient, error)
}

type multiRaftClient struct {
	cc *grpc.ClientConn
}

func NewMultiRaftClient(cc *grpc.ClientConn) MultiRaftClient {
	return &multiRaftClient{cc}
}

func (c *multiRaftClient) RaftMessageBatch(ctx context.Context, opts ...grpc.CallOption) (MultiRaft_RaftMessageBatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MultiRaft_serviceDesc.Streams[0], "/cockroach.storage.MultiRaft/RaftMessageBatch", opts...)
	if err != nil {
		return nil, err
	}
	x := &multiRaftRaftMessageBatchClient{stream}
	return x, nil
}

type MultiRaft_RaftMessageBatchClient interface {
	Send(*RaftMessageRequestBatch) error
	Recv() (*RaftMessageResponse, error)
	grpc.ClientStream
}

type multiRaftRaftMessageBatchClient struct {
	grpc.ClientStream
}

func (x *multiRaftRaftMessageBatchClient) Send(m *RaftMessageRequestBatch) error {
	return x.ClientStream.SendMsg(m)
}

func (x *multiRaftRaftMessageBatchClient) Recv() (*RaftMessageResponse, error) {
	m := new(RaftMessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *multiRaftClient) RaftSnapshot(ctx context.Context, opts ...grpc.CallOption) (MultiRaft_RaftSnapshotClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MultiRaft_serviceDesc.Streams[1], "/cockroach.storage.MultiRaft/RaftSnapshot", opts...)
	if err != nil {
		return nil, err
	}
	x := &multiRaftRaftSnapshotClient{stream}
	return x, nil
}

type MultiRaft_RaftSnapshotClient interface {
	Send(*SnapshotRequest) error
	Recv() (*SnapshotResponse, error)
	grpc.ClientStream
}

type multiRaftRaftSnapshotClient struct {
	grpc.ClientStream
}

func (x *multiRaftRaftSnapshotClient) Send(m *SnapshotRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *multiRaftRaftSnapshotClient) Recv() (*SnapshotResponse, error) {
	m := new(SnapshotResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MultiRaftServer is the server API for MultiRaft service.
type MultiRaftServer interface {
	RaftMessageBatch(MultiRaft_RaftMessageBatchServer) error
	RaftSnapshot(MultiRaft_RaftSnapshotServer) error
}

func RegisterMultiRaftServer(s *grpc.Server, srv MultiRaftServer) {
	s.RegisterService(&_MultiRaft_serviceDesc, srv)
}

func _MultiRaft_RaftMessageBatch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MultiRaftServer).RaftMessageBatch(&multiRaftRaftMessageBatchServer{stream})
}

type MultiRaft_RaftMessageBatchServer interface {
	Send(*RaftMessageResponse) error
	Recv() (*RaftMessageRequestBatch, error)
	grpc.ServerStream
}

type multiRaftRaftMessageBatchServer struct {
	grpc.ServerStream
}

func (x *multiRaftRaftMessageBatchServer) Send(m *RaftMessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *multiRaftRaftMessageBatchServer) Recv() (*RaftMessageRequestBatch, error) {
	m := new(RaftMessageRequestBatch)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MultiRaft_RaftSnapshot_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MultiRaftServer).RaftSnapshot(&multiRaftRaftSnapshotServer{stream})
}

type MultiRaft_RaftSnapshotServer interface {
	Send(*SnapshotResponse) error
	Recv() (*SnapshotRequest, error)
	grpc.ServerStream
}

type multiRaftRaftSnapshotServer struct {
	grpc.ServerStream
}

func (x *multiRaftRaftSnapshotServer) Send(m *SnapshotResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *multiRaftRaftSnapshotServer) Recv() (*SnapshotRequest, error) {
	m := new(SnapshotRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MultiRaft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cockroach.storage.MultiRaft",
	HandlerType: (*MultiRaftServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RaftMessageBatch",
			Handler:       _MultiRaft_RaftMessageBatch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "RaftSnapshot",
			Handler:       _MultiRaft_RaftSnapshot_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "kv/kvserver/storage_services.proto",
}

// PerReplicaClient is the client API for PerReplica service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PerReplicaClient interface {
	CollectChecksum(ctx context.Context, in *CollectChecksumRequest, opts ...grpc.CallOption) (*CollectChecksumResponse, error)
	WaitForApplication(ctx context.Context, in *WaitForApplicationRequest, opts ...grpc.CallOption) (*WaitForApplicationResponse, error)
	WaitForReplicaInit(ctx context.Context, in *WaitForReplicaInitRequest, opts ...grpc.CallOption) (*WaitForReplicaInitResponse, error)
}

type perReplicaClient struct {
	cc *grpc.ClientConn
}

func NewPerReplicaClient(cc *grpc.ClientConn) PerReplicaClient {
	return &perReplicaClient{cc}
}

func (c *perReplicaClient) CollectChecksum(ctx context.Context, in *CollectChecksumRequest, opts ...grpc.CallOption) (*CollectChecksumResponse, error) {
	out := new(CollectChecksumResponse)
	err := c.cc.Invoke(ctx, "/cockroach.storage.PerReplica/CollectChecksum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *perReplicaClient) WaitForApplication(ctx context.Context, in *WaitForApplicationRequest, opts ...grpc.CallOption) (*WaitForApplicationResponse, error) {
	out := new(WaitForApplicationResponse)
	err := c.cc.Invoke(ctx, "/cockroach.storage.PerReplica/WaitForApplication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *perReplicaClient) WaitForReplicaInit(ctx context.Context, in *WaitForReplicaInitRequest, opts ...grpc.CallOption) (*WaitForReplicaInitResponse, error) {
	out := new(WaitForReplicaInitResponse)
	err := c.cc.Invoke(ctx, "/cockroach.storage.PerReplica/WaitForReplicaInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PerReplicaServer is the server API for PerReplica service.
type PerReplicaServer interface {
	CollectChecksum(context.Context, *CollectChecksumRequest) (*CollectChecksumResponse, error)
	WaitForApplication(context.Context, *WaitForApplicationRequest) (*WaitForApplicationResponse, error)
	WaitForReplicaInit(context.Context, *WaitForReplicaInitRequest) (*WaitForReplicaInitResponse, error)
}

func RegisterPerReplicaServer(s *grpc.Server, srv PerReplicaServer) {
	s.RegisterService(&_PerReplica_serviceDesc, srv)
}

func _PerReplica_CollectChecksum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectChecksumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerReplicaServer).CollectChecksum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cockroach.storage.PerReplica/CollectChecksum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerReplicaServer).CollectChecksum(ctx, req.(*CollectChecksumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerReplica_WaitForApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WaitForApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerReplicaServer).WaitForApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cockroach.storage.PerReplica/WaitForApplication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerReplicaServer).WaitForApplication(ctx, req.(*WaitForApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerReplica_WaitForReplicaInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WaitForReplicaInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerReplicaServer).WaitForReplicaInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cockroach.storage.PerReplica/WaitForReplicaInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerReplicaServer).WaitForReplicaInit(ctx, req.(*WaitForReplicaInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PerReplica_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cockroach.storage.PerReplica",
	HandlerType: (*PerReplicaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CollectChecksum",
			Handler:    _PerReplica_CollectChecksum_Handler,
		},
		{
			MethodName: "WaitForApplication",
			Handler:    _PerReplica_WaitForApplication_Handler,
		},
		{
			MethodName: "WaitForReplicaInit",
			Handler:    _PerReplica_WaitForReplicaInit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv/kvserver/storage_services.proto",
}

func init() {
	proto.RegisterFile("kv/kvserver/storage_services.proto", fileDescriptor_storage_services_26b72250ee725fca)
}

var fileDescriptor_storage_services_26b72250ee725fca = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xbf, 0x4e, 0xeb, 0x30,
	0x14, 0x87, 0xe3, 0x3b, 0x5c, 0x81, 0x85, 0x04, 0x58, 0xc0, 0x90, 0xc1, 0x43, 0x07, 0x40, 0x95,
	0x70, 0x4b, 0x79, 0x02, 0x5a, 0x09, 0x89, 0xa1, 0x12, 0x2a, 0x03, 0x12, 0x0b, 0x32, 0xd6, 0x69,
	0x62, 0x25, 0xc4, 0xc1, 0x76, 0xb3, 0xf0, 0x12, 0x3c, 0x56, 0xc7, 0x8e, 0x1d, 0x21, 0x5d, 0x78,
	0x06, 0x26, 0x94, 0x34, 0x16, 0xe1, 0x4f, 0xd5, 0x30, 0xe5, 0xe4, 0xf8, 0x3b, 0xbf, 0xef, 0xc8,
	0x32, 0x6e, 0x45, 0x59, 0x27, 0xca, 0x0c, 0xe8, 0x0c, 0x74, 0xc7, 0x58, 0xa5, 0x79, 0x00, 0x77,
	0xc5, 0xaf, 0x14, 0x60, 0x58, 0xaa, 0x95, 0x55, 0x64, 0x57, 0x28, 0x11, 0x69, 0xc5, 0x45, 0xc8,
	0x2a, 0xc2, 0x3f, 0xa8, 0x8f, 0x69, 0x3e, 0xb6, 0x4b, 0xd4, 0xdf, 0xaf, 0xf7, 0x79, 0x2a, 0xab,
	0xf6, 0x5e, 0xa0, 0x02, 0x55, 0x96, 0x9d, 0xa2, 0x5a, 0x76, 0x7b, 0x6f, 0x08, 0x6f, 0x0e, 0x27,
	0xb1, 0x95, 0x23, 0x3e, 0xb6, 0xc4, 0xe2, 0x9d, 0xe2, 0x3b, 0x04, 0x63, 0x78, 0x00, 0x7d, 0x6e,
	0x45, 0x48, 0x18, 0xfb, 0x54, 0x47, 0x19, 0x73, 0xc9, 0xac, 0x06, 0x8e, 0xe0, 0x71, 0x02, 0xc6,
	0x96, 0xbc, 0xdf, 0x6e, 0xc2, 0x9b, 0x54, 0x25, 0x06, 0x5a, 0xde, 0x31, 0xea, 0x22, 0x02, 0x78,
	0xab, 0x38, 0xbc, 0x4e, 0x78, 0x6a, 0x42, 0x65, 0xc9, 0xe1, 0x8a, 0x04, 0x07, 0x54, 0x3a, 0xff,
	0x68, 0x2d, 0x57, 0xd7, 0xf4, 0xde, 0xff, 0x61, 0x7c, 0x05, 0x7a, 0x04, 0x69, 0x2c, 0x05, 0x27,
	0x1a, 0x6f, 0x0f, 0x54, 0x1c, 0x83, 0xb0, 0x83, 0x10, 0x44, 0x64, 0x26, 0x0f, 0xe4, 0x64, 0x45,
	0xe0, 0x37, 0xce, 0xf9, 0x59, 0x53, 0xdc, 0xad, 0x41, 0x9e, 0x30, 0xb9, 0xe1, 0xd2, 0x5e, 0x28,
	0x7d, 0x9e, 0x96, 0x5b, 0x58, 0xa9, 0x12, 0xd2, 0x5d, 0x91, 0xf3, 0x13, 0x75, 0xe6, 0xd3, 0x3f,
	0x4c, 0xfc, 0x22, 0xaf, 0xae, 0xe0, 0x32, 0x91, 0x76, 0x9d, 0xbc, 0x86, 0x36, 0x94, 0x7f, 0x99,
	0x70, 0xf2, 0x7e, 0x7b, 0xfa, 0x4a, 0xbd, 0x69, 0x4e, 0xd1, 0x2c, 0xa7, 0x68, 0x9e, 0x53, 0xf4,
	0x92, 0x53, 0xf4, 0xbc, 0xa0, 0xde, 0x6c, 0x41, 0xbd, 0xf9, 0x82, 0x7a, 0xb7, 0x1b, 0x2e, 0xe8,
	0xfe, 0x7f, 0xf9, 0x34, 0xcf, 0x3e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x83, 0x9e, 0xc7, 0x08, 0x18,
	0x03, 0x00, 0x00,
}
