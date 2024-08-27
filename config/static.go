package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetConfig(path string, v interface{}) {
	configString, err := os.ReadFile(filepath.Join(BasePath, path))
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := json.Unmarshal(configString, &v); err != nil {
		logger.Fatal(err.Error())
	}
}
