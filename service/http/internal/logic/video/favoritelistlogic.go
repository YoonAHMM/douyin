package video

import (
	"context"

	"douyin/service/http/internal/config"
	"douyin/service/http/internal/svc"
	"douyin/service/http/internal/types"
	"douyin/service/jwt/Jwt"
	"douyin/service/video/rpc/videorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListReq) (resp *types.FavoriteListResp, err error) {
	token, err := l.svcCtx.JwtRpc.ParseToken(l.ctx, &Jwt.ParseTokenReq{Token: req.Token})
	if err != nil {
		return &types.FavoriteListResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.STATUS_FAIL_TOKEN_MSG,
			},
			VideoList: nil,
		}, nil
	}
	userid := token.UserID

	r, err := l.svcCtx.VideoRpc.GetFavoriteList(l.ctx, &videorpc.FavoriteListReq{UserId: userid, QueryUserId: req.UserId})
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
	return &types.FavoriteListResp{
		Status: types.Status{
			StatusCode: r.StatusCode,
			StatusMsg:  r.StatusMsg,
		},
		VideoList: videoList,
	}, nil

	
}
