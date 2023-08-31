package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/xiongjsh/learn_tiktok_project/config"
	
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.ConfigInfo.RDB.IP, config.ConfigInfo.RDB.Port),
		Password: "",
		DB: config.ConfigInfo.RDB.Database,
	})
}

type ProxyIndexMap struct {}

var proxyIndexOperation ProxyIndexMap

func NewProxyIndexMap() *ProxyIndexMap {
	return &proxyIndexOperation
}

func (p *ProxyIndexMap) UpdateUser2Video(userId, videoId int64, isLike bool) {
	key := fmt.Sprintf("%d %s", userId, "like")
	if isLike {
		rdb.SAdd(ctx, key, videoId)
		return
	}
	rdb.SRem(ctx, key, videoId)
}

func (p *ProxyIndexMap) IsUserLikeVideo(userId, videoId int64) bool {
	key := fmt.Sprintf("%d %s", userId, "like")
	rbc := rdb.SIsMember(ctx, key, videoId)
	return rbc.Val()
}

func (p *ProxyIndexMap) UpdateFolloweRelation(followerId, followeeId int64, isFollow bool) {
	key := fmt.Sprintf("%d %s", followerId, "follow")
	if isFollow {
		rdb.SAdd(ctx, key, followeeId)
		return
	}
	rdb.SRem(ctx, key, followeeId)
}

func (p *ProxyIndexMap) GetFollowRelation(followerId, followeeId int64) bool {
	key := fmt.Sprintf("%d %s", followerId, "follow")
	rbc := rdb.SIsMember(ctx, key, followeeId)
	return rbc.Val()
}