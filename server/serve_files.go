package server

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/khlieng/dispatch/assets"
)

type File struct {
	Path         string
	Asset        string
	ContentType  string
	CacheControl string
	Gzip         bool
}

var (
	files = []*File{
		&File{
			Path:         "bundle.js",
			Asset:        "bundle.js.gz",
			ContentType:  "text/javascript",
			CacheControl: "max-age=31536000",
			Gzip:         true,
		},
		&File{
			Path:         "bundle.css",
			Asset:        "bundle.css.gz",
			ContentType:  "text/css",
			CacheControl: "max-age=31536000",
			Gzip:         true,
		},
	}

	contentTypes = map[string]string{
		".woff2": "font/woff2",
		".woff":  "application/font-woff",
		".ttf":   "application/x-font-ttf",
	}

	hstsHeader string
	cspEnabled bool
)

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

		fonts, err := assets.AssetDir("font")
		if err != nil {
			log.Fatal(err)
		}

		for _, font := range fonts {
			p := strings.TrimSuffix(font, ".gz")

			files = append(files, &File{
				Path:         path.Join("font", p),
				Asset:        path.Join("font", font),
				ContentType:  contentTypes[filepath.Ext(p)],
				CacheControl: "max-age=31536000",
				Gzip:         strings.HasSuffix(font, ".gz"),
			})
		}

		if viper.GetBool("https.hsts.enabled") && viper.GetBool("https.enabled") {
			hstsHeader = "max-age=" + viper.GetString("https.hsts.max_age")

			if viper.GetBool("https.hsts.include_subdomains") {
				hstsHeader += "; includeSubDomains"
			}
			if viper.GetBool("https.hsts.preload") {
				hstsHeader += "; preload"
			}
		}

		cspEnabled = true
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
			serveFile(w, r, file)
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

	if cspEnabled {
		var connectSrc string
		if r.TLS != nil {
			connectSrc = "wss://" + r.Host
		} else {
			connectSrc = "ws://" + r.Host
		}

		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self'; connect-src "+connectSrc)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	if hstsHeader != "" {
		w.Header().Set("Strict-Transport-Security", hstsHeader)
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")

		gzw := gzip.NewWriter(w)
		renderIndex(gzw, getIndexData(r, session))
		gzw.Close()
	} else {
		renderIndex(w, getIndexData(r, session))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, file *File) {
	info, err := assets.AssetInfo(file.Asset)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if !modifiedSince(w, r, info.ModTime()) {
		return
	}

	data, err := assets.Asset(file.Asset)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if file.CacheControl != "" {
		w.Header().Set("Cache-Control", file.CacheControl)
	}

	w.Header().Set("Content-Type", file.ContentType)

	if file.Gzip && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	} else if !file.Gzip {
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
