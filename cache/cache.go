package cache

import (
	"sync"
	"time"
)

type CacheBlockDataType map[string]interface{}

type Cache struct {
	Exp    int64
	blocks map[string]*CacheBlock
}

type CacheBlock struct {
	Data        CacheBlockDataType
	createdTime time.Time
	rwLock      sync.RWMutex
}

func NewCache() *Cache {
	c := &Cache{
		Exp:    time.Hour.Milliseconds(),
		blocks: map[string]*CacheBlock{},
	}

	return c
}

func (c *Cache) Set(key string, data CacheBlockDataType) {
	cb, ok := c.blocks[key]
	if !ok {
		cb = &CacheBlock{}
	}
	cb.createdTime = time.Now()
	cb.Data = data

	c.blocks[key] = cb
}

func (c *Cache) Get(key string) CacheBlockDataType {
	cb, ok := c.blocks[key]
	if !ok {
		return nil
	}

	cb.rwLock.RLock()
	defer cb.rwLock.RUnlock()

	if c.Exp+cb.createdTime.UnixMilli() > time.Now().UnixMilli() {
		return cb.Data
	} else {
		delete(c.blocks, key)
		return nil
	}
}
