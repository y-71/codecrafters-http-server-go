package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Header struct{
	HttpMethod string
	Path string
	HttpVersion string
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
	
	if len(headers) < 3 {
		fmt.Println("Invalid number of headers")
	}
	httpMethod := headers[0]
	path := headers[1]
	httpVersion := headers[2]

	header := Header{
				HttpMethod: httpMethod,
				Path: path,
				HttpVersion: httpVersion,
			}
	if (header.Path == "/"){
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n "))
	}else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n "))
	}
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
