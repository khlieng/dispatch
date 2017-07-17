package links

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	Client = &http.Client{
		Timeout: 15 * time.Second,
	}

	ErrContentType = errors.New("Unsupported Content-Type")
)

type Meta struct {
	URL         string `json:"URL"`
	SiteName    string `json:"siteName,omitempty"`
	Color       string `json:"color,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageURL,omitempty"`
	VideoURL    string `json:"videoURL,omitempty"`
}

func Fetch(url string) (*Meta, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: Image links
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return nil, ErrContentType
	}

	return ExtractMeta(resp.Body, url)
}

func ExtractMeta(body io.Reader, url string) (*Meta, error) {
	meta := Meta{URL: url}
	var currentNode atom.Atom

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return &meta, nil
			}
			return nil, z.Err()

		case html.TextToken:
			if currentNode == atom.Title && meta.Title == "" {
				meta.Title = string(z.Text())
			}

		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
			name, hasAttr := z.TagName()
			node := atom.Lookup(name)

			if node == atom.Meta && hasAttr {
				var key, val []byte
				var name, content string
				for hasAttr {
					key, val, hasAttr = z.TagAttr()
					switch atom.String(key) {
					case "name":
						name = string(val)

					case "property":
						name = string(val)

					case "content":
						content = string(val)
					}
				}

				if content != "" {
					switch name {
					case "og:site_name":
						meta.SiteName = content

					case "theme-color", "msapplication-TileColor":
						meta.Color = content

					case "og:title", "twitter:title", "title":
						meta.Title = content

					case "og:description", "twitter:description":
						meta.Description = content

					case "description":
						if meta.Description == "" {
							meta.Description = content
						}

					case "og:image", "og:image:secure_url", "twitter:image":
						if !strings.HasPrefix(meta.ImageURL, "https:") {
							meta.ImageURL = content
						}

					case "og:video:url", "og:video:secure_url", "twitter:player":
						if !strings.HasPrefix(meta.VideoURL, "https:") {
							meta.VideoURL = content
						}
					}
				}

				continue
			}

			if tt == html.StartTagToken {
				currentNode = node
			} else {
				currentNode = 0
			}

			if (node == atom.Head && tt == html.EndTagToken) || node == atom.Body {
				return &meta, nil
			}
		}
	}
}
