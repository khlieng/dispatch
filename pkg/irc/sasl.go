package irc

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"hash"
	"strings"

	"github.com/xdg-go/scram"
)

var DefaultSASLMechanisms = []string{
	"EXTERNAL",
	//"SCRAM-SHA-512",
	"SCRAM-SHA-256",
	//"SCRAM-SHA-1",
	"PLAIN",
}

type SASL interface {
	Name() string
	Step(response string) (string, error)
}

type SASLPlain struct {
	Username string
	Password string
}

func (s *SASLPlain) Name() string {
	return "PLAIN"
}

func (s *SASLPlain) Step(string) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString(s.Username)
	buf.WriteByte(0x0)
	buf.WriteString(s.Username)
	buf.WriteByte(0x0)
	buf.WriteString(s.Password)

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

type SASLExternal struct{}

func (s *SASLExternal) Name() string {
	return "EXTERNAL"
}

func (s *SASLExternal) Step(string) (string, error) {
	return "+", nil
}

var (
	scramHashes = map[string]scram.HashGeneratorFcn{
		"SHA-512": func() hash.Hash { return sha512.New() },
		"SHA-256": func() hash.Hash { return sha256.New() },
		"SHA-1":   func() hash.Hash { return sha1.New() },
	}

	ErrUnsupportedHash = errors.New("unsupported hash algorithm")
)

type SASLScram struct {
	Username string
	Password string
	Hash     string
	conv     *scram.ClientConversation
}

func (s *SASLScram) Name() string {
	return "SCRAM-" + s.Hash
}

func (s *SASLScram) Step(response string) (string, error) {
	if s.conv == nil {
		if hash, ok := scramHashes[s.Hash]; ok {
			client, err := hash.NewClient(s.Username, s.Password, "")
			if err != nil {
				return "", err
			}
			s.conv = client.NewConversation()
		} else {
			return "", ErrUnsupportedHash
		}
	}

	challenge := ""
	if response != "+" {
		b, err := base64.StdEncoding.DecodeString(response)
		if err != nil {
			return "", err
		}
		challenge = string(b)
	}

	res, err := s.conv.Step(challenge)
	if err != nil {
		return "", err
	}

	if s.conv.Done() {
		s.conv = nil
		return "+", nil
	} else {
		return base64.StdEncoding.EncodeToString([]byte(res)), nil
	}
}

func (c *Client) currentSASL() SASL {
	return c.saslMechanisms[c.currentSASLIndex]
}

func (c *Client) nextSASL() SASL {
	if c.currentSASLIndex < len(c.saslMechanisms)-1 {
		c.currentSASLIndex++
		return c.saslMechanisms[c.currentSASLIndex]
	}
	return nil
}

func (c *Client) handleSASL(msg *Message) {
	switch msg.Command {
	case AUTHENTICATE:
		auth, err := c.currentSASL().Step(msg.LastParam())
		if err != nil {
			c.finishCAP()
			return
		}

		for len(auth) >= 400 {
			c.authenticate(auth)
			auth = auth[400:]
		}
		if len(auth) > 0 {
			c.authenticate(auth)
		} else {
			c.authenticate("+")
		}

	case RPL_SASLMECHS:
		if len(msg.Params) > 1 {
			supportedMechs := strings.Split(msg.Params[1], ",")

			for i, mech := range c.saslMechanisms {
				for _, supported := range supportedMechs {
					if mech.Name() == supported && i > c.currentSASLIndex {
						c.currentSASLIndex = i
						c.authenticate(mech.Name())
						return
					}
				}
			}
		}
		c.finishCAP()

	case ERR_SASLFAIL:
		if next := c.nextSASL(); next != nil {
			c.authenticate(next.Name())
		} else {
			c.finishCAP()
		}

	case RPL_SASLSUCCESS, RPL_LOGGEDIN:
		c.finishCAP()

	case ERR_NICKLOCKED, ERR_SASLTOOLONG, ERR_SASLABORTED:
		c.finishCAP()
	}
}
