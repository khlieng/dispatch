package server

import (
	"compress/gzip"
	"io"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func getGzipWriter(w io.Writer) *gzip.Writer {
	gzw := gzipWriterPool.Get().(*gzip.Writer)
	gzw.Reset(w)
	return gzw
}

func putGzipWriter(gzw *gzip.Writer) {
	gzw.Close()
	gzipWriterPool.Put(gzw)
}
