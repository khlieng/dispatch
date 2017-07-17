// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package byteutil provides various operations on bytes and byte strings.
package byteutil

var (
	digit     [256]bool
	hexdigit  [256]bool
	letter    [256]bool
	uppercase [256]bool
	lowercase [256]bool
	alphanum  [256]bool
	tolower   [256]byte
	toupper   [256]byte
)

func init() {
	for _, b := range "0123456789" {
		digit[b] = true
	}
	for _, b := range "0123456789abcdefABCDEF" {
		hexdigit[b] = true
	}
	for _, b := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		letter[b] = true
	}
	for _, b := range "abcdefghijklmnopqrstuvwxyz" {
		lowercase[b] = true
	}
	for _, b := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		uppercase[b] = true
	}
	for _, b := range "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		alphanum[b] = true
	}
	for i := 0; i < 256; i++ {
		tolower[i] = byte(i)
		toupper[i] = byte(i)
	}
	for _, b := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		tolower[b] = byte(b) - 'A' + 'a'
	}
	for _, b := range "abcdefghijklmnopqrstuvwxyz" {
		toupper[b] = byte(b) - 'a' + 'A'
	}
}

func IsDigit(b byte) bool {
	return digit[b]
}

func IsHexDigit(b byte) bool {
	return hexdigit[b]
}

func IsLetter(b byte) bool {
	return letter[b]
}

func IsLowercaseLetter(b byte) bool {
	return lowercase[b]
}

func IsUppercaseLetter(b byte) bool {
	return uppercase[b]
}

func IsAlphaNum(b byte) bool {
	return alphanum[b]
}

func ToLower(s string) string {
	if s == "" {
		return ""
	}

	hasUpper := false
	for i := 0; i < len(s); i++ {
		if uppercase[s[i]] {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return s
	}

	buf := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		buf[i] = tolower[s[i]]
	}
	return string(buf)
}

func ToUpper(s string) string {
	if s == "" {
		return ""
	}

	hasLower := false
	for i := 0; i < len(s); i++ {
		if lowercase[s[i]] {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return s
	}

	buf := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		buf[i] = toupper[s[i]]
	}
	return string(buf)
}

func ByteToLower(b byte) byte {
	return tolower[b]
}

func ByteToUpper(b byte) byte {
	return toupper[b]
}

func IndexAny(s, chars string) int {
	var t [256]bool
	for i := 0; i < len(chars); i++ {
		t[chars[i]] = true
	}
	for i := 0; i < len(s); i++ {
		if t[s[i]] {
			return i
		}
	}
	return -1
}

func IndexAnyTable(s string, t *[256]bool) int {
	for i := 0; i < len(s); i++ {
		if t[s[i]] {
			return i
		}
	}
	return -1
}

func Unhex(d byte) byte {
	switch {
	case digit[d]:
		return d - '0'
	case uppercase[d]:
		return d - 'A' + 10
	case lowercase[d]:
		return d - 'a' + 10
	}
	panic("unhex: not hex digit")
}
