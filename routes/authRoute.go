package routes

import (
	"simplerestapi/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetAuthRoute(r *gin.Engine,db *gorm.DB) {
	
	authRepo := controllers.NewAuthRepo(db)
	userRepo := controllers.NewUserRepo(db)
	authGroup := r.Group("/")
	authGroup.POST("/login",authRepo.Login)
	authGroup.POST("/register",userRepo.CreateUser)
	// userGroup.POST("/users", userRepo.CreateUser)
	// userGroup.GET("/users", userRepo.GetUsers)
	// userGroup.GET("/users/:id", userRepo.GetUser)
	// userGroup.PUT("/users/:id", userRepo.UpdateUser)
	// userGroup.DELETE("/users/:id", userRepo.DeleteUser)
}