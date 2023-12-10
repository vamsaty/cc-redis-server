package command

import "github.com/vamsaty/cc-redis-server/store"

// RunPing returns PONG
func RunPing(_ *store.Cache, args ...string) string { return "PONG" }

// RunEcho returns the first argument sent to it
func RunEcho(_ *store.Cache, args ...string) string { return args[0] }
