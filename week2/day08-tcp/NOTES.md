# Day 8 Notes — TCP from Scratch

## Questions to answer after completing the exercise

- What is the difference between `net.Listen` and `net.Dial`?

net.Dial is for clients to initiate a connection to a server with a specific protocol and port while net.Listen is for a server to listen for a particulat protocol and port

- Why does `io.Copy(conn, conn)` implement an echo server? Walk through what it does.
it copmletely copies the io stream of one io reader/writer to another , so ones input is also the others input

- What problem does framing solve, and why doesn't TCP solve it for you?
tcp doesnt provide specific framing sized by default so determining actual payload and metadata or anything else is really difficult so even if the order is preserved the splitting of content isnt reliable without framing

- What is a connection deadline, and when would you use one in production?
a connection deadline is a trigger set on the connection which will throw a timeout error after a certain time  has passed. It could be used to limit time allocated to each client either in a fixed time format or a inactivity like format, incase of inactivity format we need to keep renewing the timeout each time the client interacts ( send payload )

- What does the chat server teach you about shared mutable state across goroutines?

---

## Observations

<!-- fill in as you work through each part -->
