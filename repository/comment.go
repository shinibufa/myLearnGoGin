package repository

import (
	"errors"
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	Id int64 `json:"-"`
	UserInfoId int64 `json:"-"`
	VideoId int64 `json:"-"`
	User *UserInfo `json:"user" gorm:"-"`
	Content string `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

type CommentDao struct {}

var commentDao *CommentDao


func NewCommentDao() *CommentDao {
	commentDao = &CommentDao{}
	return commentDao
}

func (d *CommentDao) AddComment(comment *Comment) error {
	if comment == nil {
		return errors.New("comment pointer is nil")
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE video v SET v.comment_count = v.comment+1 WHERE v.id = ?", comment.VideoId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *CommentDao) DeleteCommentById(commentId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE comment WHERE id = ?", commentId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE video v SET v.comment_count = v.comment_count-1 WHERE v.id = ? AND v.comment_count > 0", videoId).Error; err != nil {
			return err
		}
		return nil

	})
}

func (d *CommentDao) QueryCommentById(commentId int64, comment *Comment) error {
	if comment == nil {
		return errors.New("comment is a nil pointer")
	}
	if err := DB.Where("id = ?", commentId).First(comment).Error; err != nil {
		return err
	}
	return nil
}

func QueryCommentsByVideoId(videoId int64, comments []*Comment) error {
	if comments == nil {
		return errors.New("comments is a nil pointer")
	}
	if err := DB.Model(&Comment{}).Where("video_id = ?", videoId).Find(comments).Error; err != nil {
		return err
	}
	return nil
}