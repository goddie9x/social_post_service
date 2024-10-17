package repositories

import (
	"post_service/internal/models"
	"post_service/internal/requests"
	"post_service/pkg/exceptions"
	"post_service/pkg/middlewares"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(middlewares.UserAuth, *models.Post) (*models.Post, exceptions.CommonExceptionInterface)
	Update(middlewares.UserAuth, *models.Post) (*models.Post, exceptions.CommonExceptionInterface)
	GetByIdIfUserCanView(currentUser middlewares.UserAuth, id string) (*models.Post, exceptions.CommonExceptionInterface)
	DeleteById(currentUser middlewares.UserAuth, id string) exceptions.CommonExceptionInterface
	GetPostsByTagWithPagination(middlewares.UserAuth, requests.GetPostByTagWithPaginationRequest) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface)
	GetPostsWithPagination(requests.GetPostWithPaginationInterface, ...func(*gorm.DB) *gorm.DB) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface)
	GetPostsForUserProfile(middlewares.UserAuth, requests.GetPostForUserWithPagination) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface)
}
