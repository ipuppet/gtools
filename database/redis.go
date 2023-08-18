package database

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func SetRedisCache(key string, value string) error {
	key = rdbKeyPrefix + key
	err := rdb.Set(rdbCtx, key, value, rdbExpiration).Err()
	if err != redis.Nil && err != nil {
		return err
	}

	return nil
}

func GetRedisCache(key string) (string, error) {
	key = rdbKeyPrefix + key
	val, err := rdb.Get(rdbCtx, key).Result()

	if err == redis.Nil {
		return "", errors.New("no such key: " + key)
	}
	if err != nil {
		return "", err
	}

	return val, nil
}
