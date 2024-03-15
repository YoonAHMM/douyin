package video

import (
	"context"
	"encoding/json"
	"net/http"

	"douyin/service/http/internal/config"
	"douyin/service/http/internal/svc"
	"douyin/service/http/internal/types"
	"douyin/service/jwt/client/jwtrpc"

	"mime/multipart"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hashicorp/go-uuid"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
type MsgInfo struct {
	Title        string `json:"title"`
	OssObjectKey string `json:"ossObjectKey"`
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionReq, formFile multipart.File) (resp *types.PublishActionResp, err error) {
	defer formFile.Close()
	// 1. token 鉴权
	token, err := l.svcCtx.JwtRpc.ParseToken(l.ctx, &jwtrpc.ParseTokenReq{Token: req.Token})
	if err != nil {
		return nil, err
	}

	// 2. 检测文件是否为空
	if formFile == nil {
		return &types.PublishActionResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.FILE_EMPTY_ERROR,
			}}, nil
	}

	// 3. 检测文件类型是否为 mp4
	ok, err := IsFileTypeMP4(formFile)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &types.PublishActionResp{
			Status: types.Status{
				StatusCode: config.STATUS_FAIL,
				StatusMsg:  config.FILE_TYPE_ERROR,
			},
		}, nil
	}

	userid := token.UserID
	title := req.Title

	// 4. 将文件上传至 OSS
	ossObjKey, err := l.UploadFormFile(formFile)
	if err != nil {
		return nil, err
	}

	// 5. 将视频转码请求写入 kafka，随后可以马上返回客户端了（所以返回后客户端会延迟一段时间，等待服务处理完成后才可看到新视频）
	m := MsgInfo{
		Title:        title,
		OssObjectKey: ossObjKey,
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(userid),
		Value: marshal,
	}
	err = l.svcCtx.KafkaWriter.WriteMessages(l.ctx, kafkaMsg)
	if err != nil {
		return nil, err
	}
	return &types.PublishActionResp{
		Status: types.Status{
			StatusCode: config.STATUS_SUCCESS,
			StatusMsg:  config.STATUS_SUCCESS_MSG,
		},
	}, nil
}

// UploadFormFile 将文件上传到OSS
func (l *PublishActionLogic) UploadFormFile(formFile multipart.File) (string, error) {
	endpoint := l.svcCtx.Config.AliyunOss.Endpoint
	accessKeyId := l.svcCtx.Config.AliyunOss.AccessKeyId
	accessKeySecret := l.svcCtx.Config.AliyunOss.AccessKeySecret

	cli, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return "", err
	}

	bucket, err := cli.Bucket(l.svcCtx.Config.AliyunOss.VideoBucket)
	if err != nil {
		return "", err
	}
	fileName, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	ossObjKey := l.svcCtx.Config.AliyunOss.VideoPath + fileName

	err = bucket.PutObject(ossObjKey, formFile, oss.Checkpoint(true, ""))
	if err != nil {
		return "", err
	}
	return ossObjKey, nil
}

// IsFileTypeMP4 检查上传文件的前 512 字节来判断文件类型是否为 MP4
// bool 返回值为 true 则表示文件为 mp4 格式
func IsFileTypeMP4(formFile multipart.File) (bool, error) {
	buffer := make([]byte, 512)
	_, err := formFile.Read(buffer)
	_, err = formFile.Seek(0, 0)
	if err != nil {
		return false, err
	}
	contentType := http.DetectContentType(buffer)
	if contentType != config.MP4_TYPE {
		return false, nil
	}
	return true, nil
}

