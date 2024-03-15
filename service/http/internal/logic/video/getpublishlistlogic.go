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

type GetPublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPublishListLogic) GetPublishList(req *types.PublishListReq) (resp *types.PublishListResp, err error) {
	token, err := l.svcCtx.JwtRpc.ParseToken(l.ctx, &Jwt.ParseTokenReq{Token: req.Token})
	if err != nil {
		return &types.PublishListResp{
			Status:  types.Status{StatusCode: config.STATUS_FAIL, StatusMsg: config.STATUS_FAIL_TOKEN_MSG},
			VideoList: nil,
		}, nil
	}
	userid := token.UserID
	r, err := l.svcCtx.VideoRpc.GetPublishList(l.ctx, &videorpc.PublishListReq{UserID: userid, QueryId: req.UserId})
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

	return &types.PublishListResp{
		Status:  types.Status{StatusCode: r.StatusCode, StatusMsg: r.StatusMsg},
		VideoList: videoList,
	}, nil

	
}
