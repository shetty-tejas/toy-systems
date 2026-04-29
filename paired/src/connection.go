package src

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn,
		bufio.NewReader(conn),
	}
}

func (c *Connection) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) Close() {
	if err := c.conn.Close(); err != nil {
		panic(err)
	}
}

func (c *Connection) ReadLine() ([]byte, bool) {
	var eof bool
	data, err := c.reader.ReadBytes('\n')

	if err != nil {
		if errors.Is(err, io.EOF) {
			eof = true
		} else {
			panic(err)
		}
	}

	return data, eof
}

func (c *Connection) Write(data string) {
	i, err := fmt.Fprint(c.conn, data)
	if err != nil {
		panic(err)
	}

	if i != len(data) {
		panic("not all bytes written")
	}
}
