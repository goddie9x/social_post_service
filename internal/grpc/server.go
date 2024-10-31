package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"post_service/internal/grpc/converters"
	pb "post_service/internal/grpc/post_service/proto"
	"post_service/internal/repositories"
	"post_service/internal/requests"
	"post_service/internal/services"
	pkg_request "post_service/pkg/requests"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPostServiceServer
	postService repositories.PostRepository
}

func NewServer() *server {
	return &server{
		postService: services.NewPostService(),
	}
}

func (s *server) Create(ctx context.Context, req *pb.PostWithAuthRequest) (*pb.PostResponse, error) {
	response := s.postService.Create(converters.ConvertPostWithAuthRequestFromGRPC(req))
	if response.Ex == nil {
		return &pb.PostResponse{
			Post: converters.ConvertPostToGRPC(*response.Post),
		}, nil
	}
	return nil, response.Ex
}

func (s *server) Update(ctx context.Context, req *pb.PostWithAuthRequest) (*pb.PostResponse, error) {
	response := s.postService.Update(converters.ConvertPostWithAuthRequestFromGRPC(req))
	if response.Ex == nil {
		return &pb.PostResponse{
			Post: converters.ConvertPostToGRPC(*response.Post),
		}, nil
	}
	return nil, response.Ex
}

func (s *server) GetByIdIfUserCanView(ctx context.Context, req *pb.RequestWithAuthAndId) (*pb.PostResponse, error) {
	response := s.postService.GetByIdIfUserCanView(requests.RequestWithAuthAndId{
		User: converters.ConvertUserAuthFromGRPC(req.User),
		Id:   req.Id,
	})

	if response.Ex == nil {
		return &pb.PostResponse{
			Post: converters.ConvertPostToGRPC(*response.Post),
		}, nil
	}
	return nil, response.Ex
}

func (s *server) DeleteById(ctx context.Context, req *pb.RequestWithAuthAndId) (*pb.PostResponse, error) {
	response := s.postService.DeleteById(requests.RequestWithAuthAndId{
		User: converters.ConvertUserAuthFromGRPC(req.User),
		Id:   req.Id,
	})
	if response.Ex == nil {
		return &pb.PostResponse{
			Post: converters.ConvertPostToGRPC(*response.Post),
		}, nil
	}
	return nil, response.Ex
}

func (s *server) GetPostsByTagWithPagination(ctx context.Context, req *pb.GetPostByTagWithPaginationRequest) (*pb.ListPostWithPaginationResponse, error) {
	response := s.postService.GetPostsByTagWithPagination(requests.GetPostByTagWithPaginationRequest{
		PaginationRequest: pkg_request.PaginationRequest{
			Page: int(req.Pagination.Page),
			Size: int(req.Pagination.PageSize),
		},
		Tag:  req.Tag,
		User: converters.ConvertUserAuthFromGRPC(req.User),
	})
	if response.Ex != nil {
		return nil, response.Ex
	}

	return converters.ConvertListPostWithPaginationResponseToGRPCListPostWithPaginationResponse(response), nil
}

func (s *server) GetPostByMentionWithPagination(ctx context.Context, req *pb.GetPostByMentionWithPaginationRequest) (*pb.ListPostWithPaginationResponse, error) {
	response := s.postService.GetPostByMentionWithPagination(requests.GetPostByMentionWithPaginationRequest{
		Mention: req.Mention,
		PaginationRequest: pkg_request.PaginationRequest{
			Page: int(req.Pagination.Page),
			Size: int(req.Pagination.PageSize),
		},
	})
	if response.Ex != nil {
		return nil, response.Ex
	}
	return converters.ConvertListPostWithPaginationResponseToGRPCListPostWithPaginationResponse(response), nil
}

func (s *server) GetPostsForUserProfile(ctx context.Context, req *pb.GetPostForUserWithPagination) (*pb.ListPostWithPaginationResponse, error) {
	response := s.postService.GetPostsForUserProfile(requests.GetPostForUserWithPagination{
		User: converters.ConvertUserAuthFromGRPC(req.User),
		PaginationRequest: pkg_request.PaginationRequest{
			Page: int(req.Pagination.Page),
			Size: int(req.Pagination.PageSize),
		},
		UserId: req.UserId,
	})
	if response.Ex != nil {
		return nil, response.Ex
	}
	return converters.ConvertListPostWithPaginationResponseToGRPCListPostWithPaginationResponse(response), nil
}

func StartGRPCServer(portStr string) {
	address := fmt.Sprintf(":%s", portStr)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", portStr, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterPostServiceServer(grpcServer, NewServer())

	log.Printf("Server is running on port %s", portStr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
