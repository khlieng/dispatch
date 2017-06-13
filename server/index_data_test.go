package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTabFromPath(t *testing.T) {
	cases := []struct {
		input           string
		expectedServer  string
		expectedChannel string
	}{
		{
			"/chat.freenode.net/%23r%2Fstuff%2F/",
			"chat.freenode.net",
			"#r/stuff/",
		}, {
			"/chat.freenode.net/%23r%2Fstuff%2F",
			"chat.freenode.net",
			"#r/stuff/",
		}, {
			"/chat.freenode.net/%23r%2Fstuff",
			"chat.freenode.net",
			"#r/stuff",
		}, {
			"/chat.freenode.net/%23stuff",
			"chat.freenode.net",
			"#stuff",
		}, {
			"/chat.freenode.net/%23stuff/cake",
			"",
			"",
		},
	}

	for _, tc := range cases {
		server, channel := getTabFromPath(tc.input)
		assert.Equal(t, tc.expectedServer, server)
		assert.Equal(t, tc.expectedChannel, channel)
	}
}
