package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Day 9: HTTP Request Lifecycle (Raw)
//
// Read EXERCISE.md before starting.
// No net/http allowed today — raw TCP only.
//
// Steps:
// 1. Accept CLI args (just "server" for now)
// 2. Start a TCP listener (same pattern as day 8)
// 3. For each connection, spawn a goroutine that:
//    - calls ParseRequest
//    - routes to the right handler based on method + path
//    - calls WriteResponse

type Stats struct {
	startTime time.Time
	reqCount int
	mu       sync.RWMutex
}

type Status struct {
	Uptime string `json:"uptime"`
	RequestCount int     `json:"request_count"`
}

var serverStats Stats = Stats{
	startTime: time.Time{},
	reqCount: 0,
	mu: sync.RWMutex{},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage <server>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "server":
		err := startServer(":8080")
		if err != nil {
			fmt.Println("There was an error starting the server")
		}
	default:
		os.Exit(0)
	}
}

func startServer(addr string) error {
	serverStats.startTime = time.Now()
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	req,err := ParseRequest(conn)
	if err != nil {
		fmt.Printf("Error parsing-%v\n", err)
		err := WriteResponse(conn, 400, nil, []byte("Bad request"))
		if err != nil {
			fmt.Printf("Invalid  -%v", err)
			return
		}
		return
	}
	serverStats.mu.Lock()
	serverStats.reqCount++
	serverStats.mu.Unlock()
	switch req.Path {
	case "/":
		handlePathRoot(conn, req)
	case "/hello":
		handlePathHello(conn, req)
	case "/echo":
		handlePathEcho(conn, req)
	case "/status":
		handlePathStatus(conn, req)
	default:
		err := WriteResponse(conn, 404, nil, []byte("Request Path not found"))
		if err != nil {
			fmt.Printf("Invalid  -%v", err)
			return
		}
	}

	fmt.Printf("The request structure is - %v\n",req)
}

func handlePathRoot(conn net.Conn, req HTTPRequest) {
	switch req.Method {
	case "GET":
		err := WriteResponse(conn, 200, nil, []byte("Welcome"))
		if err != nil {
			fmt.Printf("There was an error sending response inside root path -%v", err)
			return
		}
	default:
		err := WriteResponse(conn, 405, nil, []byte("Method not allowed"))
		if err != nil {
			fmt.Printf("Invalid  -%v", err)
			return
		}
	}
}

func handlePathHello(conn net.Conn, req HTTPRequest) {
	switch req.Method {
	case "GET":
		name, ok := req.Query["name"]
		if !ok {
			err := WriteResponse(conn, 400, nil, []byte("Request missing query params"))
			if err != nil {
				fmt.Printf("There was an error sending response inside hello path -%v", err)
				return
			}
			return
		}
		err := WriteResponse(conn, 200, nil, []byte(fmt.Sprintf("Hello, %s!", name)))
		if err != nil {
			fmt.Printf("There was an error sending response inside hello path -%v", err)
			return
		}
	default:
		err := WriteResponse(conn, 405, nil, []byte("Method not allowed"))
		if err != nil {
			fmt.Printf("Invalid  -%v", err)
			return
		}
	}
}

func handlePathEcho(conn net.Conn, req HTTPRequest) {
	switch req.Method {
	case "GET":
		jsonHeaders, err := json.MarshalIndent(req.Headers, "", "  ")
		if err != nil {
			err := WriteResponse(conn, 500, nil, []byte("Something went wrong"))
			if err != nil {
				fmt.Printf("There was an error sending response inside hello path -%v", err)
				return
			}
			return
		}
		
		err = WriteResponse(conn, 200, map[string]string{"Content-Type": "application/json"}, jsonHeaders)
		if err != nil {
			fmt.Printf("There was an error sending response inside hello path -%v", err)
			return
		}
	case "POST":
		err := WriteResponse(conn, 200, map[string]string{"Content-Type": "application/json"}, req.Body)
		if err != nil {
			fmt.Printf("There was an error sending response inside hello path -%v", err)
			return
		}
	default:
		err := WriteResponse(conn, 405, nil, []byte("Method not allowed"))
		if err != nil {
			fmt.Printf("Invalid  -%v", err)
			return
		}
	}
}

func handlePathStatus(conn net.Conn, req HTTPRequest) {
	switch req.Method {
	case "GET":
		serverStats.mu.RLock()
		status := Status{
			Uptime: time.Since(serverStats.startTime).String(),
			RequestCount: serverStats.reqCount,
		}
		serverStats.mu.RUnlock()
		jsonStatus, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			err := WriteResponse(conn, 500, nil, []byte("Something went wrong"))
			if err != nil {
				fmt.Printf("There was an error sending response inside status path -%v", err)
				return
			}
		}
		err = WriteResponse(conn, 200, map[string]string{"Content-Type": "application/json"}, jsonStatus )
		if err != nil {
			fmt.Printf("There was an error sending response inside status path -%v", err)
			return
		}
	default:
		err := WriteResponse(conn, 405, nil, []byte("Method not allowed"))
		if err != nil {
			fmt.Printf("Invalid  -%v", err)
			return
		}
	}
}


