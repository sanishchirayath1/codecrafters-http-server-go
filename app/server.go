package main

import (
	"fmt"
	"net"
	"os"
)

var (
	IP      = "0.0.0.0"
	PORT    = "4221"
	IP_PORT = IP + ":" + PORT
)

func main() {
	l, err := net.Listen("tcp", IP_PORT)
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

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Unable to close connection: ", err)
		}
	}(l)

	fmt.Println("Listening on " + IP_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error while sending response ", err)
		}
		err = conn.Close()
		if err != nil {
			fmt.Println("Unable to close connection ", err)
		}
	}
}
