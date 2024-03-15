package logic

import (
	"context"

	"douyin/service/user/internal/config"
	"douyin/service/user/db/model"
	"douyin/service/user/internal/logic/utils"
	"douyin/service/user/internal/svc"
	"douyin/service/user/user"
 
	snowflake "github.com/ncghost1/snowflake-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//注册
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	username := in.Username
	password := in.Password

	//向数据库查询user是否存在
	exist := int64(0)
	err := l.svcCtx.Db.Model(&model.User{}).Where(&model.User{Username: username}).Count(&exist).Error
	if err != nil {
		return nil, err
	}

	if exist > 0 {
		//存在匹配用户
		return &user.RegisterResp{
			StatusCode: config.STATUS_FAIL,
			StatusMsg:  config.STATUS_USER_EXISTS_MSG,
			UserID:     0,
		}, nil
	}

	sf, err := snowflake.New(l.svcCtx.Config.WorkerId)
	if err != nil {
		return nil, err
	}
	//得到一个uid
	uuid, err := sf.Generate() 
	if err != nil {
		return nil, err
	}

	userInfo := &model.User{
		Id:       uuid,
		Username: username,
		Password: utils.BcryptHash(password),
	}

	// 向数据库写入
	err = l.svcCtx.Db.Create(&userInfo).Error
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		StatusCode: config.STATUS_SUCCESS,
		StatusMsg:  config.STATUS_SUCCESS_MSG,
		UserID:     userInfo.Id,
	}, nil
}
