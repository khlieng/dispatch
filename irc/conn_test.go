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
	buf.WriteString("001\r\n")
	c.reader = bufio.NewReader(buf)

	c.ready.Add(1)
	c.sendRecv.Add(2)
	go c.send()
	go c.recv()

	assert.Equal(t, "PONG :test\r\n", <-conn.hook)
	assert.Equal(t, &Message{Command: "CMD"}, <-c.Messages)
}

func TestRecvTriggersReconnect(t *testing.T) {
	c := testClient()
	c.conn = &mockConn{}
	c.ready.Add(1)
	c.reader = bufio.NewReader(&bytes.Buffer{})
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
MIIB0zCCAX2gAwIBAgIJAI/M7BYjwB+uMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTIwOTEyMjE1MjAyWhcNMTUwOTEyMjE1MjAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJ
hPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wok/4xIA+ui35/MmNa
rtNuC+BdZ1tMuVCPFZcCAwEAAaNQME4wHQYDVR0OBBYEFJvKs8RfJaXTH08W+SGv
zQyKn0H8MB8GA1UdIwQYMBaAFJvKs8RfJaXTH08W+SGvzQyKn0H8MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADQQBJlffJHybjDGxRMqaRmDhX0+6v02TUKZsW
r5QuVbpQhH6u+0UgcW0jp9QwpxoPTLTWGXEWBBBurxFwiCBhkQ+V
-----END CERTIFICATE-----`)

var testKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBANLJhPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wo
k/4xIA+ui35/MmNartNuC+BdZ1tMuVCPFZcCAwEAAQJAEJ2N+zsR0Xn8/Q6twa4G
6OB1M1WO+k+ztnX/1SvNeWu8D6GImtupLTYgjZcHufykj09jiHmjHx8u8ZZB/o1N
MQIhAPW+eyZo7ay3lMz1V01WVjNKK9QSn1MJlb06h/LuYv9FAiEA25WPedKgVyCW
SmUwbPw8fnTcpqDWE3yTO3vKcebqMSsCIBF3UmVue8YU3jybC3NxuXq3wNm34R8T
xVLHwDXh/6NJAiEAl2oHGGLz64BuAfjKrqwz7qMYr9HCLIe/YsoWq/olzScCIQDi
D2lWusoe2/nEqfDVVWGWlyJ7yOmqaVm/iNUN9B2N2g==
-----END RSA PRIVATE KEY-----`)
