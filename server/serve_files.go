package server

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dsnet/compress/brotli"
	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/assets"
)

const longCacheControl = "public, max-age=31536000, immutable"
const disabledCacheControl = "no-cache, no-store, must-revalidate"

type File struct {
	Path         string
	Asset        string
	GzipAsset    []byte
	Hash         string
	ContentType  string
	CacheControl string
	Compressed   bool
}

var (
	files = []*File{
		&File{
			Path:         "bundle.js",
			Asset:        "bundle.js.br",
			ContentType:  "text/javascript",
			CacheControl: longCacheControl,
			Compressed:   true,
		},
		&File{
			Path:         "bundle.css",
			Asset:        "bundle.css.br",
			ContentType:  "text/css",
			CacheControl: longCacheControl,
			Compressed:   true,
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

func (d *Dispatch) initFileServer() {
	if !viper.GetBool("dev") {
		data, err := assets.Asset(files[0].Asset)
		if err != nil {
			log.Fatal(err)
		}

		hash := md5.Sum(data)
		files[0].Hash = base64.RawURLEncoding.EncodeToString(hash[:])[:8]
		files[0].Path = "bundle." + files[0].Hash + ".js"

		br, err := brotli.NewReader(bytes.NewReader(data), nil)
		if err != nil {
			log.Fatal(err)
		}

		buf := &bytes.Buffer{}
		gzw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(gzw, br)
		gzw.Close()
		files[0].GzipAsset = buf.Bytes()

		data, err = assets.Asset(files[1].Asset)
		if err != nil {
			log.Fatal(err)
		}

		hash = md5.Sum(data)
		files[1].Hash = base64.RawURLEncoding.EncodeToString(hash[:])[:8]
		files[1].Path = "bundle." + files[1].Hash + ".css"

		br.Reset(bytes.NewReader(data))
		buf = &bytes.Buffer{}
		gzw.Reset(buf)

		io.Copy(gzw, br)
		gzw.Close()
		files[1].GzipAsset = buf.Bytes()

		fonts, err := assets.AssetDir("font")
		if err != nil {
			log.Fatal(err)
		}

		for _, font := range fonts {
			p := strings.TrimSuffix(font, ".br")

			file := &File{
				Path:         path.Join("font", p),
				Asset:        path.Join("font", font),
				ContentType:  contentTypes[filepath.Ext(p)],
				CacheControl: longCacheControl,
				Compressed:   strings.HasSuffix(font, ".br"),
			}

			if file.Compressed {
				data, err = assets.Asset(file.Asset)
				if err != nil {
					log.Fatal(err)
				}

				br.Reset(bytes.NewReader(data))
				buf = &bytes.Buffer{}
				gzw.Reset(buf)

				io.Copy(gzw, br)
				gzw.Close()
				file.GzipAsset = buf.Bytes()
			}

			files = append(files, file)
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

func (d *Dispatch) serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		d.serveIndex(w, r)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(r.URL.Path, file.Path) {
			d.serveFile(w, r, file)
			return
		}
	}

	d.serveIndex(w, r)
}

func (d *Dispatch) serveIndex(w http.ResponseWriter, r *http.Request) {
	state := d.handleAuth(w, r, false)

	if cspEnabled {
		var connectSrc string
		if r.TLS != nil {
			connectSrc = "wss://" + r.Host
		} else {
			connectSrc = "ws://" + r.Host
		}

		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self'; img-src data:; connect-src "+connectSrc)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", disabledCacheControl)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	if hstsHeader != "" {
		w.Header().Set("Strict-Transport-Security", hstsHeader)
	}

	if pusher, ok := w.(http.Pusher); ok {
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": r.Header["Accept-Encoding"],
			},
		}

		cookie, err := r.Cookie("push")
		if err != nil {
			pusher.Push("/"+files[1].Path, options)
			pusher.Push("/"+files[0].Path, options)
			setPushCookie(w, r)
		} else {
			pushed := false

			if files[1].Hash != cookie.Value[8:] {
				pusher.Push("/"+files[1].Path, options)
				pushed = true
			}
			if files[0].Hash != cookie.Value[:8] {
				pusher.Push("/"+files[0].Path, options)
				pushed = true
			}

			if pushed {
				setPushCookie(w, r)
			}
		}
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")

		gzw := gzip.NewWriter(w)
		IndexTemplate(gzw, getIndexData(r, state), files[1].Path, files[0].Path)
		gzw.Close()
	} else {
		IndexTemplate(w, getIndexData(r, state), files[1].Path, files[0].Path)
	}
}

func setPushCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "push",
		Value:    files[0].Hash + files[1].Hash,
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})
}

func (d *Dispatch) serveFile(w http.ResponseWriter, r *http.Request, file *File) {
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

	if file.Compressed && strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		w.Header().Set("Content-Encoding", "br")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	} else if file.Compressed && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", strconv.Itoa(len(file.GzipAsset)))
		w.Write(file.GzipAsset)
	} else if !file.Compressed {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	} else {
		gzr, err := gzip.NewReader(bytes.NewReader(file.GzipAsset))
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
