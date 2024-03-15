package svc

import (
	"douyin/service/user/userrpc"
	"douyin/service/video/rpc/internal/config"
	"douyin/service/video/rpc/model"
	"douyin/service/video/rpc/model/redisCache"

	"log"

	"github.com/hibiken/asynq"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	UserRpc     userrpc.UserRpc
	Redis       *redisCache.RedisPool
	Db          *gorm.DB
	AsynqClient  *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := model.InitGorm(c.DbConfig)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	pool := redisCache.NewRedisPool(c)
	conn := pool.NewRedisConn()
	_, err = conn.Do("PING")
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return &ServiceContext{
		Config:  c,
		UserRpc: userrpc.NewUserRpc(zrpc.MustNewClient(c.UserRpc)),
		Redis:   pool,
		Db:      db,
	}
}
