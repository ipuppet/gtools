package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	loggerMap map[string]*log.Logger
	LogPath   string
)

func init() {
	loggerMap = make(map[string]*log.Logger, 10)
}

func Logger(name string) *log.Logger {
	if loggerMap[name] == nil {
		var out io.Writer
		if PathExists(LogPath) {
			file := filepath.Join(LogPath, "main.log")
			logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
			if err != nil {
				panic(err.Error())
			}
			out = logFile
		} else {
			out = os.Stderr
			log.Panicln("path " + LogPath + " does not exist")
		}
		// 默认输出到 os.Stderr
		loggerMap[name] = log.New(out, "["+name+"] ", log.LstdFlags|log.LUTC)
	}

	return loggerMap[name]
}
