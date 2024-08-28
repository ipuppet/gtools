package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetConfig(path string, v interface{}) {
	configBytes, err := os.ReadFile(filepath.Join(GetBasePath(), path))
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := json.Unmarshal(configBytes, &v); err != nil {
		logger.Fatal(err.Error())
	}
}
