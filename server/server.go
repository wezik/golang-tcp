package main

import (
	"encoding/binary"
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
        VERSION byte = 1
        HEADER_SIZE = 4
)

type TCPCommand struct {
        Command byte
        Data []byte
}

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
        bytes := make([]byte, 1024)
        _, err := conn.Read(bytes)
        handleErr(err)

        tcpCommand, err := unmarshalBinary(bytes)
        handleErr(err)

        // Read and respond
        if tcpCommand.Command == 1 {
                timeNow := time.Now().Format(time.RFC3339)
                response := fmt.Sprintf("Received command: %d, Data: %s, Received time: %s", tcpCommand.Command, string(tcpCommand.Data), timeNow)

                responseData, err := marshalBinary(TCPCommand{Command: 1, Data: []byte(response)})
                handleErr(err)

                _, err = conn.Write(responseData)
                handleErr(err)
        }
        
        conn.Close()
}

func unmarshalBinary(bytes []byte) (TCPCommand, error) {
        if bytes[0] != VERSION {
                return TCPCommand{}, fmt.Errorf("version mismatch %d != %d", bytes[0], VERSION)
        }

        length := int(binary.BigEndian.Uint16(bytes[2:]))
        end := HEADER_SIZE + length
        
        if len(bytes) < end {
                return TCPCommand{}, fmt.Errorf("invalid length %d < %d", len(bytes), end)
        }

        t := TCPCommand{}

        t.Command = bytes[1]
        t.Data = bytes[HEADER_SIZE:end]

        return t, nil
}

func marshalBinary(t TCPCommand) ([]byte, error) {
        length := uint16(len(t.Data))
        lengthBytes := make([]byte, 2)
        binary.BigEndian.PutUint16(lengthBytes, length)

        bytes := make([]byte, 0, HEADER_SIZE + length)
        bytes = append(bytes, VERSION)
        bytes = append(bytes, t.Command)
        bytes = append(bytes, lengthBytes...)
        return append(bytes, t.Data...), nil
}

func handleErr(err error) {
        if (err != nil) {
                log.Fatal(err)
                os.Exit(1)
        }
}
