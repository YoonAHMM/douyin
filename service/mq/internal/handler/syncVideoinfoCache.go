package handler

import (
	"context"
	"strconv"

	"douyin/service/video/rpc/model"
	"douyin/service/video/rpc/video"

	"github.com/hibiken/asynq"
)

func (l *AsynqServer) syncVideoInfoCacheHandler(ctx context.Context, t *asynq.Task) error {
	res, err := l.svcCtx.Redis.LRange(ctx, model.VideoCacheKeyPrefix , 0, -1).Result()
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}

	for _, v := range res {
		videoId,_ := strconv.ParseUint(v, 10, 64)
		// 读取缓存
		vi:=model.Video{}
		videoInfo, err := l.svcCtx.Redis.HGetAll(ctx, vi.CacheKey(videoId)).Result()
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}

		// 更新视频信息
		
		FC,_:=strconv.ParseInt(videoInfo["FavoriteCount"], 10, 64)
		CC,_:=strconv.ParseInt(videoInfo["CommentCount"], 10, 64)
		
		_, err = l.svcCtx.VideoRpc.UpdateVideo(ctx, &video.UpdateVideoReq{
			Video: &video.VideoInfo{
				Id:            int64(videoId),
				FavoriteCount: FC,
				CommentCount:  CC,
			},
		})
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
