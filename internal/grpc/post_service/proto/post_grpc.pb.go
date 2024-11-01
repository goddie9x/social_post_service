// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0--rc2
// source: pkg/grpc/post.proto

package post_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	PostService_Create_FullMethodName                         = "/post_service.PostService/Create"
	PostService_Update_FullMethodName                         = "/post_service.PostService/Update"
	PostService_GetByIdIfUserCanView_FullMethodName           = "/post_service.PostService/GetByIdIfUserCanView"
	PostService_DeleteById_FullMethodName                     = "/post_service.PostService/DeleteById"
	PostService_GetPostsByTagWithPagination_FullMethodName    = "/post_service.PostService/GetPostsByTagWithPagination"
	PostService_GetPostsWithPagination_FullMethodName         = "/post_service.PostService/GetPostsWithPagination"
	PostService_GetPostByMentionWithPagination_FullMethodName = "/post_service.PostService/GetPostByMentionWithPagination"
	PostService_GetPostsForUserProfile_FullMethodName         = "/post_service.PostService/GetPostsForUserProfile"
)

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	Create(ctx context.Context, in *PostWithAuthRequest, opts ...grpc.CallOption) (*PostResponse, error)
	Update(ctx context.Context, in *PostWithAuthRequest, opts ...grpc.CallOption) (*PostResponse, error)
	GetByIdIfUserCanView(ctx context.Context, in *RequestWithAuthAndId, opts ...grpc.CallOption) (*PostResponse, error)
	DeleteById(ctx context.Context, in *RequestWithAuthAndId, opts ...grpc.CallOption) (*PostResponse, error)
	GetPostsByTagWithPagination(ctx context.Context, in *GetPostByTagWithPaginationRequest, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error)
	GetPostsWithPagination(ctx context.Context, in *PaginationRequest, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error)
	GetPostByMentionWithPagination(ctx context.Context, in *GetPostByMentionWithPaginationRequest, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error)
	GetPostsForUserProfile(ctx context.Context, in *GetPostForUserWithPagination, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) Create(ctx context.Context, in *PostWithAuthRequest, opts ...grpc.CallOption) (*PostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostResponse)
	err := c.cc.Invoke(ctx, PostService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) Update(ctx context.Context, in *PostWithAuthRequest, opts ...grpc.CallOption) (*PostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostResponse)
	err := c.cc.Invoke(ctx, PostService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetByIdIfUserCanView(ctx context.Context, in *RequestWithAuthAndId, opts ...grpc.CallOption) (*PostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostResponse)
	err := c.cc.Invoke(ctx, PostService_GetByIdIfUserCanView_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) DeleteById(ctx context.Context, in *RequestWithAuthAndId, opts ...grpc.CallOption) (*PostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostResponse)
	err := c.cc.Invoke(ctx, PostService_DeleteById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetPostsByTagWithPagination(ctx context.Context, in *GetPostByTagWithPaginationRequest, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPostWithPaginationResponse)
	err := c.cc.Invoke(ctx, PostService_GetPostsByTagWithPagination_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetPostsWithPagination(ctx context.Context, in *PaginationRequest, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPostWithPaginationResponse)
	err := c.cc.Invoke(ctx, PostService_GetPostsWithPagination_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetPostByMentionWithPagination(ctx context.Context, in *GetPostByMentionWithPaginationRequest, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPostWithPaginationResponse)
	err := c.cc.Invoke(ctx, PostService_GetPostByMentionWithPagination_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetPostsForUserProfile(ctx context.Context, in *GetPostForUserWithPagination, opts ...grpc.CallOption) (*ListPostWithPaginationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPostWithPaginationResponse)
	err := c.cc.Invoke(ctx, PostService_GetPostsForUserProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility.
type PostServiceServer interface {
	Create(context.Context, *PostWithAuthRequest) (*PostResponse, error)
	Update(context.Context, *PostWithAuthRequest) (*PostResponse, error)
	GetByIdIfUserCanView(context.Context, *RequestWithAuthAndId) (*PostResponse, error)
	DeleteById(context.Context, *RequestWithAuthAndId) (*PostResponse, error)
	GetPostsByTagWithPagination(context.Context, *GetPostByTagWithPaginationRequest) (*ListPostWithPaginationResponse, error)
	GetPostsWithPagination(context.Context, *PaginationRequest) (*ListPostWithPaginationResponse, error)
	GetPostByMentionWithPagination(context.Context, *GetPostByMentionWithPaginationRequest) (*ListPostWithPaginationResponse, error)
	GetPostsForUserProfile(context.Context, *GetPostForUserWithPagination) (*ListPostWithPaginationResponse, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPostServiceServer struct{}

func (UnimplementedPostServiceServer) Create(context.Context, *PostWithAuthRequest) (*PostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPostServiceServer) Update(context.Context, *PostWithAuthRequest) (*PostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedPostServiceServer) GetByIdIfUserCanView(context.Context, *RequestWithAuthAndId) (*PostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByIdIfUserCanView not implemented")
}
func (UnimplementedPostServiceServer) DeleteById(context.Context, *RequestWithAuthAndId) (*PostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteById not implemented")
}
func (UnimplementedPostServiceServer) GetPostsByTagWithPagination(context.Context, *GetPostByTagWithPaginationRequest) (*ListPostWithPaginationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostsByTagWithPagination not implemented")
}
func (UnimplementedPostServiceServer) GetPostsWithPagination(context.Context, *PaginationRequest) (*ListPostWithPaginationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostsWithPagination not implemented")
}
func (UnimplementedPostServiceServer) GetPostByMentionWithPagination(context.Context, *GetPostByMentionWithPaginationRequest) (*ListPostWithPaginationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostByMentionWithPagination not implemented")
}
func (UnimplementedPostServiceServer) GetPostsForUserProfile(context.Context, *GetPostForUserWithPagination) (*ListPostWithPaginationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostsForUserProfile not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}
func (UnimplementedPostServiceServer) testEmbeddedByValue()                     {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	// If the following call pancis, it indicates UnimplementedPostServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostWithAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Create(ctx, req.(*PostWithAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostWithAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Update(ctx, req.(*PostWithAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetByIdIfUserCanView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithAuthAndId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetByIdIfUserCanView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetByIdIfUserCanView_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetByIdIfUserCanView(ctx, req.(*RequestWithAuthAndId))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_DeleteById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithAuthAndId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).DeleteById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_DeleteById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).DeleteById(ctx, req.(*RequestWithAuthAndId))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetPostsByTagWithPagination_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostByTagWithPaginationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPostsByTagWithPagination(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetPostsByTagWithPagination_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPostsByTagWithPagination(ctx, req.(*GetPostByTagWithPaginationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetPostsWithPagination_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaginationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPostsWithPagination(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetPostsWithPagination_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPostsWithPagination(ctx, req.(*PaginationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetPostByMentionWithPagination_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostByMentionWithPaginationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPostByMentionWithPagination(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetPostByMentionWithPagination_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPostByMentionWithPagination(ctx, req.(*GetPostByMentionWithPaginationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetPostsForUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostForUserWithPagination)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPostsForUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetPostsForUserProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPostsForUserProfile(ctx, req.(*GetPostForUserWithPagination))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "post_service.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _PostService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _PostService_Update_Handler,
		},
		{
			MethodName: "GetByIdIfUserCanView",
			Handler:    _PostService_GetByIdIfUserCanView_Handler,
		},
		{
			MethodName: "DeleteById",
			Handler:    _PostService_DeleteById_Handler,
		},
		{
			MethodName: "GetPostsByTagWithPagination",
			Handler:    _PostService_GetPostsByTagWithPagination_Handler,
		},
		{
			MethodName: "GetPostsWithPagination",
			Handler:    _PostService_GetPostsWithPagination_Handler,
		},
		{
			MethodName: "GetPostByMentionWithPagination",
			Handler:    _PostService_GetPostByMentionWithPagination_Handler,
		},
		{
			MethodName: "GetPostsForUserProfile",
			Handler:    _PostService_GetPostsForUserProfile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/grpc/post.proto",
}