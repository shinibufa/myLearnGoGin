package repository

import (
	"errors"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	Id           int64       `json:"id,omitempty"`
	UserInfoId   int64       `json:"-"`
	Author       *UserInfo   `json:"author,omitempty" gorm:"-"`
	PlayUrl      string      `json:"play_url,omitempty"`
	CoverUrl     string      `json:"cover_url,omitempty"`
	LikeCount    int64       `json:"like_count,omitempty"`
	CommentCount int64       `json:"comment_count,omitempty"`
	IsLiked      bool        `json:"is_liked,omitempty"`
	Title        string      `json:"title,omitempty"`
	Users        []*UserInfo `json:"-" gorm:"many2many:user_like_video"`
	Comments     []*Comment  `json:"-"`
	CreateTime   time.Time   `json:"-"`
	UpdateTime   time.Time   `json:"-"`
}

type VideoDao struct{}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDao() *VideoDao {
	videoOnce.Do(func() {
		videoDao = &VideoDao{}
	})
	return videoDao
}

func (d *VideoDao) AddVideo(video *Video) error {
	if video == nil {
		return errors.New("video is a nil pointer")
	}

	if err := DB.Create(video).Error; err != nil {
		return err
	}
	return nil
}

func (d *VideoDao) QueryVideoById(videoId int64, video *Video) error {
	if video == nil {
		return errors.New("QueryVideoById: video is a nil pointer")
	}
	if err := DB.Where("id = ?", videoId).Select([]string{"id", "user_info_id", "play_url", "cover_url", "like_count", "is_liked", "title", "comment_count"}).First(video).Error; err != nil {
		return err
	}
	return nil
}

func (d *VideoDao) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId: count is a nil pointer")
	}
	return DB.Model(&Video{}).Where("user_info_id = ?", userId).Count(count).Error
}

func (d *VideoDao) QueryVideosByLimitAndTime(limit int64, latestTime time.Time, videos []*Video) error {
	if videos == nil {
		return errors.New("QueryVideosByLimitAndTime: videos is a nil pointer")
	}
	return DB.Model(&Video{}).
		Where("create_time < ?", latestTime).
		Order("create_time ASC").
		Select([]string{"id", "user_info_id", "title", "create_time", "update_time", "play_url", "cover_url", "like_count", "comment_count", "is_liked"}).
		Find(videos).Error
}

func (d *VideoDao) UserLikeVideo(userId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE video SET like_count = like_count+1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `user_like_video` (`user`, `video`) VALUES (?,?)", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *VideoDao) UserDislikeVideo(userId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE video SET like_count = lilke_count-1 WHERE like_count > 0 AND id = ?", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_like_video` WHERE `user` = ? AND `video` = ?", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *VideoDao) QueryLikeVideosByUserId(userId int64, videos []*Video) error {
	if videos == nil {
		return errors.New("QueryLikeVideosByUserId: videos is a nil pointer")
	}
	if err := DB.Raw("SELECT v.* FROM user_like_video u, video v WHERE u.user = ? AND v.id = u.video").Scan(videos).Error; err != nil {
		return err
	}
	return nil
}

func (d *VideoDao) IsVideoExist(videoId int64) bool {
	var video Video
	if err := DB.Where("id = ?", videoId).Select("id").First(&video).Error; err != nil {
		log.Println(err)
	}
	if video.Id == 0 {
		return false
	}
	return true
}
