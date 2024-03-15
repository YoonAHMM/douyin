package logic

import (
	"context"
	"errors"
	"time"

	"douyin/service/video/rpc/model"
	"douyin/service/video/rpc/internal/svc"
	"douyin/service/video/rpc/video"

	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"douyin/service/mq/common"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FavoriteAction 处理点赞/取消点赞请求
func (l *FavoriteActionLogic) FavoriteAction(in *video.FavoriteReq) (*video.FavoriteResp, error) {

	// 提取请求参数
	videoId, err := strconv.ParseUint(in.VideoId, 10, 64)
	if err != nil {
		return nil, err
	}

	userid, err := strconv.ParseUint(in.UserId, 10, 64)
	if err != nil {
		return nil, err
	}

	actionType := in.ActionType

	switch actionType{
	case FAVORITE_UPDATE:
		createTime := time.Now().Unix()
		err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
			// 先查询用户是否点赞过该视频
			f := model.Favorite{}
			err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("user_id = ? And video_id = ?", userid, videoId).
				First(&f).Error
	
			// 点赞记录已存在
			if err == nil {
				return nil
			}
	
			// 数据库查询错误
			if err != gorm.ErrRecordNotFound {
				return err
			}
	
			// 未点赞，创建记录
			f.VideoId = videoId
			f.UserId =  userid
			if err := tx.Create(&f).Error; err != nil {
				return err
			}
	
			// 缓存视频点赞量加一
			conn := l.svcCtx.Redis.NewRedisConn()
			defer conn.Close()
			err = l.svcCtx.Redis.AddFavorite(conn, videoId, userid, createTime)
			if err != nil {
				return err
			}

			//异步处理DB
			v:=model.Video{}
			task, err := common.NewAddCacheValueTask(v.CacheKey(videoId), "FavoriteCount", 1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return err
			}
			return nil
		})
		if err !=nil{
			return nil,err
		}
	case FAVORITE_DELETE:{
		createTime := time.Now().Unix()
		err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
			// 查询用户喜欢记录是否存在
			f := model.Favorite{}
			err := tx.
				Where("user_id = ? And video_id = ?", userid, videoId).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&f).Error
	
			// 点赞记录不存在
			if err == gorm.ErrRecordNotFound {
				return nil
			}
	
			if err != nil {
				return err
			}
	
			// 删除记录
			if err := tx.Where("user_id = ? And video_id = ?", userid, videoId).Delete(&f).Error; err != nil {
				return err
			}
	
			// 缓存视频点赞数减一
			if err := tx.Model(&model.Video{}).
				Where("id = ?", in.VideoId).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).
				Error; err != nil {
	
				return err
			}
			conn := l.svcCtx.Redis.NewRedisConn()
			defer conn.Close()

			err = l.svcCtx.Redis.DelFavorite(conn, videoId, userid, l.svcCtx.Config.CacheConfig.FAVORITE_CACHE_TTL,int(createTime))
			if err != nil {
				return err
			}

			//异步处理DB
			v:=model.Video{}
			task, err := common.NewDelCacheTask(v.CacheKey(videoId),)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return err
			}

			return nil
		})
	
		if err != nil {
			return nil, err
		}
	
	}
	default:
		return nil, errors.New(STATUS_FAIL_PARAM_MSG)
	}
	return &video.FavoriteResp{
		StatusCode: STATUS_SUCCESS,
		StatusMsg:  STATUS_SUCCESS_MSG,
	}, nil

}