package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/mailru/easyjson"
)

func writeJSON(w http.ResponseWriter, r *http.Request, data easyjson.Marshaler) {
	json, err := easyjson.Marshal(data)
	if err != nil {
		log.Println(err)
		fail(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(json) > 1400 && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")

		gzw := getGzipWriter(w)
		gzw.Write(json)
		putGzipWriter(gzw)
	} else {
		w.Write(json)
	}
}
