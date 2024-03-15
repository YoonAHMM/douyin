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

type FollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowerListLogic) FollowerList(in *user.FollowerListReq) (*user.FollowerListResp, error) {
	userid, err := strconv.ParseUint(in.UserId, 10, 64)
	if err != nil {
		return nil, err
	}

	toUserId, err := strconv.ParseUint(in.ToUserId, 10, 64)
	if err != nil {
		return nil, err
	}

	conn := l.svcCtx.Redis.GetRedisConn()
	defer conn.Close()

	//缓存得到关注列表
	cacheList, err := l.svcCtx.Redis.GetFollowUserList(conn, toUserId)
	if err != nil {
		if err.Error() != RedisCache.CACHE_KEY_NOT_EXISTS_MSG {
			return nil, err
		}
	}

	remain := l.svcCtx.Config.CacheConfig.FOLLOWERLIST_MAX_CACHE_SIZE - len(cacheList)
	var DbFollowerList []model.Follow

	//若缓存为空或缓存已满则需要到 DB 查询数据
	if remain == 0 || remain == l.svcCtx.Config.CacheConfig.FOLLOWERLIST_MAX_CACHE_SIZE {
		err = l.svcCtx.Db.Find(&DbFollowerList, &model.Follow{Following: toUserId}).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
	}

	var followerList []*user.User

	// 缓存用户信息
	for _, v := range cacheList {
		id, err := strconv.ParseUint(v[0], 10, 64)
		if err != nil {
			return nil, err
		}

		username := v[1]

		followCount, err := strconv.ParseInt(v[2], 10, 64)
		if err != nil {
			return nil, err
		}

		followerCount, err := strconv.ParseInt(v[3], 10, 64)
		if err != nil {
			return nil, err
		}

		isfollow := false
		count := int64(0)
		db := l.svcCtx.Db.Model(&model.Follow{}).Where(model.Follow{Follower: userid, Following: id}).Count(&count)
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			return nil, db.Error
		}
		
		if count > 0 {
			isfollow = true 
		}

		userInfo := &user.User{
			FollowCount:   followCount,
			FollowerCount: followerCount,
			ID:            id,
			IsFollow:      isfollow,
			Name:          username,
		}

		followerList = append(followerList, userInfo)
	}

	// 数据库用户信息
	for _, v := range DbFollowerList {
		id := v.Follower

		getUserLogic := NewGetUserLogic(l.ctx, l.svcCtx)
		u, err := getUserLogic.GetUser(&user.GetUserReq{
			UserID:  strconv.FormatUint(userid, 10),
			QueryID: strconv.FormatUint(id, 10),
		})
		if err != nil {
			return nil, err
		}

		userInfo := &user.User{
			FollowCount:   u.User.FollowCount,
			FollowerCount: u.User.FollowerCount,
			ID:            u.User.ID,
			IsFollow:      u.User.IsFollow,
			Name:          u.User.Name,
		}

		followerList = append(followerList, userInfo)
	}

	return &user.FollowerListResp{
		StatusCode: config.STATUS_SUCCESS,
		StatusMsg:  config.STATUS_SUCCESS_MSG,
		UserList:   followerList,
	}, nil
}
