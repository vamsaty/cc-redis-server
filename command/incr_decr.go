package command

import (
	"github.com/vamsaty/cc-redis-server/store"
	"strconv"
)

// RunIncr runs the INCR command
func RunIncr(cache *store.Cache, args ...string) string { return runAdd(cache, 1, args...) }

// RunDecr runs the DECR command
func RunDecr(cache *store.Cache, args ...string) string { return runAdd(cache, -1, args...) }

// runAdd helper function to run the INCR and DECR commands
func runAdd(cache *store.Cache, addNum int, args ...string) string {
	key := args[0]
	if x, _ := cache.Contains(key); x {
		value, _ := cache.Get(key)
		// supported int format
		if intVal, err := strconv.Atoi(value.Value); err == nil {
			value.Value = strconv.Itoa(intVal + addNum)
			if err = cache.Set(value.Key, value); err != nil {
				return "ERR " + err.Error()
			}
			return value.Value
		}
		return "ERR value is not an integer or out of range"
	}
	return "0"
}
