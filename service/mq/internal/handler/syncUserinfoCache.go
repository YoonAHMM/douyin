package handler

import (
	"context"
	"strconv"

	"douyin/service/user/db/model"
	"douyin/service/user/user"
	"github.com/hibiken/asynq"
)



func (l *AsynqServer) syncUserInfoCacheHandler(ctx context.Context, t *asynq.Task) error {
	
	res, err := l.svcCtx.Redis.LRange(ctx, model.UserCacheKeyPrefix, 0, -1).Result()
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}

	for _, v := range res {
		userId, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return  err
		}
		// 读取缓存
		u:=model.User{}
		userInfo, err := l.svcCtx.Redis.HGetAll(ctx, u.CacheKey(userId)).Result()
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
		// 更新用户信息
		followcount,_:=strconv.ParseInt(userInfo["followCount"], 10, 64)
		followercount,_:=strconv.ParseInt(userInfo["followerCount"], 10, 64)
		_, err = l.svcCtx.UserRpc.UpdateUser(ctx, &user.UpdateUserReq{
			Id:          int64(userId),
			Name:        userInfo["Name"],
			Password:    userInfo["Password"],
			FollowCount: followcount,
			FanCount:    followercount,
		})
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
