package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ipuppet/gtools/utils"
)

var (
	logger   *log.Logger
	BasePath string
)

func SetLogger(l *log.Logger) {
	logger = l
}

type Config struct {
	Path           string
	_path          string
	data           map[string]interface{}
	lastModifyTime int64
	rwLock         sync.RWMutex
	notifyList     []Notifyer
}

func New(path string) *Config {
	config := &Config{
		Path: path,
	}
	config.Init()

	return config
}

func (c *Config) Init() {
	m, err := c.parse()
	if err != nil {
		logger.Println(err)
		return
	}
	c.data = m
	go c.reload()
}

func (c *Config) path() string {
	if c._path == "" {
		c._path = filepath.Join(BasePath, c.Path)
		if !utils.PathExists(c._path) {
			log.Fatal("config dose not existe: " + c._path)
		}
	}

	return c._path
}

func (c *Config) AddNotifyer(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}

func (c *Config) parse() (m map[string]interface{}, err error) {
	m = make(map[string]interface{}, 50)
	configString, err := os.ReadFile(c.path())
	if err != nil {
		return
	}

	if err = json.Unmarshal(configString, &m); err != nil {
		return
	}

	return
}

func (c *Config) reload() {
	// 每 5 秒重新加载一次配置文件
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		func() {
			file, err := os.Open(c.path())
			if err != nil {
				logger.Printf("open %s failed, err: %v\n", c.Path, err)
				return
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				logger.Printf("stat %s failed, err: %v\n", c.Path, err)
				return
			}

			curModifyTime := fileInfo.ModTime().Unix()

			// 判断文件的修改时间是否大于最后一次修改时间
			if curModifyTime > c.lastModifyTime {
				m, err := c.parse()
				if err != nil {
					logger.Println("parse failed, err: ", err)
					return
				}

				c.rwLock.Lock()
				c.data = m
				c.rwLock.Unlock()

				c.lastModifyTime = curModifyTime

				for _, n := range c.notifyList {
					n.Callback(c)
				}
			}
		}()
	}
}

func (c *Config) Get(key string) (value interface{}, err error) {
	// 根据字符串获取
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		err = errors.New("key[" + key + "] not found")
		return
	}
	return
}

func (c *Config) ShouldGet(key string) (value interface{}) {
	// 根据字符串获取
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		value = nil
		return
	}
	return
}
