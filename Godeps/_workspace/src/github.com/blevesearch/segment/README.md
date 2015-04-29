# segment

A Go library for performing Unicode Text Segmentation
as described in [Unicode Standard Annex #29](http://www.unicode.org/reports/tr29/)

## Features

* Currently only segmentation at Word Boundaries is supported.

## License

Apache License Version 2.0

## Usage

The functionality is exposed in two ways:

1.  You can use a bufio.Scanner with the SplitWords implementation of SplitFunc.
The SplitWords function will identify the appropriate word boundaries in the input
text and the Scanner will return tokens at the appropriate place.

		scanner := bufio.NewScanner(...)
		scanner.Split(segment.SplitWords)
		for scanner.Scan() {
			tokenBytes := scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}

2.  Sometimes you would also like information returned about the type of token.
To do this we have introduce a new type named Segmenter.  It works just like Scanner
but additionally a token type is returned.

		segmenter := segment.NewWordSegmenter(...)
		for segmenter.Segment() {
			tokenBytes := segmenter.Bytes())
			tokenType := segmenter.Type()
		}
		if err := segmenter.Err(); err != nil {
			t.Fatal(err)
		}

## Generating Tables

The tables.go file is generated from the data in the Unicode Text Segmentation property data files.  Also the tables_test.go file is generated from the data in the Unicode Text Segmentation test data files.

To regenerate the files run:

         go generate

 The data generated will be based on the Unicode version set by the unicode package value ```unicode.Version```.

## Status


[![Build Status](https://travis-ci.org/blevesearch/segment.svg?branch=master)](https://travis-ci.org/blevesearch/segment)

[![Coverage Status](https://img.shields.io/coveralls/blevesearch/segment.svg)](https://coveralls.io/r/blevesearch/segment?branch=master)

[![GoDoc](https://godoc.org/github.com/blevesearch/segment?status.svg)](https://godoc.org/github.com/blevesearch/segment)