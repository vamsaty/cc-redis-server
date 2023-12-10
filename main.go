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
	"time"
)

var cache = store.GetCacheInstance()

// handleConnection handles a connection
func handleConnection(conn net.Conn) {
	fmt.Println("connection received", conn.RemoteAddr())

	var err error

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing connection", err)
		}
		fmt.Println("connection closed", conn.RemoteAddr())
	}()

	// read the request
	reader := bufio.NewReader(conn)
	var token []rune

	for {
		// set a read deadline of 10 seconds
		future := time.Now().Add(time.Second * 1)
		err = conn.SetReadDeadline(future)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("error setting read deadline", err)
			continue
		}

		// read the token
		token, err = resp.ReadToken(reader)
		if err != nil {
			if err == io.EOF {
				log.Println("end of token", err)
				return
			}
			log.Println("error reading token", err)
			continue
		}

		// the request is always an array
		if token[0] == '*' {
			items := resp.ParseArray(token, reader)
			output := command.Execute(items, cache)
			output = fmt.Sprintf("+%s\r\n", output)

			_, err = conn.Write([]byte(output))
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println("error writing response", err)
				continue
			}
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
