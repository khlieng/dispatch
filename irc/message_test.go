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
			":user CMD #chan :some message\r\n",
			&Message{
				Prefix:  "user",
				Nick:    "user",
				Command: "CMD",
				Params:  []string{"#chan", "some message"},
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
				Params:  []string{"a", "b", ""},
			},
		}, {
			"CMD a b\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"a", "b"},
			},
		}, {
			"CMD\r\n",
			&Message{
				Command: "CMD",
			},
		}, {
			"CMD :tests and stuff\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"tests and stuff"},
			},
		}, {
			":nick@host.com CMD\r\n",
			&Message{
				Prefix:  "nick@host.com",
				Nick:    "nick",
				Command: "CMD",
			},
		}, {
			":ni@ck!user!name@host!.com CMD\r\n",
			&Message{
				Prefix:  "ni@ck!user!name@host!.com",
				Nick:    "ni@ck",
				Command: "CMD",
			},
		}, {
			"CMD #cake pie  \r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie"},
			},
		}, {
			" CMD #cake pie\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie"},
			},
		}, {
			"CMD #cake ::pie\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", ":pie"},
			},
		}, {
			"CMD #cake :  pie\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "  pie"},
			},
		}, {
			"CMD #cake :pie :P <3\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie :P <3"},
			},
		}, {
			"CMD   #cake  :pie!\r\n",
			&Message{
				Command: "CMD",
				Params:  []string{"#cake", "pie!"},
			},
		}, {
			"@x=y CMD\r\n",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Command: "CMD",
			},
		}, {
			"@x=y :nick!user@host.com CMD\r\n",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Prefix:  "nick!user@host.com",
				Nick:    "nick",
				Command: "CMD",
			},
		}, {
			"@x=y :nick!user@host.com CMD :pie and cake\r\n",
			&Message{
				Tags: map[string]string{
					"x": "y",
				},
				Prefix:  "nick!user@host.com",
				Nick:    "nick",
				Command: "CMD",
				Params:  []string{"pie and cake"},
			},
		}, {
			"@x=y;a=b CMD\r\n",
			&Message{
				Tags: map[string]string{
					"x": "y",
					"a": "b",
				},
				Command: "CMD",
			},
		}, {
			"@x=y;a=\\\\\\:\\s\\r\\n CMD\r\n",
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
		assert.Equal(t, tc.expected, parseMessage(tc.input))
	}
}

func TestLastParam(t *testing.T) {
	assert.Equal(t, "some message", parseMessage(":user CMD #chan :some message\r\n").LastParam())
	assert.Equal(t, "", parseMessage("NO_PARAMS").LastParam())
}

func TestBadMessagePanic(t *testing.T) {
	parseMessage("@\r\n")
	parseMessage("@ :\r\n")
	parseMessage("@  :\r\n")
	parseMessage(":user\r\n")
	parseMessage(":\r\n")
	parseMessage(":")
	parseMessage("")
}

func TestParseISupport(t *testing.T) {
	s := newISupport()
	s.parse([]string{"bob", "CAKE=31", "PIE", ":durr"})
	assert.Equal(t, 31, s.GetInt("CAKE"))
	assert.Equal(t, "31", s.Get("CAKE"))
	assert.True(t, s.Has("CAKE"))
	assert.True(t, s.Has("PIE"))
	assert.False(t, s.Has("APPLES"))
	assert.Equal(t, "", s.Get("APPLES"))
	assert.Equal(t, 0, s.GetInt("APPLES"))

	s.parse([]string{"bob", "-PIE", ":hurr"})
	assert.False(t, s.Has("PIE"))

	s.parse([]string{"bob", "CAKE=1337", ":durr"})
	assert.Equal(t, 1337, s.GetInt("CAKE"))

	s.parse([]string{"bob", "CAKE=", ":durr"})
	assert.Equal(t, "", s.Get("CAKE"))
	assert.True(t, s.Has("CAKE"))

	s.parse([]string{"bob", "CAKE===", ":durr"})
	assert.Equal(t, "==", s.Get("CAKE"))

	s.parse([]string{"bob", "-CAKE=31", ":durr"})
	assert.False(t, s.Has("CAKE"))
}
