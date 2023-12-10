package command

import (
	"github.com/vamsaty/cc-redis-server/store"
)

// ExecMap is a map of command to function
var ExecMap = map[string]func(store.Cacher, ...string) string{
	"SET":    RunSet,
	"GET":    RunGet,
	"PING":   RunPing,
	"ECHO":   RunEcho,
	"EXISTS": RunExists,
	"DEL":    RunDel,
	"INCR":   RunIncr,
	"DECR":   RunDecr,
	// DEBUG commands
	"DEBUG": RunDebug,
}

// Execute executes the command
func Execute(args [][]string, cache store.Cacher) string {
	var data []string
	for i := range args {
		data = append(data, args[i][1])
	}

	execFunc, ok := ExecMap[data[0]]
	if !ok {
		return "ERR unknown command '" + data[0] + "'"
	}
	return execFunc(cache, data[1:]...)
}
