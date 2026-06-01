# Go by First Principles — 30 Days of Fundamentals


## Structure

```
go/
  week1/  — Language internals + execution model
  week2/  — I/O, networking, APIs
  week3/  — Concurrency
  week4/  — Systems thinking

python/
  (mirror exercises for selected days)
```

## Rules

1. No copy-paste. Type everything.
2. Break it intentionally. Then fix it.
3. Write NOTES.md after each day.
4. Use `go run -race` on every concurrent program.

## Week 1 — Language Internals

| Day | Topic |
|-----|-------|
| 1 | Variables, types, functions |
| 2 | Pointers, stack vs heap |
| 3 | Structs, methods, interfaces |
| 4 | Slices vs arrays (internals) |
| 5 | Maps, struct tags, JSON |
| 6 | Errors, custom types, wrapping |
| 7 | Debugger session (dlv) |

## Week 2 — I/O, Networking, APIs

| Day | Topic |
|-----|-------|
| 8  | TCP from scratch |
| 9  | HTTP request lifecycle (raw) |
| 10 | net/http server |
| 11 | Middleware chain |
| 12 | Serialization + streaming |
| 13 | Config loading |
| 14 | Build: REST API |

## Week 3 — Concurrency

| Day | Topic |
|-----|-------|
| 15 | Goroutines + WaitGroup |
| 16 | Channels (unbuffered) |
| 17 | Buffered channels + select |
| 18 | Mutexes + race conditions |
| 19 | Worker pool (fan-out/fan-in) |
| 20 | Context + cancellation |
| 21 | Rate limiter + job queue |

## Week 4 — Systems Thinking

| Day | Topic |
|-----|-------|
| 22 | Concurrent log processor |
| 23 | WebSocket broadcaster |
| 24 | Reverse proxy |
| 25 | Task scheduler |
| 26 | Mini Redis clone |
| 27 | Distributed worker prototype |
| 28-30 | Capstone |
