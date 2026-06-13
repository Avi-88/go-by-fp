package main

import (
	"fmt"
	"net"
	"os"
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
		return
	}
	fmt.Printf("The request structure is - %v\n",req)
}
