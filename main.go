package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn ,err := ln.Accept() // recive connection 
	if err!= nil {
		fmt.Println(err)
		return 
	}

	defer conn.Close() // close the connection once finished 

	for{
		// init a buffer to hold the incoming date from network connection
		buf := make([]byte , 1024)  
		// read the data from the connection 
		_, err = conn.Read(buf) 
		if err != nil {
			if err ==io.EOF{
				break 
			}
			fmt.Println("Error reading from client ", err.Error())
			os.Exit(1)
			return 
		}
		conn.Write([]byte("+OK\r\n"))

	}
}
