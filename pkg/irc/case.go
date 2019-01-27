package irc

import (
	"unicode/utf8"
)

const (
	// ASCII maps a-z as the lower case of A-Z
	ASCII = "ascii"
	// RFC1459 maps a-z and {, |, }, ~ as the lower case of A-Z and [, \, ], ^
	RFC1459 = "rfc1459"
	// RFC1459Strict maps a-z and {, |, } as the lower case of A-Z and [, \, ]
	RFC1459Strict = "strict-rfc1459"
)

func (c *Client) Casefold(s string) string {
	mapping := c.Features.String("CASEMAPPING")
	if mapping == "" {
		mapping = RFC1459
	}
	return Casefold(mapping, s)
}

func (c *Client) EqualFold(s1, s2 string) bool {
	mapping := c.Features.String("CASEMAPPING")
	if mapping == "" {
		mapping = RFC1459
	}
	return EqualFold(mapping, s1, s2)
}

func Casefold(mapping, s string) string {
	switch mapping {
	case ASCII:
		return toLower(s, 'Z')
	case RFC1459:
		return toLower(s, '^')
	case RFC1459Strict:
		return toLower(s, ']')
	}

	return s
}

func EqualFold(mapping, s1, s2 string) bool {
	switch mapping {
	case ASCII:
		return equalFold(s1, s2, 'Z')
	case RFC1459:
		return equalFold(s1, s2, '^')
	case RFC1459Strict:
		return equalFold(s1, s2, ']')
	}

	return s1 == s2
}

func toLower(s string, end byte) string {
	hasUpper := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if hasUpper = 'A' <= c && c <= end; hasUpper {
			break
		}
	}

	if !hasUpper {
		return s
	}

	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]

		// Skip Unicode characters
		if c >= utf8.RuneSelf {
			_, size := utf8.DecodeRuneInString(s[i:])
			for cEnd := i + size; i < cEnd; i++ {
				b[i] = s[i]
			}
			i--
			continue
		}

		if 'A' <= c && c <= end {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

func equalFold(s1, s2 string, end rune) bool {
	for s1 != "" && s2 != "" {
		var r1, r2 rune
		if s1[0] < utf8.RuneSelf {
			r1, s1 = rune(s1[0]), s1[1:]
		} else {
			r, size := utf8.DecodeRuneInString(s1)
			r1, s1 = r, s1[size:]
		}
		if s2[0] < utf8.RuneSelf {
			r2, s2 = rune(s2[0]), s2[1:]
		} else {
			r, size := utf8.DecodeRuneInString(s2)
			r2, s2 = r, s2[size:]
		}

		if r1 == r2 {
			continue
		}

		if r2 < r1 {
			r2, r1 = r1, r2
		}

		if 'A' <= r1 && r1 <= end && r2 == r1+32 {
			continue
		}

		return false
	}

	return s1 == s2
}
