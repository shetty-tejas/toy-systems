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

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")

		request, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		conn.Write(request)
		response, eof := conn.ReadLine()

		_, err = fmt.Fprint(os.Stdout, string(response))
		if err != nil {
			panic(err)
		}

		if eof {
			fmt.Println("Connection closed for", conn.Addr())
			return
		}
	}
}
