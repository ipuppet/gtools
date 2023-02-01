package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetConfig(filename string, v interface{}) {
	configString, err := os.ReadFile(filepath.Join(ConfigPath, filename))
	if err != nil {
		logger.Fatal(err)
	}

	if err := json.Unmarshal(configString, &v); err != nil {
		logger.Fatal(err)
	}
}
