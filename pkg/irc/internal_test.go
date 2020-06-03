package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlePing(t *testing.T) {
	c, out := testClientSend()
	c.handleMessage(&Message{
		Command: "PING",
		Params:  []string{"voi voi"},
	})
	assert.Equal(t, "PONG :voi voi\r\n", <-out)
}

func TestHandleNamreply(t *testing.T) {
	c, _ := testClientSend()

	c.handleMessage(&Message{
		Command: RPL_NAMREPLY,
		Params:  []string{"", "", "#chan", "a b c"},
	})
	c.handleMessage(&Message{
		Command: RPL_NAMREPLY,
		Params:  []string{"", "", "#chan", "d"},
	})

	endMsg := &Message{
		Command: RPL_ENDOFNAMES,
		Params:  []string{"", "#chan"},
	}
	c.handleMessage(endMsg)

	assert.Equal(t, []string{"a", "b", "c", "d"}, endMsg.meta)
}
