package goSam

import (
	"fmt"
	"net"
)

// AcceptI2P creates a new Client and accepts a connection on it
func (c *Client) AcceptI2P() (net.Conn, error) {
	listener, err := c.Listen()
	if err != nil {
		return nil, err
	}
	return listener.Accept()
}

// Listen creates a new Client and returns a net.listener which *must* be started
// with Accept
func (c *Client) Listen() (net.Listener, error) {
	return c.ListenI2P(c.destination)
}

// ListenI2P creates a new Client and returns a net.listener which *must* be started
// with Accept
func (c *Client) ListenI2P(dest string) (net.Listener, error) {
	var err error
	var id int32
	c.id = c.NewID()
	c.destination, err = c.CreateStreamSession(id, dest)
	if err != nil {
		return nil, err
	}

	fmt.Println("destination:", c.destination)

	c, err = c.NewClient()
	if err != nil {
		return nil, err
	}

	if c.debug {
		c.SamConn = WrapConn(c.SamConn)
	}

	return c, nil
}

// Accept accepts a connection on a listening goSam.Client(Implements net.Listener)
// or, if the connection isn't listening yet, just calls AcceptI2P for compatibility
// with older versions.
func (c *Client) Accept() (net.Conn, error) {
	if c.id == 0 {
		return c.AcceptI2P()
	}
	resp, err := c.StreamAccept(c.id)
	if err != nil {
		return nil, err
	}

	fmt.Println("Accept Resp:", resp)

	return c.SamConn, nil
}
