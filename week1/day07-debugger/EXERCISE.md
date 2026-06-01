# Day 7 — Debugger Session (dlv)

## Goal
Use the Go debugger (dlv) to inspect programs you've already written.
No new features — today is about developing the skill of *seeing inside* a running program.
This is what separates engineers who guess from engineers who know.

---

## Setup: Install dlv

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Verify:
```bash
dlv version
```

---

## How to run dlv

```bash
# Debug a package
dlv debug .

# Debug with arguments (e.g. your todo CLI)
dlv debug . -- add "buy milk"

# Attach to a running process
dlv attach <pid>
```

---

## Core dlv commands (learn these)

| Command | Short | What it does |
|---------|-------|-------------|
| `break main.main` | `b main.main` | Set breakpoint at function |
| `break main.go:42` | `b main.go:42` | Set breakpoint at line |
| `continue` | `c` | Run until next breakpoint |
| `next` | `n` | Step over (execute one line, don't enter functions) |
| `step` | `s` | Step into (enter the called function) |
| `stepout` | `so` | Step out of current function |
| `print x` | `p x` | Print value of variable |
| `locals` | | Print all local variables |
| `args` | | Print function arguments |
| `stack` | | Print call stack |
| `goroutines` | | List all goroutines |
| `list` | `l` | Show source around current line |
| `quit` | `q` | Exit debugger |

---

## Part 1: Debug the Linked List (Day 2)

Go to `week1/day02-pointers-memory/` and run:

```bash
dlv debug .
```

Set a breakpoint at `Append` and step through it:

```
b main.Append
c
```

**Questions to answer by observing in the debugger:**
1. When `Append` is called with `value=2` on a list that already has `[1]`, what is `l.head.value`? - the value is 1
2. After the for loop finishes, what is `tail`? Print it with `p tail`. - the tail value is also 1 and next is nil
3. After `tail.next = &nNode`, print `tail.next.value`. What do you see? - now the tail.next is node with value 2
4. Step through `Push` — at what point does `l.head` change? the value of l.head changes after l.head = &node

Write your observations in NOTES.md.

---

## Part 2: Debug the Calculator (Day 1)

Go to `week1/day01-variables-functions/` and run:

```bash
dlv debug . -- 10 / 3
```

Set a breakpoint at `divide`:
```
b main.divide
c
```

**Questions to answer:**
1. Print `A` and `B` with `p A`, `p B`. Are they what you expect? yes
2. Step into the return statement. What does Go return for `float64(A)/float64(B)`? it returns zero and nil
3. Now run with `dlv debug . -- 10 / 0`. Step through the error path. Print `err` after the return.

---

## Part 3: Debug the Todo CLI (Day 6)

Go to `week1/day06-errors-todo/` and run with arguments:

```bash
dlv debug . -- add "debug session"
```

Set a breakpoint at `add`:
```
b main.add
c
```

**Questions to answer:**
1. Print `currTodos` before the append. What does it look like? there are previous 2 tasks present there
2. Print `newId` — is it what you expect? yes
3. Step through `updateTodos`. At what line does the file get written? at line 141
4. Run `goroutines` — how many goroutines are running in this simple program? 6 goroutines, 5 of which seem to be related to garbage collector

---

## Part 4: Inspect Memory Addresses

Go back to `week1/day02-pointers-memory/` and set a breakpoint in `ShowAddresses`:

```
b main.ShowAddresses
c
```

1. Print `x` and `&x` — observe the address
2. Print `p` — does it match `&x`? yes it does
3. Print `*p` — does it match `x`? yes it dows
4. Step to after `*p = 100` — print `x`. Confirm it changed. it did change

This is the same exercise you did with `fmt.Printf` on Day 2, but now you're watching it happen live in the debugger.

---

## Part 5: Conditional Breakpoints

In `week1/day04-slices-arrays/`, set a breakpoint inside `DynamicArray.Push` that only triggers when capacity doubles:

```
b main.(*DynamicArray).Push
cond 1 d.length == d.capacity
c
```

Step through only the reallocation path. Print `d.capacity` before and after. it changed from 2 to 4 

---

## What to write in NOTES.md

- What is the difference between `next` and `step`?
- When is `stepout` useful?
- What did you learn by watching the linked list Append in the debugger that you couldn't see from the code alone?
- What surprised you about the number of goroutines in a simple program?
