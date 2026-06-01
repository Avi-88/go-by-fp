# Day 5 — Maps, Struct Tags, and JSON

## Goal
Understand maps as a data structure, struct tags as metadata,
and JSON encoding/decoding as the bridge between Go types and the outside world.
This is the foundation for everything in Week 2 (HTTP APIs).

---

## Part 1: Maps

Write a `WordCount` function:

```go
func WordCount(s string) map[string]int
```

- Takes a string sentence
- Returns a map of each word to how many times it appears
- Case-insensitive: "The" and "the" count as the same word
- Ignore punctuation: strip `,`, `.`, `!`, `?`, `:`, `;`

```go
counts := WordCount("the cat sat on the mat the cat")
// map[cat:2 mat:1 on:1 sat:1 the:3]
```

Then write these map utility functions:

```go
func TopN(counts map[string]int, n int) []string
// returns the top N most frequent words, sorted by frequency descending

func MergeCounts(a, b map[string]int) map[string]int
// merges two word count maps, summing counts for shared keys
```

---

## Part 2: Structs with JSON Tags

Define this struct:

```go
type Person struct {
    ID        int
    FirstName string
    LastName  string
    Email     string
    Age       int
    Active    bool
    Tags      []string
}
```

Add JSON struct tags so that:
- `ID` serializes as `"id"`
- `FirstName` serializes as `"first_name"`
- `LastName` serializes as `"last_name"`
- `Email` serializes as `"email"`
- `Age` serializes as `"age"`, and is **omitted if zero**
- `Active` serializes as `"active"`
- `Tags` serializes as `"tags"`, and is **omitted if nil/empty**

---

## Part 3: JSON Encoding and Decoding

**Encoding (Go → JSON):**

```go
p := Person{
    ID: 1, FirstName: "Ada", LastName: "Lovelace",
    Email: "ada@example.com", Age: 36, Active: true,
    Tags: []string{"engineer", "mathematician"},
}
```

1. Encode `p` to JSON using `json.Marshal` — print the raw bytes as string
2. Encode again using `json.MarshalIndent` with 2-space indent — print pretty
3. Create a `Person` with zero `Age` and nil `Tags` — encode and verify those fields are omitted

**Decoding (JSON → Go):**

Decode this JSON string back into a `Person` struct:

```json
{"id":2,"first_name":"Alan","last_name":"Turing","email":"alan@example.com","active":true,"tags":["cs","math"]}
```

Print all fields of the decoded struct to verify.

---

## Part 4: Config Loader

Build a simple config loader that reads JSON from a file.

Define:

```go
type Config struct {
    AppName  string
    Port     int
    Debug    bool
    Database struct {
        Host     string
        Port     int
        Name     string
        Password string
    }
    AllowedOrigins []string
}
```

Add appropriate JSON tags to all fields (snake_case).

1. Write a function `SaveConfig(cfg Config, path string) error` that writes the config as indented JSON to a file
2. Write a function `LoadConfig(path string) (Config, error)` that reads and decodes it back
3. In `main`: create a config, save it, load it back, print it

Use `os.Create`, `os.Open`, `json.NewEncoder`, `json.NewDecoder` — stream-based, not `Marshal`/`Unmarshal`.

---

## Concepts to understand deeply

- Map zero value is `nil` — you must `make` it before writing
- Reading a missing key returns the zero value (no panic, no error)
- The comma-ok idiom: `val, ok := m[key]` — use `ok` to distinguish "missing" from "zero value"
- Struct tags are backtick strings: `` `json:"name,omitempty"` ``
- `json.Marshal` returns `[]byte` — convert to string with `string(b)`
- `json.NewDecoder` is preferred over `json.Unmarshal` when reading from a stream/file/request body

---

## Intentional breaks (do after it works)

1. Access a key that doesn't exist in a map: `m["missing"]`. Print the result. No panic — understand why.
2. Write to a nil map: `var m map[string]int` then `m["x"] = 1`. See the panic.
3. Remove the `json:"first_name"` tag from `FirstName`. Encode to JSON. See how the field name changes.
4. Try decoding invalid JSON: `json.Unmarshal([]byte("not json"), &p)`. Check the error.
5. Add an unexported field `secret string` to `Person`. Encode to JSON — it won't appear. Understand why.

---

## When done, write NOTES.md

- What is the zero value of a map, and what happens when you read from it vs write to it?
- What does `omitempty` do, and when would you use it?
- What is the difference between `json.Marshal` and `json.NewEncoder`?
- Why can't JSON encode unexported struct fields?
