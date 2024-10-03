package services

import (
	"errors"
	"log"
	"os"
	"post_service/internal/constants"
	"post_service/internal/models"
	"post_service/internal/repositories"
	"post_service/internal/requests"
	"post_service/pkg/dotenv"
	"strconv"

	oracle "github.com/godoes/gorm-oracle"
	"github.com/google/uuid"
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
		"CONNECTION TIMEOUT": dotenv.GetEnvOrDefaultValue("DB_CONNECTION_TIMEOUT", "90"),
		"SSL":                dotenv.GetEnvOrDefaultValue("DB_SSL", "false"),
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
	db.AutoMigrate(&models.Tag{}, &models.Post{}, &models.Mention{})
	return &PostService{
		db: db,
	}
}

func (ps *PostService) Create(post *models.Post) (err error) {
	err = ps.db.Create(post).Error
	return
}

func (ps *PostService) Update(post *models.Post) error {
	target, err := ps.GetById(post.Id.String())
	if err != nil {
		return err
	}
	target.AcceptNewData(post)
	if err = ps.db.Save(target).Error; err != nil {
		return err
	}
	return nil
}
func (ps *PostService) GetById(id string) (*models.Post, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	var target *models.Post
	if err := ps.db.Preload("Mention").Preload("Tag").First(&target, id).Error; err != nil {
		return nil, err
	}
	return target, nil
}
func (ps *PostService) GetPostsByTagWithPagination(request requests.GetPostByTagsWithPaginationRequest) (posts []models.Post, amountPage int64, err error) {
	var postIds []uuid.UUID

	if err = ps.db.Table("post_tags").
		Select("post_id").
		Where("tag_name IN ?", request.Tags).
		Scan(&postIds).Error; err != nil {
		return nil, invalid_amount, err
	}
	requestGetPostsByTag := requests.GetPostWithPaginationRequest{
		PaginationRequest: request.PaginationRequest,
		PostQuery:         postIds,
	}
	return ps.GetPostsWithPagination(requestGetPostsByTag)
}
func (ps *PostService) GetPostsForUserProfile(request requests.GetPostForUserWithPagination) (posts []models.Post, amountPage int64, err error) {
	var totalCount int64

	query := ps.db.Model(&models.Post{}).
		Where("(type = ? AND owner_id = ?) OR id IN (SELECT post_id FROM mentions WHERE user_id = ? AND accepted_show_in_profile = true)",
			constants.PersonalPost, request.UserId, request.UserId)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err = query.
		Preload("Mentions").
		Preload("Tags").
		Offset(request.GetOffset()).
		Limit(request.GetSize()).
		Order("created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}
	amountPage = (totalCount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return posts, amountPage, nil
}
func (ps *PostService) GetPostsWithPagination(request requests.GetPostWithPaginationInterface) (posts []models.Post, amountPage int64, err error) {
	if posts, err = ps.FetchPosts(request); err != nil {
		return nil, invalid_amount, err
	}
	amount, err := ps.FetchAmountPosts(request)
	if err != nil {
		return nil, invalid_amount, err
	}
	amountPage = (amount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return posts, amountPage, err
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
	if id == "" {
		return errors.New("id cannot be empty")
	}
	if err := ps.db.Delete(models.Post{}, id).Error; err != nil {
		return err
	}
	return
}
