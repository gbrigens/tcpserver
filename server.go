package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "sync"
)

var (
    clients = make(map[net.Conn]struct{})
    mutex   = &sync.Mutex{}
)

func main() {
    var port string
    flag.StringVar(&port, "port", "3001", "port to listen on")
    flag.Parse()

    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    fmt.Printf("Server started on port %s\n", port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        mutex.Lock()
        clients[conn] = struct{}{}
        mutex.Unlock()

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            mutex.Lock()
            delete(clients, conn)
            mutex.Unlock()
            break
        }

        fmt.Print("Message received:", message)
        broadcastMessage(message, conn)
    }
}

func broadcastMessage(message string, origin net.Conn) {
    mutex.Lock()
    defer mutex.Unlock()
    for conn := range clients {
        if conn != origin {
            conn.Write([]byte(message))
        }
    }
}

// Navigate to the directory of server.go
// You can run `go run server.go` or to specify a different port, use the -port flag `go run server.go -port=8080`
