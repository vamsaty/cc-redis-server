package main

import "redis-server/cmd/server/server"

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	server.NewRedisServer()
	s := server.NewRedisServer()
	PanicIf(s.Start())
}
