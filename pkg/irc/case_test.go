package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCasefold(t *testing.T) {
	assert.Equal(t, "caላke[^", Casefold(ASCII, "CaላkE[^"))
	assert.Equal(t, "caላke{~", Casefold(RFC1459, "CaላkE[^"))
	assert.Equal(t, "caላke{^", Casefold(RFC1459Strict, "CaላkE[^"))
}

func TestEqualFold(t *testing.T) {
	assert.True(t, EqualFold(ASCII, "caላke[^", "CaላkE[^"))
	assert.False(t, EqualFold(ASCII, "caላke{~", "CaላkE[^"))

	assert.True(t, EqualFold(RFC1459, "caላke{~", "CaላkE[^"))
	assert.False(t, EqualFold(RFC1459, "cላke[^", "CaላkE[^"))

	assert.True(t, EqualFold(RFC1459Strict, "caላke{^", "CaላkE[^"))
	assert.False(t, EqualFold(RFC1459Strict, "caላke[~", "CaላkE[^"))

	assert.True(t, EqualFold(ASCII, "", ""))
	assert.False(t, EqualFold(ASCII, "", " "))
}
