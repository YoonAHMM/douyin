package svc

import (
	"douyin/service/http/internal/config"
	"douyin/service/jwt/Jwt"
	"douyin/service/user/userrpc"
	"douyin/service/video/rpc/videorpc"

	"github.com/segmentio/kafka-go"
)

type ServiceContext struct {
	Config config.Config
	UserRpc     userrpc.UserRpc
	VideoRpc    videorpc.VideoRpc
	JwtRpc      Jwt.JwtRpcClient
	KafkaWriter *kafka.Writer

}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
