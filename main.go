package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := ln.Accept() // recive connection
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // close the connection once finished

	for {
		resp :=  NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(value)

		_, err = conn.Write([]byte("+OK\r\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err.Error())
			continue // Continue to the next iteration instead of exiting
		}

	}
}
