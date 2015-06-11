package irc

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

var ircd *mockIrcd

func init() {
	initTestServer()
}

func initTestServer() {
	ircd = &mockIrcd{
		conn:       make(chan bool, 1),
		connClosed: make(chan bool, 1),
	}
	ircd.start()
}

type mockIrcd struct {
	conn       chan bool
	connClosed chan bool
}

func (i *mockIrcd) start() {
	ln, err := net.Listen("tcp", ":45678")
	if err != nil {
		log.Fatal(err)
	}
	go i.accept(ln)
}

func (i *mockIrcd) accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go i.handle(conn)
		i.conn <- true
	}
}

func (i *mockIrcd) handle(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			i.connClosed <- true
			return
		}
	}
}

func TestConnect(t *testing.T) {
	c.Connect("127.0.0.1:45678")
	assert.Equal(t, c.Host, "127.0.0.1")
	assert.Equal(t, c.Server, "127.0.0.1:45678")
	waitConnAndClose(t)
	initTestClient()
}

func TestConnectDefaultPorts(t *testing.T) {
	c.Connect("127.0.0.1")
	assert.Equal(t, "127.0.0.1:6667", c.Server)
	initTestClient()

	c.TLS = true
	c.Connect("127.0.0.1")
	assert.Equal(t, "127.0.0.1:6697", c.Server)
	initTestClient()
}

func TestWrite(t *testing.T) {
	c.write("test")
	assert.Equal(t, "test\r\n", <-conn.hook)
	c.Write("test")
	assert.Equal(t, "test\r\n", <-conn.hook)
	c.writef("test %d", 2)
	assert.Equal(t, "test 2\r\n", <-conn.hook)
	c.Writef("test %d", 2)
	assert.Equal(t, "test 2\r\n", <-conn.hook)
}

func TestClose(t *testing.T) {
	defer initTestClient()
	c.close()
	ok := false
	done := make(chan struct{})
	go func() {
		_, ok = <-c.out
		_, ok = <-c.Messages
		close(done)
	}()

	select {
	case <-done:
		assert.False(t, ok)
		return

	case <-time.After(100 * time.Millisecond):
		t.Error("Channels not closed")
	}
}

func waitConnAndClose(t *testing.T) {
	done := make(chan struct{})
	go func() {
		<-ircd.conn
		c.Quit()
		<-ircd.connClosed
		close(done)
	}()

	select {
	case <-done:
		return

	case <-time.After(500 * time.Millisecond):
		t.Error("Took too long")
	}
}
