// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linkify

import "unicode"

var (
	trigger    [256]bool
	unreserved [256]bool
	subdelims  [256]bool
	emailcs    [256]bool
	basicPunct [256]bool
)

func init() {
	for _, b := range "-._~" {
		unreserved[b] = true
	}
	for _, b := range "!$&'()*+,;=" {
		subdelims[b] = true
	}
	for _, b := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+/=?^_`{|}~-" {
		emailcs[b] = true
	}
	for _, b := range ".,?!;:" {
		basicPunct[b] = true
	}
}

func isAllowedInEmail(r rune) bool {
	return r < 0x7f && emailcs[r]
}

func isLetterOrDigit(r rune) bool {
	return unicode.In(r, unicode.Letter, unicode.Digit)
}

func isPunctOrSpaceOrControl(r rune) bool {
	return r == '<' || r == '>' || unicode.In(r, unicode.Punct, unicode.Space, unicode.Cc)
}

func isUnreserved(r rune) bool {
	return (r < 0x7f && unreserved[r]) || isLetterOrDigit(r)
}

func isSubDelimiter(r rune) bool {
	return r < 0x7f && subdelims[r]
}
