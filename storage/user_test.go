package storage

import (
	"io/ioutil"
	"testing"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func tempdir() string {
	f, _ := ioutil.TempDir("", "")
	return f
}

func TestUser(t *testing.T) {
	Initialize(tempdir())
	Open()

	srv := Server{
		Name:    "Freenode",
		Address: "irc.freenode.net",
		Nick:    "test",
	}
	chan1 := Channel{
		Server: srv.Address,
		Name:   "#test",
	}
	chan2 := Channel{
		Server: srv.Address,
		Name:   "#testing",
	}

	user := NewUser("unique")
	user.AddServer(srv)
	user.AddChannel(chan1)
	user.AddChannel(chan2)
	user.Close()

	users := LoadUsers()
	assert.Len(t, users, 1)

	user = users[0]
	assert.Equal(t, "unique", user.UUID)

	servers := user.GetServers()
	assert.Len(t, servers, 1)
	assert.Equal(t, srv, servers[0])

	channels := user.GetChannels()
	assert.Len(t, channels, 2)
	assert.Equal(t, chan1, channels[0])
	assert.Equal(t, chan2, channels[1])

	user.SetNick("bob", srv.Address)
	assert.Equal(t, "bob", user.GetServers()[0].Nick)

	user.RemoveChannel(srv.Address, chan1.Name)
	channels = user.GetChannels()
	assert.Len(t, channels, 1)
	assert.Equal(t, chan2, channels[0])

	user.RemoveServer(srv.Address)
	assert.Len(t, user.GetServers(), 0)
	assert.Len(t, user.GetChannels(), 0)
}
