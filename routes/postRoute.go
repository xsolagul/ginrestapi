package routes

import (
	"simplerestapi/common"
	"simplerestapi/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetPostRoute(r *gin.Engine,db *gorm.DB) {
	
	postRepo := controllers.NewPostRepo(db)
	postGroup := r.Group("/")
	authTool := common.NewTokenTool(db)
	postGroup.Use(authTool.AuthorizationMiddleware("user"))
	postGroup.POST("/posts", postRepo.CreatePost)
	postGroup.GET("/posts", postRepo.GetPosts)
	postGroup.GET("/posts/:id", postRepo.GetPost)
	postGroup.PUT("/posts/:id", postRepo.UpdatePost)
	postGroup.DELETE("/posts/:id", postRepo.DeletePost)


}