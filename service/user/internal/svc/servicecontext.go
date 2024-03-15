package svc

import (
	"douyin/service/user/db"
	"douyin/service/user/db/RedisCache"
	"douyin/service/user/internal/config"
	"log"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *RedisCache.RedisPool
	AsynqClient  *asynq.Client
	Db     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err :=db.InitGorm(c.DbConfig)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	pool := RedisCache.NewRedisPool(c)
	conn := pool.GetRedisConn()

	_, err = conn.Do("PING")
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return &ServiceContext{
		Config: c,
		Redis:  pool,
		Db:     db,
		AsynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: "redis:"+c.RedisConfig.Host, Password: c.RedisConfig.Password}),
	}
}
