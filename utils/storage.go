package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

var StoragePath string

func getStoragePath(app string, fileName string) string {
	if StoragePath == "" {
		StoragePath, _ = os.Getwd()
		StoragePath = filepath.Join(StoragePath, "storage")
	}
	return filepath.Join(StoragePath, app, fileName)
}

func GetStorageContent(app string, fileName string) (string, error) {
	// 读取文件
	content, err := os.ReadFile(getStoragePath(app, fileName))

	return string(content), err
}

func GetStorageJSON(app string, fileName string, v interface{}) error {
	// 读取文件
	content, err := os.ReadFile(getStoragePath(app, fileName))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	return err
}

func SetStorageContent(app string, fileName string, content string) error {
	// 写入文件
	return os.WriteFile(getStoragePath(app, fileName), []byte(content), 0666)
}

func AppendStorageContent(app string, fileName string, content string) error {
	file, err := os.OpenFile(getStoragePath(app, fileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(content + "\n")
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}
