package cron

import (
	"context"
	"douyin/service/cron/internal/config"
	"douyin/service/cron/internal/scheduler"
	"douyin/service/cron/internal/svc"
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/trace"
)

var configFile = flag.String("f", "etc/cron.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)
	// nolint:staticcheck
	prometheus.StartAgent(c.Prometheus)
	trace.StartAgent(c.Telemetry)

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	logx.DisableStat()

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	serviceGroup.Add(scheduler.NewAsynqServer(ctx, svcContext))
	serviceGroup.Start()
}