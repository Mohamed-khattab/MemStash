package main

var handelers = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(args []Value) Value {
    if len(args) == 0 {
        return Value{typ: "string", str: "PONG"}
    }
	return Value{typ: "string", str: "PONG " + args[0].bulk}
}

