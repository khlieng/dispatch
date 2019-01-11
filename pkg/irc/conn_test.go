package irc

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	ln, err := net.Listen("tcp", "127.0.0.1:45678")
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(testCert, testKey)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	lnTLS, err := tls.Listen("tcp", "127.0.0.1:45679", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	go i.accept(ln)
	go i.accept(lnTLS)
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
	c := testClient()
	c.Connect("127.0.0.1:45678")
	assert.Equal(t, c.Host, "127.0.0.1")
	assert.Equal(t, c.Server, "127.0.0.1:45678")
	waitConnAndClose(t, c)
}

func TestConnectTLS(t *testing.T) {
	c := testClient()
	c.TLS = true
	c.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	c.Connect("127.0.0.1:45679")
	assert.Equal(t, c.Host, "127.0.0.1")
	assert.Equal(t, c.Server, "127.0.0.1:45679")
	waitConnAndClose(t, c)
}

func TestConnectDefaultPorts(t *testing.T) {
	c := testClient()
	c.Connect("127.0.0.1")
	assert.Equal(t, "127.0.0.1:6667", c.Server)

	c = testClient()
	c.TLS = true
	c.Connect("127.0.0.1")
	assert.Equal(t, "127.0.0.1:6697", c.Server)
}

func TestWrite(t *testing.T) {
	c, out := testClientSend()
	c.write("test")
	assert.Equal(t, "test\r\n", <-out)
	c.Write("test")
	assert.Equal(t, "test\r\n", <-out)
	c.writef("test %d", 2)
	assert.Equal(t, "test 2\r\n", <-out)
	c.Writef("test %d", 2)
	assert.Equal(t, "test 2\r\n", <-out)
}

func TestRecv(t *testing.T) {
	c := testClient()
	conn := &mockConn{hook: make(chan string, 16)}
	c.conn = conn

	buf := &bytes.Buffer{}
	buf.WriteString("CMD\r\n")
	buf.WriteString("PING :test\r\n")
	buf.WriteString("001 foo\r\n")
	c.scan = bufio.NewScanner(buf)

	c.sendRecv.Add(1)
	go c.recv()

	assert.Equal(t, "PONG :test\r\n", <-conn.hook)
	assert.Equal(t, &Message{Command: "CMD"}, <-c.Messages)
	assert.Equal(t, &Message{Command: Ping, Params: []string{"test"}}, <-c.Messages)
	assert.Equal(t, &Message{Command: ReplyWelcome, Params: []string{"foo"}}, <-c.Messages)
}

func TestRecvTriggersReconnect(t *testing.T) {
	c := testClient()
	c.conn = &mockConn{}
	c.scan = bufio.NewScanner(bytes.NewBufferString("001 bob\r\n"))
	done := make(chan struct{})
	ok := false
	go func() {
		c.sendRecv.Add(1)
		c.recv()
		_, ok = <-c.reconnect
		close(done)
	}()

	select {
	case <-done:
		assert.False(t, ok)
		return

	case <-time.After(100 * time.Millisecond):
		t.Error("Reconnect not triggered")
	}
}

func TestClose(t *testing.T) {
	c := testClient()
	close(c.quit)
	ok := false
	done := make(chan struct{})
	go func() {
		_, ok = <-c.Messages
		close(done)
	}()

	c.run()

	select {
	case <-done:
		assert.False(t, ok)
		return

	case <-time.After(100 * time.Millisecond):
		t.Error("Channels not closed")
	}
}

func waitConnAndClose(t *testing.T, c *Client) {
	done := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		<-ircd.conn
		quit <- struct{}{}
		<-ircd.connClosed
		close(done)
	}()

	for {
		select {
		case <-done:
			return

		case <-quit:
			assert.True(t, c.Connected())
			c.Quit()

		case <-time.After(500 * time.Millisecond):
			t.Error("Took too long")
			return
		}
	}
}

var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`)

var testKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9
SjY1bIw4iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZB
l2+XsDulrKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQAB
AoGAGRzwwir7XvBOAy5tM/uV6e+Zf6anZzus1s1Y1ClbjbE6HXbnWWF/wbZGOpet
3Zm4vD6MXc7jpTLryzTQIvVdfQbRc6+MUVeLKwZatTXtdZrhu+Jk7hx0nTPy8Jcb
uJqFk541aEw+mMogY/xEcfbWd6IOkp+4xqjlFLBEDytgbIECQQDvH/E6nk+hgN4H
qzzVtxxr397vWrjrIgPbJpQvBsafG7b0dA4AFjwVbFLmQcj2PprIMmPcQrooz8vp
jy4SHEg1AkEA/v13/5M47K9vCxmb8QeD/asydfsgS5TeuNi8DoUBEmiSJwma7FXY
fFUtxuvL7XvjwjN5B30pNEbc6Iuyt7y4MQJBAIt21su4b3sjXNueLKH85Q+phy2U
fQtuUE9txblTu14q3N7gHRZB4ZMhFYyDy8CKrN2cPg/Fvyt0Xlp/DoCzjA0CQQDU
y2ptGsuSmgUtWj3NM9xuwYPm+Z/F84K6+ARYiZ6PYj013sovGKUFfYAqVXVlxtIX
qyUBnu3X9ps8ZfjLZO7BAkEAlT4R5Yl6cGhaJQYZHOde3JEMhNRcVFMO8dJDaFeo
f9Oeos0UUothgiDktdQHxdNEwLjQf7lJJBzV+5OtwswCWA==
-----END RSA PRIVATE KEY-----`)
