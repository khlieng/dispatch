package server

import (
	"time"

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
	var err error
	ping := time.Tick(20 * time.Second)

	for {
		select {
		case msg, ok := <-w.Out:
			if !ok {
				return
			}

			err = w.conn.WriteMessage(websocket.TextMessage, msg)

		case <-ping:
			err = w.conn.WriteJSON(WSResponse{Type: "ping"})
		}

		if err != nil {
			return
		}
	}
}

func (w *WebSocket) close() {
	close(w.Out)
}
