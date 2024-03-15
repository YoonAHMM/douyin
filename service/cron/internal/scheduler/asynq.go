package scheduler

import (
	"context"
	"douyin/service/cron/common"
	"douyin/service/cron/internal/svc"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)


type AsynqServer struct{
	ctx context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAsynqServer(ctx context.Context, svcCtx *svc.ServiceContext) *AsynqServer {
	return &AsynqServer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AsynqServer) Start() {
	fmt.Println("AsynqTask start")

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     "redis:"+l.svcCtx.Config.Redis.Host,
			Password: l.svcCtx.Config.Redis.Password},
		nil,
	)

	syncUserInfoCacheTask := common.NewSyncUserInfoCacheTask()
	entryID, err := scheduler.Register("@every 1h", syncUserInfoCacheTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	syncVideoInfoCacheTask := common.NewSyncVideoInfoCacheTask()
	entryID, err = scheduler.Register("@every 301s", syncVideoInfoCacheTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}

func (l *AsynqServer) Stop() {
	fmt.Println("AsynqTask stop")
}
