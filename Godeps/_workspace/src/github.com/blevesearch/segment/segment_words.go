//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package segment

import (
	"io"
	"unicode"
	"unicode/utf8"
)

// NewWordSegmenter returns a new Segmenter to read from r.
func NewWordSegmenter(r io.Reader) *Segmenter {
	return NewSegmenter(r)
}

// NewWordSegmenterDirect returns a new Segmenter to work directly with buf.
func NewWordSegmenterDirect(buf []byte) *Segmenter {
	return NewSegmenterDirect(buf)
}

const (
	wordCR = iota
	wordLF
	wordNewline
	wordExtend
	wordRegional_Indicator
	wordFormat
	wordKatakana
	wordHebrew_Letter
	wordALetter
	wordSingle_Quote
	wordDouble_Quote
	wordMidNumLet
	wordMidLetter
	wordMidNum
	wordNumeric
	wordExtendNumLet
	wordOther
)

// Word Types
const (
	None = iota
	Number
	Letter
	Kana
	Ideo
)

func SplitWords(data []byte, atEOF bool) (int, []byte, error) {
	advance, token, _, err := SegmentWords(data, atEOF)
	return advance, token, err
}

func SegmentWords(data []byte, atEOF bool) (advance int, token []byte, typ int, err error) {
	prevType := -1
	prevPrevType := -1
	nextType := -1
	immediateNextType := -1
	start := 0
	wordType := None
	currType := -1
	for width := 0; start < len(data); start += width {
		width = 1
		r := rune(data[start])
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRune(data[start:])
		}

		if immediateNextType > 0 {
			currType = immediateNextType
		} else {
			currType = wordSegmentProperty(r)
		}

		hasNext := false
		next := start + width
		nextToken := utf8.RuneError
		for next < len(data) {
			nextWidth := 1
			nextToken = rune(data[next])
			if nextToken >= utf8.RuneSelf {
				nextToken, nextWidth = utf8.DecodeRune(data[next:])
			}
			nextType = wordSegmentProperty(nextToken)
			if !hasNext {
				immediateNextType = nextType
			}
			hasNext = true
			if nextType != wordExtend && nextType != wordFormat {
				break
			}
			next = next + nextWidth
		}

		if start != 0 && in(currType, wordExtend, wordFormat) {
			// wb4
			// dont set prevType, prevPrevType
			// we ignore that these extended are here
			// so types should be whatever we saw before them
			continue
		} else if in(currType, wordALetter, wordHebrew_Letter) &&
			in(prevType, wordALetter, wordHebrew_Letter) {
			// wb5
			wordType = updateWordType(wordType, lookupWordType(currType))
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordMidLetter, wordMidNumLet, wordSingle_Quote) &&
			in(prevType, wordALetter, wordHebrew_Letter) &&
			hasNext && in(nextType, wordALetter, wordHebrew_Letter) {
			// wb6
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordMidLetter, wordMidNumLet, wordSingle_Quote) &&
			in(prevType, wordALetter, wordHebrew_Letter) &&
			!hasNext && !atEOF {
			// possibly wb6, need more data to know
			return 0, nil, 0, nil
		} else if in(currType, wordALetter, wordHebrew_Letter) &&
			in(prevType, wordMidLetter, wordMidNumLet, wordSingle_Quote) &&
			in(prevPrevType, wordALetter, wordHebrew_Letter) {
			// wb7
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordSingle_Quote) &&
			in(prevType, wordHebrew_Letter) {
			// wb7a
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordDouble_Quote) &&
			in(prevType, wordHebrew_Letter) &&
			hasNext && in(nextType, wordHebrew_Letter) {
			// wb7b
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordDouble_Quote) &&
			in(prevType, wordHebrew_Letter) &&
			!hasNext && !atEOF {
			// possibly wb7b, need more data
			return 0, nil, 0, nil
		} else if in(currType, wordHebrew_Letter) &&
			in(prevType, wordDouble_Quote) && in(prevPrevType, wordHebrew_Letter) {
			// wb7c
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordNumeric) &&
			in(prevType, wordNumeric) {
			// wb8
			wordType = updateWordType(wordType, Number)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordNumeric) &&
			in(prevType, wordALetter, wordHebrew_Letter) {
			// wb9
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordALetter, wordHebrew_Letter) &&
			in(prevType, wordNumeric) {
			// wb10
			wordType = updateWordType(wordType, Letter)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordNumeric) &&
			in(prevType, wordMidNum, wordMidNumLet, wordSingle_Quote) &&
			in(prevPrevType, wordNumeric) {
			// wb11
			wordType = updateWordType(wordType, Number)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordMidNum, wordMidNumLet, wordSingle_Quote) &&
			in(prevType, wordNumeric) &&
			hasNext && in(nextType, wordNumeric) {
			// wb12
			wordType = updateWordType(wordType, Number)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordMidNum, wordMidNumLet, wordSingle_Quote) &&
			in(prevType, wordNumeric) &&
			!hasNext && !atEOF {
			// possibly wb12, need more data
			return 0, nil, 0, nil
		} else if in(currType, wordKatakana) &&
			in(prevType, wordKatakana) {
			// wb13
			wordType = updateWordType(wordType, Ideo)
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordExtendNumLet) &&
			in(prevType, wordALetter, wordHebrew_Letter, wordNumeric, wordKatakana, wordExtendNumLet) {
			// wb13a
			wordType = updateWordType(wordType, lookupWordType(currType))
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordALetter, wordHebrew_Letter, wordNumeric, wordKatakana) &&
			in(prevType, wordExtendNumLet) {
			// wb13b
			wordType = updateWordType(wordType, lookupWordType(currType))
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordRegional_Indicator) &&
			in(prevType, wordRegional_Indicator) {
			// wb13c
			prevPrevType = prevType
			prevType = currType
			continue
		} else if start == 0 && in(currType, wordCR) &&
			hasNext && in(immediateNextType, wordLF) {
			prevPrevType = prevType
			prevType = currType
			continue
		} else if start == 0 && !in(currType, wordCR, wordLF, wordNewline) {
			// only first char, keep goin
			wordType = lookupWordType(currType)
			if wordType == None {
				if unicode.In(r, unicode.Katakana, unicode.Hiragana, unicode.Ideographic) {
					wordType = Ideo
				}
			}
			prevPrevType = prevType
			prevType = currType
			continue
		} else if in(currType, wordLF) && in(prevType, wordCR) {
			start += width
			break
		} else {
			// wb14
			if start == 0 {
				start = width
			}
			break
		}
	}
	if start > 0 && atEOF {
		return start, data[:start], wordType, nil
	}

	// Request more data
	return 0, nil, 0, nil

}

func wordSegmentProperty(r rune) int {
	if unicode.Is(_WordALetter, r) {
		return wordALetter
	} else if unicode.Is(_WordCR, r) {
		return wordCR
	} else if unicode.Is(_WordLF, r) {
		return wordLF
	} else if unicode.Is(_WordNewline, r) {
		return wordNewline
	} else if unicode.Is(_WordExtend, r) {
		return wordExtend
	} else if unicode.Is(_WordRegional_Indicator, r) {
		return wordRegional_Indicator
	} else if unicode.Is(_WordFormat, r) {
		return wordFormat
	} else if unicode.Is(_WordKatakana, r) {
		return wordKatakana
	} else if unicode.Is(_WordHebrew_Letter, r) {
		return wordHebrew_Letter
	} else if unicode.Is(_WordSingle_Quote, r) {
		return wordSingle_Quote
	} else if unicode.Is(_WordDouble_Quote, r) {
		return wordDouble_Quote
	} else if unicode.Is(_WordMidNumLet, r) {
		return wordMidNumLet
	} else if unicode.Is(_WordMidLetter, r) {
		return wordMidLetter
	} else if unicode.Is(_WordMidNum, r) {
		return wordMidNum
	} else if unicode.Is(_WordNumeric, r) {
		return wordNumeric
	} else if unicode.Is(_WordExtendNumLet, r) {
		return wordExtendNumLet
	} else {
		return wordOther
	}
}

func lookupWordType(tokenType int) int {
	if tokenType == wordNumeric {
		return Number
	} else if tokenType == wordALetter {
		return Letter
	} else if tokenType == wordHebrew_Letter {
		return Letter
	} else if tokenType == wordKatakana {
		return Ideo
	}

	return None
}

func updateWordType(currentWordType, newWordType int) int {
	if newWordType > currentWordType {
		return newWordType
	}
	return currentWordType
}
