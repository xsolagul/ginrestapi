package routes

import (
	"simplerestapi/common"
	"simplerestapi/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUserRoute(r *gin.Engine,db *gorm.DB) {
	
	userRepo := controllers.NewUserRepo(db)
	userGroup := r.Group("/")
	authTool := common.NewTokenTool(db)
	userGroup.Use(authTool.AuthorizationMiddleware("user"))
	userGroup.POST("/users", userRepo.CreateUser)
	userGroup.GET("/users", userRepo.GetUsers)
	userGroup.GET("/users/:id", userRepo.GetUser)
	userGroup.PUT("/users/:id", userRepo.UpdateUser)
	userGroup.DELETE("/users/:id", userRepo.DeleteUser)
}