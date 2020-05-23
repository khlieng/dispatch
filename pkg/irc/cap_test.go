package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCaps(t *testing.T) {
	cases := []struct {
		input    string
		expected map[string][]string
	}{
		{
			"sasl",
			map[string][]string{
				"sasl": nil,
			},
		}, {
			"sasl=PLAIN",
			map[string][]string{
				"sasl": {"PLAIN"},
			},
		}, {
			"cake sasl=PLAIN",
			map[string][]string{
				"cake": nil,
				"sasl": {"PLAIN"},
			},
		}, {
			"cake sasl=PLAIN pie",
			map[string][]string{
				"cake": nil,
				"sasl": {"PLAIN"},
				"pie":  nil,
			},
		}, {
			"cake sasl=PLAIN pie=BLUEBERRY,RASPBERRY",
			map[string][]string{
				"cake": nil,
				"sasl": {"PLAIN"},
				"pie":  {"BLUEBERRY", "RASPBERRY"},
			},
		}, {
			"cake sasl=PLAIN pie=BLUEBERRY,RASPBERRY cheesecake",
			map[string][]string{
				"cake":       nil,
				"sasl":       {"PLAIN"},
				"pie":        {"BLUEBERRY", "RASPBERRY"},
				"cheesecake": nil,
			},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, parseCaps(tc.input))
	}
}
