/*
The MIT License (MIT)

Copyright (c) 2014 Henry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package goSam

import (
	"io"
	//"github.com/miolini/datacounter"
)

type RWC struct {
	io.Reader
	io.Writer
	c io.Closer
}

func WrapRWC(c io.ReadWriteCloser) io.ReadWriteCloser {
	rl := NewReadLogger("<", c)
	wl := NewWriteLogger(">", c)

	return &RWC{
		Reader: rl,
		Writer: wl,
		c:      c,
	}
}

func (c *RWC) Close() error {
	return c.c.Close()
}

/*
type Counter struct {
	io.Reader
	io.Writer
	c io.Closer

	Cr *datacounter.ReaderCounter
	Cw *datacounter.WriterCounter
}

func WrapCounter(c io.ReadWriteCloser) *Counter {
	rc := datacounter.NewReaderCounter(c)
	wc := datacounter.NewWriterCounter(c)

	return &Counter{
		Reader: rc,
		Writer: wc,
		c:      c,

		Cr: rc,
		Cw: wc,
	}
}

func (c *Counter) Close() error {
	return c.c.Close()
}
*/
