package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// HTTP parsing and response writing.
//
// ParseRequest steps:
// 1. Wrap the conn in a bufio.Reader so you can read line by line
// 2. Read the first line — split it into method, full path, version
//    (hint: strings.Fields splits on whitespace)
// 3. Split the full path into path and query string on "?"
//    then split query string into key=value pairs on "&" and "="
// 4. Read headers line by line until you hit a blank line "\r\n"
//    each header line is "Key: Value\r\n" — split on ": "
//    normalize keys to lowercase (strings.ToLower)
// 5. Check for "content-length" header — if present, read exactly
//    that many bytes as the body (strconv.Atoi + io.ReadFull)
//
// WriteResponse steps:
// 1. Write status line: "HTTP/1.1 <code> <reason>\r\n"
//    use the statusText map to look up the reason phrase
// 2. Write each header as "Key: Value\r\n"
// 3. Always write "Content-Length: <len(body)>\r\n"
// 4. Write blank line "\r\n"
// 5. Write body bytes

type HTTPRequest struct {
    Method  string
    Path    string
    Query   map[string]string
    Version string
    Headers map[string]string
    Body    []byte
}

func ParseRequest(conn net.Conn) (HTTPRequest, error) {
	r := bufio.NewReader(conn)
	req := HTTPRequest{
		Method: "",
		Path: "",
		Query: make(map[string]string),
		Version: "",
		Headers: make(map[string]string),
		Body: nil,
	}
	// request line
	reqLine, err := r.ReadString('\n')
	reqLine = strings.TrimRight(reqLine, "\r\n")
	if err != nil {
		return req, fmt.Errorf("There was an error parsing the request string -%w\n", err)
	}
	fields := strings.Fields(reqLine)
	if len(fields) < 3  {
		return req, errors.New("Invalid request")
	}
	req.Method = fields[0]
	req.Version = fields[2]

	url := fields[1]
	queryString := strings.Split(url, "?")
	req.Path = queryString[0]
	if len(queryString) > 1 {
		for _, param := range strings.Split(queryString[1], "&") {
			parts := strings.SplitN(param, "=", 2)
			req.Query[parts[0]] = parts[1]
		}
	}

	// headers
	for {
		reqLine, err = r.ReadString('\n')
		reqLine = strings.TrimRight(reqLine, "\r\n")
		if reqLine == "" {
			break
		}
		if err != nil {
			return req, fmt.Errorf("There was an error parsing the request string -%w\n", err)
		}
		headers := strings.SplitN(reqLine, ":", 2)
		req.Headers[strings.ToLower(headers[0])] = strings.TrimSpace(headers[1])
	}

	// body
	val, ok := req.Headers["content-length"]
	if ok && (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") {
		ln, err := strconv.Atoi(val)
		if err != nil {
			return req, fmt.Errorf("Failed to parse Content-Length header -%w\n", err)
		}
		body := make([]byte, ln)
		_, err = io.ReadFull(r, body); if err != nil {
			return req, fmt.Errorf("Failed to parse request body -%w\n", err)
		}
		req.Body = body
	}

	return req, nil
}	

