# Day 10 — net/http Server

## Goal
Use Go's standard `net/http` package — now that you know what it does under the hood.
Today you'll see how much Day 9's manual work is replaced by the stdlib,
and learn the Handler interface, ServeMux routing, and middleware patterns.

---

## Background: The Handler Interface

Everything in `net/http` revolves around one interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Any type that implements `ServeHTTP` can be registered as a handler.
`http.HandlerFunc` is a convenience type that lets you use a plain function as a Handler:

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello")
}
// http.HandlerFunc(myHandler) satisfies Handler
```

---

## Part 1: Basic Server with ServeMux

Build the same server as Day 9, but using `net/http`:

```
GET  /           → 200, "Welcome"
GET  /hello      → 200, "Hello, {name}!" (from ?name= query param)
GET  /echo       → 200, echo back all request headers as JSON
POST /echo       → 200, echo back request body
GET  /status     → 200, JSON with server uptime and request count
```

Requirements:
- Use `http.NewServeMux()` — not the global `http.DefaultServeMux`
- Register handlers with `mux.HandleFunc("/path", handlerFunc)`
- Use `r.Method` to check the HTTP method inside each handler
- Use `r.URL.Query().Get("name")` for query params
- Use `r.Header` for headers
- Use `json.NewEncoder(w).Encode(data)` to write JSON
- Use `http.Error(w, message, statusCode)` for error responses
- Track request count with a mutex-protected counter (same as Day 9)

---

## Part 2: Middleware

Middleware wraps a `Handler` and adds behaviour before/after the inner handler runs.

Write these three middleware functions:

```go
// Logger logs: METHOD PATH → STATUS (duration)
func Logger(next http.Handler) http.Handler

// Recovery catches panics and returns 500 instead of crashing
func Recovery(next http.Handler) http.Handler

// RequestID adds a unique request ID to every request's context
// and includes it in the response as X-Request-ID header
func RequestID(next http.Handler) http.Handler
```

Each middleware has the signature `func(http.Handler) http.Handler` — it takes a handler and returns a new handler that wraps it.

Chain them together:

```go
handler := Logger(Recovery(RequestID(mux)))
http.ListenAndServe(":8080", handler)
```

### Logging middleware specifics:
- Log before: `METHOD PATH`
- Log after: `METHOD PATH → STATUS (elapsed)`
- To capture the status code written by the inner handler, wrap `ResponseWriter` in your own struct that records the status:

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}
```

### Recovery middleware specifics:
- Use `defer` + `recover()` to catch panics
- Write a 500 response: `http.Error(w, "Internal Server Error", 500)`
- Add a route `/panic` that intentionally panics — verify Recovery catches it

### RequestID middleware specifics:
- Generate a simple ID: `fmt.Sprintf("%d", time.Now().UnixNano())`
- Store it in the request context: `r.WithContext(context.WithValue(r.Context(), "request-id", id))`
- Set it in the response header: `w.Header().Set("X-Request-ID", id)`

---

## Part 3: File Server

Add a static file server route:

```go
mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
```

Create a `static/` folder with a few files and verify they're served correctly.

---

## Part 4: Compare with Day 9

Write a short comparison (in NOTES.md) of what Day 9's manual implementation required
vs what `net/http` gives you for free.

---

## Concepts to understand deeply

- `http.Handler` is an interface — anything with `ServeHTTP(w, r)` works
- `http.HandlerFunc` is a type that lets a plain function satisfy `Handler`
- `ServeMux` matches routes by longest prefix — `/api/` matches `/api/foo` and `/api/bar`
- `ResponseWriter` is also an interface — you can wrap it (as in the logger middleware)
- Middleware is just function composition — each layer adds behaviour
- `context.Context` carries request-scoped values through the call chain
- `http.ListenAndServe` blocks forever — it's the same TCP loop you wrote in Day 8

---

## Intentional breaks (do after it works)

1. Register two handlers on the same path. What does `ServeMux` do?
2. Write to `w` before calling `w.WriteHeader(200)`. What status code does the client see? What does it print on the server?
3. In Recovery, call `recover()` outside a `defer`. Does it catch the panic?
4. Remove Recovery middleware, hit `/panic`. Watch the server goroutine crash — does the whole server die or just that request?

---

## When done, write NOTES.md

- What is the `http.Handler` interface and why is it useful?
- What is `http.HandlerFunc` and how is it different from a plain function?
- What does middleware do, and how does the chaining pattern work?
- What does `ServeMux` give you that your Day 9 switch statement didn't?
- What does `http.ListenAndServe` do under the hood?
