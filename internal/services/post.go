package services

import (
	"log"
	"os"
	"post_service/internal/models"
	"post_service/internal/repositories"
	"post_service/internal/requests"
	"strconv"

	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

const (
	invalid_amount = -1
)

type PostService struct {
	db *gorm.DB
}

func NewPostService() repositories.PostRepository {
	options := map[string]string{
		"CONNECTION TIMEOUT": "90",
		"SSL":                "false",
	}
	db_port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error converting APP_PORT to integer: %v", err)
	}
	url := oracle.BuildUrl(os.Getenv("DB_HOST"), db_port, os.Getenv("DB_SERVICE"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), options)
	dialector := oracle.New(oracle.Config{
		DSN:                     url,
		IgnoreCase:              false,
		NamingCaseSensitive:     true,
		VarcharSizeIsCharLength: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to db %v", err)
	}
	db.AutoMigrate(&models.Tag{}, &models.Mention{}, &models.Post{})
	return &PostService{
		db: db,
	}
}

func (ps *PostService) Create(post *models.Post) (int64, error) {
	result := ps.db.Create(post)
	return result.RowsAffected, result.Error
}

func (ps *PostService) Update(post *models.Post) error {
	target, err := ps.GetById(post.Id.String())
	if err != nil {
		return err
	}
	target.AcceptNewData(post)
	ps.db.Save(target)
	return nil
}
func (ps *PostService) GetById(id string) (*models.Post, error) {
	var target *models.Post
	if err := ps.db.Preload("Mention").Preload("Tag").First(&target, id).Error; err != nil {
		return nil, err
	}
	return target, nil
}
func (ps *PostService) GetAllWithPagination(request requests.GetPostWithPaginationInterface) (posts []models.Post, amount int64, err error) {
	if posts, err = ps.FetchPosts(request); err != nil {
		return nil, invalid_amount, err
	}
	if amount, err = ps.FetchAmountPosts(request); err != nil {
		return nil, invalid_amount, err
	}
	return posts, amount, err
}

func (ps *PostService) FetchPosts(request requests.GetPostWithPaginationInterface) ([]models.Post, error) {
	var posts []models.Post
	err := ps.db.Model(&models.Post{}).Offset(request.GetOffset()).
		Preload("Mention").
		Preload("Tag").Limit(request.GetSize()).
		Where(request.GetQuery()).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
func (ps *PostService) FetchAmountPosts(request requests.GetPostWithPaginationInterface) (int64, error) {
	var amount int64
	err := ps.db.Model(&models.Post{}).Where(request.GetQuery()).
		Count(&amount).Error
	if err != nil {
		return invalid_amount, err
	}
	return amount, nil
}
func (ps *PostService) DeleteById(id string) (err error) {
	if err := ps.db.Delete(models.Post{}, id).Error; err != nil {
		return err
	}
	return
}
