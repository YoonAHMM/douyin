package logic

import (
	"context"
	"strconv"

	"douyin/service/user/internal/config"
	"douyin/service/user/db/RedisCache"
	"douyin/service/user/db/model"
	"douyin/service/user/internal/svc"
	"douyin/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//user 查询 query 用户信息
func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	userid, err := strconv.ParseUint(in.UserID, 10, 64)
	if err != nil {
		return &user.GetUserResp{
			StatusCode: config.STATUS_FAIL,
			StatusMsg:  config.STATUS_FAIL_PARAM_MSG,
			User:       nil,
		}, nil
	}
	queryid, err := strconv.ParseUint(in.QueryID, 10, 64)

	if err != nil {
		return &user.GetUserResp{
			StatusCode: config.STATUS_FAIL,
			StatusMsg:  config.STATUS_FAIL_PARAM_MSG,
			User:       nil,
		}, nil
	}

	conn := l.svcCtx.Redis.GetRedisConn()
	defer conn.Close()

	// 查询用户是否关注查询对象
	isfollow := false
	isfollow, err = l.svcCtx.Redis.IsFollow(conn, userid, queryid)
	if err != nil {
		return nil, err
	}

	//redis中找不到关注信息，无法确定是否关注
	if !isfollow {
		count := int64(0)
		err := l.svcCtx.Db.Model(&model.Follow{}).Where(&model.Follow{Follower: userid, Following: queryid}).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if count > 0 {
			isfollow = true 
		}
	}

	// 从缓存中获取查询对象用户名，关注数，粉丝数
	username, followCnt, followerCnt, err := l.svcCtx.Redis.GetUserInfo(conn, queryid)
	needUpdateCache := false // 是否需要更新缓存

	if err != nil {
		if err.Error() != RedisCache.CACHE_KEY_NOT_EXISTS_MSG {
			return nil, err
		}
		needUpdateCache  = true
	}

	//缓存找不到到DB找
	if username == "" {
		err = l.svcCtx.Db.Model(model.User{}).Select("username").Where(&model.User{Id: queryid}).Take(&username).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // 用户不存在
				return &user.GetUserResp{
					StatusCode: config.STATUS_FAIL,
					StatusMsg:  config.STATUS_USER_NOTEXIST_MSG,
					User:       nil,
				}, nil
			}
			return nil, err
		}
	}

	if followCnt == config.COUNT_NOT_FOUND {
		followCnt = 0
		err = l.svcCtx.Db.Model(model.User{}).Select("FollowCount").Where(&model.User{Id: queryid}).Take(&followCnt).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
	}

	if followerCnt == config.COUNT_NOT_FOUND {
		followerCnt = 0
		err = l.svcCtx.Db.Model(model.User{}).Select("FanCount").Where(&model.User{Id: queryid}).Take(&followerCnt).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}

	}


	// 将从 DB 获取到的数据写入缓存
	if needUpdateCache {
		err = l.svcCtx.Redis.SetUserInfo(conn, queryid, username, followCnt, followerCnt)
		if err != nil {
			return nil, err
		}
	}

	UserInfo := &user.User{
		FollowCount:   followCnt,
		FollowerCount: followerCnt,
		ID:            queryid,
		IsFollow:      isfollow,
		Name:          username,
	}

	return &user.GetUserResp{
		StatusCode: config.STATUS_SUCCESS,
		StatusMsg:  config.STATUS_SUCCESS_MSG,
		User:       UserInfo,
	}, nil

}
