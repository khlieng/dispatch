package irc

import (
	"bytes"
	"encoding/base64"
)

type SASL interface {
	Name() string
	Encode() string
}

type SASLPlain struct {
	Username string
	Password string
}

func (s *SASLPlain) Name() string {
	return "PLAIN"
}

func (s *SASLPlain) Encode() string {
	buf := bytes.Buffer{}
	buf.WriteString(s.Username)
	buf.WriteByte(0x0)
	buf.WriteString(s.Username)
	buf.WriteByte(0x0)
	buf.WriteString(s.Password)

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

type SASLExternal struct{}

func (s *SASLExternal) Name() string {
	return "EXTERNAL"
}

func (s *SASLExternal) Encode() string {
	return "+"
}

func (c *Client) handleSASL(msg *Message) {
	switch msg.Command {
	case AUTHENTICATE:
		auth := c.SASL.Encode()

		for len(auth) >= 400 {
			c.write("AUTHENTICATE " + auth)
			auth = auth[400:]
		}
		if len(auth) > 0 {
			c.write("AUTHENTICATE " + auth)
		} else {
			c.write("AUTHENTICATE +")
		}

	case RPL_SASLSUCCESS:
		c.write("CAP END")

	case ERR_NICKLOCKED, ERR_SASLFAIL, ERR_SASLTOOLONG, ERR_SASLABORTED, RPL_SASLMECHS:
		c.write("CAP END")
	}
}
