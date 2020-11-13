package goSam

import (
	"errors"

	"github.com/eyedeekay/sam3/i2pkeys"
)

// NewDestination generates a new I2P destination, creating the underlying
// public/private keys in the process. The public key can be used to send messages
// to the destination, while the private key can be used to reply to messages
func (c *Client) NewDestination(sigType ...string) (i2pkeys.I2PKeys, error) {
	var (
		sigtmp string
		keys   i2pkeys.I2PKeys
	)
	if len(sigType) > 0 {
		sigtmp = sigType[0]
	}
	r, err := c.sendCmd(
		"DEST GENERATE %s\n",
		sigtmp,
	)
	if err != nil {
		return keys, err
	}
	var pub, priv string
	if priv = r.Pairs["PRIV"]; priv == "" {
		return keys, errors.New("failed to generate private destination key")
	}
	if pub = r.Pairs["PUB"]; pub == "" {
		return keys, errors.New("failed to generate public destination key")
	}
	return i2pkeys.NewKeys(i2pkeys.I2PAddr(pub), priv), nil
}
