package handler

import (
	"context"
	"douyin/service/mq/internal/svc"
	"fmt"
	mqcommon "douyin/service/mq/common"
	croncommon "douyin/service/cron/common"
	"log"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)


type AsynqServer struct {
	ctx    context.Context
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

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr:"redis:"+l.svcCtx.Config.Redis.Host, Password: l.svcCtx.Config.Redis.Password},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	
	mux.HandleFunc(mqcommon.TypeDelCache, l.delCacheHandler)
	mux.HandleFunc(mqcommon.TypeAddCacheValue, l.addCacheValueHandler)

	mux.HandleFunc(croncommon.TypeSyncUserInfoCache, l.syncUserInfoCacheHandler)
	mux.HandleFunc(croncommon.TypeSyncVideoInfoCache, l.syncVideoInfoCacheHandler)
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (l *AsynqServer) Stop() {
	fmt.Println("AsynqTask stop")
}

