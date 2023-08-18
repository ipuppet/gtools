package database

import (
	"context"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type DatabaseConfig struct {
	Driver   string            `json:"driver"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Args     map[string]string `json:"args"`
}

var (
	logger        *log.Logger
	rdb           *redis.Client
	rdbKeyPrefix  string
	rdbCtx        context.Context
	rdbExpiration time.Duration
)

func SetLogger(l *log.Logger) {
	logger = l
}

func SetRedis(r *redis.Client) {
	rdb = r
	rdbKeyPrefix = "gtools-database-cache"
	rdbCtx = context.Background()
	rdbExpiration = 24 * time.Hour
}

func SetRedisWithExpiration(r *redis.Client, expiration time.Duration) {
	SetRedis(r)
	rdbExpiration = expiration
}

func CleanCache() error {
	iter := rdb.Scan(rdbCtx, 0, rdbKeyPrefix+"*", 0).Iterator()
	for iter.Next(rdbCtx) {
		err := rdb.Del(rdbCtx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
