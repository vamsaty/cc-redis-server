package command

import (
	"fmt"
	"github.com/vamsaty/cc-redis-server/store"
	"testing"
	"time"
)

func padArgs(args []string) [][]string {
	var data [][]string
	for i := range args {
		data = append(data, []string{"", args[i]})
	}
	return data
}

func TestExecMap(t *testing.T) {
	cache := store.GetCacheInstance()
	// handle set command

	type TestCase struct {
		name string
		args []string
		want string
	}

	testCases := []TestCase{
		{
			name: "SET",
			args: []string{"SET", "key", "123"},
			want: "OK",
		},
		{
			name: "GET",
			args: []string{"GET", "key"},
			want: "123",
		},
		{
			name: "PING",
			args: []string{"PING"},
			want: "PONG",
		},
		{
			name: "ECHO",
			args: []string{"ECHO", "hello"},
			want: "hello",
		},
		{
			name: "EXISTS",
			args: []string{"EXISTS", "key"},
			want: "1",
		},
		{
			name: "INCR",
			args: []string{"INCR", "key"},
			want: "124",
		},
		{
			name: "DECR",
			args: []string{"DECR", "key"},
			want: "123",
		},
		{
			name: "DEL",
			args: []string{"DEL", "key"},
			want: "1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := Execute(padArgs(tc.args), cache)
			if response != tc.want {
				t.Errorf("Execute() = %v, want %v", response, tc.want)
			}
		})
	}
}

func TestSetParsing(t *testing.T) {
	cacher := store.GetCacheInstance()
	type TestCase struct {
		name string
		args []string
		want func()
	}

	var wantFutureTLL = func() {
		if item, found := cacher.Get("ttl_key"); !found {
			t.Errorf("Expected key to be found")
		} else if item.TTLns < time.Now().UnixNano() {
			t.Errorf("Expected TTL to be in the future, got %v", item.TTLns)
		}
	}

	testCases := []TestCase{
		{
			name: "Set key value",
			args: []string{"SET", "key", "10"},
			want: nil,
		},
		{
			name: "Don't set if already exists",
			args: []string{"SET", "key", "123", "xx"},
			want: func() {
				cacher.Debug()
				if value, found := cacher.Get("key"); !found {
					t.Errorf("Expected key to be found")
				} else if value.Value != "10" {
					t.Errorf("Expected value to be 10, got %v", value.Value)
				}
			},
		},
		{
			name: "Don't set if key doesn't exist, set if exists",
			args: []string{"SET", "key", "123", "nx"},
			want: func() {
				if value, found := cacher.Get("key"); !found {
					t.Errorf("Expected key to be found")
				} else if value.Value != "123" {
					t.Errorf("Expected value to be 123, got %v", value.Value)
				}
			},
		},
		{
			name: "Set with TTL in seconds",
			args: []string{"SET", "ttl_key", "123", "ex", "1"},
			want: wantFutureTLL,
		},
		{
			name: "Set with TTL in milliseconds",
			args: []string{"SET", "ttl_key", "123", "px", "1000"},
			want: wantFutureTLL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := Execute(padArgs(tc.args), cacher)
			fmt.Println("RESPONSE", response)
			if tc.want != nil {
				tc.want()
			}
		})
	}
}
