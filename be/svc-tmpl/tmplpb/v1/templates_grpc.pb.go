// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: be/svc-tmpl/tmplpb/v1/templates.proto

package tmplpbv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TmplSvc_Save_FullMethodName    = "/tmplpb.v1.TmplSvc/Save"
	TmplSvc_GetById_FullMethodName = "/tmplpb.v1.TmplSvc/GetById"
	TmplSvc_Search_FullMethodName  = "/tmplpb.v1.TmplSvc/Search"
	TmplSvc_Render_FullMethodName  = "/tmplpb.v1.TmplSvc/Render"
)

// TmplSvcClient is the client API for TmplSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TmplSvcClient interface {
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error)
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	Render(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error)
}

type tmplSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewTmplSvcClient(cc grpc.ClientConnInterface) TmplSvcClient {
	return &tmplSvcClient{cc}
}

func (c *tmplSvcClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveResponse)
	err := c.cc.Invoke(ctx, TmplSvc_Save_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tmplSvcClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetByIdResponse)
	err := c.cc.Invoke(ctx, TmplSvc_GetById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tmplSvcClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, TmplSvc_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tmplSvcClient) Render(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RenderResponse)
	err := c.cc.Invoke(ctx, TmplSvc_Render_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TmplSvcServer is the server API for TmplSvc service.
// All implementations must embed UnimplementedTmplSvcServer
// for forward compatibility.
type TmplSvcServer interface {
	Save(context.Context, *SaveRequest) (*SaveResponse, error)
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	Render(context.Context, *RenderRequest) (*RenderResponse, error)
	mustEmbedUnimplementedTmplSvcServer()
}

// UnimplementedTmplSvcServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTmplSvcServer struct{}

func (UnimplementedTmplSvcServer) Save(context.Context, *SaveRequest) (*SaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedTmplSvcServer) GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedTmplSvcServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedTmplSvcServer) Render(context.Context, *RenderRequest) (*RenderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Render not implemented")
}
func (UnimplementedTmplSvcServer) mustEmbedUnimplementedTmplSvcServer() {}
func (UnimplementedTmplSvcServer) testEmbeddedByValue()                 {}

// UnsafeTmplSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TmplSvcServer will
// result in compilation errors.
type UnsafeTmplSvcServer interface {
	mustEmbedUnimplementedTmplSvcServer()
}

func RegisterTmplSvcServer(s grpc.ServiceRegistrar, srv TmplSvcServer) {
	// If the following call pancis, it indicates UnimplementedTmplSvcServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TmplSvc_ServiceDesc, srv)
}

func _TmplSvc_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TmplSvcServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TmplSvc_Save_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TmplSvcServer).Save(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TmplSvc_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TmplSvcServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TmplSvc_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TmplSvcServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TmplSvc_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TmplSvcServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TmplSvc_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TmplSvcServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TmplSvc_Render_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TmplSvcServer).Render(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TmplSvc_Render_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TmplSvcServer).Render(ctx, req.(*RenderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TmplSvc_ServiceDesc is the grpc.ServiceDesc for TmplSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TmplSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tmplpb.v1.TmplSvc",
	HandlerType: (*TmplSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _TmplSvc_Save_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _TmplSvc_GetById_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _TmplSvc_Search_Handler,
		},
		{
			MethodName: "Render",
			Handler:    _TmplSvc_Render_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "be/svc-tmpl/tmplpb/v1/templates.proto",
}
