package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ipuppet/gtools/flags"
)

var (
	loggerMap map[string]*log.Logger
)

func init() {
	loggerMap = make(map[string]*log.Logger, 10)
}

func Logger(name string) *log.Logger {
	if loggerMap[name] == nil {
		logPath := BasePath
		if flags.IsParse {
			if !PathExists(flags.LogPath) {
				log.Fatal("logger path not existe: " + flags.LogPath)
			}
			logPath = flags.LogPath
		}
		file := filepath.Join(logPath, "main.log")
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			panic(err)
		}
		loggerMap[name] = log.New(logFile, "["+name+"] ", log.LstdFlags|log.LUTC)
	}

	return loggerMap[name]
}
