# Day 8 — TCP from Scratch

## Goal
Understand TCP at the Go level — not just "use a library", but actually open sockets,
accept connections, read and write raw bytes, and handle the full connection lifecycle.
Everything in `net/http` and every networked service you'll ever build sits on top of this.

---

## Background: What TCP actually is

TCP gives you a **reliable, ordered, bidirectional byte stream** between two endpoints.
There is no concept of "messages" — just a stream of bytes.
It's your responsibility to frame them.

The server side:
1. Binds to a port (announces "I'm listening here")
2. Accepts incoming connections (blocks until a client connects)
3. Gets a `net.Conn` — reads and writes bytes on it
4. Closes the connection when done

The client side:
1. Dials a server (connects to its IP:port)
2. Gets a `net.Conn`
3. Writes bytes, reads the response
4. Closes the connection

---

## Part 1: Echo Server

Write a TCP server that listens on `localhost:9000` and echoes back whatever it receives.

### Server (`server.go` or inline in `main.go`)

```go
func startServer(addr string) error
```

Requirements:
- Call `net.Listen("tcp", addr)` to get a `net.Listener`
- Loop: call `listener.Accept()` to get each `net.Conn`
- Handle each connection in its own goroutine (so the server stays responsive)
- In the handler: read from the conn, write the same bytes back, repeat until EOF or error
- Use `io.Copy(conn, conn)` — this is the simplest echo implementation. Understand why it works.
- Log each new connection with the remote address: `conn.RemoteAddr()`

### Client (in `main.go`)

```go
func runClient(addr string, message string) error
```

Requirements:
- `net.Dial("tcp", addr)` to connect
- Write `message + "\n"` to the conn
- Read the response and print it
- Close the conn

Run both from `main`:
```
go run . server       → starts the echo server
go run . client "hello world"  → connects, sends, prints response
```

---

## Part 2: Framed Messages

`io.Copy` is elegant but it hides something important: **TCP has no message boundaries**.

If you send "hello" and "world" as two separate writes, the receiver might read them as
"helloworld" in one read, or "hel", "lo", "wo", "rld" in four reads.

You need a **framing protocol**. The simplest: length-prefixed messages.

Write two helper functions:

```go
// WriteMessage writes a 4-byte big-endian length header followed by the message bytes.
func WriteMessage(w io.Writer, msg []byte) error

// ReadMessage reads a 4-byte big-endian length header, then reads exactly that many bytes.
func ReadMessage(r io.Reader) ([]byte, error)
```

Use `encoding/binary` with `binary.BigEndian`.

Then update your echo server and client to use `WriteMessage`/`ReadMessage` instead of raw reads.

Verify: send multiple messages in a loop from the client, receive them all correctly.

---

## Part 3: Multi-client Chat Server

Build a TCP chat server where multiple clients can connect and broadcast messages to all others.

```
go run . chat-server       → starts on localhost:9001
go run . chat-client alice  → connects as "alice"
go run . chat-client bob    → connects as "bob"
```

When "alice" types a message and sends it, all other connected clients see:
```
[alice]: hello everyone
```

### Requirements

- Keep a global list of active connections (protected with a `sync.Mutex`)
- When a new client connects, register it; when it disconnects, deregister it
- Each client runs two goroutines: one reading from the TCP conn, one writing to it
- Use a `chan string` per client for outbound messages
- Broadcast function: iterate over all registered connections and send to each channel

```go
type Client struct {
    name string
    conn net.Conn
    send chan string
}
```

This is the hardest part of the day. Think carefully about:
- What happens when you write to a closed channel?
- What happens when a client disconnects mid-broadcast?
- How do you avoid deadlocks when holding the mutex while sending?

---

## Part 4: Timeouts and Deadlines

Connection bugs often manifest as hangs, not crashes.
Go lets you set deadlines on any `net.Conn`.

Add to your echo server handler:

```go
conn.SetDeadline(time.Now().Add(30 * time.Second))
```

Then experiment:
1. Connect with `nc localhost 9000`, type nothing for 30 seconds — observe what happens
2. `conn.SetReadDeadline` vs `conn.SetDeadline` — what's the difference?
3. After a timeout error, can you reset the deadline and keep using the conn?

---

## Concepts to understand deeply

- `net.Listen` vs `net.Dial` — server vs client side
- `listener.Accept()` blocks until a connection arrives — it returns a `net.Conn`
- `net.Conn` implements both `io.Reader` and `io.Writer` — that's why `io.Copy` works
- `Read` on a TCP conn blocks until *some* bytes arrive — not necessarily all you asked for
- `Read` returning `(0, io.EOF)` means the other side closed the connection
- Goroutine per connection is the standard Go pattern (not the only one, but the natural one)
- `sync.Mutex` protects shared state accessed from multiple goroutines
- Deadlines are absolute times (`time.Time`), not durations — set them fresh each loop iteration

---

## Intentional breaks (do after it works)

1. Remove the `go` before your connection handler. What happens when a second client connects while the first is still connected?
2. In the chat server, `close(client.send)` immediately when a client disconnects. What panic do you get? How do you fix it?
3. Write to a conn after calling `conn.Close()`. What error do you get?
4. In `ReadMessage`, don't use `io.ReadFull` — use a plain `Read`. Send a large message. Does it always work?

---

## When done, write NOTES.md

- What is the difference between `net.Listen` and `net.Dial`?
- Why does `io.Copy(conn, conn)` implement an echo server? Walk through what it does.
- What problem does framing solve, and why doesn't TCP solve it for you?
- What is a connection deadline, and when would you use one in production?
- What does the chat server teach you about shared mutable state across goroutines?
