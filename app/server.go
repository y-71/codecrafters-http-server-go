package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Header struct {
	HttpMethod  string
	Path        string
	HttpVersion string
	UserAgent   string
}

func handleConnection(conn net.Conn, directory string) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	headers := strings.Fields(strings.TrimSpace(string(buf)))

	if len(headers) < 3 {
		// Return a 400 Bad Request response for incomplete or malformed requests
		response := []byte("HTTP/1.1 400 Bad Request\r\n\r\n")
		_, err := conn.Write(response)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
		}
		return
	}

	httpMethod := headers[0]
	path := headers[1]
	httpVersion := headers[2]
	userAgent := ""

	if len(headers) > 6 {
		userAgent = headers[6]
	}

	header := Header{
		HttpMethod:  httpMethod,
		Path:        path,
		HttpVersion: httpVersion,
		UserAgent:   userAgent,
	}

	var response []byte

	switch {
	case header.Path == "/":
		response = []byte("HTTP/1.1 200 OK\r\n\r\n ")
	case header.isRoute("/echo"):
		content, _ := header.getEchoPathContent()
		response = []byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/plain\r\n"+
				"Content-Length: %d\r\n\r\n"+
				"%s", len(content), content))
	case header.isRoute("/user-agent"):
		response = []byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/plain\r\n"+
				"Content-Length: %d\r\n\r\n"+
				"%s", len(header.UserAgent), header.UserAgent))
	case strings.HasPrefix(header.Path, "/files/"):
		filePath := filepath.Join(directory, header.Path[len("/files/"):])
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		} else {
			response = []byte(fmt.Sprintf(
				"HTTP/1.1 200 OK\r\n"+
					"Content-Type: application/octet-stream\r\n"+
					"Content-Length: %d\r\n\r\n"+
					"%s", len(content), content))
		}
	default:
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

func main() {
	fmt.Println("Logs from your program will appear here!")

	if len(os.Args) < 3 || os.Args[1] != "--directory" {
		fmt.Println("Usage: ./your_server.sh --directory <directory>")
		os.Exit(1)
	}

	directory := os.Args[2]

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, directory)
	}
}

func (h Header) isRoute(route string) bool {
	return strings.HasPrefix(h.Path, route)
}

func (h Header) getEchoPathContent() (string, error) {
	if !h.isRoute("/echo") {
		return "", errors.New("invalid path format: expected path to start with '/echo/'")
	}
	if !strings.HasPrefix(h.Path, "/echo/") {
		return "", nil
	}
	return h.Path[len("/echo/"):], nil
}
