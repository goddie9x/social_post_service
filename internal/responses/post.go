package responses

import (
	"post_service/internal/models"
	"post_service/pkg/exceptions"
)

type PostResponse struct {
	Post *models.Post
	Ex   exceptions.CommonExceptionInterface
}

type ListPostWithPaginationResponse struct {
	Posts      []models.Post
	AmountPage int64
	Ex         exceptions.CommonExceptionInterface
}
