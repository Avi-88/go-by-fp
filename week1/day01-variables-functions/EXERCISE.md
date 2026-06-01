# Day 1 — Variables, Types, Functions

## Goal
Understand Go's type system, how values are passed, and how functions work.
No frameworks, no libraries beyond `fmt` and `os`.

---

## Exercise: Command-Line Calculator

Build a calculator that:
1. Takes two numbers and an operator as CLI arguments
2. Supports: `+`, `-`, `*`, `/`, `%`, `^` (power)
3. Handles bad input gracefully (wrong number of args, bad operator, divide by zero)
4. Prints the result with the full expression: `3 + 4 = 7`

### Usage
```
go run main.go 10 + 3     → 10 + 3 = 13
go run main.go 10 / 0     → error: division by zero
go run main.go 2 ^ 8      → 2 ^ 8 = 256
go run main.go             → error: usage: calc <num> <op> <num>
```

---

## Constraints (these are the point)

- Write a separate function for each operation. Do NOT use a switch inside main.
- Your `calculate` function must return BOTH a result AND an error.
- Parse the string arguments into numbers yourself using `strconv`.
- Do NOT use `math/big`. Use `float64` for everything.

---

## Concepts to notice as you write this

1. **Multiple return values** — Go functions can return (value, error). This is not optional style, it's idiomatic Go.
2. **Named types** — Try defining `type Operator string` and using it.
3. **Value semantics** — Every argument to your functions is a copy. Verify this by printing the address of a variable inside and outside a function: `fmt.Printf("%p\n", &x)`
4. **Short variable declaration** — `:=` vs `var`. Know when each is used.
5. **Zero values** — What is the zero value of `float64`? `string`? `bool`? Declare variables without initializing them and print them.

---

## Intentional breaks (do these after it works)

1. Return only the result from `calculate`, discard the error. Watch what happens on divide by zero.
2. Declare a variable with `var x float64 = "hello"` and read the compiler error carefully.
3. Call a function that returns `(float64, error)` and assign it to a single variable (no `_`). Read the error.

---

## When done, write NOTES.md

Answer these questions in your own words:
- What is the difference between `:=` and `var`?
- What does Go do when you don't handle an error return?
- What surprised you?
