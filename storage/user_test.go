package storage_test

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/khlieng/dispatch/storage"
	"github.com/khlieng/dispatch/storage/bleve"
	"github.com/khlieng/dispatch/storage/boltdb"
	"github.com/kjk/betterguid"
	"github.com/stretchr/testify/assert"
)

func tempdir() string {
	f, _ := ioutil.TempDir("", "")
	return f
}

func TestUser(t *testing.T) {
	storage.Initialize(tempdir())

	db, err := boltdb.New(storage.Path.Database())
	assert.Nil(t, err)

	user, err := storage.NewUser(db)
	assert.Nil(t, err)

	srv := storage.Server{
		Name: "Freenode",
		Host: "irc.freenode.net",
		Nick: "test",
	}
	chan1 := storage.Channel{
		Server: srv.Host,
		Name:   "#test",
	}
	chan2 := storage.Channel{
		Server: srv.Host,
		Name:   "#testing",
	}

	user.AddServer(&srv)
	user.AddChannel(&chan1)
	user.AddChannel(&chan2)

	users, err := storage.LoadUsers(db)
	assert.Nil(t, err)
	assert.Len(t, users, 1)

	user = &users[0]
	assert.Equal(t, uint64(1), user.ID)

	servers, err := user.GetServers()
	assert.Len(t, servers, 1)
	assert.Equal(t, srv, servers[0])

	channels, err := user.GetChannels()
	assert.Len(t, channels, 2)
	assert.Equal(t, chan1, channels[0])
	assert.Equal(t, chan2, channels[1])

	user.SetNick("bob", srv.Host)
	servers, err = user.GetServers()
	assert.Equal(t, "bob", servers[0].Nick)

	user.SetServerName("cake", srv.Host)
	servers, err = user.GetServers()
	assert.Equal(t, "cake", servers[0].Name)

	user.RemoveChannel(srv.Host, chan1.Name)
	channels, err = user.GetChannels()
	assert.Len(t, channels, 1)
	assert.Equal(t, chan2, channels[0])

	user.RemoveServer(srv.Host)
	servers, err = user.GetServers()
	assert.Len(t, servers, 0)
	channels, err = user.GetChannels()
	assert.Len(t, channels, 0)

	user.Remove()
	_, err = os.Stat(storage.Path.User(user.Username))
	assert.True(t, os.IsNotExist(err))

	users, err = storage.LoadUsers(db)
	assert.Nil(t, err)

	for i := range users {
		assert.NotEqual(t, user.ID, users[i].ID)
	}
}

func TestMessages(t *testing.T) {
	storage.Initialize(tempdir())

	db, err := boltdb.New(storage.Path.Database())
	assert.Nil(t, err)

	user, err := storage.NewUser(db)
	assert.Nil(t, err)

	os.MkdirAll(storage.Path.User(user.Username), 0700)

	search, err := bleve.New(storage.Path.Index(user.Username))
	assert.Nil(t, err)

	user.SetMessageStore(db)
	user.SetMessageSearchProvider(search)

	messages, hasMore, err := user.GetMessages("irc.freenode.net", "#go-nuts", 10, "6")
	assert.Nil(t, err)
	assert.False(t, hasMore)
	assert.Len(t, messages, 0)

	messages, hasMore, err = user.GetLastMessages("irc.freenode.net", "#go-nuts", 10)
	assert.Nil(t, err)
	assert.False(t, hasMore)
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

	messages, hasMore, err = user.GetMessages("irc.freenode.net", "#go-nuts", 10, ids[4])
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message3", messages[3].Content)
	assert.Nil(t, err)
	assert.False(t, hasMore)
	assert.Len(t, messages, 4)

	messages, hasMore, err = user.GetMessages("irc.freenode.net", "#go-nuts", 10, betterguid.New())
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message4", messages[4].Content)
	assert.Nil(t, err)
	assert.False(t, hasMore)
	assert.Len(t, messages, 5)

	messages, hasMore, err = user.GetMessages("irc.freenode.net", "#go-nuts", 10, ids[2])
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message1", messages[1].Content)
	assert.Nil(t, err)
	assert.False(t, hasMore)
	assert.Len(t, messages, 2)

	messages, hasMore, err = user.GetLastMessages("irc.freenode.net", "#go-nuts", 10)
	assert.Equal(t, "message0", messages[0].Content)
	assert.Equal(t, "message4", messages[4].Content)
	assert.Nil(t, err)
	assert.False(t, hasMore)
	assert.Len(t, messages, 5)

	messages, hasMore, err = user.GetLastMessages("irc.freenode.net", "#go-nuts", 4)
	assert.Equal(t, "message1", messages[0].Content)
	assert.Equal(t, "message4", messages[3].Content)
	assert.Nil(t, err)
	assert.True(t, hasMore)
	assert.Len(t, messages, 4)

	messages, err = user.SearchMessages("irc.freenode.net", "#go-nuts", "message")
	assert.Nil(t, err)
	assert.True(t, len(messages) > 0)

	db.Close()
}
