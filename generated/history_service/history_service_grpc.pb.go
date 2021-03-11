// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package history_service

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// HistoryServiceClient is the client API for HistoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HistoryServiceClient interface {
	GetBySymbolUuid(ctx context.Context, in *GetBySymbolUuidRequest, opts ...grpc.CallOption) (*GetBySymbolUuidResponse, error)
}

type historyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHistoryServiceClient(cc grpc.ClientConnInterface) HistoryServiceClient {
	return &historyServiceClient{cc}
}

func (c *historyServiceClient) GetBySymbolUuid(ctx context.Context, in *GetBySymbolUuidRequest, opts ...grpc.CallOption) (*GetBySymbolUuidResponse, error) {
	out := new(GetBySymbolUuidResponse)
	err := c.cc.Invoke(ctx, "/v1.history_service.HistoryService/GetSymbolHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HistoryServiceServer is the server API for HistoryService service.
// All implementations must embed UnimplementedHistoryServiceServer
// for forward compatibility
type HistoryServiceServer interface {
	GetBySymbolUuid(context.Context, *GetBySymbolUuidRequest) (*GetBySymbolUuidResponse, error)
	mustEmbedUnimplementedHistoryServiceServer()
}

// UnimplementedHistoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHistoryServiceServer struct {
}

func (UnimplementedHistoryServiceServer) GetBySymbolUuid(context.Context, *GetBySymbolUuidRequest) (*GetBySymbolUuidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSymbolHistory not implemented")
}
func (UnimplementedHistoryServiceServer) mustEmbedUnimplementedHistoryServiceServer() {}

// UnsafeHistoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HistoryServiceServer will
// result in compilation errors.
type UnsafeHistoryServiceServer interface {
	mustEmbedUnimplementedHistoryServiceServer()
}

func RegisterHistoryServiceServer(s *grpc.Server, srv HistoryServiceServer) {
	s.RegisterService(&_HistoryService_serviceDesc, srv)
}

func _HistoryService_GetBySymbolUuid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBySymbolUuidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryServiceServer).GetBySymbolUuid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.history_service.HistoryService/GetSymbolHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryServiceServer).GetBySymbolUuid(ctx, req.(*GetBySymbolUuidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HistoryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.history_service.HistoryService",
	HandlerType: (*HistoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSymbolHistory",
			Handler:    _HistoryService_GetBySymbolUuid_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "history_service.proto",
}