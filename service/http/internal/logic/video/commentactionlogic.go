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

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionReq) (resp *types.CommentActionResp, err error) {
	token, err := l.svcCtx.JwtRpc.ParseToken(l.ctx, &Jwt.ParseTokenReq{Token: req.Token})
	if err != nil {
		return &types.CommentActionResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.STATUS_FAIL_TOKEN_MSG,
			},
			Comment: nil,
		}, nil
	}
	userid := token.UserID
	content := ""
	commentId := ""
	if req.CommentText != nil {
		content = *req.CommentText
	}
	if req.CommentId != nil {
		commentId = *req.CommentId
	}

	r, err := l.svcCtx.VideoRpc.CommentAction(l.ctx, &videorpc.CommentReq{
		VideoId:    req.VideoId,
		UserId:     userid,
		ActionType: req.ActionType,
		Content:    content,
		CommentId:  commentId,
	})
	if err != nil {
		return nil, err
	}

	var comment *types.Comment
	comment = nil
	if r.Comment != nil {
		comment = &types.Comment{
			Content:    r.Comment.Content,
			CreateDate: r.Comment.CreateDate,
			ID:         r.Comment.ID,
			User: types.User{
				FollowCount:   r.Comment.User.FollowCount,
				FollowerCount: r.Comment.User.FollowerCount,
				ID:            r.Comment.User.ID,
				IsFollow:      r.Comment.User.IsFollow,
				Name:          r.Comment.User.Name,
			},
		}
	}

	return &types.CommentActionResp{
		Status: types.Status{
			StatusCode: config.STATUS_SUCCESS,
			StatusMsg:  config.STATUS_SUCCESS_MSG,
		},
		Comment: comment,
	}, nil

}
