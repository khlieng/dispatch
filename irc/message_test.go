package irc

import (
	"testing"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	cases := []struct {
		input    string
		expected *Message
	}{
		{
			":user CMD #chan :some message\r\n",
			&Message{
				Prefix:   "user",
				Nick:     "user",
				Command:  "CMD",
				Params:   []string{"#chan", "some message"},
				Trailing: "some message",
			},
		}, {
			":nick!user@host.com CMD a b\r\n",
			&Message{
				Prefix:  "nick!user@host.com",
				Nick:    "nick",
				Command: "CMD",
				Params:  []string{"a", "b"},
			},
		}, {
			"CMD a b :\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"a", "b"},
			},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, parseMessage(tc.input))
	}
}
