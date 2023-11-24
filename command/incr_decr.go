package command

import (
	"github.com/vamsaty/cc-redis-server/store"
	"strconv"
)

func RunIncr(cache *store.Cache, args ...string) string {
	key := args[0]
	if x, _ := cache.Contains(key); x {
		value, _ := cache.Get(key)
		// supported int format
		if intVal, err := strconv.Atoi(value.Value); err == nil {
			value.Value = strconv.Itoa(intVal + 1)
			cache.Set(value.Key, value)
			return value.Value
		}
		return "ERR value is not an integer or out of range"
	}
	return "0"
}

func RunDecr(cache *store.Cache, args ...string) string {
	key := args[0]
	if x, _ := cache.Contains(key); x {
		value, _ := cache.Get(key)
		// supported int format
		if intVal, err := strconv.Atoi(value.Value); err == nil {
			value.Value = strconv.Itoa(intVal - 1)
			cache.Set(value.Key, value)
			return value.Value
		}
		return "ERR value is not an integer or out of range"
	}
	return "0"
}
