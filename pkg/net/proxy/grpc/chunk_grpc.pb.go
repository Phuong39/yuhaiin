// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: chunk.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StreamClient is the client API for Stream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamClient interface {
	Conn(ctx context.Context, opts ...grpc.CallOption) (Stream_ConnClient, error)
}

type streamClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamClient(cc grpc.ClientConnInterface) StreamClient {
	return &streamClient{cc}
}

func (c *streamClient) Conn(ctx context.Context, opts ...grpc.CallOption) (Stream_ConnClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[0], "/stream.stream/conn", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamConnClient{stream}
	return x, nil
}

type Stream_ConnClient interface {
	Send(*wrapperspb.BytesValue) error
	Recv() (*wrapperspb.BytesValue, error)
	grpc.ClientStream
}

type streamConnClient struct {
	grpc.ClientStream
}

func (x *streamConnClient) Send(m *wrapperspb.BytesValue) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamConnClient) Recv() (*wrapperspb.BytesValue, error) {
	m := new(wrapperspb.BytesValue)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamServer is the server API for Stream service.
// All implementations must embed UnimplementedStreamServer
// for forward compatibility
type StreamServer interface {
	Conn(Stream_ConnServer) error
	mustEmbedUnimplementedStreamServer()
}

// UnimplementedStreamServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServer struct {
}

func (UnimplementedStreamServer) Conn(Stream_ConnServer) error {
	return status.Errorf(codes.Unimplemented, "method Conn not implemented")
}
func (UnimplementedStreamServer) mustEmbedUnimplementedStreamServer() {}

// UnsafeStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServer will
// result in compilation errors.
type UnsafeStreamServer interface {
	mustEmbedUnimplementedStreamServer()
}

func RegisterStreamServer(s grpc.ServiceRegistrar, srv StreamServer) {
	s.RegisterService(&Stream_ServiceDesc, srv)
}

func _Stream_Conn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServer).Conn(&streamConnServer{stream})
}

type Stream_ConnServer interface {
	Send(*wrapperspb.BytesValue) error
	Recv() (*wrapperspb.BytesValue, error)
	grpc.ServerStream
}

type streamConnServer struct {
	grpc.ServerStream
}

func (x *streamConnServer) Send(m *wrapperspb.BytesValue) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamConnServer) Recv() (*wrapperspb.BytesValue, error) {
	m := new(wrapperspb.BytesValue)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Stream_ServiceDesc is the grpc.ServiceDesc for Stream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stream.stream",
	HandlerType: (*StreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "conn",
			Handler:       _Stream_Conn_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chunk.proto",
}
