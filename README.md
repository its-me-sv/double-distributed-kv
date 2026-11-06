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

## Recent update: B+ tree groundwork (in progress)

The program described above still works the same way. Behind the scenes, the store keeps everything in a Go map in memory ([`internal/store`](internal/store/store.go)).

Work has started on a B+ tree in [`internal/datastructures/bplustree`](internal/datastructures/bplustree). A B+ tree keeps keys sorted and splits data into fixed-size pages, which is what you need to eventually save data to disk instead of keeping it all in memory. **It is not connected to the REPL or the store yet**, so the commands above don't use it.

What exists so far:

- **Page-sized nodes** ([`node.go`](internal/datastructures/bplustree/node.go)): each node is a 4096-byte block. Leaf nodes hold key-value pairs and internal nodes hold pointers to child pages. A node is laid out like this:

  ```text
  | type | nkeys |  pointers  |  offsets   | key-values | unused |
  |  2B  |   2B  | nkeys × 8B | nkeys × 2B |     ...    |        |
  ```

  Each key-value entry stores the key length (2 bytes), the value length (2 bytes), then the key and value bytes. The offsets array means you can jump straight to the nth entry without reading the ones before it. Keys can be up to 1000 bytes and values up to 3000 bytes.
- **Leaf edits** ([`utils.go`](internal/datastructures/bplustree/utils.go)): `LeafInsert`, `LeafUpdate` and `LeafSet` never change a node in place. They copy the old node into a new one and add or replace the key on the way, a copy-on-write approach that suits on-disk storage. `NodeLookupLS` scans the node in order to find where a key is, or where it would go.
- **Node splitting**: when an edit makes a node bigger than one page, `NodeSplitInto2` and `NodeSplitInto3` split it into two or three nodes that each fit in a page.
- **Tree skeleton** ([`bplustree.go`](internal/datastructures/bplustree/bplustree.go)): a `Tree` type that holds the root page number and callbacks to read, allocate and free pages. Tree-level operations (lookup, insert and delete from the root down) haven't been written yet.

Next steps: add tree-level insert, lookup and delete, store pages in a file, and swap the map in `internal/store` for the B+ tree.
