package svc

import (
	"context"
	"douyin/service/mq/internal/config"
	"douyin/service/user/userrpc"
	"douyin/service/video/rpc/videorpc"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	Redis      *redis.Client
	VideoRpc   videorpc.VideoRpc
	UserRpc    userrpc.UserRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Redis:      initRedis(c),
		VideoRpc:   videorpc.NewVideoRpc(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:    userrpc.NewUserRpc(zrpc.MustNewClient(c.UserRpc)),
	}
}

func initRedis(c config.Config) *redis.Client {
	fmt.Println("connect Redis ...")
	db := redis.NewClient(&redis.Options{
		Addr:     "redis:"+c.Redis.Host,
		Password: c.Redis.Password,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  3 * time.Second,
	})
	_, err := db.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("connect Redis failed")
		panic(err)
	}
	fmt.Println("connect Redis success")
	return db
}
