package src

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	return &Server{listener: l}
}

func (s *Server) Close() {
	if err := s.listener.Close(); err != nil {
		panic(err)
	}
}

func (s *Server) Accept() *Connection {
	conn, err := s.listener.Accept()
	if err != nil {
		panic(err)
	}

	return NewConnection(conn)
}

func (s *Server) ProcessConnection(c *Connection) {
	defer c.Close()

	for {
		data, eof := c.ReadLine()

		for len(data) > 0 {
			c.Write(fmt.Sprintf("received: %s", data))

			data = data[len(data):]
		}

		if eof {
			fmt.Println("Connection closed for ", c.Addr())
			return
		}
	}
}
