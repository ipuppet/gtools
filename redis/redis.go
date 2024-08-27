package redis

import (
	"context"
	"errors"
	"time"

	"github.com/ipuppet/gtools/utils"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Md5Key   bool   `json:"md5key" binding:"-"`
}

type RedisClient struct {
	Config  *RedisConfig
	rdb     *redis.Client
	Context context.Context
}

func NewRedisClient(config *RedisConfig) *RedisClient {
	return &RedisClient{Config: config}
}

func (rc *RedisClient) GetClient() *redis.Client {
	if rc.rdb == nil {
		rc.rdb = redis.NewClient(&redis.Options{
			Addr:     rc.Config.Addr,
			Password: rc.Config.Password, // 没有密码，默认值
			DB:       rc.Config.DB,       // 默认DB 0
		})
		if rc.Context == nil {
			rc.Context = context.Background()
		}
	}
	return rc.rdb
}

func (rc *RedisClient) Set(key string, value string, expiration time.Duration) error {
	if rc.Config.Md5Key {
		key = utils.MD5(key)
	}
	err := rc.GetClient().Set(rc.Context, key, value, expiration).Err()

	if err != redis.Nil && err != nil {
		return err
	}

	return nil
}

func (rc *RedisClient) Get(key string) (string, error) {
	if rc.Config.Md5Key {
		key = utils.MD5(key)
	}
	val, err := rc.GetClient().Get(rc.Context, key).Result()

	if err == redis.Nil {
		return "", errors.New("no such key: " + key)
	}
	if err != nil {
		return "", err
	}

	return val, nil
}
