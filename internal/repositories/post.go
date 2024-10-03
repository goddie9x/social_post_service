package repositories

import (
	"post_service/internal/models"
	"post_service/internal/requests"
)

type PostRepository interface {
	Create(post *models.Post) error
	Update(post *models.Post) error
	GetById(id string) (*models.Post, error)
	GetPostsByTagWithPagination(request requests.GetPostByTagsWithPaginationRequest) (posts []models.Post, amountPage int64, err error)
	GetPostsWithPagination(request requests.GetPostWithPaginationInterface) (posts []models.Post, amountPage int64, err error)
	GetPostsForUserProfile(request requests.GetPostForUserWithPagination) ([]models.Post, int64, error)
	DeleteById(id string) error
}
