package server

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/khlieng/name_pending/assets"
)

var files = []File{
	File{"bundle.js", "text/javascript"},
	File{"bundle.css", "text/css"},
	File{"font/fontello.woff", "application/font-woff"},
	File{"font/fontello.ttf", "application/x-font-ttf"},
	File{"font/fontello.eot", "application/vnd.ms-fontobject"},
	File{"font/fontello.svg", "image/svg+xml"},
}

type File struct {
	Path        string
	ContentType string
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveFile(w, r, "index.html.gz", "text/html")
		return
	}

	for _, file := range files {
		if strings.HasSuffix(r.URL.Path, file.Path) {
			serveFile(w, r, file.Path+".gz", file.ContentType)
			return
		}
	}

	serveFile(w, r, "index.html.gz", "text/html")
}

func serveFile(w http.ResponseWriter, r *http.Request, path, contentType string) {
	info, err := assets.AssetInfo(path)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if !modifiedSince(w, r, info.ModTime()) {
		return
	}

	data, err := assets.Asset(path)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	} else {
		gzr, err := gzip.NewReader(bytes.NewReader(data))
		buf, err := ioutil.ReadAll(gzr)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
		w.Write(buf)
	}
}

func modifiedSince(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
	t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))

	if err == nil && modtime.Before(t.Add(1*time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return false
	}

	w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	return true
}
