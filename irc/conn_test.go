package irc

import (
	"testing"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	c.write("test")
	assert.Equal(t, "test\r\n", <-conn.hook)
	c.Write("test")
	assert.Equal(t, "test\r\n", <-conn.hook)
	c.writef("test %d", 2)
	assert.Equal(t, "test 2\r\n", <-conn.hook)
	c.Writef("test %d", 2)
	assert.Equal(t, "test 2\r\n", <-conn.hook)
}
