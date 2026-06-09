package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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
		if err == io.EOF {
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