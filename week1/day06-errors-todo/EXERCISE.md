# Day 6 — Errors, Custom Types, and Wrapping

## Goal
Understand Go's error model deeply — not just how to return errors,
but how to create meaningful error types, wrap context into them,
and let callers inspect and handle them precisely.
This directly applies to every function you'll write in Week 2+.

---

## Part 1: The Basic Error Patterns

Write a function that demonstrates all three common error patterns:

```go
func divide(a, b float64) (float64, error)
func sqrt(n float64) (float64, error)
func parseAge(s string) (int, error)
```

Rules:
- `divide`: error if `b == 0`
- `sqrt`: error if `n < 0`
- `parseAge`: error if not a valid number, or if age < 0 or age > 150

Use `errors.New` for simple errors and `fmt.Errorf` for errors with context.

---

## Part 2: Custom Error Types

Define two custom error types:

```go
type ValidationError struct {
    Field   string
    Message string
}

type NotFoundError struct {
    Resource string
    ID       int
}
```

Both must implement the `error` interface (i.e. have an `Error() string` method).

Then write:

```go
func validateUser(name string, age int) error
// returns ValidationError if name is empty or age is invalid

func findUser(id int) (string, error)
// returns NotFoundError if id is not in a hardcoded map of users
```

In `main`, call both and use `errors.As` to extract the concrete error type and print its fields:

```go
err := validateUser("", 200)
var ve ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Validation failed on field '%s': %s\n", ve.Field, ve.Message)
}
```

---

## Part 3: Error Wrapping and Unwrapping

Write a three-layer call chain that wraps errors with context:

```go
func readConfig(path string) (string, error)
// tries to open a file — if it fails, wrap: fmt.Errorf("readConfig: %w", err)

func loadApp(configPath string) error
// calls readConfig — if it fails, wrap: fmt.Errorf("loadApp: %w", err)

func startServer(configPath string) error
// calls loadApp — if it fails, wrap: fmt.Errorf("startServer: %w", err)
```

Call `startServer("missing.json")` from `main` and:
1. Print the full error message — observe how the context chain reads
2. Use `errors.Is` to check if the root cause is `os.ErrNotExist`
3. Use `errors.Unwrap` manually to peel back the layers one by one

---

## Part 4: The Todo CLI (Mini Project)

Build a CLI todo app that persists to a JSON file.

Commands:
```
go run . add "buy groceries"
go run . list
go run . done 1
go run . delete 2
```

### Data model
```go
type Todo struct {
    ID        int    `json:"id"`
    Text      string `json:"text"`
    Done      bool   `json:"done"`
    CreatedAt string `json:"created_at"`
}
```

### Requirements
- Load todos from `todos.json` on startup (if file doesn't exist, start empty)
- Save todos to `todos.json` after every change
- `add`: creates a new todo with auto-incremented ID and current timestamp
- `list`: prints all todos with status (✓ or ✗), ID, and text
- `done`: marks a todo as complete by ID — error if ID not found
- `delete`: removes a todo by ID — error if ID not found
- Use your Day 5 JSON skills (`json.NewEncoder`/`json.NewDecoder`)
- Return proper errors everywhere — no `log.Fatal` except in `main`

---

## Concepts to understand deeply

- `error` is just an interface: `type error interface { Error() string }`
- `errors.New` creates a simple sentinel error
- `fmt.Errorf("context: %w", err)` wraps an error — `%w` (not `%v`) makes it unwrappable
- `errors.Is(err, target)` checks if any error in the chain matches `target`
- `errors.As(err, &target)` checks if any error in the chain matches the type of `target`
- The difference between `%w` and `%v`: `%v` formats as string only, `%w` preserves the chain

---

## Intentional breaks (do after it works)

1. Use `%v` instead of `%w` in `fmt.Errorf`. Then try `errors.Is` — it won't match. Understand why.
2. Return a `*ValidationError` (pointer) from `validateUser` but do `errors.As(err, &ve)` where `ve` is `ValidationError` (not a pointer). See what happens.
3. In the todo app, try loading a corrupted JSON file (manually break the JSON). Handle the error gracefully rather than panicking.

---

## When done, write NOTES.md

- What is the difference between `errors.Is` and `errors.As`?
- What does `%w` do that `%v` does not?
- When would you use a custom error type vs `errors.New` vs `fmt.Errorf`?
- What does it mean for an error to be "wrapped"?
