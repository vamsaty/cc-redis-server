package store

import (
	"fmt"
	"sync"
	"time"
)

var once = sync.Once{}
var cacheInstance *Cache

func GetCacheInstance() *Cache {
	once.Do(func() {
		if cacheInstance == nil {
			cacheInstance = &Cache{
				data: make(map[string]Item),
				lock: sync.Mutex{},
			}
		}
	})
	return cacheInstance
}

/* ------------ Cache ------------ */

type Cache struct {
	lock sync.Mutex
	data map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, found := c.data[key]
	if found && item.Expired() {
		delete(c.data, key)
		found = false
	}
	return item, found
}

func (c *Cache) Set(key string, value Item) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = value
	return nil
}

func (c *Cache) Contains(key string) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, found := c.data[key]
	return found, nil
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, key)
}

/* ------------ Item ------------ */

type Item struct {
	Key   string
	Value string
	TTL   int64
}

func (it *Item) Expired() bool { return it.TTL > 0 && time.Now().UnixNano() > it.TTL }

func (it *Item) Serialize() string { return fmt.Sprintf("$%d\r\n%s\r\n", len(it.Value), it.Value) }
