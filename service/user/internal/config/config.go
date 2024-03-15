package config

import (
	"strconv"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DbConfig    DbConfig
	RedisConfig struct {
		Host        string
		Port        int
		Username    string
		Password    string
		Auth        bool //认证
		MaxIdle     int
		Active      int
		IdleTimeout int
	}
	CacheConfig CacheConfig
	WorkerId    uint32
}

type DbConfig struct{
	Path         string `json:"path" yaml:"path"`                     // 服务器地址
	Port         int    `json:"port" yaml:"port"`                     // 端口
	Config       string `json:"config" yaml:"config"`                 // 高级配置
	Dbname       string `json:"db-name" yaml:"db-name"`               // 数据库名
	Username     string `json:"username" yaml:"username"`             // 数据库用户名
	Password     string `json:"password" yaml:"password"`             // 数据库密码
	MaxIdleConns int    `json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
}

type Mysql struct {
	DbConfig
}


func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + strconv.FormatInt(int64(m.Port), 10) + ")/" + m.Dbname + "?" + m.Config
}


type CacheConfig struct {
	FOLLOW_CACHE_TTL            int //关注数据缓存的过期时间
	FOLLOW_COUNT_CACHE_TTL      int //关注计数缓存的过期时间
	FOLLOW_COUNT_THRESHOLD      int //关注计数超过该阈值才会缓存
	USER_CACHE_TTL              int //用户数据缓存的过期时间
	USER_CACHE_INIT_SIZE        int //用户缓存创建时的大小
	FOLLOWLIST_MAX_CACHE_SIZE   int //关注列表缓存的最大容量
	FOLLOWERLIST_MAX_CACHE_SIZE int //粉丝列表缓存的最大容量
}



const (
	STATUS_SUCCESS            = "0"							//操作成功
	STATUS_SUCCESS_MSG        = "OK"						//操作成功信息，常用于 "OK"
	STATUS_FAIL               = "1"							//操作失败
	STATUS_FAIL_PARAM_MSG     = "param incorrect" 			//请求参数格式不正确
	STATUS_USER_EXISTS_MSG    = "username already exists"	//用户名已经被注册
	STATUS_USER_NOTEXIST_MSG  = "user not exist"			//用户名不存在
	STATUS_WRONG_PASSWORD_MSG = "wrong Password"			//密码错误
	COUNT_NOT_FOUND           = int64(-1)					//找不到
	OP_FOLLOW                 = "1"							//关注
	OP_CANCEL_FOLLOW          = "2"							//取消关注	
)