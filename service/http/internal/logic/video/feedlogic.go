package video

import (
	"context"
	"strconv"
	"time"

	"douyin/service/http/internal/config"
	"douyin/service/http/internal/svc"
	"douyin/service/http/internal/types"
	"douyin/service/jwt/Jwt"
	"douyin/service/video/rpc/videorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//不限制登录状态
func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedResp, err error) {
	var userid *string = new(string)

	if req.Token == nil {
		*userid = config.USER_NO_LOGIN
	} else {
		token, err := l.svcCtx.JwtRpc.ParseToken(l.ctx, &Jwt.ParseTokenReq{Token: *req.Token})
		if err != nil {
			return &types.FeedResp{
				Status:  types.Status{StatusCode: config.STATUS_FAIL, StatusMsg: config.STATUS_FAIL_TOKEN_MSG},
				NextTime:  nil,
				VideoList: nil,
			}, nil
		}
		userid = &token.UserID
	}

	latest := req.LatestTime
	latestTs := time.Now().Unix()

	if latest != nil {
		latestTs, err = strconv.ParseInt(*latest, 10, 64)
		if err != nil {
			return &types.FeedResp{
					Status: types.Status{
						StatusCode: config.STATUS_FAIL,
						StatusMsg:  config.STATUS_FAIL_PARAM_MSG,
				},
				NextTime:  nil,
				VideoList: nil,
			}, nil
		}
	}

	r, err := l.svcCtx.VideoRpc.GetFeed(l.ctx, &videorpc.FeedReq{
		UserId:     *userid,
		LatestTime: latestTs,
		Limit:      l.svcCtx.Config.FeedLimit,
	})
	if err != nil {
		return nil, err
	}

	videoList := make([]types.Video, len(r.VideoList))
	for i, v := range r.VideoList {
		videoList[i] = types.Video{
			Author: types.User{
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
				ID:            v.Author.ID,
				IsFollow:      v.Author.IsFollow,
				Name:          v.Author.Name,
			},
			CommentCount:  v.CommentCount,
			CoverURL:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			ID:            v.ID,
			IsFavorite:    v.IsFavorite,
			PlayURL:       v.PlayURL,
			Title:         v.Title,
		}
	}

	return &types.FeedResp{
		Status: types.Status{
			StatusCode: r.StatusCode,
			StatusMsg:  r.StatusMsg,
		},
		NextTime:  &r.NextTime,
		VideoList: videoList,
	}, nil
	
}
