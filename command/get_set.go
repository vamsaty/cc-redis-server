package command

import (
	"github.com/vamsaty/cc-redis-server/store"
	"strconv"
	"time"
)

// RunSet sets the value of the key
func RunSet(cache *store.Cache, args ...string) string {
	item := store.Item{
		Key:   args[0],
		Value: args[1],
		TTL:   0,
	}
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "xx":
			// don't set if key already exists
			if _, found := cache.Get(args[0]); found {
				return "OK"
			}
		case "nx":
			// don't set if key doesn't exist
			if _, found := cache.Get(args[0]); !found {
				return "OK"
			}
		case "ex":
			x, _ := strconv.Atoi(args[i+1])
			item.TTL = time.Now().Add(time.Duration(x) * time.Second).UnixNano()
			i++
		case "px":
			x, _ := strconv.Atoi(args[i+1])
			item.TTL = time.Now().Add(time.Duration(x) * time.Millisecond).UnixNano()
			i++
		case "pxat":
			pxat, _ := strconv.Atoi(args[i+1])
			item.TTL = int64(pxat) * 1000000
			i++
		case "exat":
			exat, _ := strconv.Atoi(args[i+1])
			item.TTL = int64(exat) * 1000000000
			i++
		}
	}
	cache.Set(args[0], item)
	return "OK"
}

// RunGet returns the value of the key if present, else "(nil)"
func RunGet(cache *store.Cache, args ...string) string {
	if item, found := cache.Get(args[0]); found { // is present
		return item.Value
	}
	return "(nil)"
}
