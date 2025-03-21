package routes

import (
	"blog-api/controllers"
	"blog-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	auth := r.Group("/").Use(middlewares.AuthMiddleware())
	auth.POST("/posts", controllers.CreatePost)
	r.GET("/posts", controllers.GetPosts)
}
