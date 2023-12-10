# Write Your Own Redis Server

## Description
The challenge was to use build your own Redis Server.
It runs on port 6379 by default.
It currently supports basic Redis commands like -
1. SET
2. GET
3. PING
4. ECHO
5. DEL
6. EXISTS
7. INCR
8. DECR
---
## Usage
Steps to build the binary and execute it -
```
go build -o redis-server ./cmd/main.go && ./redis-server
``` 
---
## File Structure

### command
Contains the code for the commands executed by the server. Each file can contain multiple commands.
The commands implemented in the file are mentioned in the file's name.
e.g. del_exists.go -> contains DEL and EXISTS commands.

### resp
Contains the code for the RESP protocol. This is used by the server to parse the commands sent by the client.

### store
Contains the code for the data store. This is used by the server to store/fetch the data.

---

###### PS: Humble request to provide feedback on the code and the challenge. Kindly raise a bug in case of any issues.
###### Thanks!