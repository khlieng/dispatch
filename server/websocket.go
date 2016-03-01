package server

import (
	"time"

	"github.com/gorilla/websocket"
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

			err = c.conn.WriteJSON(res)

		case <-ping:
			err = c.conn.WriteJSON(WSResponse{Type: "ping"})
		}

		if err != nil {
			return
		}
	}
}

func (c *wsConn) recv() {
	var req WSRequest

	for {
		err := c.conn.ReadJSON(&req)
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
