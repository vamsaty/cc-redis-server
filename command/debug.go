package command

import (
	"fmt"
	"github.com/vamsaty/cc-redis-server/store"
)

func RunDebug(cache store.Cacher, _ ...string) string {
	m := cache.Debug()
	fmt.Println("DEBUG", m)
	return fmt.Sprintf("%v", m)
}
