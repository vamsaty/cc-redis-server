package command

import "github.com/vamsaty/cc-redis-server/store"

func RunPing(_ *store.Cache, args ...string) string {
	return "PONG"
}

func RunEcho(_ *store.Cache, args ...string) string {
	return args[0]
}
