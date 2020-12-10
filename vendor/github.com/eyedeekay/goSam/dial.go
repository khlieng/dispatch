package goSam

import (
	"context"
	"log"
	"net"
	"strings"
)

// DialContext implements the net.DialContext function and can be used for http.Transport
func (c *Client) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	errCh := make(chan error, 1)
	connCh := make(chan net.Conn, 1)
	go func() {
		if conn, err := c.DialContextFree(network, addr); err != nil {
			errCh <- err
		} else if ctx.Err() != nil {
			log.Println(ctx)
			errCh <- ctx.Err()
		} else {
			connCh <- conn
		}
	}()
	select {
	case err := <-errCh:
		return nil, err
	case conn := <-connCh:
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) Dial(network, addr string) (net.Conn, error) {
	return c.DialContext(context.TODO(), network, addr)
}

// Dial implements the net.Dial function and can be used for http.Transport
func (c *Client) DialContextFree(network, addr string) (net.Conn, error) {
	portIdx := strings.Index(addr, ":")
	if portIdx >= 0 {
		addr = addr[:portIdx]
	}
	addr, err := c.Lookup(addr)
	if err != nil {
		log.Printf("LOOKUP DIALER ERROR %s %s", addr, err)
		return nil, err
	}

	c.destination, err = c.CreateStreamSession(c.id, c.destination)
	if err != nil {
		c.Close()
		d, err := c.NewClient(c.id + 1) /**/
		if err != nil {
			return nil, err
		}
		d.destination, err = d.CreateStreamSession(d.id, c.destination)
		if err != nil {
			return nil, err
		}
		d, err = d.NewClient(d.id)
		if err != nil {
			return nil, err
		}
		//	  d.lastaddr = addr
		err = d.StreamConnect(d.id, addr)
		if err != nil {
			return nil, err
		}
		c = d
		return d.SamConn, nil
	}
	c, err = c.NewClient(c.id)
	if err != nil {
		return nil, err
	}
	err = c.StreamConnect(c.id, addr)
	if err != nil {
		return nil, err
	}
	return c.SamConn, nil
}
