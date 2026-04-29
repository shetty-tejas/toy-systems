package src

import "net"

type Connection struct {
	conn net.Conn
}

func (c *Connection) Close() {
	if err := c.conn.Close(); err != nil {
		panic(err)
	}
}

func (c *Connection) Process() {
	// TODO: Complete

	_, err := c.conn.Write([]byte("Hello, World!"))
	if err != nil {
		panic(err)
	}
}

func (c *Connection) ProcessAndClose() {
	defer c.Close()

	c.Process()
}
