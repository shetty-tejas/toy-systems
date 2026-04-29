package src

import (
	"errors"
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
		panic(errors.Join(err, conn.Close()))
	}

	return &Connection{conn}
}
