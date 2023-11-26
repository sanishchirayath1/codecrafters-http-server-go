package main

import (
	"fmt"
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
			os.Exit(0)
			continue
		}
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// defer conn.Close()
	fmt.Println("Received connection request", conn)
	response := "HTTP/1.1 200 OK\r\n\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(0)
	}

}
