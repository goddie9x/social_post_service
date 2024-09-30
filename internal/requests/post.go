package requests

import (
	"post_service/internal/models"
	pkg_request "post_service/pkg/requests"
)

type GetPostWithPaginationInterface interface {
	pkg_request.PaginationRequestInterface
	GetQuery() interface{}
}

type GetPostWithPaginationRequest struct {
	pagination pkg_request.PaginationRequest
	postQuery  models.Post
}

func (r *GetPostWithPaginationRequest) GetQuery() models.Post {
	return r.postQuery
}
