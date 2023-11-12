# MemStash: Redis-Like Server Documentation

This documentation provides an overview of MemStash, a Redis-like server implemented in Go. MemStash supports a variety of common Redis commands, offers serialization and deserialization of data in the Redis Serialization Protocol (RESP) format, and includes persistent logging of executed commands in an Append-Only File (AOF). Additionally, it emphasizes future plans for further development, including the addition of commands and a testing section.


## Extended Redis Commands Support

MemStash currently supports the following common Redis commands:

- **Basic Commands:**
  - PING
  - SET
  - GET
  - DEL

- **Hash Commands:**
  - HSET
  - HGET
  - HGETALL

## Project Structure Summary

### `main.go`

- Main server entry point.
- Listens on port 6379.
- Handles basic Redis commands and logs data in an Append-Only File (AOF).

### `deserializer.go`

- Converts RESP format to Go data structures.
- Supports types: string, error, integer, bulk, array.

### `serializer.go`

- Converts Go data structures to RESP format.
- Supports types: string, error, integer, bulk, array, null.

### `handler.go`

- Defines handlers for Redis-like commands.
- Handles PING, SET, GET, DEL, HSET, HGET, HGETALL.
- Utilizes in-memory structures and mutexes for thread safety.

### `aof.go`

- Implements Append-Only File (AOF) for persistent command logging.
- Uses buffered reader and file pointer.
- Periodically syncs AOF file with storage (`dump.aof` file).


## Execution

1. **Starting MemStash:**
   - Run MemStash by executing `make run` and it will run all the fo files .
### Video Walkthrough

