package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
        HOST = "localhost"
        PORT = "8080"
        TYPE = "tcp"
)

func main() {

        listen, err := net.Listen(TYPE, HOST+":"+PORT)
        handleErr(err)
        defer listen.Close()

        fmt.Printf("Server is listening on %s:%s\n", HOST, PORT)
        for {
                conn, err := listen.Accept()
                handleErr(err)
                go handleRequest(conn)
        }
}

func handleRequest(conn net.Conn) {
        buffer := make([]byte, 1024)
        _, err := conn.Read(buffer)
        handleErr(err)

        time := time.Now().Format(time.ANSIC)
        response := fmt.Sprintf("Hello, %s, received time: %s", string(buffer), time)
        _, err = conn.Write([]byte(response))
        handleErr(err)
        
        conn.Close()
}

func handleErr(err error) {
        if (err != nil) {
                log.Fatal(err)
                os.Exit(1)
        }
}
