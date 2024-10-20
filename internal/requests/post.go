package requests

import (
	"post_service/internal/constants"
	pkg_request "post_service/pkg/requests"
)

type GetPostWithPaginationInterface interface {
	pkg_request.PaginationRequestInterface
	GetQuery() map[string]interface{}
}

type GetPostWithPaginationRequest struct {
	pkg_request.PaginationRequest
	PostQuery map[string]interface{}
}

func (r GetPostWithPaginationRequest) GetQuery() map[string]interface{} {
	return r.PostQuery
}

type GetPostInGroupWithPaginationRequest struct {
	pkg_request.PaginationRequest
	TargetId string `form:"targetId"`
}

func (r GetPostInGroupWithPaginationRequest) GetQuery() map[string]interface{} {
	return map[string]interface{}{
		"target_id": r.TargetId,
		"type":      constants.GroupPost,
	}
}

type GetPostByTagWithPaginationRequest struct {
	pkg_request.PaginationRequest
	Tag string `form:"tag"`
}
type GetPostByMentionWithPaginationRequest struct {
	pkg_request.PaginationRequest
	Mention string `form:"mention"`
}

func (r GetPostByMentionWithPaginationRequest) GetQuery() map[string]interface{} {
	return map[string]interface{}{"mention": r.Mention}
}

type GetPostForUserWithPagination struct {
	pkg_request.PaginationRequest
	UserId string `form:"userId"`
}
type GetPostOfUserInWithPagination struct {
	pkg_request.PaginationRequest
	UserId string `form:"userId"`
}

type GetPostByUserInGroupRequestWithPagination struct {
	pkg_request.PaginationRequest
	OwnerId  string
	TargetId string
}

func (r GetPostByUserInGroupRequestWithPagination) GetQuery() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":  r.OwnerId,
		"target_id": r.TargetId,
		"type":      constants.GroupPost,
	}
}
