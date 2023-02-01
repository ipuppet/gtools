package flags

import (
	"flag"
)

var (
	ConfigPath string
	LogPath    string
	IsParse    bool = false
)

const (
	configPathDefault = "./config"
	configPathUsage   = "Set config file path."

	logPathDefault = "./var"
	logPathUsage   = "Set log file path."
)

func init() {
	flag.StringVar(&ConfigPath, "config", configPathDefault, configPathUsage)

	flag.StringVar(&LogPath, "log", logPathDefault, logPathUsage)
}

func Parse() {
	flag.Parse()
	IsParse = true
}
