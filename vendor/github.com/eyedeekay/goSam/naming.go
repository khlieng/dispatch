package goSam

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
)

// Lookup askes SAM for the internal i2p address from name
func (c *Client) Lookup(name string) (string, error) {
	r, err := c.sendCmd("NAMING LOOKUP NAME=%s\n", name)
	if err != nil {
		return "", nil
	}

	// TODO: move check into sendCmd()
	if r.Topic != "NAMING" || r.Type != "REPLY" {
		return "", fmt.Errorf("Naming Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return "", ReplyError{result, r}
	}

	if r.Pairs["NAME"] != name {
		// somehow different on i2pd
		if r.Pairs["NAME"] != "ME" {
			return "", fmt.Errorf("Lookup() Replyed to another name.\nWanted:%s\nGot: %+v\n", name, r)
		}
		fmt.Fprintln(os.Stderr, "WARNING: Lookup() Replyed to another name. assuming i2pd c++ fluke")
	}

	return r.Pairs["VALUE"], nil
}

func (c *Client) forward(client, conn net.Conn) {
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func (c *Client) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	if c.lastaddr == "invalid" || c.lastaddr != name {
		client, err := c.DialContext(ctx, "", name)
		if err != nil {
			return ctx, nil, err
		}
		ln, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			return ctx, nil, err
		}
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					fmt.Println(err.Error())
				}
				go c.forward(client, conn)
			}
		}()
	}
	return ctx, nil, nil
}
