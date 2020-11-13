package goSam

import (
	"fmt"
	//	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CreateStreamSession creates a new STREAM Session.
// Returns the Id for the new Client.
func (c *Client) CreateStreamSession(id int32, dest string) (string, error) {
	if dest == "" {
		dest = "TRANSIENT"
	}
	c.id = id
	r, err := c.sendCmd(
		"SESSION CREATE STYLE=STREAM ID=%d DESTINATION=%s %s %s %s %s \n",
		c.id,
		dest,
		c.from(),
		c.to(),
		c.sigtype(),
		c.allOptions(),
	)
	if err != nil {
		return "", err
	}

	// TODO: move check into sendCmd()
	if r.Topic != "SESSION" || r.Type != "STATUS" {
		return "", fmt.Errorf("Session Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return "", ReplyError{ResultKeyNotFound, r}
	}
	c.destination = r.Pairs["DESTINATION"]
	return c.destination, nil
}
