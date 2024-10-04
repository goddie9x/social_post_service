package controllers

import (
	"net/http"
	"post_service/internal/models"
	"post_service/internal/repositories"
	request_internal "post_service/internal/requests"
	"post_service/pkg/exceptions"

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
	GetPostByUserWithPagination(c *gin.Context)
	GetPostForUserWithPagination(c *gin.Context)
}

func CreatePostController(ps repositories.PostRepository) PostControllerInterface {
	return &PostController{
		ps: ps,
	}
}
func (pc *PostController) Create(c *gin.Context) {
	post, ex := pc.ps.Create(c)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": post})
}
func (pc *PostController) Update(c *gin.Context) {
	post, ex := pc.ps.Update(c)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"post": post})
}
func (pc *PostController) GetById(c *gin.Context) {
	post, ex := pc.ps.GetById(c)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}
func (pc *PostController) DeleteById(c *gin.Context) {
	if ex := pc.ps.DeleteById(c); ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted post"})
}
func (pc *PostController) GetPostByTagWithPagination(c *gin.Context) {

	posts, amountPage, ex := pc.ps.GetPostsByTagWithPagination(c)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}
func (pc *PostController) GetPostForUserWithPagination(c *gin.Context) {

	posts, amountPage, ex := pc.ps.GetPostsForUserProfile(c)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}
func (pc *PostController) GetPostInGroupWithPagination(c *gin.Context) {
	posts, amountPage, ex := pc.ps.GetPostsInGroupWithPagination(c)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage})
}
func (pc *PostController) GetPostByUserWithPagination(c *gin.Context) {
	var request request_internal.GetPostWithPaginationRequest
	var postQuery models.Post

	if err := c.ShouldBindQuery(&request); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	if err := c.ShouldBindQuery(&postQuery); err != nil {
		exceptions.HandleExceptionByGin(c, exceptions.NewBadRequestException(err.Error()))
	}
	request.PostQuery = models.Post{
		OwnerId: postQuery.OwnerId,
	}
	posts, amountPage, ex := pc.ps.GetPostsWithPagination(request, nil)
	if ex != nil {
		exceptions.HandleExceptionByGin(c, ex)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage, "page": request.GetPage()})
}
