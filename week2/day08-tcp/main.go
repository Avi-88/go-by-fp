package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"errors"
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
		err := runClient(":3000",os.Stdin)
		if err != nil {
			fmt.Printf("Error with client-%v\n", err)
			os.Exit(1)
		}
		return
	case "chat-server":
		startChatServer(":4000")
		return
	case "chat-client":
		if len(os.Args) < 3 {
			fmt.Println("usage: go run . <chat-client> [message...]")
			os.Exit(1)
		}
		err := runChatClient(":4000", string(os.Args[2]), os.Stdin)
		if err != nil {
			fmt.Printf("Error with client-%v\n", err)
			os.Exit(1)
		}
		return
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runClient(addr string, input io.Reader) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("There was a error connecting to the server -%w\n", err)
	}

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		message := scanner.Text()
	
		err = WriteMessage(conn, []byte(message))
		if err != nil {
			return fmt.Errorf("There was a error sending payload to the server -%w\n", err)
		}

	}

	return nil
}


func runChatClient(addr string,name string, input io.Reader) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("There was a error connecting to the server -%w\n", err)
	}
	scanner := bufio.NewScanner(input)
	err = WriteMessage(conn, []byte(name))
	if err != nil {
		return fmt.Errorf("There was a error sending client name to the server -%w\n", err)
	}

	go printMessage(conn)

	for scanner.Scan() {
		message := scanner.Text()

		err = WriteMessage(conn, []byte(message))
		if err != nil {
			return fmt.Errorf("There was a error sending payload to the server -%w\n", err)
		}
	}
	return nil
}

func printMessage(conn net.Conn) {
	for {
		msg, err := ReadMessage(conn)
		if errors.Is(err, io.EOF) {
			conn.Close()
			return 
		}
		if err != nil {
			fmt.Printf("There was a error reading payload from the server -%v\n", err)
			return
		}
		fmt.Println(string(msg))
	}
}
