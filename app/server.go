package main

import (
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
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	request := string(buf[:n])

	lines := strings.Split(request, "\n")
	firstLine := strings.Fields(lines[0])

	if len(firstLine) < 3 {
		// Return a 400 Bad Request response for incomplete or malformed requests
		response := []byte("HTTP/1.1 400 Bad Request\r\n\r\n")
		_, err := conn.Write(response)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
		}
		return
	}

	httpMethod := firstLine[0]
	path := firstLine[1]
	httpVersion := firstLine[2]
	userAgent := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "User-Agent:") {
			userAgent = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent:"))
			break
		}
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
	case strings.HasPrefix(header.Path, "/files/") && header.HttpMethod == "POST":
		filename := header.Path[len("/files/"):]
		filePath := filepath.Join(directory, filename)
		body := strings.SplitN(request, "\n\n", 2)[1]

		err := ioutil.WriteFile(filePath, []byte(body), 0644)
		if err != nil {
			response = []byte("HTTP/1.1 500 Internal Server Error\r\n\r\n")
		} else {
			response = []byte("HTTP/1.1 201 Created\r\n\r\n")
		}
	default:
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
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

	var directory string
	if len(os.Args) > 2 && os.Args[1] == "--directory" {
		directory = os.Args[2]
	}

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
