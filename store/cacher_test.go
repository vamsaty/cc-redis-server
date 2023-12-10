package store

import (
	"testing"
	"time"
)

func TestItem(t *testing.T) {
	type TestCase struct {
		name string
		item Item
		want string
	}

	testCases := []TestCase{
		{
			name: "Test Item debug",
			item: Item{
				Key:   "key",
				Value: "test",
			},
			want: "{key: test}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.item.Debug()
			if got != tc.want {
				t.Errorf("Item.Debug() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCache(t *testing.T) {
	cache := GetCacheInstance()

	// no expiry
	cache.Set("key", Item{
		Key:   "key",
		Value: "test",
		TTL:   time.Now().UnixNano() + 1*time.Second.Nanoseconds(),
	})
	var found bool

	_, found = cache.Get("key")
	if !found {
		t.Errorf("Item should exist")
	}

	time.Sleep(2 * time.Second)
	_, found = cache.Get("key")
	if found {
		t.Errorf("Item should have expired")
	}
}
