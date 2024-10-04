package repositories

import (
	"post_service/internal/models"
	"post_service/internal/requests"
	"post_service/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(c *gin.Context) (*models.Post, exceptions.CommonExceptionInterface)
	Update(c *gin.Context) (*models.Post, exceptions.CommonExceptionInterface)
	GetById(c *gin.Context) (*models.Post, exceptions.CommonExceptionInterface)
	DeleteById(c *gin.Context) exceptions.CommonExceptionInterface
	GetPostsByTagWithPagination(c *gin.Context) ([]models.Post, int64, exceptions.CommonExceptionInterface)
	GetPostsInGroupWithPagination(c *gin.Context) ([]models.Post, int64, exceptions.CommonExceptionInterface)
	GetPostsWithPagination(request requests.GetPostWithPaginationInterface, additionalWhereClause func(*gorm.DB) *gorm.DB) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface)
	GetPostsForUserProfile(c *gin.Context) ([]models.Post, int64, exceptions.CommonExceptionInterface)
}
