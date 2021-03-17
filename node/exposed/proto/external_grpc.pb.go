// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package exposed

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

// ExternalListenerClient is the client API for ExternalListener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExternalListenerClient interface {
	Upload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*IPResponse, error)
	Download(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	CheckIP(ctx context.Context, in *Request, opts ...grpc.CallOption) (*IPResponse, error)
}

type externalListenerClient struct {
	cc grpc.ClientConnInterface
}

func NewExternalListenerClient(cc grpc.ClientConnInterface) ExternalListenerClient {
	return &externalListenerClient{cc}
}

func (c *externalListenerClient) Upload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*IPResponse, error) {
	out := new(IPResponse)
	err := c.cc.Invoke(ctx, "/grpc.ExternalListener/Upload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalListenerClient) Download(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.ExternalListener/Download", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalListenerClient) CheckIP(ctx context.Context, in *Request, opts ...grpc.CallOption) (*IPResponse, error) {
	out := new(IPResponse)
	err := c.cc.Invoke(ctx, "/grpc.ExternalListener/CheckIP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExternalListenerServer is the server API for ExternalListener service.
// All implementations must embed UnimplementedExternalListenerServer
// for forward compatibility
type ExternalListenerServer interface {
	Upload(context.Context, *UploadRequest) (*IPResponse, error)
	Download(context.Context, *Request) (*Response, error)
	CheckIP(context.Context, *Request) (*IPResponse, error)
	mustEmbedUnimplementedExternalListenerServer()
}

// UnimplementedExternalListenerServer must be embedded to have forward compatible implementations.
type UnimplementedExternalListenerServer struct {
}

func (UnimplementedExternalListenerServer) Upload(context.Context, *UploadRequest) (*IPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedExternalListenerServer) Download(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedExternalListenerServer) CheckIP(context.Context, *Request) (*IPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIP not implemented")
}
func (UnimplementedExternalListenerServer) mustEmbedUnimplementedExternalListenerServer() {}

// UnsafeExternalListenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExternalListenerServer will
// result in compilation errors.
type UnsafeExternalListenerServer interface {
	mustEmbedUnimplementedExternalListenerServer()
}

func RegisterExternalListenerServer(s grpc.ServiceRegistrar, srv ExternalListenerServer) {
	s.RegisterService(&ExternalListener_ServiceDesc, srv)
}

func _ExternalListener_Upload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalListenerServer).Upload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalListener/Upload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalListenerServer).Upload(ctx, req.(*UploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalListener_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalListenerServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalListener/Download",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalListenerServer).Download(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalListener_CheckIP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalListenerServer).CheckIP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalListener/CheckIP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalListenerServer).CheckIP(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// ExternalListener_ServiceDesc is the grpc.ServiceDesc for ExternalListener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExternalListener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.ExternalListener",
	HandlerType: (*ExternalListenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upload",
			Handler:    _ExternalListener_Upload_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _ExternalListener_Download_Handler,
		},
		{
			MethodName: "CheckIP",
			Handler:    _ExternalListener_CheckIP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/external.proto",
}
