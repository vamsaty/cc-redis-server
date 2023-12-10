package store

import (
	"fmt"
	"sync"
	"time"
)

type Cacher interface {
	Get(key string) (Item, bool)
	Set(key string, value Item) error
	Contains(key string) (bool, error)
	Delete(key string)
	Debug() map[string]string
}

var once = sync.Once{}
var cacheInstance Cacher

// GetCacheInstance returns a singleton instance of the cache
func GetCacheInstance() Cacher {
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

// Cache satisfies the Cacher interface
type Cache struct {
	lock sync.Mutex
	data map[string]Item
}

// Get returns an item from the cache
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

// Set adds an item to the cache
func (c *Cache) Set(key string, value Item) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = value
	return nil
}

// Contains returns true if the key exists in the cache
func (c *Cache) Contains(key string) (bool, error) {
	_, found := c.Get(key)
	return found, nil
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, key)
}

// Debug is a helper function to print the cache
func (c *Cache) Debug() map[string]string {
	data := make(map[string]string)
	for key, item := range c.data {
		fmt.Printf("DEBUG: key=(%s), item=(%s), expired=(%v)\n", key, item.Debug(), item.Expired())
		data[key] = item.Debug()
	}
	return data
}

/* ------------ Item ------------ */

type Item struct {
	Key   string
	Value string
	TTLns int64 // stored as nanoseconds
}

// Expired returns true if the item has expired
func (it *Item) Expired() bool {
	return it.TTLns > 0 && time.Now().UnixNano() > it.TTLns
}

// Serialize returns a string representation of the item
func (it *Item) Serialize() string { return fmt.Sprintf("$%d\r\n%s\r\n", len(it.Value), it.Value) }

// Debug returns a string representation of the item for debugging purposes
func (it *Item) Debug() string { return fmt.Sprintf("{%s: %s}", it.Key, it.Value) }
