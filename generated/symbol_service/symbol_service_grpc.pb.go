// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package symbol_service

import (
	context "context"
	proto_models "github.com/vectorman1/analysis/analysis-api/generated/proto_models"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SymbolServiceClient is the client API for SymbolService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SymbolServiceClient interface {
	GetPaged(ctx context.Context, in *GetPagedRequest, opts ...grpc.CallOption) (*GetPagedResponse, error)
	Overview(ctx context.Context, in *SymbolOverviewRequest, opts ...grpc.CallOption) (*SymbolOverview, error)
	Get(ctx context.Context, in *SymbolRequest, opts ...grpc.CallOption) (*proto_models.Symbol, error)
	StartUpdateJob(ctx context.Context, in *StartUpdateJobRequest, opts ...grpc.CallOption) (*StartUpdateJobResponse, error)
}

type symbolServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSymbolServiceClient(cc grpc.ClientConnInterface) SymbolServiceClient {
	return &symbolServiceClient{cc}
}

func (c *symbolServiceClient) GetPaged(ctx context.Context, in *GetPagedRequest, opts ...grpc.CallOption) (*GetPagedResponse, error) {
	out := new(GetPagedResponse)
	err := c.cc.Invoke(ctx, "/v1.symbol_service.SymbolService/GetPaged", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *symbolServiceClient) Overview(ctx context.Context, in *SymbolOverviewRequest, opts ...grpc.CallOption) (*SymbolOverview, error) {
	out := new(SymbolOverview)
	err := c.cc.Invoke(ctx, "/v1.symbol_service.SymbolService/Overview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *symbolServiceClient) Get(ctx context.Context, in *SymbolRequest, opts ...grpc.CallOption) (*proto_models.Symbol, error) {
	out := new(proto_models.Symbol)
	err := c.cc.Invoke(ctx, "/v1.symbol_service.SymbolService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *symbolServiceClient) StartUpdateJob(ctx context.Context, in *StartUpdateJobRequest, opts ...grpc.CallOption) (*StartUpdateJobResponse, error) {
	out := new(StartUpdateJobResponse)
	err := c.cc.Invoke(ctx, "/v1.symbol_service.SymbolService/StartUpdateJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SymbolServiceServer is the server API for SymbolService service.
// All implementations must embed UnimplementedSymbolServiceServer
// for forward compatibility
type SymbolServiceServer interface {
	GetPaged(context.Context, *GetPagedRequest) (*GetPagedResponse, error)
	Overview(context.Context, *SymbolOverviewRequest) (*SymbolOverview, error)
	Get(context.Context, *SymbolRequest) (*proto_models.Symbol, error)
	StartUpdateJob(context.Context, *StartUpdateJobRequest) (*StartUpdateJobResponse, error)
	mustEmbedUnimplementedSymbolServiceServer()
}

// UnimplementedSymbolServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSymbolServiceServer struct {
}

func (UnimplementedSymbolServiceServer) GetPaged(context.Context, *GetPagedRequest) (*GetPagedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaged not implemented")
}
func (UnimplementedSymbolServiceServer) Overview(context.Context, *SymbolOverviewRequest) (*SymbolOverview, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Overview not implemented")
}
func (UnimplementedSymbolServiceServer) Get(context.Context, *SymbolRequest) (*proto_models.Symbol, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSymbolServiceServer) StartUpdateJob(context.Context, *StartUpdateJobRequest) (*StartUpdateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartUpdateJob not implemented")
}
func (UnimplementedSymbolServiceServer) mustEmbedUnimplementedSymbolServiceServer() {}

// UnsafeSymbolServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SymbolServiceServer will
// result in compilation errors.
type UnsafeSymbolServiceServer interface {
	mustEmbedUnimplementedSymbolServiceServer()
}

func RegisterSymbolServiceServer(s *grpc.Server, srv SymbolServiceServer) {
	s.RegisterService(&_SymbolService_serviceDesc, srv)
}

func _SymbolService_GetPaged_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPagedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SymbolServiceServer).GetPaged(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.symbol_service.SymbolService/GetPaged",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SymbolServiceServer).GetPaged(ctx, req.(*GetPagedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SymbolService_Overview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SymbolOverviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SymbolServiceServer).Overview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.symbol_service.SymbolService/Overview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SymbolServiceServer).Overview(ctx, req.(*SymbolOverviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SymbolService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SymbolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SymbolServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.symbol_service.SymbolService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SymbolServiceServer).Get(ctx, req.(*SymbolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SymbolService_StartUpdateJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartUpdateJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SymbolServiceServer).StartUpdateJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.symbol_service.SymbolService/StartUpdateJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SymbolServiceServer).StartUpdateJob(ctx, req.(*StartUpdateJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SymbolService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.symbol_service.SymbolService",
	HandlerType: (*SymbolServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPaged",
			Handler:    _SymbolService_GetPaged_Handler,
		},
		{
			MethodName: "Overview",
			Handler:    _SymbolService_Overview_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SymbolService_Get_Handler,
		},
		{
			MethodName: "StartUpdateJob",
			Handler:    _SymbolService_StartUpdateJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "symbol_service.proto",
}
