package server

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
)

type wsConn struct {
	conn *websocket.Conn
	in   chan WSRequest
	out  chan WSResponse
}

func newWSConn(conn *websocket.Conn) *wsConn {
	return &wsConn{
		conn: conn,
		in:   make(chan WSRequest, 32),
		out:  make(chan WSResponse, 32),
	}
}

func (c *wsConn) send() {
	var err error
	ping := time.Tick(20 * time.Second)

	for {
		select {
		case res, ok := <-c.out:
			if !ok {
				return
			}

			err = c.writeJSON(res)

		case <-ping:
			err = c.writeJSON(WSResponse{Type: "ping"})
		}

		if err != nil {
			return
		}
	}
}

func (c *wsConn) recv() {
	var req WSRequest

	for {
		err := c.readJSON(&req)
		if err != nil {
			close(c.in)
			return
		}

		c.in <- req
	}
}

func (c *wsConn) close() {
	close(c.out)
	c.conn.Close()
}

func (c *wsConn) readJSON(v easyjson.Unmarshaler) error {
	_, r, err := c.conn.NextReader()
	if err != nil {
		return err
	}

	return easyjson.UnmarshalFromReader(r, v)
}

func (c *wsConn) writeJSON(v easyjson.Marshaler) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err1 := easyjson.MarshalToWriter(v, w)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
