package goSam

import (
	"fmt"
)

// StreamConnect asks SAM for a TCP-Like connection to dest, has to be called on a new Client
func (c *Client) StreamConnect(id int32, dest string) error {
	if dest == "" {
		return nil
	}
	r, err := c.sendCmd("STREAM CONNECT ID=%d DESTINATION=%s %s %s\n", id, dest, c.from(), c.to())
	if err != nil {
		return err
	}

	// TODO: move check into sendCmd()
	if r.Topic != "STREAM" || r.Type != "STATUS" {
		return fmt.Errorf("Stream Connect Unknown Reply: %+v\n", r)
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
		return nil, fmt.Errorf("Stream Accept Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return nil, ReplyError{result, r}
	}

	return r, nil
}
