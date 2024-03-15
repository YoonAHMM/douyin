package user

import (
	"context"

	"douyin/service/http/internal/config"
	"douyin/service/http/internal/svc"
	"douyin/service/http/internal/types"
	"douyin/service/jwt/Jwt"
	"douyin/service/user/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionReq) (resp *types.RelationActionResp, err error) {
	token, err := l.svcCtx.JwtRpc.ParseToken(l.ctx, &Jwt.ParseTokenReq{Token: req.Token})
	if err != nil {
		return &types.RelationActionResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.STATUS_FAIL_TOKEN_MSG,
			},
		}, nil
	}
	userid := token.UserID

	if userid == req.ToUserId {
		return &types.RelationActionResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.STATUS_FAIL_FOLLOW_SELF,
			},
		}, nil
	}

	r, err := l.svcCtx.UserRpc.FollowAction(l.ctx, &userrpc.FollowActionReq{UserId: userid, ToUserId: req.ToUserId, ActionType: req.ActionType})
	if err != nil {
		return nil, err
	}

	return &types.RelationActionResp{Status: types.Status{
		StatusCode: r.StatusCode,
		StatusMsg:  r.StatusMsg,
	}}, nil

	
}
