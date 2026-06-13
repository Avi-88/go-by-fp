# Day 9 — HTTP Request Lifecycle (Raw)

## Goal
Understand what HTTP actually is at the byte level — before using `net/http`.
HTTP/1.1 is a text protocol built on top of TCP. You already built the TCP layer.
Today you parse it by hand to understand what `net/http` does for you.

---

## Background: What HTTP/1.1 looks like on the wire

When a browser (or `curl`) makes a request, it sends raw text over TCP:

```
GET /hello?name=ada HTTP/1.1\r\n
Host: localhost:8080\r\n
User-Agent: curl/7.64.1\r\n
Accept: */*\r\n
\r\n
```

The response looks like:
```
HTTP/1.1 200 OK\r\n
Content-Type: text/plain\r\n
Content-Length: 13\r\n
\r\n
Hello, world!
```

Key rules:
- Request line: `METHOD PATH HTTP/VERSION\r\n`
- Each header: `Key: Value\r\n`
- Blank line (`\r\n`) separates headers from body
- Body follows immediately after (for POST/PUT)

---

## Part 1: Parse a Raw HTTP Request

Write a function that parses a raw HTTP/1.1 request from a `net.Conn`:

```go
type HTTPRequest struct {
    Method  string
    Path    string
    Query   map[string]string
    Version string
    Headers map[string]string
    Body    []byte
}

func ParseRequest(conn net.Conn) (HTTPRequest, error)
```

Requirements:
- Read the request line: extract method, path (without query), query params
- Read headers until the blank line (`\r\n\r\n`)
- If `Content-Length` header exists, read exactly that many bytes as body
- Header keys should be normalized to lowercase

Test it with `curl`:
```bash
# Terminal 1 — start your server
go run . server

# Terminal 2
curl -v http://localhost:8080/hello?name=ada
curl -v -X POST -d '{"user":"ada"}' http://localhost:8080/users
```

---

## Part 2: Write a Raw HTTP Response

Write a function:

```go
func WriteResponse(conn net.Conn, status int, headers map[string]string, body []byte) error
```

It should write:
1. Status line: `HTTP/1.1 <status> <reason>\r\n`
2. Each header as `Key: Value\r\n`
3. Always include `Content-Length` (derived from body)
4. Blank line `\r\n`
5. Body bytes

Map these status codes to reason phrases:
```go
var statusText = map[int]string{
    200: "OK",
    201: "Created",
    400: "Bad Request",
    404: "Not Found",
    405: "Method Not Allowed",
    500: "Internal Server Error",
}
```

---

## Part 3: A Minimal HTTP Server

Wire up Parts 1 and 2 into a TCP server that handles HTTP requests:

```go
func startHTTPServer(addr string) error
```

Handle these routes manually (no router library):

```
GET  /           → 200, "Welcome"
GET  /hello      → 200, "Hello, {name}!" (from ?name= query param)
GET  /echo       → 200, echo back all request headers as JSON
POST /echo       → 200, echo back request body
GET  /status     → 200, JSON with server uptime and request count
unknown route    → 404, "Not Found"
wrong method     → 405, "Method Not Allowed"
```

Keep a request counter and server start time as package-level variables (protected with a mutex since requests are handled concurrently).

---

## Part 4: Test with curl

Run your server and test each route:

```bash
# Basic GET
curl http://localhost:8080/

# Query param
curl "http://localhost:8080/hello?name=Ada"

# Echo headers
curl -H "X-Custom: foo" http://localhost:8080/echo

# Echo body
curl -X POST -d "hello from curl" http://localhost:8080/echo

# Status
curl http://localhost:8080/status

# 404
curl http://localhost:8080/missing

# Wrong method
curl -X DELETE http://localhost:8080/hello
```

---

## Concepts to understand deeply

- HTTP is just text over TCP — `net/http` parses exactly what you're parsing by hand
- `\r\n` (CRLF) is the HTTP line terminator — not just `\n`
- Headers and body are separated by a blank line — `\r\n\r\n`
- `Content-Length` tells the receiver exactly how many body bytes to read
- HTTP/1.1 keeps connections alive by default (`Connection: keep-alive`) — your server can ignore this for now
- Every `net.Conn` from `Accept()` is one HTTP request (in HTTP/1.0) or potentially many (keep-alive)

---

## Intentional breaks (do after it works)

1. Remove the `\r\n` after your status line. Use `curl -v` to see the response — observe it breaks.
2. Send a `Content-Length` that's wrong (too small). See how curl handles truncated body.
3. Don't send `Content-Length` at all. See what curl does — it may work but complain.
4. Try parsing a request with `bufio.Scanner` using `\n` instead of `\r\n` as separator. Works on most clients but is technically wrong — understand why.

---

## When done, write NOTES.md

- What is the structure of an HTTP/1.1 request in plain text?
- What does `Content-Length` do, and what happens without it?
- What is the difference between the path and query string?
- Why does `curl -v` show both the request and response headers?
- After doing this by hand, what does `net/http` give you for free?
