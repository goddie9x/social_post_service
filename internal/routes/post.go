package routes

import (
	"post_service/internal/controllers"
	"post_service/internal/services"
	pkg "post_service/pkg/eureka"
	"post_service/pkg/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func MappingRoute(r *gin.Engine) {
	start_time := time.Now()
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			posts := v1.Group("/posts")
			{
				pc := controllers.CreatePostController(services.NewPostService())
				posts.GET("/by-tag", pc.GetPostByTagWithPagination)
				posts.GET("/by-mention", pc.GetPostByUserWithPagination)
				posts.GET("/for-user", pc.GetPostForUserWithPagination)
				posts.Use(middlewares.PutAuthToContext)
				posts.GET("/:id", pc.GetById)
				posts.POST("/create", pc.Create)
				posts.PATCH("/update", pc.Update)
				posts.DELETE("/delete/:id", pc.DeleteById)
			}
		}
	}
	global := r.Group("/")
	{
		global.GET("health", pkg.Health)
		global.GET("status", pkg.Status(start_time))
	}
}
