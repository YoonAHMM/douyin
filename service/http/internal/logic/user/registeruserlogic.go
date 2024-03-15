package user

import (
	"context"
	"strconv"

	"douyin/service/http/internal/config"
	"douyin/service/http/internal/svc"
	"douyin/service/http/internal/types"
	"douyin/service/jwt/Jwt"
	"douyin/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterUserLogic) RegisterUser(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	
	if len(req.Username) > 32 || len(req.Password) > 32 {
		return &types.RegisterResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.STATUS_FAIL_TOOLONG_MSG,
			},
			Token:  nil,
			UserID: 0,
		}, err
	}
	
	r, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	var token *string
	token = nil
	if r.StatusCode == config.STATUS_SUCCESS {
		tokenResp, err := l.svcCtx.JwtRpc.CreateToken(l.ctx, &Jwt.CreateTokenReq{
			UserID:       strconv.FormatUint(r.UserID, 10),
			AccessExpire: l.svcCtx.Config.JwtConfig.AccessExpire,
		})
		if err != nil {
			return nil, err
		}
		token = &tokenResp.Token
	}

	return &types.RegisterResp{
		Status: types.Status{
			StatusCode: r.StatusCode,
			StatusMsg:  r.StatusMsg,
		},
		Token:  token,
		UserID: r.UserID,
	}, nil

	
}
