package userinfo

import (
	"errors"

	"github.com/xiongjsh/learn_tiktok_project/repository"
	"github.com/xiongjsh/learn_tiktok_project/cache"
)

const (
	FOLLOW = 1
	UNFOLLOW = 2
)

type PostFollowActionFlow struct {
	followerId int64
	followeeId int64
	actionType int
}

func PostFollowAction(followerId, followeeId int64, actionType int) error {
	return NewPostFollowActionFlow(followerId, followeeId, actionType).Do()
}

func NewPostFollowActionFlow(followerId, followeeId int64, actionType int) *PostFollowActionFlow {
	return &PostFollowActionFlow{
		followerId: followerId,
		followeeId: followeeId,
		actionType: actionType,
	}
}

func (f *PostFollowActionFlow) Do() error {
	if err := f.checkPara(); err != nil {
		return err
	}
	if err := f.publish(); err != nil {
		return err
	} 
	return nil
}

func (f *PostFollowActionFlow) checkPara() error {
	if !repository.NewUserInfoDao().CheckExistById(f.followeeId) {
		return errors.New("post_follow_action checkPara:followee doesn't exist")
	}
	if f.followerId == f.followeeId {
		return errors.New("you can not follow yourself")
	}
	return nil
}

func (f *PostFollowActionFlow) publish() error {
	userInfoDao := repository.NewUserInfoDao()
	var err error
	switch f.actionType {
	case FOLLOW:
		err = userInfoDao.AddFollow(f.followerId, f.followeeId)
		cache.NewProxyIndexMap().UpdateFolloweRelation(f.followerId, f.followeeId, true)

	case UNFOLLOW:
		err = userInfoDao.CancelFollow(f.followerId, f.followeeId)
		cache.NewProxyIndexMap().UpdateFolloweRelation(f.followerId, f.followeeId, false)
	
	default:
		return errors.New("unsupported actionType")
	}
	return err
}