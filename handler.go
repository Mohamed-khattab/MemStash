var handelers = map[string]func([]Value) Value{
 "PING": ping,

}

func ping(args []Value) Value{
    return Value{typ: "string", str: "PONG"}
}
