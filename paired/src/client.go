package src

import (
	"net"
)

func NewClientConnection() *Connection {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	return NewConnection(conn)
}
