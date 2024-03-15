package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

	"douyin/service/user/internal/config"
	"douyin/service/user/db/model"
	"douyin/service/user/internal/svc"
	"douyin/service/user/user"
	"douyin/service/mq/common"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type FollowActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowActionLogic {
	return &FollowActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowActionLogic) FollowAction(in *user.FollowActionReq) (*user.FollowActionResp, error) {
	userid, err := strconv.ParseUint(in.UserId, 10, 64)
	if err !=nil{
		return nil,err
	}

	toUserId,err :=strconv.ParseUint(in.ToUserId,10,64)
	if err !=nil{
		return nil,err
	}

	conn := l.svcCtx.Redis.GetRedisConn()
	defer conn.Close()
	

	switch in.ActionType {
	//关注
	case config.OP_FOLLOW:
		createTime := time.Now().Unix()
		// 查找用户是否存在
		err = l.svcCtx.Db.Model(model.User{}).Where(&model.User{Id: toUserId}).Take(&model.User{}).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound { 
				return &user.FollowActionResp{
					StatusCode: config.STATUS_FAIL,
					StatusMsg:  config.STATUS_USER_NOTEXIST_MSG,
				}, nil
			}
			return nil, err
		}

		err = l.svcCtx.Db.Model(model.User{}).Where(&model.User{Id: userid}).Take(&model.User{}).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound { 
				return &user.FollowActionResp{
					StatusCode: config.STATUS_FAIL,
					StatusMsg:  config.STATUS_USER_NOTEXIST_MSG,
				}, nil
			}
			return nil, err
		}
		
		err := l.svcCtx.Db.Create(&model.Follow{
			Follower:   userid,
			Following:  toUserId,
			CreateTime: createTime,
		}).Error
		if err != nil {
			return nil, err
		}

		s,_:=l.svcCtx.Redis.GetFollowerCount(conn,userid)
		if model.IsPopularUser(s){
			task, err := common.NewAddCacheValueTask(model.User{}.CacheKey(userid), "FollowCount", 1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return nil,err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return nil,err
			}
			err = l.svcCtx.Redis.IncrFollowCount(conn, userid)
			if err != nil {
				return nil, err
			}
		}

		s,_=l.svcCtx.Redis.GetFollowerCount(conn,toUserId)
		if model.IsPopularUser(s){
			task, err := common.NewAddCacheValueTask(model.User{}.CacheKey(toUserId), "FanCount", 1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return nil,err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return nil,err
			}
			err = l.svcCtx.Redis.IncrFollowerCount(conn, toUserId)
			if err != nil {
				return nil, err
			}
		}
		// 更新双方用户的最新关注列表以及最新粉丝列表
		err = l.svcCtx.Redis.AddFollowUserList(conn, userid, toUserId, createTime)
		if err != nil {
			return nil, err
		}
	//取消关注
	case config.OP_CANCEL_FOLLOW:

		err := l.svcCtx.Db.Delete(&model.Follow{}, &model.Follow{
			Follower:  userid,
			Following: toUserId,
		}).Error
		if err != nil {
			return nil, err
		}

		s,_:=l.svcCtx.Redis.GetFollowerCount(conn,userid)
		if model.IsPopularUser(s){
			task, err := common.NewAddCacheValueTask(model.User{}.CacheKey(userid), "FollowCount", -1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return nil,err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return nil,err
			}
			err = l.svcCtx.Redis.DecrFollowCount(conn, userid)
			if err != nil {
				return nil, err
			}
		}

		s,_=l.svcCtx.Redis.GetFollowerCount(conn,toUserId)
		if model.IsPopularUser(s){
			task, err := common.NewAddCacheValueTask(model.User{}.CacheKey(userid), "FanCount", -1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return nil,err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return nil,err
			}
			err = l.svcCtx.Redis.DecrFollowerCount(conn, toUserId)
			if err != nil {
				return nil, err
			}
		}
		// 更新双方用户的最新关注列表以及最新粉丝列表
		err = l.svcCtx.Redis.RemFollowUserList(conn, userid, toUserId)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(config.STATUS_FAIL_PARAM_MSG)
	}
	return &user.FollowActionResp{
		StatusCode: config.STATUS_SUCCESS,
		StatusMsg:  config.STATUS_SUCCESS_MSG,
	}, nil
}
