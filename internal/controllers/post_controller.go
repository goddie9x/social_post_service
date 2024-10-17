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
	var post *models.Post
	currentUser := middlewares.GetUserAuthFromContext(c)
	if err := c.ShouldBindJSON(post); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	post, ex := pc.ps.Create(currentUser, post)

	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": post})
}
func (pc *PostController) Update(c *gin.Context) {
	var post *models.Post
	if err := c.ShouldBindJSON(post); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	currentUser := middlewares.GetUserAuthFromContext(c)

	post, ex := pc.ps.Update(currentUser, post)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"post": post})
}
func (pc *PostController) GetById(c *gin.Context) {
	var post *models.Post
	currentUser := middlewares.GetUserAuthFromContext(c)

	if err := c.ShouldBind(post); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	post, ex := pc.ps.GetByIdIfUserCanView(currentUser, post.Id)

	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}
func (pc *PostController) DeleteById(c *gin.Context) {
	var post *models.Post
	if err := c.ShouldBind(post); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	currentUser := middlewares.GetUserAuthFromContext(c)

	if ex := pc.ps.DeleteById(currentUser, post.Id); ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
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

	posts, amountPage, ex := pc.ps.GetPostsByTagWithPagination(currentUser, request)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}
func (pc *PostController) GetPostForUserWithPagination(c *gin.Context) {
	var request requests.GetPostForUserWithPagination

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	currentUser := middlewares.GetUserAuthFromContext(c)

	posts, amountPage, ex := pc.ps.GetPostsForUserProfile(currentUser, request)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}
func (pc *PostController) GetPostByMentionWithPagination(c *gin.Context) {
	var request requests.GetPostByMentionWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())

		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	posts, amountPage, ex := pc.ps.GetPostsWithPagination(request, nil)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}

func (pc *PostController) GetPostInGroupWithPagination(c *gin.Context) {
	var request requests.GetPostInGroupWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		ex := exceptions.NewBadRequestException(err.Error())
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	posts, amountPage, ex := pc.ps.GetPostsWithPagination(request, nil)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}
func (pc *PostController) GetPostsForUserProfile(c *gin.Context) {
	var request requests.GetPostForUserWithPagination
	currentUser := middlewares.GetUserAuthFromContext(c)

	if err := c.ShouldBindQuery(&request); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}

	posts, amountPage, ex := pc.ps.GetPostsForUserProfile(currentUser, request)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage, "page": request.GetPage()})
}

func (pc *PostController) GetPostByUserInGroupWithPagination(c *gin.Context) {
	var request request_internal.GetPostByUserInGroupRequestWithPagination
	var postQuery models.Post

	if err := c.ShouldBindQuery(&request); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	if err := c.ShouldBindQuery(&postQuery); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	posts, amountPage, ex := pc.ps.GetPostsWithPagination(request, nil)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage, "page": request.GetPage()})
}
