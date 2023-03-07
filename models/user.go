package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	//gorm.Model// for create default field
	ID    int `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Name  string `json:"name,omitempty" gorm:"unique;notnull"`
	Email string `json:"email,omitempty" gorm:"unique;notnull"`
	Password string `json:"-" gorm:"not null"`
	Role string `json:"-" gorm:"notnull;default:user"`
	Posts []Post `gorm:"foreignKey:UserID"` 
	CreatedAt time.Time
}

//create a user
func CreateUser(db *gorm.DB, user *User) (err error) {
	err = db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

//get users
func GetUsers(db *gorm.DB, user *[]User) (err error) {
	err = db.Find(user).Error
	if err != nil {
		return err
	}
	return nil
}

//get user by id
func GetUser(db *gorm.DB, user *User, id int) (err error) {
	err = db.Where("id = ?", id).First(user).Error
	if err != nil {
		return err
	}
	return nil
}
//get user by email
func GetUserByEmail(db *gorm.DB, user *User, email string) (err error) {
	err = db.Where("email = ?", email).First(user).Error
	if err != nil {
		return err
	}
	return nil
}
//update user
func UpdateUser(db *gorm.DB, user *User, id int) (err error) {
	var oldData User
	err = GetUser(db,&oldData,id)
	if err != nil {
		return err
	}

	
	if user.Email != "" {
		oldData.Email = user.Email
	}
	if user.Name != "" {
		oldData.Name = user.Name
	} 
	// oldData.Name = user.Name
	// oldData.Email = user.Email
	//err = db.Save(&oldData).Error
	err = db.Where("id = ?", id).Updates(&oldData).Error
	if err != nil {
		return err
	}
	*user = oldData
	return nil
}

//delete user
func DeleteUser(db *gorm.DB, user *User, id int) (err error) {
	result := db.Where("id = ?", id).Delete(user)
	err = result.Error
	if err != nil {
		return err
	}else if result.RowsAffected < 1 {
        return errors.New("id does not exist")
    }
	return nil
}