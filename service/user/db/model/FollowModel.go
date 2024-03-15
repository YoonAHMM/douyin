package model

import "strconv"

type Follow struct {
	Follower   uint64 `gorm:"column:follower_id"`
	Following  uint64 `gorm:"column:following_id"`
	CreateTime int64  `gorm:"column:create_time"`
}

const (
	FollowListCacheKeyPrefix   = "FollowList:Userid:FollowId:"
	FollowerListCacheKeyPrefix = "FollowerList:Userid:FollowId:"
)


func (Follow) TableName() string {
	return "follow"
}

//关注
func (Follow) FollowListCacheKey(userid uint64) string {
	return FollowListCacheKeyPrefix + strconv.FormatUint(userid, 10)
}

//粉丝
func (Follow) FollowerListCacheKey(userid uint64) string {
	return FollowerListCacheKeyPrefix + strconv.FormatUint(userid, 10)
}
