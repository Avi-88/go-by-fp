# Day 8 Notes — TCP from Scratch

## Questions to answer after completing the exercise

- What is the difference between `net.Listen` and `net.Dial`?

- Why does `io.Copy(conn, conn)` implement an echo server? Walk through what it does.

- What problem does framing solve, and why doesn't TCP solve it for you?

- What is a connection deadline, and when would you use one in production?

- What does the chat server teach you about shared mutable state across goroutines?

---

## Observations

<!-- fill in as you work through each part -->
