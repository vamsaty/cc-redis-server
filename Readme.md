# Write Your Own Redis Server

## Description
The challenge was to use build your own Redis Server.
It runs on port 6379 by default.
It currently supports basic Redis commands like -
1. SET
2. GET
3. PING
4. ECHO
---
## Usage

Steps to build the binary and execute it -
```
go build -o redis-server ./cmd/main.go && ./redis-server
``` 
---
## Flags

----


## File Structure

### executor
Defines the supported commands. Executes the commands in the datastore and generates a response.

### server
Contains the code for the server. Starts a listener (at 6379 port) and connection handler (concurrent).

### tokenizer
Contains the code for the tokenizer. This is used by the server to parse the commands sent by the client.

### types
Defines the RESP datatypes and helper function to convert RESP types to string & vice-versa.

---