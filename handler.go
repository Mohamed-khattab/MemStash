package main

import (
	"strconv"
	"sync"
)

var handelers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"DEL":     del,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"MGET":    mget,
	"MSET":    mset,
	"APPEND":  append_,
	"INCR":    incr,
	"INCRBY":  incrby,
	"DECR":    decr,
	"DECRBY":  decrby,
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
func mget(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'mget' command"}
	}
	var values []Value
	SETSMutex.RLock()
	defer SETSMutex.RUnlock()

	for _, key := range args {
		val, ok := SETs[key.bulk]
		if !ok {
			values = append(values, Value{typ: "null"})
		} else {
			values = append(values, Value{typ: "bulk", bulk: val})
		}
	}
	return Value{typ: "array", array: values}

}
func mset(args []Value) Value {
	if len(args)%2 != 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'mset' command"}
	}

	SETSMutex.Lock()
	defer SETSMutex.Unlock()
	for i := 0; i < len(args); i += 2 {
		key := args[i].bulk
		val := args[i+1].bulk
		SETs[key] = val
	}
	return Value{typ: "string", str: "OK"}
}

func append_(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'append' command "}
	}
	key := args[0].bulk
	val := args[1].bulk
	SETSMutex.Lock()
	defer SETSMutex.Unlock()
	_, ok := SETs[key]
	if !ok {
		SETs[key] = val
		return Value{typ: "integer", num: len(val)}
	}
	SETs[key] += val
	return Value{typ: "integer", num: len(SETs[key])}
}

func incr(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'incr' command "}
	}
	key := args[0].bulk
	val, ok := SETs[key]  // okay - if err nill if not now i wann  check  for the if err 
	if ok {
		SETs[key] = "1"
		return Value{typ: "integer", num: 1}
	}
	valInt, err := strconv.Atoi(val)
	if err == nil {
		SETs[key] = strconv.Itoa(valInt + 1)
		return Value{typ: "integer", num: valInt + 1}
	}
	return Value{typ: "error", str: "value is not an integer or out of range"}
}

func decr(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'decr' command "}
	}
	key := args[0].bulk
	val, ok := SETs[key]
	if !ok {
		SETs[key] = "-1"
		return Value{typ: "integer", num: -1}
	}
	valInt, err := strconv.Atoi(val)
	if err == nil {
		SETs[key] = strconv.Itoa(valInt - 1)
		return Value{typ: "integer", num: valInt - 1}
	}
	return Value{typ: "error", str: "value is not an integer or out of range"}
}

func incrby(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'incrby' command "}
	}
	key := args[0].bulk
	increment, err1 := strconv.Atoi(args[1].bulk)

	val, ok := SETs[key]
	if !ok {
		SETs[key] = strconv.Itoa(increment)
		return Value{typ: "integer", num: increment}
	}
	valInt, err2 := strconv.Atoi(val)
	if err1 != nil || err2 != nil {
		return Value{typ: "error", str: "value of the key or the increment is not an integer or maybe out of range"}
	}
	SETs[key] = strconv.Itoa(valInt + increment)
	return Value{typ: "integer", num: valInt + increment}
}
func decrby(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'decrby' command "}
	}
	key := args[0].bulk
	decrement, err1 := strconv.Atoi(args[1].bulk)

	val, ok := SETs[key]
	if !ok {
		SETs[key] = strconv.Itoa(-1 * decrement)
		return Value{typ: "integer", num: -1 * decrement}
	}
	valInt, err2 := strconv.Atoi(val)
	if err1 != nil || err2 != nil {
		return Value{typ: "error", str: "value of the key or the decrement is not an integer or out of range"}
	}
	SETs[key] = strconv.Itoa(valInt - decrement)
	return Value{typ: "integer", num: valInt - decrement}
}
