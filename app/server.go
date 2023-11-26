package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	defer l.Close()

	fmt.Println("Waiting for connections...")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Received connection request", conn)
	response := "HTTP/1.1 200 OK\r\n\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(0)
	}
	// Wait for the client to close the connection
    buf := make([]byte, 1024)
    for {
        _, err := conn.Read(buf)
        if err != nil {
            if err != io.EOF {
                fmt.Println("Error reading from connection: ", err.Error())
            }
            break
        }
    }
    fmt.Println("Client has closed the connection")

}
