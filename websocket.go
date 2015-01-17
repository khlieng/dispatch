package main

import (
	"golang.org/x/net/websocket"
)

type WebSocket struct {
	conn *websocket.Conn

	In chan []byte
}

func NewWebSocket(ws *websocket.Conn) *WebSocket {
	return &WebSocket{
		conn: ws,
		In:   make(chan []byte, 32),
	}
}

func (w *WebSocket) write() {
	for data := range w.In {
		w.conn.Write(data)
	}
}
