package controllers

import (
	"net/http"
	"post_service/internal/constants"
	"post_service/internal/models"
	"post_service/internal/repositories"
	request_internal "post_service/internal/requests"
	"post_service/pkg/errors"

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
	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	err := pc.ps.Create(&post)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": post})
}
func (pc *PostController) Update(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := pc.ps.Update(&post); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"post": post})
}
func (pc *PostController) GetById(c *gin.Context) {
	var post *models.Post

	if err := c.ShouldBind(post); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	post, err := pc.ps.GetById(post.Id.String())
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}
func (pc *PostController) DeleteById(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindUri(&post); err != nil {
		errors.HandleError(c, http.StatusBadRequest, "Id not exist or invalid Id")
		return
	}
	if err := pc.ps.DeleteById(post.Id.String()); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted post"})
}
func (pc *PostController) GetPostByTagWithPagination(c *gin.Context) {
	var request request_internal.GetPostByTagsWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(request.Tags) < 1 {
		errors.HandleError(c, http.StatusBadRequest, "Tags not provided")
		return
	}
	posts, amountPage, err := pc.ps.GetPostsByTagWithPagination(request)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage, "page": request.GetPage()})
}
func (pc *PostController) GetPostInGroupWithPagination(c *gin.Context) {
	var request request_internal.GetPostInGroupWithPaginationRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	getPostInGroupRequest := request_internal.GetPostWithPaginationRequest{
		PaginationRequest: request.PaginationRequest,
		PostQuery: models.Post{
			TargetId: request.TargetId,
			Type:     constants.GroupPost,
		},
	}
	posts, amountPage, err := pc.ps.GetPostsWithPagination(getPostInGroupRequest)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage, "page": request.GetPage()})
}
func (pc *PostController) GetPostByUserWithPagination(c *gin.Context) {
	var request request_internal.GetPostWithPaginationRequest
	var postQuery models.Post

	if err := c.ShouldBindQuery(&request); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
	}
	if err := c.ShouldBindQuery(&postQuery); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
	}
	request.PostQuery = models.Post{
		OwnerId: postQuery.OwnerId,
	}
	posts, amountPage, err := pc.ps.GetPostsWithPagination(request)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "amountPage": amountPage, "page": request.GetPage()})
}
func (pc *PostController) GetPostForUserWithPagination(c *gin.Context) {
	var request request_internal.GetPostForUserWithPagination

	if err := c.ShouldBindQuery(&request); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	userId := c.Query("userId")
	if userId == "" {
		errors.HandleError(c, http.StatusBadRequest, "userId is required")
		return
	}

	posts, amountPage, err := pc.ps.GetPostsForUserProfile(request)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":      posts,
		"amountPage": amountPage,
		"page":       request.GetPage(),
	})
}
