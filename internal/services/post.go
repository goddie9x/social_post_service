package services

import (
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

	"github.com/gin-gonic/gin"
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

func (ps *PostService) additionQueryPostBaseOnUser(currentUser middlewares.UserAuth) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if currentUser.Role == pkg_constants.User {
			query.Where("owner_id = ?", currentUser.UserId)
		}
		return query
	}
}

func (ps *PostService) Create(c *gin.Context) (post *models.Post, ex exceptions.CommonExceptionInterface) {
	currentUser := middlewares.GetUserAuthFromContext(c)
	if err := c.ShouldBindJSON(&post); err != nil {
		return nil, exceptions.NewBadRequestException(err.Error())
	}
	post.OwnerId = currentUser.UserId
	if err := ps.db.Create(&post).Error; err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return post, nil
}

func (ps *PostService) Update(c *gin.Context) (post *models.Post, ex exceptions.CommonExceptionInterface) {
	if err := c.ShouldBindJSON(&post); err != nil {
		return nil, exceptions.NewBadRequestException(err.Error())
	}
	target, ex := ps.getById(post.Id)
	if ex != nil {
		return nil, ex
	}
	currentUser := middlewares.GetUserAuthFromContext(c)
	if !validates.CanModifyTarget(currentUser, target.OwnerId) {
		return nil, exceptions.NewNotHavePermissionException()
	}
	if err := ps.db.Save(target).Error; err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return target, nil
}
func (ps *PostService) GetById(c *gin.Context) (*models.Post, exceptions.CommonExceptionInterface) {
	var post *models.Post

	if err := c.ShouldBind(post); err != nil {
		return nil, exceptions.NewBadRequestException(err.Error())
	}
	id := post.Id
	if id == "" {
		return nil, exceptions.NewBadRequestException("Post id must provided")
	}
	post, ex := ps.getById(id)
	if ex != nil {
		return nil, ex
	}
	currentUser := middlewares.GetUserAuthFromContext(c)
	if post.Privacy == constants.Public {
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
	var target *models.Post
	if err := ps.db.Preload("Mention").Preload("Tag").First(&target, id).Error; err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return target, nil
}

func (ps *PostService) DeleteById(c *gin.Context) exceptions.CommonExceptionInterface {
	var post *models.Post

	if err := c.ShouldBind(post); err != nil {
		return exceptions.NewBadRequestException(err.Error())
	}
	id := post.Id
	if id == "" {
		return exceptions.NewBadRequestException("Post id must provided")
	}
	currentUser := middlewares.GetUserAuthFromContext(c)
	post, ex := ps.getById(id)
	if ex != nil {
		return ex
	}
	if !validates.CanModifyTarget(currentUser, post.OwnerId) {
		return exceptions.NewNotHavePermissionException()
	}
	if err := ps.db.Delete(models.Post{}, id).Error; err != nil {
		return exceptions.NewInternalErrorException("Cannot delete post, please try again later")
	}
	return nil
}

func (ps *PostService) GetPostsByTagWithPagination(c *gin.Context) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	var request requests.GetPostByTagsWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		return nil, invalid_amount, exceptions.NewBadRequestException("Cannot get info from query")
	}
	if len(request.Tags) < 1 {
		return nil, invalid_amount, exceptions.NewBadRequestException("Tags not provided")
	}
	var postIds []uuid.UUID

	if err := ps.db.Table("post_tags").
		Select("post_id").
		Where("tag_name IN ?", request.Tags).
		Scan(&postIds).Error; err != nil {
		return nil, invalid_amount, exceptions.NewInternalErrorException(err.Error())
	}
	currentUser := middlewares.GetUserAuthFromContext(c)

	requestGetPostsByTag := requests.GetPostWithPaginationRequest{
		PaginationRequest: request.PaginationRequest,
		PostQuery:         postIds,
	}
	return ps.GetPostsWithPagination(requestGetPostsByTag, ps.additionQueryPostBaseOnUser(currentUser))
}
func (ps *PostService) GetPostsForUserProfile(c *gin.Context) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	var request requests.GetPostForUserWithPagination

	if err := c.ShouldBindQuery(&request); err != nil {
		return nil, invalid_amount, exceptions.NewBadRequestException(err.Error())
	}
	currentUser := middlewares.GetUserAuthFromContext(c)
	userId := c.Query("userId")
	if userId == "" {
		return nil, invalid_amount, exceptions.NewBadRequestException("userId is required")
	}
	var totalCount int64
	whereClause := "(type = ? AND owner_id = ?) OR id IN (SELECT post_id FROM mentions WHERE user_id = ? AND accepted_show_in_profile = true)"
	query := ps.db.Model(&models.Post{})
	if !validates.CanModifyTarget(currentUser, userId) {
		query.Where(whereClause,
			constants.PersonalPost, request.UserId, request.UserId)
	} else {
		query.Where(whereClause,
			constants.PersonalPost, request.UserId, request.UserId)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, exceptions.NewInternalErrorException(err.Error())
	}

	err := query.
		Preload("Mentions").
		Preload("Tags").
		Offset(request.GetOffset()).
		Limit(request.GetSize()).
		Order("created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, 0, exceptions.NewInternalErrorException(err.Error())
	}
	amountPage = (totalCount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return posts, amountPage, nil
}
func (ps *PostService) GetPostsWithPagination(request requests.GetPostWithPaginationInterface, additionalWhereClause func(*gorm.DB) *gorm.DB) (posts []models.Post, amountPage int64, ex exceptions.CommonExceptionInterface) {
	if posts, ex = ps.FetchPosts(request, additionalWhereClause); ex != nil {
		return nil, invalid_amount, ex
	}
	amount, ex := ps.FetchAmountPosts(request, additionalWhereClause)
	if ex != nil {
		return nil, invalid_amount, ex
	}
	amountPage = (amount + int64(request.GetSize()) - 1) / int64(request.GetSize())

	return posts, amountPage, nil
}

func (ps *PostService) FetchPosts(request requests.GetPostWithPaginationInterface, additionOptionsWrapper func(*gorm.DB) *gorm.DB) ([]models.Post, exceptions.CommonExceptionInterface) {
	var posts []models.Post
	query := ps.db.Model(&models.Post{}).Offset(request.GetOffset()).
		Preload("Mention").
		Preload("Tag").Limit(request.GetSize())
	if additionOptionsWrapper != nil {
		query = additionOptionsWrapper(query)
	}
	err := query.
		Where(request.GetQuery()).
		Where(additionOptionsWrapper).Find(&posts).Error
	if err != nil {
		return nil, exceptions.NewInternalErrorException(err.Error())
	}
	return posts, nil
}
func (ps *PostService) FetchAmountPosts(request requests.GetPostWithPaginationInterface, additionOptionsWrapper func(*gorm.DB) *gorm.DB) (int64, exceptions.CommonExceptionInterface) {
	var amount int64
	query := ps.db.Model(&models.Post{}).
		Where(request.GetQuery())

	if additionOptionsWrapper != nil {
		query = additionOptionsWrapper(query)
	}
	err := query.
		Count(&amount).Error
	if err != nil {
		return invalid_amount, exceptions.NewInternalErrorException(err.Error())
	}
	return amount, nil
}
func (ps *PostService) GetPostsInGroupWithPagination(c *gin.Context) ([]models.Post, int64, exceptions.CommonExceptionInterface) {
	var request requests.GetPostInGroupWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		return nil, invalid_amount, exceptions.NewBadRequestException(err.Error())
	}
	getPostInGroupRequest := requests.GetPostWithPaginationRequest{
		PaginationRequest: request.PaginationRequest,
		PostQuery: models.Post{
			TargetId: request.TargetId,
			Type:     constants.GroupPost,
		},
	}
	return ps.GetPostsWithPagination(getPostInGroupRequest, nil)
}
