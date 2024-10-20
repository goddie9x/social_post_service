package services

import (
	"fmt"
	"log"
	"os"
	"post_service/internal/constants"
	"post_service/internal/models"
	"post_service/internal/repositories"
	"post_service/internal/requests"
	pkg_constants "post_service/pkg/constants"
	"post_service/pkg/dotenv"
	"post_service/pkg/exceptions"
	"post_service/pkg/middlewares"
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
func (ps *PostService) additionQueryPostBaseOnPostIds(Ids []string) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if len(Ids) == 0 {
			return query.Where("1=0")
		}
		return query.Where(Ids)
	}
}
func (ps *PostService) additionQueryPostBaseOnUser(currentUser middlewares.UserAuth) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if currentUser.Role == pkg_constants.User {
			query.Where(`"owner_id" = ?`, currentUser.UserId)
		}
		return query
	}
}

func (ps *PostService) Create(currentUser middlewares.UserAuth, post models.Post) (*models.Post, exceptions.CommonExceptionInterface) {
	post.OwnerId = currentUser.UserId
	if *post.Type == constants.PersonalPost && post.TargetId == currentUser.UserId {
		post.Approved = true
	}
	if err := ps.db.Create(&post).Error; err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return &post, nil
}

func (ps *PostService) Update(currentUser middlewares.UserAuth, post models.Post) (*models.Post, exceptions.CommonExceptionInterface) {
	target, ex := ps.getById(post.Id)
	if ex != nil {
		return nil, ex
	}
	if !validates.CanModifyTarget(currentUser, target.OwnerId) {
		return nil, exceptions.NewNotHavePermissionException()
	}
	if target == nil {
		return nil, exceptions.NewTargetNotExistException("Post not exist")
	}
	ex = target.AcceptNewDataForUpdate(&post)
	if ex != nil {
		return nil, ex
	}
	if err := ps.db.Save(&target).Error; err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return target, nil
}
func (ps *PostService) GetByIdIfUserCanView(currentUser middlewares.UserAuth, id string) (*models.Post, exceptions.CommonExceptionInterface) {
	post, ex := ps.getById(id)

	if ex != nil {
		return nil, ex
	}
	if *post.Privacy == constants.Public {
		return post, nil
	}
	//TODO: add check friend of owner by grpc
	if validates.CanModifyTarget(currentUser, post.OwnerId) {
		return post, nil
	} else {
		return nil, exceptions.NewNotHavePermissionException()
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

func (ps *PostService) DeleteById(currentUser middlewares.UserAuth, id string) exceptions.CommonExceptionInterface {
	post, ex := ps.getById(id)
	if ex != nil {
		return ex
	}
	if !validates.CanModifyTarget(currentUser, post.OwnerId) {
		return exceptions.NewNotHavePermissionException()
	}
	if err := ps.db.Delete(&post).Error; err != nil {
		return exceptions.NewInternalErrorException("Cannot delete post, please try again later")
	}
	return nil
}

func (ps *PostService) GetPostsByTagWithPagination(currentUser middlewares.UserAuth, request requests.GetPostByTagWithPaginationRequest) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	if request.Tag == "" {
		return nil, invalid_amount, exceptions.NewBadRequestException("Tag not provided")
	}
	var postIds []string

	if err := ps.db.Table("post_tags").
		Select(`"post_id"`).
		Where(`"tag_name" = ?`, request.Tag).
		Scan(&postIds).Error; err != nil {
		return nil, invalid_amount, exceptions.NewInternalErrorException(err.Error())
	}
	requestGetPostsByTag := requests.GetPostWithPaginationRequest{
		PaginationRequest: request.PaginationRequest,
	}
	additionQueryBaseOnUser := ps.additionQueryPostBaseOnUser(currentUser)
	additionQueryBaseOnPostIds := ps.additionQueryPostBaseOnPostIds(postIds)
	return ps.GetPostsWithPagination(requestGetPostsByTag, additionQueryBaseOnUser, additionQueryBaseOnPostIds)
}
func (ps *PostService) GetPostByMentionWithPagination(currentUser middlewares.UserAuth, request requests.GetPostByMentionWithPaginationRequest) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	if request.Mention == "" {
		return nil, invalid_amount, exceptions.NewBadRequestException("Mention not provided")
	}
	var postIds []string

	if err := ps.db.Table("mentions").
		Select(`"post_id"`).
		Where(`"user_id" = ?`, request.Mention).
		Scan(&postIds).Error; err != nil {
		return nil, invalid_amount, exceptions.NewInternalErrorException(err.Error())
	}
	requestGetPostsByMention := requests.GetPostWithPaginationRequest{
		PaginationRequest: request.PaginationRequest,
	}
	additionQueryBaseOnPostIds := ps.additionQueryPostBaseOnPostIds(postIds)
	return ps.GetPostsWithPagination(requestGetPostsByMention, additionQueryBaseOnPostIds)
}
func (ps *PostService) GetPostsForUserProfile(currentUser middlewares.UserAuth, request requests.GetPostForUserWithPagination) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	if request.UserId == "" {
		return nil, invalid_amount, exceptions.NewBadRequestException("userId is required")
	}
	var totalCount int64
	whereClause := `("type" = ? AND "owner_id" = ? `
	query := ps.db.Model(&models.Post{})

	if !validates.CanModifyTarget(currentUser, request.UserId) {
		whereClause += `AND "privacy" = ` + fmt.Sprintf("%v", constants.Private)
	}
	whereClause += `) OR "id" IN (SELECT "post_id" FROM "mentions" WHERE "user_id" = ? AND "accepted_show_in_profile" = 1)`
	query.Where(whereClause,
		constants.PersonalPost, request.UserId, request.UserId)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, exceptions.NewInternalErrorException(err.Error())
	}

	err := query.
		Preload("Mentions").
		Preload("Tags").
		Offset(request.GetOffset()).
		Limit(request.GetSize()).
		Order(`"created_at" DESC`).
		Find(&posts).Error

	if err != nil {
		return nil, 0, exceptions.NewInternalErrorException(err.Error())
	}
	amountPage = (totalCount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return posts, amountPage, nil
}
func (ps *PostService) GetPostsWithPagination(request requests.GetPostWithPaginationInterface, additionalWhereClauses ...func(*gorm.DB) *gorm.DB) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	if posts, ex = ps.FetchPosts(request, additionalWhereClauses...); ex != nil {
		return nil, invalid_amount, ex
	}
	amount, ex := ps.FetchAmountPosts(request, additionalWhereClauses...)
	if ex != nil {
		return nil, invalid_amount, ex
	}
	amountPage = (amount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return posts, amountPage, nil
}

func (ps *PostService) FetchPosts(request requests.GetPostWithPaginationInterface, additionOptionsWrappers ...func(*gorm.DB) *gorm.DB) ([]models.Post, exceptions.CommonExceptionInterface) {
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
func (ps *PostService) FetchAmountPosts(request requests.GetPostWithPaginationInterface, additionOptionsWrappers ...func(*gorm.DB) *gorm.DB) (int64, exceptions.CommonExceptionInterface) {
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
