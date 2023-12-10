package main

import (
	"bufio"
	"fmt"
	"github.com/vamsaty/cc-redis-server/command"
	"github.com/vamsaty/cc-redis-server/resp"
	"github.com/vamsaty/cc-redis-server/store"
	ccUtils "github.com/vamsaty/cc-utils"
	"io"
	"log"
	"net"
)

var cache = store.GetCacheInstance()

// handleConnection handles a connection
func handleConnection(conn net.Conn) {
	fmt.Println("connection received")

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing connection", err)
		}
		fmt.Println("connection closed")
	}()

	// read the request
	reader := bufio.NewReader(conn)
	for {
		token, err := resp.ReadToken(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			ccUtils.PanicIf(err)
			continue
		}
		// the request is always an array
		if token[0] == '*' {
			items := resp.ParseArray(token, reader)
			output := command.Execute(items, cache)
			output = fmt.Sprintf("+%s\r\n", output)
			_, err = conn.Write([]byte(output))
			ccUtils.PanicIf(err)
		}
	}
}

func main() {
	fmt.Println("Starting redis server...")
	listener, err := net.Listen("tcp", ":6379")
	ccUtils.PanicIf(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}

}
