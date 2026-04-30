package src

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn}
}

func (c *Connection) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) Close() {
	if err := c.conn.Close(); err != nil {
		panic(err)
	}
}

func (c *Connection) ReadMessage() ([]byte, bool) {
	header := make([]byte, 4)
	i, err := io.ReadFull(c.conn, header)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return []byte{}, true
		}

		panic(err)
	}

	if i != len(header) {
		panic("header read seems to be wrong")
	}

	var length uint32

	i, err = binary.Decode(header, binary.BigEndian, &length)
	if err != nil {
		panic(err)
	}

	if i != len(header) {
		panic("header decode seems to be wrong")
	}

	message := make([]byte, length)
	i, err = io.ReadFull(c.conn, message)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return message, true
		}

		panic(err)
	}

	if i != len(message) {
		panic("message seems to be incomplete")
	}

	return message, false
}

func (c *Connection) Write(data []byte) {
	size := uint32(len(data))

	header := make([]byte, 4)

	i, err := binary.Encode(header, binary.BigEndian, size)
	if err != nil {
		panic(err)
	}

	if i != len(header) {
		panic("encoding header went wrong")
	}

	data = append(header, data...)

	i, err = c.conn.Write(data)
	if err != nil {
		panic(err)
	}

	if i != len(data) {
		panic("not all bytes written")
	}
}
