package main

import (
    "bufio"
    "flag"
    "log"
    "net"
    "os"
    "sync"
    "time"
)

var (
    maxConnections int
    port           string
    clients        = make(map[net.Conn]struct{})
    rwMutex        = &sync.RWMutex{}
    logger         = log.New(os.Stdout, "server: ", log.LstdFlags)
)

func init() {
    flag.IntVar(&maxConnections, "max", 300, "maximum number of concurrent connections")
    flag.StringVar(&port, "port", "12345", "port to listen on")
}

func main() {
    flag.Parse()

    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        logger.Fatalf("Failed to start server on port %s: %s", port, err)
    }
    defer listener.Close()
    logger.Printf("Server started on port %s, max connections: %d\n", port, maxConnections)

    var connections int
    for {
        if connections >= maxConnections {
            logger.Println("Maximum connections reached, new connections will be temporarily refused.")
            time.Sleep(10 * time.Second)
            continue
        }

        conn, err := listener.Accept()
        if err != nil {
            logger.Printf("Error accepting connection: %s", err)
            continue
        }
        conn.SetDeadline(time.Now().Add(15 * time.Minute))

        rwMutex.Lock()
        clients[conn] = struct{}{}
        connections = len(clients)
        rwMutex.Unlock()

        go handleConnection(conn, &connections)
    }
}

func handleConnection(conn net.Conn, connections *int) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            logger.Printf("Connection error from %v: %s", conn.RemoteAddr(), err)
            rwMutex.Lock()
            delete(clients, conn)
            *connections = len(clients)
            rwMutex.Unlock()
            return
        }

        logger.Printf("Message received from %v: %s", conn.RemoteAddr(), message)
        broadcastMessage(message, conn)
    }
}

func broadcastMessage(message string, origin net.Conn) {
    rwMutex.Lock()
    defer rwMutex.Unlock()
    for conn := range clients {
        if conn != origin {
            go func(c net.Conn) {
                _, err := c.Write([]byte(message))
                if err != nil {
                    logger.Printf("Broadcast error to %v: %s", c.RemoteAddr(), err)
                    delete(clients, c)
                }
            }(conn)
        }
    }
}
