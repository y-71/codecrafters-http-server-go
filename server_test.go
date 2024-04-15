package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestConcurrentClient(t *testing.T) {
	// Number of concurrent requests
	numRequests := 10

	// Address of the server
	serverAddr := "localhost:4221"

	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Perform concurrent requests
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			// Connect to the server
			conn, err := net.Dial("tcp", serverAddr)
			if err != nil {
				t.Errorf("Error connecting to server: %v", err)
				return
			}
			defer conn.Close()

			// Construct the HTTP request with headers
			request := "GET / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: GoClient\r\n" +
				"Connection: close\r\n\r\n"

			// Send the HTTP request
			_, err = conn.Write([]byte(request))
			if err != nil {
				t.Errorf("Error sending request: %v", err)
				return
			}

			// Receive and print the response
			response := make([]byte, 1024)
			n, err := conn.Read(response)
			if err != nil {
				t.Errorf("Error receiving response: %v", err)
				return
			}
			fmt.Printf("Response received: %s\n", response[:n])
		}()
	}

	// Wait for all requests to complete
	wg.Wait()
}
