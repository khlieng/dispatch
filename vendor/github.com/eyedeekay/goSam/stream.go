package goSam

import (
	"fmt"
)

// StreamConnect asks SAM for a TCP-Like connection to dest, has to be called on a new Client
func (c *Client) StreamConnect(id int32, dest string) error {
	r, err := c.sendCmd("STREAM CONNECT ID=%d %s %s DESTINATION=%s\n", id, c.from(), c.to(), dest)
	if err != nil {
		return err
	}

	// TODO: move check into sendCmd()
	if r.Topic != "STREAM" || r.Type != "STATUS" {
		return fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return ReplyError{result, r}
	}

	return nil
}

// StreamAccept asks SAM to accept a TCP-Like connection
func (c *Client) StreamAccept(id int32) (*Reply, error) {
	r, err := c.sendCmd("STREAM ACCEPT ID=%d SILENT=false\n", id)
	if err != nil {
		return nil, err
	}

	// TODO: move check into sendCmd()
	if r.Topic != "STREAM" || r.Type != "STATUS" {
		return nil, fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return nil, ReplyError{result, r}
	}

	return r, nil
}
