# Day 4 — Slices vs Arrays (Internals)

## Goal
Understand the difference between arrays and slices at the memory level.
Build your own versions of common slice operations to understand what Go does under the hood.

---

## Part 1: Arrays vs Slices — Feel the Difference

Write a `main` that demonstrates:

1. Declare an array: `var arr [5]int = [5]int{1, 2, 3, 4, 5}`
2. Declare a slice from it: `s := arr[1:4]`
3. Modify `s[0] = 99` — then print both `arr` and `s`. What changed?
4. Print `len(s)` and `cap(s)` — understand what each means
5. Now `append` to `s` beyond its capacity and print again

The key question: when does a slice stop sharing memory with its original array?

---

## Part 2: Implement These Functions From Scratch

Do NOT use built-in slice tricks for the logic. Write the loops yourself.

```go
func Map(s []int, f func(int) int) []int
// applies f to every element, returns new slice

func Filter(s []int, f func(int) bool) []int
// returns new slice with only elements where f returns true

func Reduce(s []int, initial int, f func(int, int) int) int
// folds the slice into a single value

func Contains(s []int, val int) bool
// returns true if val is in s

func Reverse(s []int) []int
// returns a new reversed slice (do NOT modify the original)

func Unique(s []int) []int
// returns a new slice with duplicates removed, preserving order
```

### Test them in main with:
```go
nums := []int{1, 2, 3, 4, 5, 2, 3, 1}

doubled  := Map(nums, func(n int) int { return n * 2 })
evens    := Filter(nums, func(n int) bool { return n%2 == 0 })
sum      := Reduce(nums, 0, func(acc, n int) int { return acc + n })
reversed := Reverse(nums)
unique   := Unique(nums)

fmt.Println(doubled)   // [2 4 6 8 10 4 6 2]
fmt.Println(evens)     // [2 4 2]
fmt.Println(sum)       // 21
fmt.Println(reversed)  // [1 3 2 5 4 3 2 1]
fmt.Println(unique)    // [1 2 3 4 5]
```

---

## Part 3: Grow Your Own Slice

Implement a simple dynamic array (like how Go slices work internally):

```go
type DynamicArray struct {
    data     []int
    length   int
    capacity int
}

func NewDynamicArray() *DynamicArray
func (d *DynamicArray) Push(val int)    // double capacity when full
func (d *DynamicArray) Get(i int) (int, error)
func (d *DynamicArray) Len() int
func (d *DynamicArray) Cap() int
func (d *DynamicArray) Print()
```

When `Push` is called and `length == capacity`:
1. Allocate a new slice with double the capacity (`make([]int, 0, newCap)`)
2. Copy existing elements over
3. Add the new element

Start with an initial capacity of 2. Print the capacity every time it grows so you can see the doubling.

---

## Concepts to understand deeply

- An array's size is part of its type: `[3]int` and `[5]int` are different types
- A slice is `{pointer, len, cap}` — 3 words (24 bytes on 64-bit)
- `len` = how many elements you can access
- `cap` = how many elements before a new array is allocated
- `append` returns a new slice — always capture the return value
- Two slices can share the same underlying array — modifying one affects the other

---

## Intentional breaks (do after it works)

1. Do `s2 := s` (slice assignment) then `s2[0] = 999`. Check what happened to `s`. Why?
2. In `Reverse`, modify the original slice instead of making a copy. Call `Reverse(nums)` then print `nums`. Observe.
3. In your `DynamicArray.Push`, forget to update `d.length`. See what breaks.
4. Try `append` on a nil slice: `var s []int` then `s = append(s, 1)`. It works — understand why.

---

## When done, write NOTES.md

- What is the difference between `len` and `cap`?
- When does `append` allocate a new array vs reuse the existing one?
- Why is it dangerous for two slices to share the same underlying array?
- What does `make([]int, 3, 5)` create?
