package main

import (
	"github.com/khlieng/name_pending/Godeps/_workspace/src/golang.org/x/net/websocket"
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
	for data := range w.Out {
		w.conn.Write(data)
	}
}
