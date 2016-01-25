package server

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/khlieng/dispatch/assets"
)

var files = []File{
	File{
		Path:         "bundle.js",
		Asset:        "bundle.js.gz",
		ContentType:  "text/javascript",
		CacheControl: "max-age=31536000",
	},
	File{
		Path:         "bundle.css",
		Asset:        "bundle.css.gz",
		ContentType:  "text/css",
		CacheControl: "max-age=31536000",
	},
	File{
		Path:        "font/fontello.woff",
		Asset:       "font/fontello.woff.gz",
		ContentType: "application/font-woff",
	},
	File{
		Path:        "font/fontello.ttf",
		Asset:       "font/fontello.ttf.gz",
		ContentType: "application/x-font-ttf",
	},
	File{
		Path:        "font/fontello.eot",
		Asset:       "font/fontello.eot.gz",
		ContentType: "application/vnd.ms-fontobject",
	},
	File{
		Path:        "font/fontello.svg",
		Asset:       "font/fontello.svg.gz",
		ContentType: "image/svg+xml",
	},
}

type File struct {
	Path         string
	Asset        string
	ContentType  string
	CacheControl string
}

func initFileServer() {
	if !viper.GetBool("dev") {
		data, err := assets.Asset(files[0].Asset)
		if err != nil {
			log.Fatal(err)
		}

		hash := md5.Sum(data)
		files[0].Path = "bundle." + base64.RawURLEncoding.EncodeToString(hash[:]) + ".js"

		data, err = assets.Asset(files[1].Asset)
		if err != nil {
			log.Fatal(err)
		}

		hash = md5.Sum(data)
		files[1].Path = "bundle." + base64.RawURLEncoding.EncodeToString(hash[:]) + ".css"
	}
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveIndex(w, r)
		return
	}

	if strings.HasSuffix(r.URL.Path, "favicon.ico") {
		w.WriteHeader(404)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(r.URL.Path, file.Path) {
			if file.CacheControl != "" {
				w.Header().Set("Cache-Control", file.CacheControl)
			}

			serveFile(w, r, file.Asset, file.ContentType)
			return
		}
	}

	serveIndex(w, r)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	session := handleAuth(w, r)
	if session == nil {
		log.Println("[Auth] No session")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html")

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		renderIndex(gzw, session)
		gzw.Close()
	} else {
		renderIndex(w, session)
	}
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
