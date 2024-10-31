package converters

import (
	proto "post_service/internal/grpc/post_service/proto"
	"post_service/internal/models"
	"post_service/internal/requests"
	"post_service/internal/responses"
	"post_service/pkg/constants"
	"post_service/pkg/middlewares"
	pkg_request "post_service/pkg/requests"
	"time"
)

func ConvertPaginationFromGRPCPagination(pagination *proto.PaginationRequest) pkg_request.PaginationRequest {
	return pkg_request.PaginationRequest{
		Page: int(pagination.Page),
		Size: int(pagination.PageSize),
	}
}
func ConvertUserAuthToGRPC(userAuth middlewares.UserAuth) *proto.UserAuth {
	return &proto.UserAuth{
		UserId: userAuth.UserId,
		Role:   int32(userAuth.Role),
	}
}

func ConvertUserAuthFromGRPC(grpcUserAuth *proto.UserAuth) middlewares.UserAuth {
	return middlewares.UserAuth{
		UserId:   grpcUserAuth.UserId,
		UserName: grpcUserAuth.Username,
		Role:     constants.Role(grpcUserAuth.Role),
	}
}

func ConvertPostToGRPC(post models.Post) *proto.Post {
	return &proto.Post{
		Id:        post.Id,
		OwnerId:   post.OwnerId,
		Type:      proto.PostType(*post.Type),
		TargetId:  post.TargetId,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Unix(),
		UpdatedAt: post.UpdatedAt.Unix(),
		Privacy:   proto.PrivacyType(*post.Privacy),
		Approved:  post.Approved,
		BlobIds:   post.BlobIds,
		Tags:      convertTagsToGRPC(post.Tags),
		Mentions:  convertMentionsToGRPC(post.Mentions),
	}
}

func ConvertPostFromGRPC(grpcPost *proto.Post) models.Post {
	postType := constants.PostType(grpcPost.Type)
	privacy := constants.PrivacyType(grpcPost.Privacy)

	return models.Post{
		Id:        grpcPost.Id,
		OwnerId:   grpcPost.OwnerId,
		Type:      &postType,
		TargetId:  grpcPost.TargetId,
		Content:   grpcPost.Content,
		CreatedAt: time.Unix(grpcPost.CreatedAt, 0),
		UpdatedAt: time.Unix(grpcPost.UpdatedAt, 0),
		Privacy:   &privacy,
		Approved:  grpcPost.Approved,
		BlobIds:   grpcPost.BlobIds,
		Tags:      convertTagsFromGRPC(grpcPost.Tags),
		Mentions:  convertMentionsFromGRPC(grpcPost.Mentions),
	}
}

func convertTagsToGRPC(tags []*models.Tag) []*proto.Tag {
	grpcTags := make([]*proto.Tag, len(tags))
	for i, tag := range tags {
		grpcTags[i] = &proto.Tag{Name: tag.Name}
	}
	return grpcTags
}

func convertTagsFromGRPC(grpcTags []*proto.Tag) []*models.Tag {
	tags := make([]*models.Tag, len(grpcTags))
	for i, grpcTag := range grpcTags {
		tags[i] = &models.Tag{Name: grpcTag.Name}
	}
	return tags
}

func convertMentionsToGRPC(mentions []*models.Mention) []*proto.Mention {
	grpcMentions := make([]*proto.Mention, len(mentions))
	for i, mention := range mentions {
		grpcMentions[i] = &proto.Mention{
			PostId:                mention.PostId,
			UserId:                mention.UserId,
			AcceptedShowInProfile: mention.AcceptedShowInProfile,
		}
	}
	return grpcMentions
}

func convertMentionsFromGRPC(grpcMentions []*proto.Mention) []*models.Mention {
	mentions := make([]*models.Mention, len(grpcMentions))
	for i, grpcMention := range grpcMentions {
		mentions[i] = &models.Mention{
			PostId:                grpcMention.PostId,
			UserId:                grpcMention.UserId,
			AcceptedShowInProfile: grpcMention.AcceptedShowInProfile,
		}
	}
	return mentions
}
func ConvertGRPCPostResponseToPostResponse(protoPostResponse *proto.PostResponse) responses.PostResponse {
	post := ConvertPostFromGRPC(protoPostResponse.Post)
	return responses.PostResponse{
		Post: &post,
	}
}

func ConvertPostResponseToGRPCPostResponse(postResponse responses.PostResponse) *proto.PostResponse {
	return &proto.PostResponse{
		Post: ConvertPostToGRPC(*postResponse.Post),
	}
}

func ConvertListPostGRPCToListPost(postsProto []*proto.Post) []models.Post {
	posts := make([]models.Post, len(postsProto))
	for i, protoPost := range postsProto {
		posts[i] = ConvertPostFromGRPC(protoPost)
	}
	return posts
}

func ConvertListPostToListPostGRPC(posts []models.Post) []*proto.Post {
	postsProto := make([]*proto.Post, len(posts))
	for i, post := range posts {
		postsProto[i] = ConvertPostToGRPC(post)
	}
	return postsProto
}

func ConvertListPostWithPaginationGRPCResponseToListPostWithPaginationResponse(res *proto.ListPostWithPaginationResponse) responses.ListPostWithPaginationResponse {
	return responses.ListPostWithPaginationResponse{
		Posts:      ConvertListPostGRPCToListPost(res.Posts),
		AmountPage: res.AmountPage,
	}
}

func ConvertListPostWithPaginationResponseToGRPCListPostWithPaginationResponse(res responses.ListPostWithPaginationResponse) *proto.ListPostWithPaginationResponse {

	return &proto.ListPostWithPaginationResponse{
		Posts:      ConvertListPostToListPostGRPC(res.Posts),
		AmountPage: res.AmountPage,
	}
}

func ConvertPostWithAuthRequestToGRPC(request requests.PostWithAuthRequest) *proto.PostWithAuthRequest {
	return &proto.PostWithAuthRequest{
		User: ConvertUserAuthToGRPC(request.User),
		Post: ConvertPostToGRPC(request.Post),
	}
}
func ConvertPostWithAuthRequestFromGRPC(grpcRequest *proto.PostWithAuthRequest) requests.PostWithAuthRequest {
	return requests.PostWithAuthRequest{
		User: ConvertUserAuthFromGRPC(grpcRequest.User),
		Post: ConvertPostFromGRPC(grpcRequest.Post),
	}
}
