package server

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dsnet/compress/brotli"
	"github.com/khlieng/dispatch/assets"
	"github.com/khlieng/dispatch/pkg/cookie"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

const longCacheControl = "public, max-age=31536000, immutable"
const disabledCacheControl = "no-cache, no-store, must-revalidate"

type File struct {
	Data         []byte
	Length       string
	GzipData     []byte
	GzipLength   string
	Hash         string
	ContentType  string
	CacheControl string
	Compressed   bool
}

type h2PushAsset struct {
	path string
	hash string
}

func newH2PushAsset(name string) h2PushAsset {
	return h2PushAsset{
		path: name,
		hash: strings.Split(name, ".")[1],
	}
}

var (
	files = map[string]*File{}

	indexPage          []byte
	indexPageLen       string
	inlineScriptSha256 string
	serviceWorker      []byte

	h2PushAssets      []h2PushAsset
	h2PushCookieValue string

	contentTypes = map[string]string{
		".js":    "text/javascript",
		".css":   "text/css",
		".woff2": "font/woff2",
		".woff":  "application/font-woff",
		".ttf":   "application/x-font-ttf",
		".png":   "image/png",
		".ico":   "image/x-icon",
		".json":  "application/json",
	}

	robots = []byte("User-agent: *\nDisallow: /")

	hstsHeader string
	cspEnabled bool
)

func (d *Dispatch) initFileServer() {
	cfg := d.Config()

	if cfg.Dev {
		renderIndexPage(indexTemplateData{
			Scripts: []string{"/boot.js", "/main.js"},
		})
	} else {
		bootloader := decompressedAsset(findAssetName("boot*.js"))
		runtime := decompressedAsset(findAssetName("runtime*.js"))

		inlineScript := string(bootloader) + string(runtime)

		hash := sha256.New()
		hash.Write(bootloader)
		hash.Write(runtime)
		inlineScriptSha256 = base64.StdEncoding.EncodeToString(hash.Sum(nil))

		indexStylesheet := "/" + findAssetName("main*.css")
		indexScripts := []string{
			"/" + findAssetName("vendors~main*.js"),
			"/" + findAssetName("main*.js"),
		}

		h2PushAssets = []h2PushAsset{
			newH2PushAsset(indexStylesheet),
			newH2PushAsset(indexScripts[0]),
			newH2PushAsset(indexScripts[1]),
			newH2PushAsset("/" + findAssetName("vendors~connect*.js")),
			newH2PushAsset("/" + findAssetName("connect*.js")),
		}

		for _, asset := range h2PushAssets {
			h2PushCookieValue += asset.hash
		}

		ignoreAssets := []string{
			findAssetName("runtime*.js"),
			findAssetName("boot*.js"),
			"sw.js",
		}

	outer:
		for _, asset := range assets.AssetNames() {
			assetName := strings.TrimSuffix(asset, ".br")

			for _, ignored := range ignoreAssets {
				if ignored == assetName {
					continue outer
				}
			}

			file := &File{
				ContentType:  contentTypes[filepath.Ext(assetName)],
				CacheControl: longCacheControl,
				Compressed:   strings.HasSuffix(asset, ".br"),
			}

			data, err := assets.Asset(asset)
			fatalErr(err)
			file.Data = data
			file.Length = strconv.Itoa(len(data))

			if file.Compressed {
				file.GzipData = gzipAsset(data)
				file.GzipLength = strconv.Itoa(len(file.GzipData))
			}

			files["/"+assetName] = file
		}

		renderIndexPage(indexTemplateData{
			Stylesheet:   indexStylesheet,
			InlineScript: inlineScript,
			Scripts:      indexScripts,
		})

		serviceWorker = decompressedAsset("sw.js")
		hash.Reset()
		hash.Write(indexPage)
		indexHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))

		serviceWorker = bytes.Replace(serviceWorker, []byte("__INDEX_REVISON__"), []byte(indexHash), 1)

		if cfg.HTTPS.HSTS.Enabled && cfg.HTTPS.Enabled {
			hstsHeader = "max-age=" + cfg.HTTPS.HSTS.MaxAge

			if cfg.HTTPS.HSTS.IncludeSubdomains {
				hstsHeader += "; includeSubDomains"
			}
			if cfg.HTTPS.HSTS.Preload {
				hstsHeader += "; preload"
			}
		}

		cspEnabled = true
	}
}

func renderIndexPage(data indexTemplateData) {
	tmpl, err := template.New("").Parse(indexTemplate)
	fatalErr(err)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	buf := &bytes.Buffer{}
	gzw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	fatalErr(err)
	mw := m.Writer("text/html", gzw)

	fatalErr(tmpl.Execute(mw, data))

	fatalErr(mw.Close())
	fatalErr(gzw.Close())

	indexPage = buf.Bytes()
	indexPageLen = strconv.Itoa(len(indexPage))
}

func findAssetName(glob string) string {
	for _, assetName := range assets.AssetNames() {
		assetName = strings.TrimSuffix(assetName, ".br")

		if m, _ := filepath.Match(glob, assetName); m {
			return assetName
		}
	}
	return ""
}

func decompressAsset(data []byte) []byte {
	br, err := brotli.NewReader(bytes.NewReader(data), nil)
	fatalErr(err)

	buf := &bytes.Buffer{}
	io.Copy(buf, br)
	return buf.Bytes()
}

func decompressedAsset(name string) []byte {
	asset, err := assets.Asset(name + ".br")
	fatalErr(err)
	return decompressAsset(asset)
}

func gzipAsset(data []byte) []byte {
	br, err := brotli.NewReader(bytes.NewReader(data), nil)
	fatalErr(err)

	buf := &bytes.Buffer{}
	gzw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	fatalErr(err)

	io.Copy(gzw, br)
	gzw.Close()
	return buf.Bytes()
}

func (d *Dispatch) serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		d.serveIndex(w, r)
	} else if file, ok := files[r.URL.Path]; ok {
		d.serveFile(w, r, file)
	} else if r.URL.Path == "/sw.js" {
		w.Header().Set("Cache-Control", disabledCacheControl)
		w.Header().Set("Content-Type", "text/javascript")
		w.Header().Set("Content-Length", strconv.Itoa(len(serviceWorker)))
		w.Write(serviceWorker)
	} else if r.URL.Path == "/robots.txt" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(robots)))
		w.Write(robots)
	} else {
		d.serveIndex(w, r)
	}
}

func (d *Dispatch) serveIndex(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": r.Header["Accept-Encoding"],
			},
		}

		cookie, err := r.Cookie(cookie.Name(r, "push"))
		if err != nil {
			for _, asset := range h2PushAssets {
				pusher.Push(asset.path, options)
			}

			setPushCookie(w, r)
		} else {
			pushed := false

			i := 0
			for _, asset := range h2PushAssets {
				if len(cookie.Value) >= i+len(asset.hash) &&
					asset.hash != cookie.Value[i:i+len(asset.hash)] {
					pusher.Push(asset.path, options)
					pushed = true
				}
				i += len(asset.hash)
			}

			if pushed {
				setPushCookie(w, r)
			}
		}
	}

	if cspEnabled {
		var wsSrc string
		if r.TLS != nil {
			wsSrc = "wss://" + r.Host
		} else {
			wsSrc = "ws://" + r.Host
		}

		csp := []string{
			"default-src 'none'",
			"script-src 'self' 'sha256-" + inlineScriptSha256 + "'",
			"style-src 'self' 'unsafe-inline'",
			"font-src 'self'",
			"img-src 'self'",
			"manifest-src 'self'",
			"connect-src 'self' " + wsSrc,
			"worker-src 'self'",
		}

		w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", disabledCacheControl)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-XSS-Protection", "0")
	w.Header().Set("Referrer-Policy", "same-origin")

	if hstsHeader != "" {
		w.Header().Set("Strict-Transport-Security", hstsHeader)
	}

	for k, v := range d.Config().Headers {
		w.Header().Set(k, v)
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", indexPageLen)
		w.Write(indexPage)
	} else {
		serveDecompressed(w, indexPage)
	}
}

func setPushCookie(w http.ResponseWriter, r *http.Request) {
	cookie.Set(w, r, &http.Cookie{
		Name:    "push",
		Value:   h2PushCookieValue,
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

func (d *Dispatch) serveFile(w http.ResponseWriter, r *http.Request, file *File) {
	w.Header().Set("Cache-Control", file.CacheControl)
	w.Header().Set("Content-Type", file.ContentType)

	if file.Compressed && strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		w.Header().Set("Content-Encoding", "br")
		w.Header().Set("Content-Length", file.Length)
		w.Write(file.Data)
	} else if file.Compressed && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", file.GzipLength)
		w.Write(file.GzipData)
	} else if !file.Compressed {
		w.Header().Set("Content-Length", file.Length)
		w.Write(file.Data)
	} else {
		serveDecompressed(w, file.GzipData)
	}
}

func serveDecompressed(w http.ResponseWriter, asset []byte) {
	gzr, err := gzip.NewReader(bytes.NewReader(asset))
	buf, err := ioutil.ReadAll(gzr)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	w.Write(buf)
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
