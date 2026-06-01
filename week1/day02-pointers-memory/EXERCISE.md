# Day 2 — Pointers and Memory

## Goal
Understand where values live, when copies happen, and when references happen.
This is the most important day of Week 1. Everything else builds on this.

---

## Part 1: Linked List

Implement a singly linked list **from scratch** with no libraries.

### The Node type
```
Node:
  value int
  next  *Node   ← pointer to next node
```

### Functions to implement

```
NewList() *List                     — create an empty list
(l *List) Push(val int)             — add to front
(l *List) Pop() (int, error)        — remove from front, error if empty
(l *List) Append(val int)           — add to end
(l *List) Len() int                 — count of nodes
(l *List) Print()                   — print all values: [1 -> 2 -> 3]
```

### Usage
```
list := NewList()
list.Append(1)
list.Append(2)
list.Append(3)
list.Push(0)
list.Print()         // [0 -> 1 -> 2 -> 3]
fmt.Println(list.Len()) // 4
v, _ := list.Pop()
fmt.Println(v)       // 0
list.Print()         // [1 -> 2 -> 3]
```

---

## Part 2: Address Printing

Add a function `ShowAddresses` that:
1. Creates an `int` variable `x = 42`
2. Creates a pointer `p` pointing to `x`
3. Prints the value of `x`, the address of `x`, the value of `p` (same address), and the value at `p` (dereference)
4. Modifies `x` through the pointer: `*p = 100`
5. Prints `x` again to show it changed

Expected output:
```
x value:    42
x address:  0xc000...
p value:    0xc000...   (same address)
*p value:   42
after *p = 100, x = 100
```

---

## Part 3: Pass by Value vs Pass by Reference

Write two functions:

```go
func doubleByValue(n int)
func doubleByPointer(n *int)
```

Call both from main. Print the value of `n` after each call.
Observe: only `doubleByPointer` changes the original.

Then write a third function:
```go
func doubleSlice(s []int)
```

Call it, then print the slice after. Observe: slices are passed "by reference" even without a pointer. Understand why (hint: a slice is a struct with a pointer inside it).

---

## Concepts to understand deeply

- `&x` gives you the address of `x` (a pointer)
- `*p` dereferences a pointer (gives you the value at that address)
- Function arguments are always **copies** in Go
- Passing a pointer copies the pointer, but both point to the same memory
- A slice is NOT an array — it's a 3-word struct: `{pointer, length, capacity}`

---

## Intentional breaks (do these after it works)

1. In `Pop()`, try returning the zero value `0` instead of an error when the list is empty. Watch how the caller can't tell the difference between a real `0` and "empty list".
2. Change `doubleByPointer` to take `n int` instead of `n *int`. Watch the compiler error when you try to pass `&x`.
3. In your linked list, make `next` a `Node` (not `*Node`) — try to compile. Read the error: "invalid recursive type".

---

## When done, write NOTES.md

- What is the difference between a pointer and a value?
- When you pass a slice to a function, is it a copy or a reference? Why?
- What does "invalid recursive type" mean?
- Draw a diagram (in ASCII) of your linked list in memory after `Append(1)`, `Append(2)`, `Append(3)`.
