# Double Distributed KV

A simple program for saving and looking up text while it is running. You use it by typing commands, similar to a chat prompt.

## Get started

1. Install [Go](https://go.dev/dl/) (version 1.27.1 or newer).
2. Open a terminal in this project folder.
3. Start the program:

   ```sh
   go run ./cmd/repl
   ```

## Try it

Type a command after the `double-distributed-kv>` prompt:

```text
SET name "Ada Lovelace"
GET name
EXISTS name
DEL name
EXIT
```

`SET` saves text, `GET` shows it, `EXISTS` checks whether it is saved, and `DEL` removes it. You can put quotes around text that contains spaces. Commands can be typed in upper or lower case.

## What to know

Your saved text is temporary: it disappears when you quit the program. This is a small local demo, not a full Redis server, and it does not connect to other computers.
