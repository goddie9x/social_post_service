package requests

import (
	"post_service/internal/models"
	"post_service/pkg/constants"
	pkg_constants "post_service/pkg/constants"
	"post_service/pkg/middlewares"
	pkg_request "post_service/pkg/requests"

	"gorm.io/gorm"
)

type PostWithAuthRequest struct {
	User middlewares.UserAuth
	Post models.Post
}

type RequestWithAuthAndId struct {
	User middlewares.UserAuth
	Id   string
}

type GetPostWithPaginationInterface interface {
	pkg_request.PaginationRequestInterface
	GetQuery() map[string]interface{}
	GetAdditionWhereClause() []func(*gorm.DB) *gorm.DB
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

func (r GetPostInGroupWithPaginationRequest) GetAdditionWhereClause() []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{}
}

func (r GetPostInGroupWithPaginationRequest) GetQuery() map[string]interface{} {
	return map[string]interface{}{
		"target_id": r.TargetId,
		"type":      constants.GroupPost,
	}
}

type GetPostByTagWithPaginationRequest struct {
	User middlewares.UserAuth `binding:"-"`
	pkg_request.PaginationRequest
	Tag     string   `form:"tag"`
	PostIds []string `binding:"-"`
}

func (r GetPostByTagWithPaginationRequest) GetQuery() map[string]interface{} {
	return map[string]interface{}{}
}
func (r GetPostByTagWithPaginationRequest) GetAdditionWhereClause() []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		additionQueryPostBaseOnPostIds(r.PostIds),
		additionQueryPostBaseOnUser(r.User),
	}
}
func additionQueryPostBaseOnPostIds(Ids []string) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if len(Ids) == 0 {
			return query.Where("1=0")
		}
		return query.Where(Ids)
	}
}
func additionQueryPostBaseOnUser(currentUser middlewares.UserAuth) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if currentUser.Role == pkg_constants.User {
			query.Where(`"owner_id" = ?`, currentUser.UserId)
		}
		return query
	}
}

type GetPostByMentionWithPaginationRequest struct {
	User middlewares.UserAuth
	pkg_request.PaginationRequest
	Mention string   `form:"mention"`
	PostIds []string `binding:"-"`
}

func (r GetPostByMentionWithPaginationRequest) GetAdditionWhereClause() []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		additionQueryPostBaseOnPostIds(r.PostIds),
	}
}
func (r GetPostByMentionWithPaginationRequest) GetQuery() map[string]interface{} {
	return map[string]interface{}{"mention": r.Mention}
}

type GetPostForUserWithPagination struct {
	User middlewares.UserAuth `binding:"-"`
	pkg_request.PaginationRequest
	UserId string `form:"userId"`
}

func (r GetPostForUserWithPagination) GetAdditionWhereClause() []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{}
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

func (r GetPostByUserInGroupRequestWithPagination) GetAdditionWhereClause() []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{}
}

func (r GetPostByUserInGroupRequestWithPagination) GetQuery() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":  r.OwnerId,
		"target_id": r.TargetId,
		"type":      constants.GroupPost,
	}
}
