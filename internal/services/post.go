package services

import (
	"fmt"
	"log"
	"os"
	"post_service/internal/constants"
	"post_service/internal/models"
	"post_service/internal/repositories"
	"post_service/internal/requests"
	"post_service/internal/responses"
	"post_service/pkg/dotenv"
	"post_service/pkg/exceptions"
	"post_service/pkg/validates"
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
	db.AutoMigrate(&models.Tag{}, &models.Post{}, &models.Mention{}, &models.PostTag{})
	return &PostService{
		db: db,
	}
}

func (ps *PostService) Create(request requests.PostWithAuthRequest) responses.PostResponse {
	userId := request.User.UserId
	request.Post.OwnerId = userId
	if *request.Post.Type == constants.PersonalPost && request.Post.TargetId == userId {
		request.Post.Approved = true
	}
	if err := ps.db.Create(&request.Post).Error; err != nil {
		return responses.PostResponse{
			Ex: exceptions.NewInternalErrorException(err.Error()),
		}
	}
	return responses.PostResponse{
		Post: &request.Post,
	}
}

func (ps *PostService) Update(request requests.PostWithAuthRequest) responses.PostResponse {
	target, ex := ps.getById(request.Post.Id)
	if ex != nil {
		return responses.PostResponse{
			Ex: ex,
		}
	}
	if !validates.CanModifyTarget(request.User, target.OwnerId) {
		return responses.PostResponse{
			Ex: exceptions.NewNotHavePermissionException(),
		}
	}
	if target == nil {
		return responses.PostResponse{
			Ex: exceptions.NewTargetNotExistException("Post not exist"),
		}
	}
	ex = target.AcceptNewDataForUpdate(&request.Post)
	if ex != nil {
		return responses.PostResponse{
			Ex: ex,
		}
	}
	if err := ps.db.Save(&target).Error; err != nil {
		return responses.PostResponse{
			Ex: exceptions.NewInternalErrorException(err.Error()),
		}
	}
	return responses.PostResponse{
		Post: target,
	}
}
func (ps *PostService) GetByIdIfUserCanView(request requests.RequestWithAuthAndId) responses.PostResponse {
	post, ex := ps.getById(request.Id)
	if ex != nil {
		return responses.PostResponse{
			Ex: ex,
		}
	}
	if *post.Privacy == constants.Public {
		return responses.PostResponse{
			Post: post,
		}
	}
	//TODO: add check friend of owner by grpc
	if validates.CanModifyTarget(request.User, post.OwnerId) {
		return responses.PostResponse{
			Post: post,
		}
	} else {
		return responses.PostResponse{
			Ex: exceptions.NewNotHavePermissionException(),
		}
	}
}

func (ps *PostService) getById(id string) (*models.Post, exceptions.CommonExceptionInterface) {
	if id == "" {
		return nil, exceptions.NewBadRequestException("Post id is required")
	}
	target := models.Post{
		Id: id,
	}
	if err := ps.db.Preload("Mentions").Preload("Tags").First(&target).Error; err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return &target, nil
}

func (ps *PostService) DeleteById(request requests.RequestWithAuthAndId) responses.PostResponse {
	post, ex := ps.getById(request.Id)
	if ex != nil {
		return responses.PostResponse{
			Ex: ex,
		}
	}
	if !validates.CanModifyTarget(request.User, post.OwnerId) {
		return responses.PostResponse{
			Ex: exceptions.NewNotHavePermissionException(),
		}
	}
	if err := ps.db.Delete(&post).Error; err != nil {
		return responses.PostResponse{
			Ex: exceptions.NewInternalErrorException("Cannot delete post, please try again later"),
		}
	}
	return responses.PostResponse{
		Post: post,
	}
}

func (ps *PostService) GetPostsByTagWithPagination(request requests.GetPostByTagWithPaginationRequest) responses.ListPostWithPaginationResponse {
	if request.Tag == "" {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewBadRequestException("Tag not provided"),
		}
	}
	var postIds []string

	if err := ps.db.Table("post_tags").
		Select(`"post_id"`).
		Where(`"tag_name" = ?`, request.Tag).
		Scan(&postIds).Error; err != nil {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewInternalErrorException(err.Error()),
		}
	}
	request.PostIds = postIds
	return ps.GetPostsWithPagination(request)
}
func (ps *PostService) GetPostByMentionWithPagination(request requests.GetPostByMentionWithPaginationRequest) responses.ListPostWithPaginationResponse {
	if request.Mention == "" {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewBadRequestException("Mention not provided"),
		}
	}
	var postIds []string

	if err := ps.db.Table("mentions").
		Select(`"post_id"`).
		Where(`"user_id" = ?`, request.Mention).
		Scan(&postIds).Error; err != nil {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewInternalErrorException(err.Error()),
		}
	}
	request.PostIds = postIds

	return ps.GetPostsWithPagination(request)
}
func (ps *PostService) GetPostsForUserProfile(request requests.GetPostForUserWithPagination) responses.ListPostWithPaginationResponse {
	if request.UserId == "" {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewBadRequestException("userId is required"),
		}
	}
	var totalCount int64
	whereClause := `("type" = ? AND "owner_id" = ? `
	query := ps.db.Model(&models.Post{})

	if !validates.CanModifyTarget(request.User, request.UserId) {
		whereClause += `AND "privacy" = ` + fmt.Sprintf("%v", constants.Private)
	}
	whereClause += `) OR "id" IN (SELECT "post_id" FROM "mentions" WHERE "user_id" = ? AND "accepted_show_in_profile" = 1)`
	query.Where(whereClause,
		constants.PersonalPost, request.UserId, request.UserId)

	if err := query.Count(&totalCount).Error; err != nil {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewInternalErrorException(err.Error()),
		}
	}
	var posts []models.Post
	err := query.
		Preload("Mentions").
		Preload("Tags").
		Offset(request.GetOffset()).
		Limit(request.GetSize()).
		Order(`"created_at" DESC`).
		Find(&posts).Error

	if err != nil {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         exceptions.NewInternalErrorException(err.Error()),
		}
	}
	amountPage := (totalCount + int64(request.GetSize()) - 1) / int64(request.GetSize())
	return responses.ListPostWithPaginationResponse{
		Posts:      posts,
		AmountPage: amountPage,
	}
}
func (ps *PostService) GetPostsWithPagination(request requests.GetPostWithPaginationInterface) responses.ListPostWithPaginationResponse {
	var posts []models.Post
	additionQuery := request.GetAdditionWhereClause()
	posts, ex := ps.FetchPosts(request, additionQuery)

	if ex != nil {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         ex,
		}
	}
	amount, ex := ps.FetchAmountPosts(request, additionQuery)

	if ex != nil {
		return responses.ListPostWithPaginationResponse{
			AmountPage: invalid_amount,
			Ex:         ex,
		}
	}
	amountPage := (amount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return responses.ListPostWithPaginationResponse{
		Posts:      posts,
		AmountPage: amountPage,
	}
}

func (ps *PostService) FetchPosts(request requests.GetPostWithPaginationInterface, additionOptionsWrappers []func(*gorm.DB) *gorm.DB) ([]models.Post, exceptions.CommonExceptionInterface) {
	var posts []models.Post
	query := ps.db.Model(&models.Post{}).Offset(request.GetOffset()).
		Preload("Mentions").
		Preload("Tags").Limit(request.GetSize())
	for _, additionQuery := range additionOptionsWrappers {
		if additionQuery != nil {
			query = additionQuery(query)
		}
	}
	err := query.
		Where(request.GetQuery()).Find(&posts).Error
	if err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return posts, nil
}
func (ps *PostService) FetchAmountPosts(request requests.GetPostWithPaginationInterface, additionOptionsWrappers []func(*gorm.DB) *gorm.DB) (int64, exceptions.CommonExceptionInterface) {
	var amount int64

	query := ps.db.Model(&models.Post{}).
		Where(request.GetQuery())
	for _, additionQuery := range additionOptionsWrappers {
		if additionQuery != nil {
			query = additionQuery(query)
		}
	}
	err := query.
		Count(&amount).Error
	if err != nil {
		return invalid_amount, exceptions.NewInternalErrorException(err.Error())
	}
	return amount, nil
}
