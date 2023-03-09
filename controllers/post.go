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
	"gorm.io/gorm"
)

type PostRepo struct {
	Db *gorm.DB
}


func NewPostRepo(db *gorm.DB) *PostRepo {
	//db := database.InitDb()//guess do one time
	db.AutoMigrate(&models.Post{})
	return &PostRepo{Db: db}
}

// create user
func (repository *PostRepo) CreatePost(c *gin.Context) {
	var post models.Post
	err := c.BindJSON(&post)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	value := c.MustGet("userId")
	post.UserID = uint(value.(int))
	err = models.CreatePost(repository.Db, &post)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, post)
}

// get users
func (repository *PostRepo) GetPosts(c *gin.Context) {
	var posts []models.Post
	err := models.GetPosts(repository.Db, &posts)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// get user by id
func (repository *PostRepo) GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var post models.Post
	err := models.GetPost(repository.Db, &post, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, post)
}

// update user
func (repository *PostRepo) UpdatePost(c *gin.Context) {
	var post models.Post
	id, _ := strconv.Atoi(c.Param("id"))

	c.BindJSON(&post)
	post.ID = id
	errC := checkPostUserId(c,repository,id)
	if errC != nil {
		if errC.Error() == "userId not match" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errC.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errC})
		return
	}
	err := models.UpdatePost(repository.Db, &post , id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, post)
}

// delete user
func (repository *PostRepo) DeletePost(c *gin.Context) {
	var post models.Post
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.GetPost(repository.Db, &post, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	errC := checkPostUserId(c,repository,id)
	if errC != nil {
		if errC.Error() == "userId not match" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errC.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errC})
		return
	}
	err = models.DeletePost(repository.Db, &post, id)
	if err != nil {
		if err.Error() == "id does not exist" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
func checkPostUserId(c *gin.Context,repository *PostRepo,id int ) error {
	var post models.Post
	err := models.GetPost(repository.Db, &post, id)
	if err != nil {
		return err 
	}
	//change to return error
	value := c.MustGet("userId")
	userId := uint(value.(int))
	if userId != post.UserID {
		return errors.New("userId not match")
	}
	return nil
}