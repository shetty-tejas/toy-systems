package main

import (
	"bufio"
	"fmt"
	"os"
	"paired/src"
)

func main() {
	fmt.Println("Paired: Client Version 0.6.9")
	conn := src.NewClientConnection()
	defer conn.Close()

	fmt.Printf("Connection established with %s\n", conn.Addr())

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")

		if scanner.Scan() {
			request := scanner.Bytes()

			conn.Write(request)
			response, eof := conn.ReadMessage()

			_, err := fmt.Fprintln(os.Stdout, string(response))
			if err != nil {
				panic(err)
			}

			if eof {
				fmt.Println("Connection closed for", conn.Addr())
				return
			}
		}
	}
}
