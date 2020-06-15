package server

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/khlieng/dispatch/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetTabFromPath(t *testing.T) {
	cases := []struct {
		input       *http.Request
		expectedTab Tab
	}{
		{
			&http.Request{
				URL:    &url.URL{Path: "/init"},
				Header: http.Header{"Referer": []string{"/chat.freenode.net/%23r%2Fstuff%2F"}},
			},
			Tab{storage.Tab{Network: "chat.freenode.net", Name: "#r/stuff/"}},
		}, {
			&http.Request{
				URL:    &url.URL{Path: "/init"},
				Header: http.Header{"Referer": []string{"/chat.freenode.net/%23r%2Fstuff"}},
			},
			Tab{storage.Tab{Network: "chat.freenode.net", Name: "#r/stuff"}},
		}, {
			&http.Request{
				URL:    &url.URL{Path: "/init"},
				Header: http.Header{"Referer": []string{"/chat.freenode.net/%23stuff"}},
			},
			Tab{storage.Tab{Network: "chat.freenode.net", Name: "#stuff"}},
		}, {
			&http.Request{
				URL:    &url.URL{Path: "/init"},
				Header: http.Header{"Referer": []string{"/chat.freenode.net/stuff"}},
			},
			Tab{storage.Tab{Network: "chat.freenode.net", Name: "stuff"}},
		}, {
			&http.Request{
				URL:    &url.URL{Path: "/init"},
				Header: http.Header{"Referer": []string{"/data/chat.freenode.net/%23apples"}},
			},
			Tab{},
		}, {
			&http.Request{
				URL: &url.URL{Path: "/ws/chat.freenode.net"},
			},
			Tab{storage.Tab{Network: "chat.freenode.net"}},
		},
	}

	for _, tc := range cases {
		tab, err := tabFromRequest(tc.input)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedTab, tab)
	}
}
