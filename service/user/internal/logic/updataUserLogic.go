package logic

import (
	"context"
	"douyin/service/user/db/model"
	"douyin/service/user/internal/config"
	"douyin/service/user/internal/svc"
	"douyin/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
)



type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//向数据库更新缓存中user的关注数和粉丝数
func(l *UpdateUserLogic)UpdateUser(in *user.UpdateUserReq)(*user.UpdateUserResp,error){
	//开启事务
	tx := l.svcCtx.Db.Begin()

	var newUser *model.User
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.Id).First(&newUser).Error
	if err != nil {
		tx.Rollback()//失败回滚
		return nil, err
	}

	//更新
	newUser.FollowCount = in.FollowCount
	newUser.FanCount = in.FanCount

	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Save(&newUser).Error
	if err != nil {
		tx.Rollback()//失败回滚
		return nil, err
	}

	tx.Commit()//成功提交
	return &user.UpdateUserResp{
		StatusCode: config.STATUS_SUCCESS,
		StatusMsg:  config.STATUS_SUCCESS_MSG,
	}, nil
}
