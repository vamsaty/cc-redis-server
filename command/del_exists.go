package command

import (
	"fmt"
	"github.com/vamsaty/cc-redis-server/store"
)

// RunExists checks if the key exists in the cache
func RunExists(cache *store.Cache, args ...string) string {
	count := 0
	for i := range args {
		if x, _ := cache.Contains(args[i]); x {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

// RunDel deletes the key from the cache
func RunDel(cache *store.Cache, args ...string) string {
	count := 0
	for i := range args {
		if x, _ := cache.Contains(args[i]); x {
			cache.Delete(args[i])
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}
