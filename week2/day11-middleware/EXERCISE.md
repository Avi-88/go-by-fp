# Day 11 — Middleware Chain

## Goal
Go deeper on middleware — build production-quality middleware patterns
that handle authentication, rate limiting, CORS, and request validation.
Day 10 introduced the shape; today you build real middleware you'd actually use.

---

## Background: Middleware Stack

The middleware chain is evaluated inside-out:

```
Request  → A → B → C → Handler
Response ← A ← B ← C ← Handler
```

If you write:
```go
handler := A(B(C(mux)))
```

A runs first on the way in, last on the way out.
Order matters — put Logger outermost, Auth innermost near the handler.

---

## Part 1: Auth Middleware (API Key)

Write middleware that checks for an `Authorization` header with a valid API key:

```go
func APIKeyAuth(validKeys []string) func(http.Handler) http.Handler
```

- If the header is missing → 401 Unauthorized
- If the key is invalid → 403 Forbidden
- If valid → call `next.ServeHTTP`
- Store the key in the request context so handlers can read which key was used

Test:
```bash
curl -H "Authorization: Bearer secret123" http://localhost:8080/protected
curl http://localhost:8080/protected                          # should 401
curl -H "Authorization: Bearer wrong" http://localhost:8080/protected  # should 403
```

---

## Part 2: Rate Limiter Middleware

Write a simple per-IP rate limiter:

```go
func RateLimit(requestsPerSecond int) func(http.Handler) http.Handler
```

- Track request counts per IP address with a `map[string]*rateLimitEntry`
- Each entry has: count, window start time
- If a client exceeds `requestsPerSecond` in a 1-second window → 429 Too Many Requests
- Reset the count when a new second begins
- Protect the map with a `sync.Mutex`
- Include a `Retry-After: 1` header on 429 responses

Test by hitting the server rapidly:
```bash
for i in $(seq 1 20); do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/; done
```

---

## Part 3: CORS Middleware

Write middleware that adds CORS headers to every response:

```go
func CORS(allowedOrigins []string) func(http.Handler) http.Handler
```

- For all responses: add `Access-Control-Allow-Methods` and `Access-Control-Allow-Headers`
- If `Origin` header matches an allowed origin: add `Access-Control-Allow-Origin: <origin>`
- Handle `OPTIONS` preflight requests: respond 204 immediately (don't call next)
- If origin not in allowed list: don't add `Allow-Origin` header (implicit denial)

Test:
```bash
curl -H "Origin: http://localhost:3000" -v http://localhost:8080/
curl -X OPTIONS -H "Origin: http://localhost:3000" -H "Access-Control-Request-Method: POST" http://localhost:8080/
```

---

## Part 4: Request Validation Middleware

Write a middleware that validates request bodies for POST/PUT endpoints:

```go
func ValidateJSON(next http.Handler) http.Handler
```

- Only applies to `POST` and `PUT` requests
- Checks `Content-Type: application/json` header — if missing, 415 Unsupported Media Type
- Reads the body, checks it's valid JSON — if invalid, 400 Bad Request
- Re-attaches the body to the request (since reading drains it) using `io.NopCloser`
- Passes the validated request to `next`

Hint: after reading `r.Body`, replace it with:
```go
r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
```

---

## Part 5: Wire It All Together

Build a small API with the full middleware stack:

```go
Routes:
  GET  /                    → public, "Welcome"
  GET  /protected           → requires API key
  POST /protected/data      → requires API key + valid JSON body
  GET  /status              → public, request count + uptime

Middleware stack (outermost first):
  Logger → RateLimit(5) → CORS(["http://localhost:3000"]) → mux

Protected routes only:
  APIKeyAuth(["secret123", "dev-key"]) → ValidateJSON (POST only)
```

Note: apply `APIKeyAuth` and `ValidateJSON` only to specific routes, not the whole mux.
You can do this by wrapping individual handlers:

```go
mux.Handle("/protected", APIKeyAuth(keys)(http.HandlerFunc(handleProtected)))
mux.Handle("/protected/data", APIKeyAuth(keys)(ValidateJSON(http.HandlerFunc(handleProtectedData))))
```

---

## Concepts to understand deeply

- Middleware order determines what wraps what — think carefully about where auth vs logging go
- `io.NopCloser` lets you wrap a `bytes.Reader` to satisfy `io.ReadCloser` interface
- Rate limiting state must be goroutine-safe — each request runs in its own goroutine
- CORS preflight (`OPTIONS`) must be handled before the actual handler runs
- Context is the idiomatic way to pass per-request data through middleware to handlers

---

## Intentional breaks (do after it works)

1. In `RateLimit`, remove the mutex. Run the concurrent curl loop. Use `go run -race` — does the race detector catch it?
2. In `ValidateJSON`, don't re-attach the body. Call the handler — watch it read an empty body.
3. Put `APIKeyAuth` outside `Logger` instead of inside. What happens to the logs for rejected requests?
4. In CORS, don't handle `OPTIONS` separately. Send a POST from a browser (or simulate with curl). What does the browser do?

---

## When done, write NOTES.md

- Why does middleware order matter? Give a concrete example.
- What is `io.NopCloser` and why is it needed?
- What is a CORS preflight request and why does it exist?
- What is the difference between 401 and 403?
- How would you apply middleware to only some routes, not all?
