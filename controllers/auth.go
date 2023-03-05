package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"simplerestapi/common"
	"simplerestapi/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRepo struct {
	Db *gorm.DB
}
type userInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{Db: db}
}
func (repository *AuthRepo) Login(c *gin.Context) {
	var input userInput
	var user models.User
	c.BindJSON(&input)
	err := models.GetUserByEmail(repository.Db, &user, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if input.Password != user.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "password not correct"})
		return
	}else {
		token,err := common.GenerateToken(user)
		if err != nil{
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK,gin.H{"token":token})
	}
}
