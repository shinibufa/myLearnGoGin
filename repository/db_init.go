package repository

import (
	"github.com/xiongjsh/learn_tiktok_project/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB, err := gorm.Open(mysql.Open(config.GetDBConStr()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&UserLogin{}, &UserInfo{}, &Comment{}, &Video{})
	if err != nil {
		panic(err)
	}

}
