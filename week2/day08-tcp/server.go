package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"errors"
	"sync"
	"time"
)

func startServer(addr string) error {
	ls,err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("There was a error starting the server - %w\n", err)
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			return fmt.Errorf("There was a error accepting client connection - %w\n", err)
		}
		fmt.Printf("New connection -%v\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		msg, err := ReadMessage(conn)
		if errors.Is(err, io.EOF) {
			conn.Close()
			return
		}
		if err != nil {
			fmt.Printf("There was an error while parsing data -%v\n", err)
		}
		
		err = WriteMessage(conn, msg)
		if err != nil {
			fmt.Printf("There was an error while sending data -%v\n", err)
		}
	}
}

// WriteMessage writes a 4-byte big-endian length header followed by the message bytes.
func WriteMessage(w io.Writer, msg []byte) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(msg)))

	if _, err := w.Write(header); err != nil {
		return fmt.Errorf("There was a error adding the length header -%w", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("There was a error sending payload -%w", err)
	}
	return nil
}

// ReadMessage reads a 4-byte big-endian length header, then reads exactly that many bytes.
func ReadMessage(r io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("There was a error reading the length header -%w", err)
	}
	ln := binary.BigEndian.Uint32(header)
	msg := make([]byte, ln)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, fmt.Errorf("There was a error parsing the payload -%w", err)
	}

	return msg, nil
}

type Client struct {
    name string
    conn net.Conn
    send chan string
}

type ActiveConn struct {
	connections map[string]Client
	mu sync.RWMutex
}

var ac ActiveConn


func startChatServer(addr string) error {
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("There was a error starting the server - %w\n", err)
	}
	ac = ActiveConn{connections: make(map[string]Client)}
	for {
		conn, err := ls.Accept()
		if err != nil {
			return fmt.Errorf("There was a error establishing connection to the server - %w\n", err)
		}
		go handleChatConnection(conn)
	}
}

func handleChatConnection(conn net.Conn) {
	name, err := ReadMessage(conn)
	if err != nil {
		fmt.Printf("There was a error receiving payload")
		return
	}
	client := Client{name: string(name), conn: conn, send: make(chan string, 32)}
	sendBroadCast(string(name) + " joined the chat", nil)
	ac.mu.Lock()
	ac.connections[string(name)] = client
	ac.mu.Unlock()

	go receiveBroadCast(client)
	for {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		msg, err := ReadMessage(conn)
		if errors.Is(err, io.EOF) {
			ac.mu.Lock()
			delete(ac.connections, string(name))
			ac.mu.Unlock()
			sendBroadCast(string(name) + " left the chat", nil)
			close(client.send)
			conn.Close()
			return
		}
		if err != nil {
			fmt.Printf("Error in sending broadcast -%v\n", err)
			ac.mu.Lock()
			delete(ac.connections, string(name))
			ac.mu.Unlock()
			sendBroadCast(string(name) + " left the chat", nil)
			close(client.send)
			conn.Close()
			return 
		}
		sendBroadCast(string(msg), &client)
	}

}

func receiveBroadCast(client Client) {
	for msg:= range client.send {
		WriteMessage(client.conn, []byte(msg))
	}
}

func sendBroadCast(message string, exclude *Client) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	for _,cl := range ac.connections {
		if exclude != nil && cl.name == exclude.name {
			continue
		} else if exclude != nil {
			cl.send <- string("[" + exclude.name + "]" + " : "+ message)
		}else {
			cl.send <- message
		}
		
	}
	
}