package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)


type Config struct {
	service.ServiceConf

	Redis      RedisConfig

	VideoRpc   zrpc.RpcClientConf
	UserRpc    zrpc.RpcClientConf
}

type RedisConfig struct {
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
	Username    string `yaml:"Username"`
	Password    string `yaml:"Password"`
	Auth        bool   `yaml:"Auth"`
	MaxIdle     int    `yaml:"MaxIdle"`
	Active      int    `yaml:"Active"`
	IdleTimeout int    `yaml:"IdleTimeout"`
}
