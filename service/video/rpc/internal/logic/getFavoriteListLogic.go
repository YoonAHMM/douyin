package logic

import (
	"context"
	"douyin/service/user/userrpc"
	"douyin/service/video/rpc/model"
	"douyin/service/video/rpc/internal/svc"
	"douyin/service/video/rpc/video"
	

	"encoding/json"
	"errors"
	"strconv"
	"time"


	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFavoriteList 获取用户点赞过的视频列表
func (l *GetFavoriteListLogic) GetFavoriteList(in *video.FavoriteListReq) (*video.FavoriteListResp, error) {
	userid, err := strconv.ParseUint(in.UserId, 10, 64)
	if err != nil {
		return nil, err
	}
	queryId, err := strconv.ParseUint(in.QueryUserId, 10, 64)
	if err != nil {
		return nil, err
	}

	conn := l.svcCtx.Redis.NewRedisConn()
	defer conn.Close()

	var modelVideoList []*model.Video
	latestTime := time.Now().Unix()
	list, exists, err := l.svcCtx.Redis.GetExFavoriteList(conn, queryId, l.svcCtx.Config.CacheConfig.FAVORITE_CACHE_TTL)
	if err != nil {
		return nil, err
	}

	if exists {
		for i, v := range list {
			if i == len(list)-1 {
				latestTime, err = strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					return nil, err
				}
				break
			}
			var vid *model.Video
			err = json.Unmarshal(v, &vid)
			if err != nil {
				return nil, err
			}
			modelVideoList = append(modelVideoList, vid)
		}
	}

	// 如果不存在该用户的缓存，或者缓存列表已满，则需要到数据库中查找该用户是否还有点赞的视频
	remain := l.svcCtx.Config.CacheConfig.VIDEO_FAVORITE_MAX_CACHE_SIZE - len(modelVideoList)
	var favoriteList []model.Favorite
	if !exists || remain == l.svcCtx.Config.CacheConfig.VIDEO_FAVORITE_MAX_CACHE_SIZE {
		err = l.svcCtx.Db.Where("user_id = ? and create_time <= ?", queryId, latestTime).Find(&favoriteList).Error
		if err != nil {
			return nil, err
		}

		for i, v := range favoriteList {
			if i <= remain {
				err = l.svcCtx.Redis.SendAddFavorList(conn, queryId, v.VideoId, v.CreateTime, l.svcCtx.Config.CacheConfig.FAVORITE_CACHE_TTL)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	for _, v := range favoriteList {

		// 使用视频 id 查视频信息（先查缓存再查数据库）
		info, exists, err := l.svcCtx.Redis.GetExVideoInfo(conn, v.VideoId, l.svcCtx.Config.CacheConfig.VIDEO_CACHE_TTL)
		var videoInfo *model.Video
		if err != nil {
			return nil, err
		}

		if exists {
			// 缓存查找成功
			err = json.Unmarshal(info, &videoInfo)
			if err != nil {
				return nil, err
			}

		} else { 
			// 缓存查找失败，则需要到数据库查询
			err = l.svcCtx.Db.Where(&model.Video{Id: v.VideoId}).Take(&videoInfo).Error
			if err != nil {
				return nil, err
			}
		}
	}

	videoList := make([]*video.Video, len(modelVideoList))

	for i, v := range modelVideoList {
		// 使用视频作者 id 查作者信息
		var userInfo *video.User
		r, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userrpc.GetUserReq{
			UserID:  in.QueryUserId,
			QueryID: strconv.FormatUint(v.UserId, 10),
		})
		if err != nil {
			return nil, err
		}

	
		userInfo = &video.User{
			FollowCount:   r.User.FollowCount,
			FollowerCount: r.User.FollowerCount,
			ID:            r.User.ID,
			IsFollow:      r.User.IsFollow,
			Name:          r.User.Name,
		}

		if r.StatusCode == STATUS_FAIL {
			return nil, errors.New(r.StatusMsg)
		}


		// 一次性获取点赞数，评论数，用户是否点赞过视频的缓存数据
		favCount, comCount, isFavor, err := l.svcCtx.Redis.GetExFavComCountIsFavor(conn, v.Id, userid,
			l.svcCtx.Config.CacheConfig.FAVORITE_CACHE_TTL,
			l.svcCtx.Config.CacheConfig.COMMENT_CACHE_TTL,
		)
		if err != nil {
			return nil, err
		}

		// 如果缓存未找到，还需要查库
		if favCount == COUNT_NOT_FOUND {
			err := l.svcCtx.Db.Model(&model.Favorite{}).Where(&model.Favorite{VideoId: v.Id}).Count(&favCount).Error
			if err != nil {
				return nil, err
			}
			err = l.svcCtx.Redis.SendSetExFavorCount(conn, v.Id, favCount, l.svcCtx.Config.CacheConfig.FAVORITE_CACHE_TTL)
			if err != nil {
				return nil, err
			}
		}

		if comCount == COUNT_NOT_FOUND {
			err := l.svcCtx.Db.Model(&model.Comment{}).Where(&model.Comment{VideoId: v.Id}).Count(&comCount).Error
			if err != nil {
				return nil, err
			}
			err = l.svcCtx.Redis.SendSetExCommentCount(conn, v.Id, comCount, l.svcCtx.Config.CacheConfig.COMMENT_CACHE_TTL)
			if err != nil {
				return nil, err
			}
		}

		if !isFavor {
			var cnt int64
			err = l.svcCtx.Db.Model(&model.Favorite{}).Where(&model.Favorite{UserId: userid, VideoId: v.Id}).Count(&cnt).Error
			if err != nil {
				return nil, err
			}
			if cnt > 0 {
				isFavor = true
			}
		}

		vid := &video.Video{
			Author: &video.User{
				FollowCount:   userInfo.FollowCount,
				FollowerCount: userInfo.FollowerCount,
				ID:            userInfo.ID,
				IsFollow:      userInfo.IsFollow,
				Name:          userInfo.Name,
			},
			CommentCount:  comCount,
			CoverURL:      v.CoverUrl,
			FavoriteCount: favCount,
			ID:            v.Id,
			IsFavorite:    isFavor,
			PlayURL:       v.PlayUrl,
			Title:         v.Title,
		}
		videoList[i] = vid
	}

	err = conn.Flush()
	if err != nil {
		return nil, err
	}

	return &video.FavoriteListResp{
		StatusCode: STATUS_SUCCESS,
		StatusMsg:  STATUS_SUCCESS_MSG,
		VideoList:  videoList,
	}, nil
}
