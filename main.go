package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := newAOF("dump.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.close()

	conn, err := ln.Accept() // recive connection
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // close the connection once finished

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}
		if value.typ != "array" {
			fmt.Println("Invalid Request, Expected an array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]
		writer := newWriter(conn)

		handler , ok := handelers[command]

		if !ok{
			fmt.Println("Unknown command", command)
			writer.Write(Value{typ: "error", str: "Unknown command"})
			continue
		}
		if command == "SET"  || command == "HSET" {
			aof.write(value)
		}
		result :=handler(args)
		writer.Write(result)

	}
}
