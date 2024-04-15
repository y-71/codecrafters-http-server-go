package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type Header struct {
	HttpMethod  string
	Path        string
	HttpVersion string
	UserAgent 	string
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Read data from the connection
	buf := make([]byte, 1024) // Create a buffer to hold the incoming data
	_, err = conn.Read(buf)   // Read data from the connection into the buffer
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	headers := strings.Fields(strings.TrimSpace(string(buf)))

	if len(headers) < 6 {
		fmt.Println("Invalid number of headers")
	}
	httpMethod := headers[0]
	path := headers[1]
	httpVersion := headers[2]
	userAgent := headers[6]

	header := Header{
		HttpMethod:  httpMethod,
		Path:        path,
		HttpVersion: httpVersion,
		UserAgent: userAgent,
	}
	var response []byte
	if header.Path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n ")
	} else if header.isRoute("/echo") {

		content, _ := header.getEchoPathContent()
		
		response = []byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n\r\n"+
			"%s", len(content), content))
	} else if header.isRoute("/user-agent"){
		response = []byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n\r\n"+
			"%s", len(header.UserAgent), header.UserAgent))
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n ")
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

	err = conn.Close()
	if err != nil {
		fmt.Println("Error closing connection: ", err.Error())
		os.Exit(1)
	}

}

func (h Header) isRoute(route string) bool {
	return strings.HasPrefix(h.Path, route)
}
func (h Header) getEchoPathContent() (string, error) {
	if !h.isRoute("/echo") {
		return "", errors.New("invalid path format: expected path to start with '/echo/'")
	}
	if !strings.HasPrefix(h.Path, "/echo/"){
		return "", nil
	}
	return h.Path[len("/echo/"):], nil
}