package server

import (
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"
)

type WebSocket struct {
	conn *websocket.Conn

	Out chan []byte
}

func NewWebSocket(ws *websocket.Conn) *WebSocket {
	return &WebSocket{
		conn: ws,
		Out:  make(chan []byte, 32),
	}
}

func (w *WebSocket) write() {
	for {
		err := w.conn.WriteMessage(websocket.TextMessage, <-w.Out)
		if err != nil {
			return
		}
	}
}

func (w *WebSocket) close() {
	close(w.Out)
}
