package main

import (
	"paired/src"
)

func main() {
	server := src.NewServer()
	defer server.Close()

	for {
		conn := server.Accept()

		go conn.ProcessAndClose()
	}
}
