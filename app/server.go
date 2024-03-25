package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	IP      = "localhost"
	PORT    = "4221"
	IP_PORT = IP + ":" + PORT
)

const CRLF = "\r\n"
const HTTP_OK = "HTTP/1.1 200 OK"
const HTTP_NOT_FOUND = "HTTP/1.1 404 Not Found"

func main() {
	filesDirectory := flag.String("directory", "", "Directory to serve files from")

	flag.Parse()

	fmt.Println("Serving files from: ", *filesDirectory)

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

		go handleConnection(conn, *filesDirectory)

	}
}

func loadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func handleConnection(conn net.Conn, filesDirectory string) {
	request := make([]byte, 4096)
	reqSize, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}

	reqUrl := extractPath(string(request[:reqSize]))

	/**
	Req Method = 0, (GET / POST/ ETC)
	Req Url = 1, (http://example.com)
	Proto Type = 2 (Http/1.1)
	*/
	// reqUrl := httpProperties[0][1]
	// userAgent := httpProperties[2][1]
	response := HTTP_OK + CRLF + CRLF

	if reqUrl != "/" {
		response = HTTP_NOT_FOUND + CRLF + CRLF

	}

	if reqUrl != "/" && strings.HasPrefix(reqUrl, "/echo/") {
		body := reqUrl[6:]
		headers := HTTP_OK + CRLF + "Content-Type: text/plain" + CRLF + "Content-Length: " + fmt.Sprint(len(body)) + CRLF + CRLF

		response = headers + body + CRLF + CRLF
	}

	if reqUrl != "/" && (reqUrl == "/user-agent") {
		userAgent := extractUserAgent(request)
		headers := HTTP_OK + CRLF + "Content-Type: text/plain" + CRLF + "Content-Length: " + fmt.Sprint(len(userAgent)) + CRLF + CRLF

		response = headers + userAgent + CRLF + CRLF
	}

	if reqUrl != "/" && strings.Contains(reqUrl, "/files") {
		fileName := strings.Split(reqUrl, "/files/")[1]

		fileBytes, err := loadFile(fmt.Sprintf("%s/%s", filesDirectory, fileName))

		if err != nil {
			response = HTTP_NOT_FOUND + CRLF + CRLF
		} else {
			headers := HTTP_OK + CRLF + "Content-Type: application/octet-stream" + CRLF + "Content-Length: " + fmt.Sprint(len(fileBytes)) + CRLF + CRLF
			response = headers + string(fileBytes) + CRLF + CRLF
		}
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error responding")
	}

	err = conn.Close()
	if err != nil {
		fmt.Println("Unable to close connection ", err)
	}
}

// func extractHttpProperties(reqBuffer []byte, reqSize int) [][]string {
// 	if reqSize == 0 {
// 		return make([][]string, 0)
// 	}
// 	req := strings.Split(string(reqBuffer[:reqSize]), CRLF)
// 	reqProperties := req[0]
// 	host := req[1]
// 	agent := req[2]
// 	// return
// 	// return a map
// 	return [][]string{strings.Split(reqProperties, " "), strings.Split(host, " "), strings.Split(agent, " ")}
// }

func extractPath(req string) string {
	var path string
	start := strings.Index(req, " ") + 1
	end := strings.Index(req[start:], " ") + start
	if start > 0 && end > start {
		path = req[start:end]
	}
	return path
}
func extractUserAgent(req []byte) string {
	lines := strings.Split(string(req), "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "User-Agent: ") {
			return strings.TrimPrefix(line, "User-Agent: ")
		}
	}
	return ""
}
