package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run . <server|client|chat-server|chat-client> [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		startServer(":3000")
		return
	case "client":
		if len(os.Args) != 3 {
			fmt.Printf("Insufficient arguments passed - Usage client <message>")
			return
		}
		message := os.Args[2]
		err := runClient(":3000",message)
		if err != nil {
			fmt.Printf("Error with client-%v\n", err)
			os.Exit(1)
		}
		return
	case "chat-server":
		return
	case "chat-client":
		return
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runClient(addr string, message string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("There was a error connecting to the server -%w\n", err)
	}
	for {
		payload := []byte(message)

		err = WriteMessage(conn, payload)
		if err != nil {
			return fmt.Errorf("There was a error sending payload to the server -%w\n", err)
		}
		msg, err := ReadMessage(conn)
		if err == io.EOF {
			conn.Close()
			return nil
		}
		if err != nil {
			return fmt.Errorf("There was a error reading payload from the server -%w\n", err)
		}
		fmt.Println(string(msg))
	}
}
