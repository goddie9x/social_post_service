package repositories

import (
	"post_service/internal/models"
	"post_service/internal/requests"
)

type PostRepository interface {
	Create(post *models.Post) (int64, error)
	Update(post *models.Post) error
	GetById(id string) (*models.Post, error)
	GetAllWithPagination(request requests.GetPostWithPaginationInterface) (posts []models.Post, amount int64, err error)
	DeleteById(id string) error
}
