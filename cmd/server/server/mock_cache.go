package server

import (
	"sync"
)

type RedisCacher interface {
	Get(key string) (CacheItem, bool)
	Set(key string, value CacheItem) error
	Contains(key string) (bool, error)
	Remove(key string) error
}

type RedisCacherImpl struct {
	store map[string]CacheItem
}

var cacherInstance RedisCacher
var once = &sync.Once{}

func GetCacherInstance() RedisCacher {
	once.Do(func() {
		cacherInstance = NewRedisCacherImpl()
	})
	return cacherInstance
}

func (r *RedisCacherImpl) Get(key string) (CacheItem, bool) {
	item, found := r.store[key]
	return item, found
}

func (r *RedisCacherImpl) Set(key string, value CacheItem) error {
	r.store[key] = value
	return nil
}

func (r *RedisCacherImpl) Contains(key string) (bool, error) {
	return false, nil
}

func (r *RedisCacherImpl) Remove(key string) error {
	return nil
}

func NewRedisCacherImpl() *RedisCacherImpl {
	return &RedisCacherImpl{
		store: make(map[string]CacheItem),
	}
}
