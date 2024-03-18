package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)
type Config struct {
	rest.RestConf
	UserRpc   zrpc.RpcClientConf
	JwtRpc    zrpc.RpcClientConf
	VideoRpc  zrpc.RpcClientConf

	JwtConfig struct {
		AccessExpire int64
	}
	AliyunOss struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		VideoBucket     string
		VideoPath       string
	}
	FeedLimit int64
}


const (

	STATUS_SUCCESS          = "0"  //成功
	STATUS_FAIL             = "1"  //失败
	STATUS_FAIL_TOKEN_MSG   = "Token is invalid" //令牌无效
	STATUS_FAIL_TOOLONG_MSG = "Username or Password must less than 32 characters" //用户名或密码长度超过 32 个字符
	STATUS_SUCCESS_MSG      = "OK" //成功消息
	STATUS_FAIL_PARAM_MSG   = "Request parameter error"//请求参数错误

	USER_NO_LOGIN           = "0" // 需保证不出现 id 为 0 的用户,用户未登录

	STATUS_FAIL_FOLLOW_SELF = "Follow yourself is not allowed"//不允许关注自己
	FILE_EMPTY_ERROR        = "upload file is empty" //上传文件为空
	FILE_TYPE_ERROR         = "upload file type error"//上传文件类型错误

	MP4_TYPE                = "video/mp4"// MP4 视频文件类型
)











