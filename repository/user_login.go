package repository

import (
	"errors"
	"log"
	"sync"
)

type UserLogin struct {
	Id         int64  `gorm:"primary_key"`
	UserInfoId int64  `gorm:"colunm:userinfo_id"`
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"csize:200;notnull"`
}

type UserLoginDao struct {
}

var userLoginDao *UserLoginDao
var userLoginOnce sync.Once

func NewUserLoginDao() *UserLoginDao {
	userLoginOnce.Do(func() {
		userLoginDao = &UserLoginDao{}
	})
	return userLoginDao
}

func (d *UserLoginDao) CheckExistByName(username string) bool {
	var userLogin UserLogin
	err := DB.Where("username = ?", username).Find(&userLogin).Error
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	if userLogin.Id == 0 {
		return false
	}
	return true
}

func (d *UserLoginDao) QueryUserLogin(username, password string, userLogin *UserLogin) error {
	if userLogin == nil {
		return errors.New("userLogin is nil")
	}
	err := DB.Where("username = ?", username).First(userLogin).Error
	if err != nil {
		return err
	}
	if userLogin.Id == 0 {
		return errors.New("Your username doesn't exist.")
	}
	err = DB.Where("username = ? AND password = ?", username, password).First(userLogin).Error
	if err != nil {
		return err
	}
	if userLogin.Id == 0 {
		return errors.New("Your oassword is wrong.")
	}
	return nil

}
