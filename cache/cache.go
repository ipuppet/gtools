package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	Exp    int64
	MaxLen int
	keys   *list.List
	blocks map[string]*CacheBlock
}

type CacheBlock struct {
	Data        interface{}
	createdTime time.Time
	rwLock      sync.RWMutex
}

func New() *Cache {
	c := &Cache{
		Exp:    time.Hour.Milliseconds(), // default 1h
		MaxLen: 20,
		keys:   list.New(),
		blocks: map[string]*CacheBlock{},
	}

	return c
}

func (c *Cache) Clean() {
	c.blocks = map[string]*CacheBlock{}
	c.keys = list.New()
}

func (c *Cache) Delete(key string) {
	delete(c.blocks, key)

	for k := c.keys.Front(); k != nil; k = k.Next() {
		if k.Value == key {
			c.keys.Remove(k)
			break
		}
	}
}

func (c *Cache) Set(key string, data interface{}) {
	cb, ok := c.blocks[key]
	if !ok {
		cb = &CacheBlock{}
		for c.keys.Len() >= c.MaxLen {
			c.Delete(c.keys.Front().Value.(string))
		}
		c.keys.PushBack(key)
	}

	cb.createdTime = time.Now()
	cb.Data = data

	c.blocks[key] = cb
}

func (c *Cache) Get(key string) interface{} {
	cb, ok := c.blocks[key]
	if !ok {
		return nil
	}

	cb.rwLock.RLock()
	defer cb.rwLock.RUnlock()

	if c.Exp+cb.createdTime.UnixMilli() > time.Now().UnixMilli() {
		return cb.Data
	} else {
		c.Delete(key)
		return nil
	}
}
