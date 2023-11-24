package server

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"net"
)

type RedisServer interface {
	Start() error
	handleConnection(conn net.Conn) error
}

type RedisServerImpl struct {
	RedisExecutor
	RedisTokenizer
	*zap.Logger
}

func NewRedisServer() *RedisServerImpl {
	logger, _ := zap.NewProduction()
	return &RedisServerImpl{
		RedisExecutor:  NewRedisExecutorImpl(),
		RedisTokenizer: DefaultTokenizer(),
		Logger:         logger,
	}
}

func (rs *RedisServerImpl) Start() error {
	rs.Info("starting server")
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		return err
	}
	for {
		c, err := listener.Accept()
		rs.Info("accepted connection", zap.String("remote_addr", c.RemoteAddr().String()))
		if err != nil {
			return err
		}
		go rs.handleConnection(c)
	}
}

// handleConnection parses the request from a client, generates a command, executes the
// command on redis executor and sends the response back to the client
func (rs *RedisServerImpl) handleConnection(conn net.Conn) (err error) {
	connId := zap.String("remote_addr", conn.RemoteAddr().String())
	rs.Info("handling connection", connId)

	defer func() {
		err = conn.Close()
		rs.Info("connection closed", connId)
	}()

	var tokens []string
	reader := bufio.NewReader(conn)

	var sendResponse = func(response string) {
		_, err = conn.Write([]byte(response))
		if err != nil {
			rs.Error("error while writing response", zap.Error(err))
		}
	}

	for {
		// read the request and parse the tokens
		if tokens, err = rs.GetTokens(reader); err != nil {
			if err == io.EOF {
				rs.Info("client closed connection", connId)
				return nil
			}
			rs.Warn("error while reading request", zap.Error(err))
			continue
		}
		if len(tokens) == 0 {
			continue
		}

		// create command from tokens
		execCmd := CreateCommandFromTokens(tokens[1:])
		if execCmd.IsExit() {
			continue
		}

		// execute the command and generate response
		rs.Info("executing command", zap.String("command", execCmd.String()))
		response := rs.Execute(execCmd)

		// send the response to the client
		rs.Info("writing request", zap.String("command", response.Serialize()))

		sendResponse(response.Serialize())
	}
}
