package main

import "sync"

var handelers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"DEL":  del,
    "HSET": hset,
    "HGET": hget,
    "HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: "PONG " + args[0].bulk}
}

// init some data structure we gonna need
// for set and get commands
var SETs = map[string]string{}
var SETSMutex = sync.RWMutex{}

// for hset and hget commands
var HSETS = map[string]map[string]string{}
var HSETSMutex = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].bulk
	val := args[1].bulk
	SETSMutex.Lock()
	defer SETSMutex.Unlock()
	SETs[key] = val
	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].bulk
	SETSMutex.RLock()

	defer SETSMutex.RUnlock()
	val, ok := SETs[key]
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}
func del(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}
	key := args[0].bulk
	SETSMutex.Lock()
	defer SETSMutex.Unlock()
	_, ok := SETs[key]
	if !ok {
		return Value{typ: "integer", num: 0}
	}
	delete(SETs, key)
	return Value{typ: "string", str: "OK"}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	key := args[0].bulk
	field := args[1].bulk
	val := args[2].bulk
	defer HSETSMutex.Unlock()
	HSETSMutex.Lock()
	// Check if the key exists in the hash map,
	// if not, create a new map for the key
	if _, ok := HSETS[key]; !ok {
		HSETS[key] = make(map[string]string)
	}
	HSETS[key][field] = val
	return Value{typ: "string", str: "OK"}
}
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}
	key := args[0].bulk
	field := args[1].bulk
	HSETSMutex.RLock()
	defer HSETSMutex.RUnlock()
	// Check if the key exists in the hash map,
	// if not, return nil
	if _, ok := HSETS[key]; !ok {
		return Value{typ: "null"}
	}
	val, ok := HSETS[key][field]
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}
func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}
	key := args[0].bulk
	// Check if the key exists in the hash map,
	// if not, return nil

	HSETSMutex.RLock()
	val, ok := HSETS[key]

	HSETSMutex.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}

	var values []Value
	for k, v := range val {
		values = append(values, Value{typ: "bulk", bulk: k}, Value{typ: "bulk", bulk: v})
	}
	return Value{typ: "array", array: values}
}
