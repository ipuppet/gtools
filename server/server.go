package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/middleware"
)

var (
	logger *log.Logger
)

func SetLogger(l *log.Logger) {
	logger = l
}

func GetServer(addr string, handle func(engine *gin.Engine)) *http.Server {
	engine := gin.Default()

	engine.Use(middleware.ErrorHandler())

	handle(engine)

	logger.Println("server listening on: " + addr)

	return &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
