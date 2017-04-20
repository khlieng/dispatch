package storage

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/kjk/betterguid"
	"github.com/stretchr/testify/assert"
)

func tempdir() string {
	f, _ := ioutil.TempDir("", "")
	return f
}

func TestUser(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	Initialize(tempdir())
	Open()

	srv := Server{
		Name: "Freenode",
		Host: "irc.freenode.net",
		Nick: "test",
	}
	chan1 := Channel{
		Server: srv.Host,
		Name:   "#test",
	}
	chan2 := Channel{
		Server: srv.Host,
		Name:   "#testing",
	}

	user, err := NewUser()
	assert.Nil(t, err)
	user.AddServer(srv)
	user.AddChannel(chan1)
	user.AddChannel(chan2)
	user.closeMessageLog()

	users := LoadUsers()
	assert.Len(t, users, 1)

	user = users[0]
	assert.Equal(t, uint64(1), user.ID)

	servers := user.GetServers()
	assert.Len(t, servers, 1)
	assert.Equal(t, srv, servers[0])

	channels := user.GetChannels()
	assert.Len(t, channels, 2)
	assert.Equal(t, chan1, channels[0])
	assert.Equal(t, chan2, channels[1])

	user.SetNick("bob", srv.Host)
	assert.Equal(t, "bob", user.GetServers()[0].Nick)

	user.RemoveChannel(srv.Host, chan1.Name)
	channels = user.GetChannels()
	assert.Len(t, channels, 1)
	assert.Equal(t, chan2, channels[0])

	user.RemoveServer(srv.Host)
	assert.Len(t, user.GetServers(), 0)
	assert.Len(t, user.GetChannels(), 0)

	user.Remove()
	_, err = os.Stat(Path.User(user.Username))
	assert.True(t, os.IsNotExist(err))

	for _, storedUser := range LoadUsers() {
		assert.NotEqual(t, user.ID, storedUser.ID)
	}
}

func TestMessages(t *testing.T) {
	Initialize(tempdir())
	Open()

	user, err := NewUser()
	assert.Nil(t, err)

	messages, err := user.GetMessages("irc.freenode.net", "#go-nuts", 10, "6")
	assert.Nil(t, err)
	assert.Len(t, messages, 0)

	messages, err = user.GetLastMessages("irc.freenode.net", "#go-nuts", 10)
	assert.Nil(t, err)
	assert.Len(t, messages, 0)

	messages, err = user.SearchMessages("irc.freenode.net", "#go-nuts", "message")
	assert.Nil(t, err)
	assert.Len(t, messages, 0)

	ids := []string{}
	for i := 0; i < 5; i++ {
		id := betterguid.New()
		ids = append(ids, id)
		err = user.LogMessage(id, "irc.freenode.net", "nick", "#go-nuts", "message"+strconv.Itoa(i))
		assert.Nil(t, err)
	}

	messages, err = user.GetMessages("irc.freenode.net", "#go-nuts", 10, ids[4])
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message3", messages[3].Content)
	assert.Nil(t, err)
	assert.Len(t, messages, 4)

	messages, err = user.GetMessages("irc.freenode.net", "#go-nuts", 10, betterguid.New())
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message4", messages[4].Content)
	assert.Nil(t, err)
	assert.Len(t, messages, 5)

	messages, err = user.GetMessages("irc.freenode.net", "#go-nuts", 10, ids[2])
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message1", messages[1].Content)
	assert.Nil(t, err)
	assert.Len(t, messages, 2)

	messages, err = user.GetLastMessages("irc.freenode.net", "#go-nuts", 10)
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message2", messages[2].Content)
	assert.Nil(t, err)
	assert.Len(t, messages, 5)

	messages, err = user.GetLastMessages("irc.freenode.net", "#go-nuts", 3)
	assert.Equal(t, "message2", messages[0].Content)
	assert.Equal(t, "message4", messages[2].Content)
	assert.Nil(t, err)
	assert.Len(t, messages, 3)

	messages, err = user.SearchMessages("irc.freenode.net", "#go-nuts", "message")
	assert.Nil(t, err)
	assert.True(t, len(messages) > 0)

	Close()
}
