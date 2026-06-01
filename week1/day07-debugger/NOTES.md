# Day 7 Notes

## What is the difference between `next` and `step`?

next is used to get to the next line that will be executed, if its a function it wont automatically traverse inside it. Step is used to traverse inside the function directly

## When is `stepout` useful?

stepout can be useful if once you are done debugging the logic of a function from inside and want to debug the parent caller further you can resume the debugging from there 

## What did you learn watching linked list Append in the debugger?

(your answer here)

## What surprised you about goroutines in a simple program?

There are a few ( 5 i guess ) that will always exist alongside the ones that your program has and are mainly related to garbage collection in golang

## Part 1 observations — Linked List

1. When `Append` is called with `value=2` on a list that already has `[1]`, what is `l.head.value`? - the value is 1
2. After the for loop finishes, what is `tail`? Print it with `p tail`. - the tail value is also 1 and next is nil
3. After `tail.next = &nNode`, print `tail.next.value`. What do you see? - now the tail.next is node with value 2
4. Step through `Push` — at what point does `l.head` change? - the value of l.head changes after l.head = &node

## Part 2 observations — Calculator

1. Print `A` and `B` with `p A`, `p B`. Are they what you expect? - yes
2. Step into the return statement. What does Go return for `float64(A)/float64(B)`? -  it returns zero and nil
3. Now run with `dlv debug . -- 10 / 0`. Step through the error path. Print `err` after the return. - Done

## Part 3 observations — Todo CLI

1. Print `currTodos` before the append. What does it look like? - there are previous 2 tasks present there
2. Print `newId` — is it what you expect? - yes
3. Step through `updateTodos`. At what line does the file get written? - at line 141
4. Run `goroutines` — how many goroutines are running in this simple program? - 6 goroutines, 5 of which seem to be related to garbage collector

## Part 4 observations — Memory Addresses

1. Print `x` and `&x` — observe the address - Done
2. Print `p` — does it match `&x`? - yes it does
3. Print `*p` — does it match `x`? - yes it dows
4. Step to after `*p = 100` — print `x`. Confirm it changed. - it did change


