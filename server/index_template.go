package server

import (
	"encoding/json"
	"io"
)

var (
	index_0 = []byte(`<!DOCTYPE html><html lang=en><head><meta charset=UTF-8><meta name=viewport content="width=device-width,initial-scale=1"><title>Dispatch</title><link href=/`)
	index_1 = []byte(` rel=stylesheet></head><body><div id=root></div><script>window.__ENV__=`)
	index_2 = []byte(`</script><script src=/`)
	index_3 = []byte(`></script></body></html>`)
)

func renderIndex(w io.Writer, data interface{}) {
	w.Write(index_0)
	w.Write([]byte(files[1].Path))
	w.Write(index_1)

	json.NewEncoder(w).Encode(data)

	w.Write(index_2)
	w.Write([]byte(files[0].Path))
	w.Write(index_3)
}
