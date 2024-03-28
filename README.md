# Go TCP Server

## Overview

 A Go TCP server that accepts connections, receives messages from clients, and broadcasts them to all other connected clients, excluding the sender.

## Running the Server

1. Clone the repo and navigate to the server directory.
2. Run `go run server.go` to start the server on the default port (3001).
3. Use `-port=8080` to specify a different port `go run server.go -port=8080`.
4. Use `-max=100` to set the maximum connections `go run server.go -port=8080 -max=100`.

## Connecting as a Client

- Implement or use a TCP client to connect to the server's IP and port.
- Messages should end with a newline character for correct processing.

