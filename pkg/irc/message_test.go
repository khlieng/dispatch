package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	cases := []struct {
		input    string
		expected *Message
	}{
		{
			":user CMD #chan :some message",
			&Message{
				Sender:  "user",
				Command: "CMD",
				Params:  []string{"#chan", "some message"},
			},
		}, {
			":nick!user@host.com CMD a b",
			&Message{
				Sender:  "nick",
				Ident:   "user",
				Host:    "host.com",
				Command: "CMD",
				Params:  []string{"a", "b"},
			},
		}, {
			"CMD a b :",
			&Message{
				Command: "CMD",
				Params:  []string{"a", "b", ""},
			},
		}, {
			"CMD a b",
			&Message{
				Command: "CMD",
				Params:  []string{"a", "b"},
			},
		}, {
			"CMD",
			&Message{
				Command: "CMD",
			},
		}, {
			"CMD :tests and stuff",
			&Message{
				Command: "CMD",
				Params:  []string{"tests and stuff"},
			},
		}, {
			":nick@host.com CMD",
			&Message{
				Sender:  "nick",
				Host:    "host.com",
				Command: "CMD",
			},
		}, {
			":ni@ck!user!name@host!.com  CMD",
			&Message{
				Sender:  "ni@ck",
				Ident:   "user!name",
				Host:    "host!.com",
				Command: "CMD",
			},
		}, {
			"CMD #cake pie  ",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie"},
			},
		}, {
			" CMD #cake pie",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie"},
			},
		}, {
			"CMD #cake ::pie",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", ":pie"},
			},
		}, {
			"CMD #cake :  pie",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "  pie"},
			},
		}, {
			"CMD #cake :pie :P <3",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie :P <3"},
			},
		}, {
			"CMD   #cake  :pie!",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie!"},
			},
		}, {
			"@x=y CMD",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Command: "CMD",
			},
		}, {
			"@x=y :nick!user@host.com CMD",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Sender:  "nick",
				Ident:   "user",
				Host:    "host.com",
				Command: "CMD",
			},
		}, {
			"@x=y  :nick!user@host.com    CMD     :pie and cake",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Sender:  "nick",
				Ident:   "user",
				Host:    "host.com",
				Command: "CMD",
				Params:  []string{"pie and cake"},
			},
		}, {
			"@x=y  :nick!user@host.com    CMD beans  rainbows :pie and cake",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Sender:  "nick",
				Ident:   "user",
				Host:    "host.com",
				Command: "CMD",
				Params:  []string{"beans", "rainbows", "pie and cake"},
			},
		},
		{
			"@x=y;a=b CMD",
			&Message{
				Tags: map[string]string{
					"x": "y",
					"a": "b",
				},
				Command: "CMD",
			},
		}, {
			"@x=y;a=\\\\\\:\\s\\r\\n CMD",
			&Message{
				Tags: map[string]string{
					"x": "y",
					"a": "\\; \r\n",
				},
				Command: "CMD",
			},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, ParseMessage(tc.input))
	}
}

func BenchmarkParseMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessage("@x=y  :nick!user@host.com    CMD beans  rainbows :pie and cake")
	}
}

func TestLastParam(t *testing.T) {
	assert.Equal(t, "some message", ParseMessage(":user CMD #chan :some message").LastParam())
	assert.Equal(t, "", ParseMessage("NO_PARAMS").LastParam())
}

func TestBadMessage(t *testing.T) {
	assert.Nil(t, ParseMessage("@"))
	assert.Nil(t, ParseMessage("@ :"))
	assert.Nil(t, ParseMessage("@  :"))
	assert.Nil(t, ParseMessage("@   :"))
	assert.Nil(t, ParseMessage(":user"))
	assert.Nil(t, ParseMessage(":"))
	assert.Nil(t, ParseMessage(""))
}
