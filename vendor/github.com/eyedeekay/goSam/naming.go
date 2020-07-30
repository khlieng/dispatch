package goSam

import (
	"fmt"
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
		return "", fmt.Errorf("Unknown Reply: %+v\n", r)
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
