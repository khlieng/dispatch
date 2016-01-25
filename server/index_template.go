package server

import (
	"encoding/json"
	"io"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"
)

var (
	index_start = []byte(`<!DOCTYPE html><html lang=en><head><meta charset=UTF-8><meta name=viewport content="width=device-width,initial-scale=1"><title>Dispatch</title><link href="https://fonts.googleapis.com/css?family=Montserrat:400,700|Roboto+Mono:400,700" rel=stylesheet><link href=/bundle.css rel=stylesheet></head><body><div id=root></div><script>window.__ENV__=`)
	index_end   = []byte(`;</script><script src=/bundle.js></script></body></html>`)
)

type connectDefaults struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Channels []string `json:"channels"`
	Password string   `json:"password"`
	SSL      bool     `json:"ssl"`
}

type indexData struct {
	Defaults connectDefaults `json:"defaults"`
}

func renderIndex(w io.Writer, session *Session) {
	w.Write(index_start)

	json.NewEncoder(w).Encode(indexData{
		Defaults: connectDefaults{
			Name:     viper.GetString("defaults.name"),
			Address:  viper.GetString("defaults.address"),
			Channels: viper.GetStringSlice("defaults.channels"),
			Password: viper.GetString("defaults.password"),
			SSL:      viper.GetBool("defaults.ssl"),
		},
	})

	w.Write(index_end)
}
