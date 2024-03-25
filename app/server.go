package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	IP      = "0.0.0.0"
	PORT    = "4221"
	IP_PORT = IP + ":" + PORT
)

const CRLF = "\r\n"
const HTTP_OK = "HTTP/1.1 200 OK" + CRLF + CRLF
const HTTP_NOT_FOUND = "HTTP/1.1 404 Not Found" + CRLF + CRLF

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

		request := make([]byte, 4096)
		reqSize, err := conn.Read(request)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
		}

		httpProperties := extractHttpProperties(request, reqSize)

		/**
		Req Method = 0, (GET / POST/ ETC)
		Req Url = 1, (http://example.com)
		Proto Type = 2 (Http/1.1)
		*/
		reqUrl := httpProperties[1]
		response := HTTP_OK

		if reqUrl != "/" {
			response = HTTP_NOT_FOUND
		}

		if reqUrl != "/" && strings.HasPrefix(reqUrl, "/echo/") {
			body := reqUrl[6:]
			headers := "Content-Type: text/plain" + CRLF + "Content-Length: " + fmt.Sprint(len(body)) + CRLF

			response = HTTP_OK + headers + CRLF + CRLF + body + CRLF + CRLF
		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error responding")
		}

		if err != nil {
			fmt.Println("Error while sending response ", err)
		}

		err = conn.Close()
		if err != nil {
			fmt.Println("Unable to close connection ", err)
		}
	}
}

func extractHttpProperties(reqBuffer []byte, reqSize int) []string {
	if reqSize == 0 {
		return make([]string, 0)
	}
	req := strings.Split(string(reqBuffer[:reqSize]), CRLF)
	reqProperties := req[0]
	return strings.Split(reqProperties, " ")
}
