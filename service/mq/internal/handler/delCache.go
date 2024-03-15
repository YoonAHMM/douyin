package handler

import (
	"context"
	"douyin/service/mq/common"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)


func (l *AsynqServer) delCacheHandler(ctx context.Context, t *asynq.Task) error {
	var p common.DelCachePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		l.Logger.Errorf("json.Unmarshal failed: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := l.svcCtx.Redis.Del(ctx, p.Key).Err(); err != nil {
		l.Logger.Errorf("redis.Del failed: %v", err)
		return err
	}

	time.Sleep(1 * time.Second)

	if err := l.svcCtx.Redis.Del(ctx, p.Key).Err(); err != nil {
		l.Logger.Errorf("redis.Del failed: %v", err)
		return err
	}

	return nil
}
