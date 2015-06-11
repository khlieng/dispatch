package irc

import (
	"net"
	"testing"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

var c *Client
var conn *mockConn

func init() {
	initTestClient()
}

func initTestClient() {
	c = NewClient("test", "testing")
	conn = &mockConn{hook: make(chan string, 1)}
	c.conn = conn
	go c.send()
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
	c.Pass("pass")
	assert.Equal(t, "PASS pass\r\n", <-conn.hook)
}

func TestNick(t *testing.T) {
	c.Nick("test2")
	assert.Equal(t, "test2", c.GetNick())
	assert.Equal(t, "NICK test2\r\n", <-conn.hook)
}

func TestUser(t *testing.T) {
	c.User("user", "rn")
	assert.Equal(t, "USER user 0 * :rn\r\n", <-conn.hook)
}

func TestOper(t *testing.T) {
	c.Oper("name", "pass")
	assert.Equal(t, "OPER name pass\r\n", <-conn.hook)
}

func TestMode(t *testing.T) {
	c.Mode("#chan", "+o", "user")
	assert.Equal(t, "MODE #chan +o user\r\n", <-conn.hook)
}

func TestQuit(t *testing.T) {
	c.connected = true
	c.Quit()
	assert.Equal(t, "QUIT\r\n", <-conn.hook)
	_, ok := <-c.quit
	assert.Equal(t, false, ok)

	initTestClient()
}

func TestJoin(t *testing.T) {
	c.Join("#a")
	assert.Equal(t, "JOIN #a\r\n", <-conn.hook)
	c.Join("#b", "#c")
	assert.Equal(t, "JOIN #b,#c\r\n", <-conn.hook)
}

func TestPart(t *testing.T) {
	c.Part("#a")
	assert.Equal(t, "PART #a\r\n", <-conn.hook)
	c.Part("#b", "#c")
	assert.Equal(t, "PART #b,#c\r\n", <-conn.hook)
}

func TestTopic(t *testing.T) {
	c.Topic("#chan")
	assert.Equal(t, "TOPIC #chan\r\n", <-conn.hook)
}

func TestInvite(t *testing.T) {
	c.Invite("user", "#chan")
	assert.Equal(t, "INVITE user #chan\r\n", <-conn.hook)
}

func TestKick(t *testing.T) {
	c.Kick("#chan", "user")
	assert.Equal(t, "KICK #chan user\r\n", <-conn.hook)
	c.Kick("#chan", "a", "b")
	assert.Equal(t, "KICK #chan a,b\r\n", <-conn.hook)
}

func TestPrivmsg(t *testing.T) {
	c.Privmsg("user", "the message")
	assert.Equal(t, "PRIVMSG user :the message\r\n", <-conn.hook)
}

func TestNotice(t *testing.T) {
	c.Notice("user", "the message")
	assert.Equal(t, "NOTICE user :the message\r\n", <-conn.hook)
}

func TestWhois(t *testing.T) {
	c.Whois("user")
	assert.Equal(t, "WHOIS user\r\n", <-conn.hook)
}

func TestAway(t *testing.T) {
	c.Away("not here")
	assert.Equal(t, "AWAY :not here\r\n", <-conn.hook)
}
