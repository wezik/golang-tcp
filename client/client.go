package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
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
        fmt.Println("[CLIENT] Connecting to server...")

        tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
        handleErr(err)

        scanner := bufio.NewScanner(os.Stdin)

        fmt.Println("[CLIENT] Client is running...")
        for {
                fmt.Printf("[CLIENT] 1. Send | x. Exit : ")

                var input string
                fmt.Scanln(&input)

                if input == "1" {
                        fmt.Printf("[CLIENT] Enter data to send : ")

                        scanner.Scan()
                        handleErr(scanner.Err())
                        handleErr(err)

                        tcpConn, err := net.DialTCP(TYPE, nil, tcpServer)
                        handleErr(err)

                        sendData(tcpConn, []byte(scanner.Text()))
                        handleErr(err)

                        fmt.Println("[CLIENT] Data sent to server")
                        t := readResponse(tcpConn);
                        tcpConn.Close()

                        fmt.Println("[SERVER]", string(t.Data))
                } else if input == "x" {
                        fmt.Println("[CLIENT] Exiting...")
                        break
                }

        }
}

func sendData(tcpConn *net.TCPConn, data []byte) {
        t := TCPCommand{1, data}
        bytes, err := marshalBinary(t)
        handleErr(err)
        _, err = tcpConn.Write(bytes)
        handleErr(err)
}

func readResponse(tcpConn *net.TCPConn) TCPCommand {
        buffer := make([]byte, 1024)
        _, err := tcpConn.Read(buffer)
        handleErr(err)
        t, err := unmarshalBinary(buffer)
        handleErr(err)
        return t
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
        fmt.Println("[CLIENT] Sending payload")
        length := uint16(len(t.Data))
        lengthBytes := make([]byte, 2)
        binary.BigEndian.PutUint16(lengthBytes, length)

        bytes := make([]byte, 0, HEADER_SIZE + length)
        bytes = append(bytes, VERSION)
        bytes = append(bytes, t.Command)
        bytes = append(bytes, lengthBytes...)
        fmt.Println("[CLIENT] Header :", bytes)
        fmt.Println("[CLIENT] Data :", t.Data)
        return append(bytes, t.Data...), nil
}

func handleErr(err error) {
        if (err != nil) {
                log.Fatal(err)
                os.Exit(1)
        }
}
