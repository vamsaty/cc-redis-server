package server

import (
	"fmt"
	"go.uber.org/zap"
)

// RedisExecutor is the interface for executing commands on Redis server
type RedisExecutor interface {
	Execute(cmd *Cmd) *RedisResponse
}

// RedisExecutorImpl executes the commands on Redis datastore
type RedisExecutorImpl struct {
	*zap.Logger
	RedisCacher
	requestChan chan *Cmd
}

func NewRedisExecutorImpl() *RedisExecutorImpl {
	logger, _ := zap.NewProduction()
	return &RedisExecutorImpl{
		RedisCacher: GetCacherInstance(),
		requestChan: make(chan *Cmd, 1000),
		Logger:      logger,
	}
}

// Execute executes the command on Redis datastore
func (re *RedisExecutorImpl) Execute(cmd *Cmd) *RedisResponse {
	var argValue interface{}

	response := &RedisResponse{
		Error: fmt.Errorf("invalid command received"),
	}

	defer func() {
		re.Info("RESPONSE", zap.String("command", response.Serialize()))
	}()
	re.Info("executing command", zap.String("command", cmd.Name()))

	// TODO : Use enums for Cmd type instead
	switch cmd.Name() {

	case "get":
		argValue = cmd.GetArg("key")
		key := argValue.(string)

		response = NotFoundResponse(key)
		if value, found := re.Get(key); found {
			response = &RedisResponse{Item: value}
		}

	case "set":
		symbol := cmd.GetArg("type").(byte)
		respType, _ := GetRespType(symbol)
		item := CacheItem{
			Key:      cmd.GetArg("key").(string),
			Value:    fmt.Sprintf("%v", cmd.GetArg("value")),
			DataType: respType,
		}

		response = OKResponse()
		if err := re.Set(item.GetKey(), item); err != nil {
			response = ErrorResponse(err)
		}

	case "ping":
		response = &RedisResponse{
			Item: CacheItem{
				Key:   "ping",
				Value: "PONG",
			},
		}

	case "echo":
		response = &RedisResponse{
			Item: CacheItem{
				Key:   "echo",
				Value: cmd.GetArg("value").(string),
			},
		}
	}

	return response
}
