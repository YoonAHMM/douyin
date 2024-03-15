package model

import "strconv"

// User 表结构
type User struct {
	Id       uint64 `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	FollowCount int64 `gorm:"column:FollowCount"`
	FanCount int64    `gorm:"column:FanCount"`
}

const (
	UserCacheKeyPrefix = "User:Userid:UserInfo:Hash"
	UsernameField      = "username"
	FollowCountField   = "followCount"
	FollowerCountField = "followerCount"
)

func (User) TableName() string {
	return "user"
}


func (User) CacheKey(userid uint64) string {
	return UserCacheKeyPrefix + strconv.FormatUint(userid, 10)
}

const (
	PopularUserStandard = 1000 // 拥有超过 1000 个粉丝的用户成为大V，有特殊处理
)

func IsPopularUser(fanCount int64) bool {
	return fanCount >= PopularUserStandard
}