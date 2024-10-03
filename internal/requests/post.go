package requests

import (
	pkg_request "post_service/pkg/requests"
)

type GetPostWithPaginationInterface interface {
	pkg_request.PaginationRequestInterface
	GetQuery() interface{}
}

type GetPostWithPaginationRequest struct {
	pkg_request.PaginationRequest
	PostQuery interface{}
}

func (r GetPostWithPaginationRequest) GetQuery() interface{} {
	return r.PostQuery
}

type GetPostInGroupWithPaginationRequest struct {
	pkg_request.PaginationRequest
	TargetId string `form:"targetId"`
}

func (r GetPostInGroupWithPaginationRequest) GetQuery() interface{} {
	return r.TargetId
}

type GetPostByTagsWithPaginationRequest struct {
	pkg_request.PaginationRequest
	Tags []string `form:"tags"`
}
