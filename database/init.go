package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ipuppet/gtools/cache"
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
	logger  *log.Logger
	dbCache *cache.Cache
)

func init() {
	dbCache = cache.New()
}

func SetLogger(l *log.Logger) {
	logger = l
}

func CleanCache() {
	dbCache.Clean()
}
