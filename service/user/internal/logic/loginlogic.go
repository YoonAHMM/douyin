package logic

import (
	"context"

	"douyin/service/user/internal/config"
	"douyin/service/user/db/model"
	"douyin/service/user/internal/logic/utils"
	"douyin/service/user/internal/svc"
	"douyin/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//登入
func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	var userInfo *model.User	
	err := l.svcCtx.Db.Where(&model.User{Username: in.Username}).Take(&userInfo).Error
	
	if err != nil {
		//记录不存在
		if err == gorm.ErrRecordNotFound {
			return &user.LoginResp{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.STATUS_USER_NOTEXIST_MSG,
				UserID:     0,
			}, nil
		}
		return nil, err
	}

	//验证密码
	if ok := utils.BcryptCheck(in.Password, userInfo.Password); !ok {
		return &user.LoginResp{
			StatusCode: config.STATUS_FAIL,
			StatusMsg:  config.STATUS_WRONG_PASSWORD_MSG,
			UserID:     0,
		}, nil
	}
	
	return &user.LoginResp{
		StatusCode: config.STATUS_SUCCESS,
		StatusMsg:  config.STATUS_SUCCESS_MSG,
		UserID:     userInfo.Id,
	}, nil

}
