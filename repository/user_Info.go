package repository

import (
	"errors"
	"log"
	"sync"

	"gorm.io/gorm"
)

type UserInfo struct {
	Id        int64      `json:"id" gorm:"primary_key"`
	Name      string     `json:"name" gorm:"name"`
	UserLogin *UserLogin `gorm:"column:user_login"`
	IsFollow  bool       `json:"is_follow" gorm:"is_follow"`
	//Followed []*UserInfo `gorm:"column:followed"`
	Following     []*UserInfo `gorm:"column:fllowing;many2many:follower_2_followee" json:"-"`
	LikeVideos    []*Video    `gorm:"colunm:like_vedios;many2many:user_like_video" json:"-"`
	OwnVideos     []*Video    `gorm:"column:own_vedios" json:"own_videos" json:"-"`
	FollowerCount int64       `gorm:"column:follower_count,omitempty" json:"follower_count"`
	FolloweeCount int64       `gorm:"colunm:followcount_count,omitempty" json:"followee_count"`
}

type UserInfoDao struct{}

var userInfoDao *UserInfoDao
var userInfoOnce sync.Once

func NewUserInfoDao() *UserInfoDao {
	userInfoOnce.Do(func() {
		userInfoDao = &UserInfoDao{}
	})
	return userInfoDao
}

func (d *UserInfoDao) AddUserInfo(userInfo *UserInfo) error {
	if userInfo == nil {
		return errors.New("AddUserInfo: userInfo is a nil pointer")
	}
	return DB.Create(userInfo).Error
}

func (d *UserInfoDao) CheckExistById(userInfoId int64) bool {
	var userInfo UserInfo
	err := DB.Where("id = ?", userInfoId).Select("id").First(&userInfo).Error
	if err != nil {
		log.Println(err)
	}
	if userInfo.Id == 0 {
		return false
	}
	return true
}

func (d *UserInfoDao) QueryUserInfoById(userInfoId int64, userInfo *UserInfo) error {
	if userInfo == nil {
		return errors.New("userInfo is a nil pointer")
	}
	DB.Where("id = ?", userInfoId).Select([]string{"id", "name", "is_follow", "follower_count", "followee_count"}).First(userInfo)
	if userInfo.Id == 0 {
		return errors.New("user doesn't exist")
	}
	return nil
}

func (d *UserInfoDao) AddFollow(followerId, followeeId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_info SET followee_count = followee_count+1 WHERE id = ?", followerId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_info SET follower_count = follower_count+1 WHERE id = ?", followeeId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `follower_2_followee` (`follower`, `followee`) VALUES (?,?)", followerId, followeeId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *UserInfoDao) CancelFollow(followerId, followeeId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_info SET followee_count = followee_count-1 WHERE id = ?", followerId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_info SET follower_count = follower_count-1 WHERE id = ?", followeeId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_info` WHERE follower = ? AND followee = ?", followerId, followeeId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *UserInfoDao) GetFollowersById(followeeId int64, userList []*UserInfo) error {
	if userList == nil {
		return errors.New("userList is a nil pointer")
	}
	if err := DB.Raw("SELECT u.* FROM follower_2_followee f, user_info u WHERE f.follower_id = ? AND u.id = f.followee_id", followeeId).Scan(userList).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserInfoDao) GetFolloweesById(followerId int64, userList []*UserInfo) error {
	if userList == nil {
		return errors.New("userList is a nil pointer")
	}
	if err := DB.Raw("SELECT u.* FROM follower_2_followee f, user_info u WHERE f.followee_id = ? AND u.id = f.follower_id").Scan(userList).Error; err != nil {
		return err
	}
	return nil
}
