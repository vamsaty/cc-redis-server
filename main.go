package main

import (
	"bufio"
	"fmt"
	"github.com/vamsaty/cc-redis-server/command"
	"github.com/vamsaty/cc-redis-server/resp"
	"github.com/vamsaty/cc-redis-server/store"
	. "github.com/vamsaty/cc-redis-server/utils"
	"io"
	"log"
	"net"
)

var cache = store.GetCacheInstance()

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing connection", err)
		}
	}()

	// read the request
	reader := bufio.NewReader(conn)
	for {
		token, err := resp.ReadToken(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			PanicIf(err)
			continue
		}
		if token[0] == '*' {
			items := resp.ParseArray(token, reader)
			output := command.Execute(items, cache)
			output = fmt.Sprintf("+%s\r\n", output)
			_, err = conn.Write([]byte(output))
			PanicIf(err)
		}
	}
}

func main() {

	// listener
	listener, err := net.Listen("tcp", ":6379")
	PanicIf(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}

}
