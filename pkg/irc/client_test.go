package irc

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testClientSend() (*Client, chan string) {
	c := NewClient(Config{})
	conn := &mockConn{hook: make(chan string, 16)}
	c.conn = conn
	c.sendRecv.Add(1)
	go c.send()
	return c, conn.hook
}

type mockConn struct {
	hook chan string
	net.Conn
}

func (c *mockConn) Write(b []byte) (int, error) {
	c.hook <- string(b)
	return len(b), nil
}

func (c *mockConn) Close() error {
	return nil
}

func TestPass(t *testing.T) {
	c, out := testClientSend()
	c.writePass("pass")
	assert.Equal(t, "PASS pass\r\n", <-out)
}

func TestNick(t *testing.T) {
	c, out := testClientSend()
	c.Nick("test2")
	assert.Equal(t, "NICK test2\r\n", <-out)

	c.writeNick("nick")
	assert.Equal(t, "NICK nick\r\n", <-out)
}

func TestUser(t *testing.T) {
	c, out := testClientSend()
	c.writeUser("user", "rn")
	assert.Equal(t, "USER user 0 * :rn\r\n", <-out)
}

func TestOper(t *testing.T) {
	c, out := testClientSend()
	c.Oper("name", "pass")
	assert.Equal(t, "OPER name pass\r\n", <-out)
}

func TestMode(t *testing.T) {
	c, out := testClientSend()
	c.Mode("#chan", "+o", "user")
	assert.Equal(t, "MODE #chan +o user\r\n", <-out)
}

func TestQuit(t *testing.T) {
	c, out := testClientSend()
	c.connected = true
	c.Quit()
	assert.Equal(t, "QUIT\r\n", <-out)
	_, ok := <-c.quit
	assert.Equal(t, false, ok)
}

func TestJoin(t *testing.T) {
	c, out := testClientSend()
	c.Join("#a")
	assert.Equal(t, "JOIN #a\r\n", <-out)
	c.Join("#b", "#c")
	assert.Equal(t, "JOIN #b,#c\r\n", <-out)
}

func TestPart(t *testing.T) {
	c, out := testClientSend()
	c.Part("#a")
	assert.Equal(t, "PART #a\r\n", <-out)
	c.Part("#b", "#c")
	assert.Equal(t, "PART #b,#c\r\n", <-out)
}

func TestTopic(t *testing.T) {
	c, out := testClientSend()
	c.Topic("#chan")
	assert.Equal(t, "TOPIC #chan\r\n", <-out)
	c.Topic("#chan", "apple pie")
	assert.Equal(t, "TOPIC #chan :apple pie\r\n", <-out)
	c.Topic("#chan", "")
	assert.Equal(t, "TOPIC #chan :\r\n", <-out)
}

func TestInvite(t *testing.T) {
	c, out := testClientSend()
	c.Invite("user", "#chan")
	assert.Equal(t, "INVITE user #chan\r\n", <-out)
}

func TestKick(t *testing.T) {
	c, out := testClientSend()
	c.Kick("#chan", "user")
	assert.Equal(t, "KICK #chan user\r\n", <-out)
	c.Kick("#chan", "a", "b")
	assert.Equal(t, "KICK #chan a,b\r\n", <-out)
}

func TestPrivmsg(t *testing.T) {
	c, out := testClientSend()
	c.Privmsg("user", "the message")
	assert.Equal(t, "PRIVMSG user :the message\r\n", <-out)
}

func TestNotice(t *testing.T) {
	c, out := testClientSend()
	c.Notice("user", "the message")
	assert.Equal(t, "NOTICE user :the message\r\n", <-out)
}

func TestReplyCTCP(t *testing.T) {
	c, out := testClientSend()
	c.ReplyCTCP("user", "PING", "PONG")
	assert.Equal(t, "NOTICE user :\x01PING PONG\x01\r\n", <-out)
}

func TestWhois(t *testing.T) {
	c, out := testClientSend()
	c.Whois("user")
	assert.Equal(t, "WHOIS user\r\n", <-out)
}

func TestAway(t *testing.T) {
	c, out := testClientSend()
	c.Away("not here")
	assert.Equal(t, "AWAY :not here\r\n", <-out)
}

func TestRegister(t *testing.T) {
	c, out := testClientSend()
	c.Config.Nick = "nick"
	c.Config.Username = "user"
	c.Config.Realname = "rn"
	t.Log(c.Config)
	c.register()
	assert.Equal(t, "CAP LS 302\r\n", <-out)
	assert.Equal(t, "NICK nick\r\n", <-out)
	assert.Equal(t, "USER user 0 * :rn\r\n", <-out)

	c.Config.Password = "pass"
	c.register()
	assert.Equal(t, "CAP LS 302\r\n", <-out)
	assert.Equal(t, "PASS pass\r\n", <-out)
	assert.Equal(t, "NICK nick\r\n", <-out)
	assert.Equal(t, "USER user 0 * :rn\r\n", <-out)
}

func TestFlushChannels(t *testing.T) {
	c, out := testClientSend()
	c.addChannel("#chan1")
	c.flushChannels()
	assert.Equal(t, <-out, "JOIN #chan1\r\n")
	c.addChannel("#chan2")
	c.addChannel("#chan4")
	c.removeChannels("#chan4")
	c.addChannel("#chan3")
	c.flushChannels()
	assert.Equal(t, <-out, "JOIN #chan2,#chan3\r\n")
}
