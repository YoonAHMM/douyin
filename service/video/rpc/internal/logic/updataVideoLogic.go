package logic

import (
	"context"
	"douyin/service/video/rpc/model"
	"douyin/service/video/rpc/internal/svc"
	"douyin/service/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
)


type UpdateVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVideoLogic {
	return &UpdateVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateVideoLogic) UpdateVideo(in *video.UpdateVideoReq) (*video.UpdateVideoResp, error) {
	// 开启事务
	tx := l.svcCtx.Db.Begin()

	var newVideo model.Video
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.Video.Id).First(&newVideo).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	newVideo.CommentCount = in.Video.CommentCount
	newVideo.FavoriteCount = in.Video.FavoriteCount

	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Save(&newVideo).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &video.UpdateVideoResp{
		StatusCode: STATUS_SUCCESS,
		StatusMsg:  STATUS_SUCCESS_MSG,
	}, nil
}
