package controllers

import (
	"net/http"
	"post_service/internal/models"
	"post_service/internal/repositories"
	"post_service/internal/requests"
	request_internal "post_service/internal/requests"
	"post_service/pkg/exceptions"
	"post_service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	ps repositories.PostRepository
}

type PostControllerInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	GetById(c *gin.Context)
	DeleteById(c *gin.Context)
	GetPostByTagWithPagination(c *gin.Context)
	GetPostInGroupWithPagination(c *gin.Context)
	GetPostByMentionWithPagination(c *gin.Context)
	GetPostByUserInGroupWithPagination(c *gin.Context)
	GetPostForUserWithPagination(c *gin.Context)
}

func CreatePostController(ps repositories.PostRepository) PostControllerInterface {
	return &PostController{
		ps: ps,
	}
}
func (pc *PostController) Create(c *gin.Context) {
	var post models.Post
	currentUser := middlewares.GetUserAuthFromContext(c)
	if err := c.ShouldBindJSON(&post); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	request := request_internal.PostWithAuthRequest{User: currentUser, Post: post}
	responses := pc.ps.Create(request)

	if responses.Ex != nil {
		exceptions.HandleExceptionByGin(c, responses.Ex)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": responses.Post})
}
func (pc *PostController) Update(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	currentUser := middlewares.GetUserAuthFromContext(c)
	request := request_internal.PostWithAuthRequest{
		User: currentUser,
		Post: post,
	}
	response := pc.ps.Update(request)

	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"post": response.Post})
}
func (pc *PostController) GetById(c *gin.Context) {
	postId := c.Param("id")
	currentUser := middlewares.GetUserAuthFromContext(c)
	request := request_internal.RequestWithAuthAndId{
		User: currentUser,
		Id:   postId,
	}
	response := pc.ps.GetByIdIfUserCanView(request)

	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": response.Post})
}
func (pc *PostController) DeleteById(c *gin.Context) {
	postId := c.Param("id")
	currentUser := middlewares.GetUserAuthFromContext(c)
	request := request_internal.RequestWithAuthAndId{
		User: currentUser,
		Id:   postId,
	}

	if response := pc.ps.DeleteById(request); response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted post"})
}
func (pc *PostController) GetPostByTagWithPagination(c *gin.Context) {
	var request requests.GetPostByTagWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException("Cannot get info from query")
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	currentUser := middlewares.GetUserAuthFromContext(c)
	request.User = currentUser
	response := pc.ps.GetPostsByTagWithPagination(request)

	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": response.Posts, "amountPage": response.AmountPage})
}
func (pc *PostController) GetPostForUserWithPagination(c *gin.Context) {
	var request requests.GetPostForUserWithPagination

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	request.User = middlewares.GetUserAuthFromContext(c)
	response := pc.ps.GetPostsForUserProfile(request)

	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": response.Posts, "amountPage": response.AmountPage})
}
func (pc *PostController) GetPostByMentionWithPagination(c *gin.Context) {
	var request requests.GetPostByMentionWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())

		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	request.User = middlewares.GetUserAuthFromContext(c)
	response := pc.ps.GetPostByMentionWithPagination(request)

	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": response.Posts, "amountPage": response.AmountPage})
}

func (pc *PostController) GetPostInGroupWithPagination(c *gin.Context) {
	var request requests.GetPostInGroupWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	response := pc.ps.GetPostsWithPagination(request)
	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": response.Posts, "amountPage": response.AmountPage})
}
func (pc *PostController) GetPostsForUserProfile(c *gin.Context) {
	var request requests.GetPostForUserWithPagination
	request.User = middlewares.GetUserAuthFromContext(c)

	if err := c.ShouldBindQuery(&request); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	response := pc.ps.GetPostsForUserProfile(request)

	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": response.Ex, "amountPage": response.AmountPage})
}

func (pc *PostController) GetPostByUserInGroupWithPagination(c *gin.Context) {
	var request request_internal.GetPostByUserInGroupRequestWithPagination

	if err := c.ShouldBindQuery(&request); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	response := pc.ps.GetPostsWithPagination(request)
	if response.Ex != nil {
		exceptions.HandleExceptionByGin(c, response.Ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": response.Posts, "amountPage": response.AmountPage})
}
