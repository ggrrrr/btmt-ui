// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: be/svc-auth/authpb/v1/passwd.proto

// Authentication service for users of the systems
// Users in this case are people who interact with the system via UI.

package authpbv1

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
	AuthSvc_UserCreate_FullMethodName       = "/authpb.v1.AuthSvc/UserCreate"
	AuthSvc_UserList_FullMethodName         = "/authpb.v1.AuthSvc/UserList"
	AuthSvc_UserUpdate_FullMethodName       = "/authpb.v1.AuthSvc/UserUpdate"
	AuthSvc_UserChangePasswd_FullMethodName = "/authpb.v1.AuthSvc/UserChangePasswd"
	AuthSvc_LoginPasswd_FullMethodName      = "/authpb.v1.AuthSvc/LoginPasswd"
	AuthSvc_TokenValidate_FullMethodName    = "/authpb.v1.AuthSvc/TokenValidate"
	AuthSvc_TokenRefresh_FullMethodName     = "/authpb.v1.AuthSvc/TokenRefresh"
)

// AuthSvcClient is the client API for AuthSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthSvcClient interface {
	UserCreate(ctx context.Context, in *UserCreateRequest, opts ...grpc.CallOption) (*UserCreateResponse, error)
	UserList(ctx context.Context, in *UserListRequest, opts ...grpc.CallOption) (*UserListResponse, error)
	UserUpdate(ctx context.Context, in *UserUpdateRequest, opts ...grpc.CallOption) (*UserUpdateResponse, error)
	UserChangePasswd(ctx context.Context, in *UserChangePasswdRequest, opts ...grpc.CallOption) (*UserChangePasswdResponse, error)
	LoginPasswd(ctx context.Context, in *LoginPasswdRequest, opts ...grpc.CallOption) (*LoginPasswdResponse, error)
	TokenValidate(ctx context.Context, in *TokenValidateRequest, opts ...grpc.CallOption) (*TokenValidateResponse, error)
	TokenRefresh(ctx context.Context, in *TokenRefreshRequest, opts ...grpc.CallOption) (*TokenRefreshResponse, error)
}

type authSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthSvcClient(cc grpc.ClientConnInterface) AuthSvcClient {
	return &authSvcClient{cc}
}

func (c *authSvcClient) UserCreate(ctx context.Context, in *UserCreateRequest, opts ...grpc.CallOption) (*UserCreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserCreateResponse)
	err := c.cc.Invoke(ctx, AuthSvc_UserCreate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authSvcClient) UserList(ctx context.Context, in *UserListRequest, opts ...grpc.CallOption) (*UserListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserListResponse)
	err := c.cc.Invoke(ctx, AuthSvc_UserList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authSvcClient) UserUpdate(ctx context.Context, in *UserUpdateRequest, opts ...grpc.CallOption) (*UserUpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserUpdateResponse)
	err := c.cc.Invoke(ctx, AuthSvc_UserUpdate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authSvcClient) UserChangePasswd(ctx context.Context, in *UserChangePasswdRequest, opts ...grpc.CallOption) (*UserChangePasswdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserChangePasswdResponse)
	err := c.cc.Invoke(ctx, AuthSvc_UserChangePasswd_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authSvcClient) LoginPasswd(ctx context.Context, in *LoginPasswdRequest, opts ...grpc.CallOption) (*LoginPasswdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginPasswdResponse)
	err := c.cc.Invoke(ctx, AuthSvc_LoginPasswd_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authSvcClient) TokenValidate(ctx context.Context, in *TokenValidateRequest, opts ...grpc.CallOption) (*TokenValidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenValidateResponse)
	err := c.cc.Invoke(ctx, AuthSvc_TokenValidate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authSvcClient) TokenRefresh(ctx context.Context, in *TokenRefreshRequest, opts ...grpc.CallOption) (*TokenRefreshResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenRefreshResponse)
	err := c.cc.Invoke(ctx, AuthSvc_TokenRefresh_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthSvcServer is the server API for AuthSvc service.
// All implementations must embed UnimplementedAuthSvcServer
// for forward compatibility.
type AuthSvcServer interface {
	UserCreate(context.Context, *UserCreateRequest) (*UserCreateResponse, error)
	UserList(context.Context, *UserListRequest) (*UserListResponse, error)
	UserUpdate(context.Context, *UserUpdateRequest) (*UserUpdateResponse, error)
	UserChangePasswd(context.Context, *UserChangePasswdRequest) (*UserChangePasswdResponse, error)
	LoginPasswd(context.Context, *LoginPasswdRequest) (*LoginPasswdResponse, error)
	TokenValidate(context.Context, *TokenValidateRequest) (*TokenValidateResponse, error)
	TokenRefresh(context.Context, *TokenRefreshRequest) (*TokenRefreshResponse, error)
	mustEmbedUnimplementedAuthSvcServer()
}

// UnimplementedAuthSvcServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthSvcServer struct{}

func (UnimplementedAuthSvcServer) UserCreate(context.Context, *UserCreateRequest) (*UserCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCreate not implemented")
}
func (UnimplementedAuthSvcServer) UserList(context.Context, *UserListRequest) (*UserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserList not implemented")
}
func (UnimplementedAuthSvcServer) UserUpdate(context.Context, *UserUpdateRequest) (*UserUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserUpdate not implemented")
}
func (UnimplementedAuthSvcServer) UserChangePasswd(context.Context, *UserChangePasswdRequest) (*UserChangePasswdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserChangePasswd not implemented")
}
func (UnimplementedAuthSvcServer) LoginPasswd(context.Context, *LoginPasswdRequest) (*LoginPasswdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginPasswd not implemented")
}
func (UnimplementedAuthSvcServer) TokenValidate(context.Context, *TokenValidateRequest) (*TokenValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenValidate not implemented")
}
func (UnimplementedAuthSvcServer) TokenRefresh(context.Context, *TokenRefreshRequest) (*TokenRefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenRefresh not implemented")
}
func (UnimplementedAuthSvcServer) mustEmbedUnimplementedAuthSvcServer() {}
func (UnimplementedAuthSvcServer) testEmbeddedByValue()                 {}

// UnsafeAuthSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthSvcServer will
// result in compilation errors.
type UnsafeAuthSvcServer interface {
	mustEmbedUnimplementedAuthSvcServer()
}

func RegisterAuthSvcServer(s grpc.ServiceRegistrar, srv AuthSvcServer) {
	// If the following call pancis, it indicates UnimplementedAuthSvcServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuthSvc_ServiceDesc, srv)
}

func _AuthSvc_UserCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).UserCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_UserCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).UserCreate(ctx, req.(*UserCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthSvc_UserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).UserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_UserList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).UserList(ctx, req.(*UserListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthSvc_UserUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).UserUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_UserUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).UserUpdate(ctx, req.(*UserUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthSvc_UserChangePasswd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserChangePasswdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).UserChangePasswd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_UserChangePasswd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).UserChangePasswd(ctx, req.(*UserChangePasswdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthSvc_LoginPasswd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginPasswdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).LoginPasswd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_LoginPasswd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).LoginPasswd(ctx, req.(*LoginPasswdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthSvc_TokenValidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).TokenValidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_TokenValidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).TokenValidate(ctx, req.(*TokenValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthSvc_TokenRefresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthSvcServer).TokenRefresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthSvc_TokenRefresh_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthSvcServer).TokenRefresh(ctx, req.(*TokenRefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthSvc_ServiceDesc is the grpc.ServiceDesc for AuthSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authpb.v1.AuthSvc",
	HandlerType: (*AuthSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserCreate",
			Handler:    _AuthSvc_UserCreate_Handler,
		},
		{
			MethodName: "UserList",
			Handler:    _AuthSvc_UserList_Handler,
		},
		{
			MethodName: "UserUpdate",
			Handler:    _AuthSvc_UserUpdate_Handler,
		},
		{
			MethodName: "UserChangePasswd",
			Handler:    _AuthSvc_UserChangePasswd_Handler,
		},
		{
			MethodName: "LoginPasswd",
			Handler:    _AuthSvc_LoginPasswd_Handler,
		},
		{
			MethodName: "TokenValidate",
			Handler:    _AuthSvc_TokenValidate_Handler,
		},
		{
			MethodName: "TokenRefresh",
			Handler:    _AuthSvc_TokenRefresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "be/svc-auth/authpb/v1/passwd.proto",
}
