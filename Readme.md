# Write Your Own Redis Server

## Description
The challenge was to use build your own Redis Server.
It runs on port 6379 by default.
It currently supports basic Redis commands like -
1. SET
2. GET
3. PING
4. ECHO

## Usage

Steps to build the binary and execute it -
```
go build -o redis-server ./cmd/main.go && ./redis-server
``` 

## Flags
None (TODO: support flag for port number)