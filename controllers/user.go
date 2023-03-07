package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	//"simplerestapi/database"
	"simplerestapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}
//struct for only create user from client 
type UserFromClient struct {
	//ID    int `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	//db := database.InitDb()//guess do one time
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

// create user
func (repository *UserRepo) CreateUser(c *gin.Context) {
	var user models.User
	var data UserFromClient
	err := c.BindJSON(&data)
	if err != nil {
		return
	}
	user.Email = data.Email
	user.Name = data.Name
	// Hashing the password with the default cost of 10
	password := []byte(data.Password)
    hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	user.Password = string(hashedPassword)
	//bind json with other struct and set to user struct after
	err = models.CreateUser(repository.Db, &user)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// get users
func (repository *UserRepo) GetUsers(c *gin.Context) {
	var user []models.User
	err := models.GetUsers(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by id
func (repository *UserRepo) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	err := models.GetUser(repository.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// update user
func (repository *UserRepo) UpdateUser(c *gin.Context) {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))
	// err := models.GetUser(repository.Db, &user, id)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		c.AbortWithStatus(http.StatusNotFound)
	// 		return
	// 	}

	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }
	c.BindJSON(&user)
	user.ID = id
	err := models.UpdateUser(repository.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete user
func (repository *UserRepo) DeleteUser(c *gin.Context) {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.DeleteUser(repository.Db, &user, id)
	if err != nil {
		if err.Error() == "id does not exist" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
