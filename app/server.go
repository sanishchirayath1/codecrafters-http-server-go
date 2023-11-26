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

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Unable to close connection: ", err)
		}
	}(l)

	fmt.Println("Listening on 4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	response := "HTTP/1.1 200 OK\r\n\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}

	err = conn.Close()
	if err != nil {
		fmt.Println("Unable to close connection ", err)

	}
}
