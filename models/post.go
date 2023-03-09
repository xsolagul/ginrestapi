package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	//gorm.Model// for create default field
	ID        int    `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Sentence  string `json:"sentence,omitempty" gorm:"notnull"`
	UserID    uint
	User      User
	CreatedAt time.Time
}

// create a post
func CreatePost(db *gorm.DB, post *Post) (err error) {
	err = db.Create(post).Error
	if err != nil {
		return err
	}
	err = db.Preload("User").First(post).Error
	if err != nil {
		return err
	}
	return nil
}

// get post
func GetPosts(db *gorm.DB, post *[]Post) (err error) {
	err = db.Preload("User",func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name")}).Find(post).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPagePosts(db *gorm.DB, post *[]Post,page int,pageSize int)(int64, error){
	var count int64 
	offset := (page - 1)*pageSize

	if err := db.Model(&Post{}).Count(&count).Error; err != nil {
        return  0, err
    }

    // retrieve users with pagination
    if err := db.Limit(pageSize).Offset(offset).Find(post).Error; err != nil {
        return  0, err
    }
	return count,nil
}

// get post by id
func GetPost(db *gorm.DB, post *Post, id int) (err error) {
	err = db.Preload("User",func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name")}).Where("id = ?", id).First(post).Error
	if err != nil {
		return err
	}
	return nil
}

// update post
func UpdatePost(db *gorm.DB, post *Post, id int) (err error) {
	var oldData Post
	err = GetPost(db, &oldData, id)
	fmt.Println(oldData)
	if err != nil {
		return err
	}

	if post.Sentence != "" {
		oldData.Sentence = post.Sentence
	}
	err = db.Where("id = ?", id).Updates(&oldData).Error
	if err != nil {
		return err
	}
	*post = oldData
	return nil
}

// delete post
func DeletePost(db *gorm.DB, post *Post, id int) (err error) {
	result := db.Where("id = ?", id).Delete(post)
	err = result.Error
	if err != nil {
		return err
	} else if result.RowsAffected < 1 {
		return errors.New("id does not exist")
	}
	return nil
}
