package main

import (
	"fmt"
	"paired/src"
)

func main() {
	fmt.Println("Paired: Server Version 0.6.9")

	server := src.NewServer()
	defer server.Close()

	fmt.Println("Server Started...")

	for {
		conn := server.Accept()
		fmt.Printf("Connection established with %s\n", conn.Addr())

		go server.ProcessConnection(conn)
	}
}
