// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthorizationClient is the client API for Authorization service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorizationClient interface {
	// User things
	GetById(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*User, error)
	Login(ctx context.Context, in *LoginCredentials, opts ...grpc.CallOption) (*User, error)
	Signup(ctx context.Context, in *SignupCredentials, opts ...grpc.CallOption) (*User, error)
	Update(ctx context.Context, in *UpdateInfo, opts ...grpc.CallOption) (*EmptyMessage, error)
	Delete(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*EmptyMessage, error)
	// Session things
	CreateSession(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error)
	DeleteSession(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error)
	CheckSession(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error)
}

type authorizationClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorizationClient(cc grpc.ClientConnInterface) AuthorizationClient {
	return &authorizationClient{cc}
}

func (c *authorizationClient) GetById(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/GetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) Login(ctx context.Context, in *LoginCredentials, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) Signup(ctx context.Context, in *SignupCredentials, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/Signup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) Update(ctx context.Context, in *UpdateInfo, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) Delete(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) CreateSession(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/CreateSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) DeleteSession(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/DeleteSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationClient) CheckSession(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, "/Authorization.Authorization/CheckSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorizationServer is the server API for Authorization service.
// All implementations must embed UnimplementedAuthorizationServer
// for forward compatibility
type AuthorizationServer interface {
	// User things
	GetById(context.Context, *UserId) (*User, error)
	Login(context.Context, *LoginCredentials) (*User, error)
	Signup(context.Context, *SignupCredentials) (*User, error)
	Update(context.Context, *UpdateInfo) (*EmptyMessage, error)
	Delete(context.Context, *UserId) (*EmptyMessage, error)
	// Session things
	CreateSession(context.Context, *Session) (*Session, error)
	DeleteSession(context.Context, *Session) (*Session, error)
	CheckSession(context.Context, *Session) (*Session, error)
	mustEmbedUnimplementedAuthorizationServer()
}

// UnimplementedAuthorizationServer must be embedded to have forward compatible implementations.
type UnimplementedAuthorizationServer struct {
}

func (UnimplementedAuthorizationServer) GetById(context.Context, *UserId) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedAuthorizationServer) Login(context.Context, *LoginCredentials) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthorizationServer) Signup(context.Context, *SignupCredentials) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}
func (UnimplementedAuthorizationServer) Update(context.Context, *UpdateInfo) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedAuthorizationServer) Delete(context.Context, *UserId) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedAuthorizationServer) CreateSession(context.Context, *Session) (*Session, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (UnimplementedAuthorizationServer) DeleteSession(context.Context, *Session) (*Session, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSession not implemented")
}
func (UnimplementedAuthorizationServer) CheckSession(context.Context, *Session) (*Session, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckSession not implemented")
}
func (UnimplementedAuthorizationServer) mustEmbedUnimplementedAuthorizationServer() {}

// UnsafeAuthorizationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorizationServer will
// result in compilation errors.
type UnsafeAuthorizationServer interface {
	mustEmbedUnimplementedAuthorizationServer()
}

func RegisterAuthorizationServer(s grpc.ServiceRegistrar, srv AuthorizationServer) {
	s.RegisterService(&Authorization_ServiceDesc, srv)
}

func _Authorization_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/GetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).GetById(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).Login(ctx, req.(*LoginCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_Signup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignupCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).Signup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/Signup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).Signup(ctx, req.(*SignupCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).Update(ctx, req.(*UpdateInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).Delete(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Session)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/CreateSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).CreateSession(ctx, req.(*Session))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_DeleteSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Session)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).DeleteSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/DeleteSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).DeleteSession(ctx, req.(*Session))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authorization_CheckSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Session)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).CheckSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authorization.Authorization/CheckSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).CheckSession(ctx, req.(*Session))
	}
	return interceptor(ctx, in, info, handler)
}

// Authorization_ServiceDesc is the grpc.ServiceDesc for Authorization service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authorization_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Authorization.Authorization",
	HandlerType: (*AuthorizationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetById",
			Handler:    _Authorization_GetById_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Authorization_Login_Handler,
		},
		{
			MethodName: "Signup",
			Handler:    _Authorization_Signup_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Authorization_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Authorization_Delete_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _Authorization_CreateSession_Handler,
		},
		{
			MethodName: "DeleteSession",
			Handler:    _Authorization_DeleteSession_Handler,
		},
		{
			MethodName: "CheckSession",
			Handler:    _Authorization_CheckSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/pkg/authorization/delivery/grpc/proto/auth.proto",
}
