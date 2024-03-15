package RedisCache

import (
	"douyin/service/user/internal/config"
	"douyin/service/user/db/model"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type RedisPool struct {
	pool *redis.Pool
}

const (
	CACHE_KEY_NOT_EXISTS_MSG = "cache key not exists or has wrong type"// 缓存键不存在或类型错误
	COUNT_NOT_FOUND          = -1
)

var cacheConfig *config.CacheConfig

func NewRedisPool(config config.Config) *RedisPool {
	return &RedisPool{&redis.Pool{
		MaxIdle:     config.RedisConfig.MaxIdle,
		MaxActive:   config.RedisConfig.Active,
		IdleTimeout: time.Duration(config.RedisConfig.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			redisUri := fmt.Sprintf("%s:%d", config.RedisConfig.Host, config.RedisConfig.Port)

			if config.RedisConfig.Auth {
				//已经认证用户
				redisConn, err := redis.Dial("tcp", redisUri,
					redis.DialUsername(config.RedisConfig.Username),
					redis.DialPassword(config.RedisConfig.Password),
				)

				if err != nil {
					return nil, err
				}
				return redisConn, nil

			} else {
				//非认证用户
				redisConn, err := redis.Dial("tcp", redisUri)
				if err != nil {
					return nil, err
				}
				return redisConn, nil

			}
		},
	}}
}

func (p *RedisPool) GetRedisConn() redis.Conn {
	return p.pool.Get()
}


//HINCRBY 
func (p *RedisPool) HIncrIntVal(conn redis.Conn, key string, field string, isIncr bool) error {
	i := 1
	if !isIncr {
		i = -1
	}

	_, err := conn.Do("EVAL", "local val = redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2]); "+
		"if (val >= ARGV[3]) then "+
		"redis.call('PERSIST', KEYS[1]); "+
		"else "+
		"redis.call('EXPIRE', ARGV[4]);"+
		"return nil; end;", 1, key, field, i)
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) IncrFollowCount(conn redis.Conn, userid uint64) error {
	key := model.User{}.CacheKey(userid)
	err := p.HIncrIntVal(conn, key, model.FollowCountField, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) DecrFollowCount(conn redis.Conn, userid uint64) error {
	key := model.User{}.CacheKey(userid)
	err := p.HIncrIntVal(conn, key, model.FollowCountField, false)
	if err != nil {
		return err
	}
	return nil
}


func (p *RedisPool) IncrFollowerCount(conn redis.Conn, userid uint64) error {
	key := model.User{}.CacheKey(userid)
	err := p.HIncrIntVal(conn, key, model.FollowerCountField, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) DecrFollowerCount(conn redis.Conn, userid uint64) error {
	key := model.User{}.CacheKey(userid)
	err := p.HIncrIntVal(conn, key, model.FollowerCountField, false)
	if err != nil {
		return err
	}
	return nil
}

//GET
func (p *RedisPool) getExHash(conn redis.Conn, key string, field string, ttl int) ([]byte, error) {

	raw, err := conn.Do("EVAL", "redis.call('EXPIRE', KEYS[1], ARGV[2]); "+
		"return redis.call('HGET', KEYS[1], ARGV[1]); ", 1, key, field, ttl)
	if err != nil {
		return nil, err
	}
	if val, ok := raw.([]byte); ok {
		return val, nil
	}
	return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
}

// 粉丝数
func (p *RedisPool) GetFollowerCount(conn redis.Conn, userid uint64) (int64, error) {
	raw, err := p.getExHash(conn, model.User{}.CacheKey(userid), model.FollowerCountField, cacheConfig.USER_CACHE_TTL)
	if err != nil {
		return COUNT_NOT_FOUND, err
	}

	cnt, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return COUNT_NOT_FOUND, err
	}

	return cnt, err
}


//SET
func (p *RedisPool) setExHashIntVal(conn redis.Conn, key string, field string, value int64, ttl int) ([]byte, error) {
	raw, err := conn.Do("EVAL", "redis.call('HSET', KEYS[1], ARGV[1], ARGV[2]); "+
		"redis.call('EXPIRE', KEYS[1], ARGV[3]); "+
		"return nil; ", 1, key, field, value, ttl)
	if err != nil {
		return nil, err
	}
	if val, ok := raw.([]byte); ok {
		return val, nil
	}
	return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
}


func (p *RedisPool) SetFollowCount(conn redis.Conn, userid uint64, count int64) error {
	_, err := p.setExHashIntVal(conn, model.User{}.CacheKey(userid), model.FollowCountField, count, cacheConfig.USER_CACHE_TTL)
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) SetFollowerCount(conn redis.Conn, userid uint64, count int64) error {
	_, err := p.setExHashIntVal(conn, model.User{}.CacheKey(userid), model.FollowerCountField, count, cacheConfig.USER_CACHE_TTL)
	if err != nil {
		return err
	}
	return nil
}

//set userinfo 如果关注或者粉丝超过阈值，加入缓存
func (p *RedisPool) SetUserInfo(conn redis.Conn, userid uint64, username string, followCount int64, followerCount int64) error {
	_, err := conn.Do("EVAL",
		"redis.call('HSETNX', KEYS[1], ARGV[4], ARGV[1]; "+
		"redis.call('HSETNX', KEYS[1], ARGV[5], ARGV[2]);"+
		"redis.call('HSETNX', KEYS[1], ARGV[6], ARGV[3]);"+
		"if (tonumber(ARGV[2]) >= tonumber(ARGV[8]) or tonumber(ARGV[3]) >= tonumber(ARGV[8])) then "+
		"redis.call('PERSIST', KEYS[1]); "+
		"else redis.call('EXPIRE', KEYS[1], ARGV[7]); end; "+
		"return nil;",
		1, model.User{}.CacheKey(userid), username, followCount, followerCount,
		model.UsernameField, model.FollowCountField, model.FollowerCountField,
		cacheConfig.USER_CACHE_TTL, cacheConfig.FOLLOW_COUNT_THRESHOLD)
	if err != nil {
		return err
	}
	return nil
}


//得到用户信息，同时设置缓存过期时间
func (p *RedisPool) GetUserInfo(conn redis.Conn, userid uint64) (username string, followCount int64, followerCount int64, err error) {
	raw, err := conn.Do("EVAL", "if (redis.call('EXISTS',KEYS[1]) ~= 1) then "+
		"return nil; end; "+
		"local cnt = redis.call('HGET', KEYS[1], ARGV[2]); "+
		"local cnt2 = redis.call('HGET', KEYS[1], ARGV[3]); "+
		"if (tonumber(cnt) >= tonumber(ARGV[4]) or tonumber(cnt2) >= tonumber(ARGV[4])) then "+
		"redis.call('PERSIST', KEYS[1]); "+
		"else redis.call('EXPIRE', KEYS[1], ARGV[5]); end; "+
		"return redis.call('HMGET',KEYS[1],ARGV[1],ARGV[2],ARGV[3]); ", 1, model.User{}.CacheKey(userid),
		model.UsernameField, model.FollowCountField, model.FollowerCountField,
		cacheConfig.FOLLOW_COUNT_THRESHOLD, cacheConfig.USER_CACHE_TTL)

	if err != nil {
		return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, err
	}

	if raw == nil {
		return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, err
	}

	if r, ok := raw.([]interface{}); ok {
		nameBytes, ok := r[0].([]byte)
		if !ok {
			return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
		}

		username = string(nameBytes)

		followCntBytes, ok := r[1].([]byte)
		if !ok {
			return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
		}

		followCount, err = strconv.ParseInt(string(followCntBytes), 10, 64)
		if err != nil {
			return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
		}

		followerCntBytes, ok := r[2].([]byte)
		if !ok {
			return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
		}

		followerCount, err = strconv.ParseInt(string(followerCntBytes), 10, 64)
		if err != nil {
			return "", COUNT_NOT_FOUND, COUNT_NOT_FOUND, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
		}
	}
	return
}


// 将每个用户id，用户名，关注数，点赞数均以 string 表示返回
func (p *RedisPool) getUserList(conn redis.Conn, key string) ([][]string, error) {
	raw, err := conn.Do("EVAL", "if (redis.call('EXISTS', KEYS[1]) ~= 1) then "+
		"return nil; "+
		"else "+
		"local zlist = redis.call('ZRANGE', KEYS[1], 0, -1, 'REV'); "+
		"local len = #zlist; "+
		"redis.call('EXPIRE', KEYS[1], ARGV[6]); "+
		"local res_array = {}; "+
		"for i, v in pairs(zlist) do "+
		"local t = {ARGV[1], v}; "+
		"local key = table.concat(t);"+
		"if(redis.call('EXISTS', key) == 0) then "+
		"return nil; "+
		"else "+
		"local followCount = redis.call('HGET', key, ARGV[3]); "+
		"local followerCount = redis.call('HGET', key, ARGV[4]); "+
		"if (tonumber(followCount) >= tonumber(ARGV[5]) or tonumber(followerCount) >= tonumber(ARGV[5])) then "+
		"redis.call('PERSIST', key); "+
		"else redis.call('EXPIRE', key, ARGV[7]); end; "+
		"local username = redis.call('HGET', key, ARGV[2]); "+
		"table.insert(res_array,{v, username, followCount, followerCount}); end; end; "+
		"return res_array; end; ", 1, key, model.UserCacheKeyPrefix,
		model.UsernameField, model.FollowCountField, model.FollowerCountField,
		cacheConfig.FOLLOW_COUNT_THRESHOLD, cacheConfig.FOLLOW_CACHE_TTL, cacheConfig.USER_CACHE_TTL)

	if err != nil {
		return nil, err
	}

	var userList [][]string

	r, ok := raw.([]interface{})
	if ok {
		for _, s := range r {
			v, ok := s.([]interface{})
			if !ok {
				return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
			}

			var userInfo []string
			IdBytes, ok := v[0].([]byte)
			if !ok {
				return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
			}
			userInfo = append(userInfo, string(IdBytes))

			nameBytes, ok := v[1].([]byte)
			if !ok {
				return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
			}
			userInfo = append(userInfo, string(nameBytes))

			followCntBytes, ok := v[2].([]byte)
			if !ok {
				return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
			}
			userInfo = append(userInfo, string(followCntBytes))

			followerCntBytes, ok := v[3].([]byte)
			if !ok {
				return nil, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
			}
			userInfo = append(userInfo, string(followerCntBytes))

			userList = append(userList,userInfo)
		}

	}

	return userList, errors.New(CACHE_KEY_NOT_EXISTS_MSG)
}


//关注列表
func (p *RedisPool) GetFollowUserList(conn redis.Conn, userid uint64) ([][]string, error) {
	IdCacheKey := model.Follow{}.FollowListCacheKey(userid)
	return p.getUserList(conn, IdCacheKey)
}

//粉丝列表
func (p *RedisPool) GetFollowerUserList(conn redis.Conn, userid uint64) ([][]string, error) {
	IdCacheKey := model.Follow{}.FollowerListCacheKey(userid)
	return p.getUserList(conn, IdCacheKey)
}

//user是否关注to user
func (p *RedisPool) IsFollow(conn redis.Conn, userid, toUserId uint64) (bool, error) {
	FollowListCacheKey := model.Follow{}.FollowListCacheKey(userid)
	FollowerListCacheKey := model.Follow{}.FollowerListCacheKey(toUserId)
	raw, err := conn.Do("EVAL",
		"if(redis.call('ZRANK', KEYS[1], ARGV[2]) ~= false or redis.call('ZRANK', KEYS[2], ARGV[1]) ~= false) "+
			"then return 1; "+
			"else return nil; end; ", 2, FollowListCacheKey, FollowerListCacheKey, userid, toUserId)
	if err != nil {
		return false, err
	}
	if raw != nil {
		return true, nil
	}
	return false, nil
}


func (p *RedisPool) addFolUserList(conn redis.Conn, userKey, toUserKey string, userid, toUserid uint64, createTime int64, ttl int) error {
	Exat := createTime + int64(ttl)
	_, err := conn.Do("EVAL", "redis.call('ZADD', KEYS[1], ARGV[3], ARGV[2]); "+
		"redis.call('EXPIREAT', KEYS[1], ARGV[4]); "+
		"redis.call('ZADD', KEYS[2], ARGV[3], ARGV[1]); "+
		"redis.call('EXPIREAT', KEYS[2], ARGV[4]); "+
		"return nil; ", 2, userKey, toUserKey, userid, toUserid, createTime, Exat)
	if err != nil {
		return err
	}
	return nil
}

//关注
func (p *RedisPool) AddFollowUserList(conn redis.Conn, userid, toUserId uint64, createTime int64) error {
	UserCacheKey := model.Follow{}.FollowListCacheKey(userid)
	ToUserCacheKey := model.Follow{}.FollowerListCacheKey(toUserId)
	return p.addFolUserList(conn, UserCacheKey, ToUserCacheKey, userid, toUserId, createTime, cacheConfig.FOLLOW_CACHE_TTL)
}


func (p *RedisPool) remExFolUserList(conn redis.Conn, userKey, toUserKey string, userid, toUserid uint64, ttl int) error {
	Exat := time.Now().Unix() + int64(ttl)
	_, err := conn.Do("EVAL", "redis.call('ZREM', KEYS[1], ARGV[2]); "+
		"redis.call('EXPIREAT', KEYS[1], ARGV[3]); "+
		"redis.call('ZREM', KEYS[2], ARGV[1]); "+
		"redis.call('EXPIREAT', KEYS[2], ARGV[3]); "+
		"return nil; ", 2, userKey, toUserKey, userid, toUserid, Exat)
	if err != nil {
		return err
	}
	return nil
}

//取消关注
func (p *RedisPool) RemFollowUserList(conn redis.Conn, userid, toUserId uint64) error {
	UserCacheKey := model.Follow{}.FollowListCacheKey(userid)
	ToUserCacheKey := model.Follow{}.FollowerListCacheKey(toUserId)
	return p.remExFolUserList(conn, UserCacheKey, ToUserCacheKey, userid, toUserId, cacheConfig.FOLLOW_CACHE_TTL)
}


// InitRedis 通过从 DB 获取数据初始化 redis（缓存预热）
func (p *RedisPool) InitRedis(conf *config.CacheConfig, db *gorm.DB) error {
	
	cacheConfig = conf
	var UserList []*model.User

	var UserCnt int64
	err := db.Model(&model.User{}).Count(&UserCnt).Error
	if err != nil {
		return err
	}

	// 用户数量大于配置的用户信息缓存初始化大小
	if UserCnt > int64(cacheConfig.USER_CACHE_INIT_SIZE) {

		// 使用粉丝数最多的（默认最多初始化 10 万条用户信息）一批用户信息用于缓存预热
		subQuery1 := db.Model(&model.Follow{}).Select("following_id,count(1) as cnt").Group("following_id").Order("cnt DESC").Limit(cacheConfig.USER_CACHE_INIT_SIZE)
		subQuery2 := db.Select("following_id as id").Table("(?) as F", subQuery1)
		err := db.Find(&UserList, "id in (?)", subQuery2).Error
		if err != nil {
			return err
		}
	} else {

		// 将所有用户信息用于缓存预热
		err := db.Find(&UserList).Error
		if err != nil {
			return err
		}
	}

	// 初始化缓存
	var FollowerCount int64  // 粉丝数
	var FollowCount int64    // 关注数
	conn := p.GetRedisConn() // Redis Conn
	defer conn.Close()

	for _, user := range UserList {
		err = db.Model(&model.Follow{}).Where(&model.Follow{Follower: user.Id}).Count(&FollowCount).Error
		if err != nil {
			return err
		}

		err = db.Model(&model.Follow{}).Where(&model.Follow{Following: user.Id}).Count(&FollowerCount).Error
		if err != nil {
			return err
		}

		err = p.SetUserInfo(conn, user.Id, user.Username, FollowCount, FollowerCount)
		if err != nil {
			return err
		}
	}

	return nil
}
