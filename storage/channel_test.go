package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetUsers(t *testing.T) {
	channelStore := NewChannelStore()
	users := []string{"a,b"}
	channelStore.SetUsers(users, "srv", "#chan")
	assert.Equal(t, channelStore.GetUsers("srv", "#chan"), users)
}

func TestAddRemoveUser(t *testing.T) {
	channelStore := NewChannelStore()
	channelStore.AddUser("user", "srv", "#chan")
	channelStore.AddUser("user", "srv", "#chan")
	assert.Len(t, channelStore.GetUsers("srv", "#chan"), 1)
	channelStore.AddUser("user2", "srv", "#chan")
	assert.Equal(t, []string{"user", "user2"}, channelStore.GetUsers("srv", "#chan"))
	channelStore.RemoveUser("user", "srv", "#chan")
	assert.Equal(t, []string{"user2"}, channelStore.GetUsers("srv", "#chan"))
}

func TestRemoveUserAll(t *testing.T) {
	channelStore := NewChannelStore()
	channelStore.AddUser("user", "srv", "#chan1")
	channelStore.AddUser("user", "srv", "#chan2")
	channelStore.RemoveUserAll("user", "srv")
	assert.Empty(t, channelStore.GetUsers("srv", "#chan1"))
	assert.Empty(t, channelStore.GetUsers("srv", "#chan2"))
}

func TestRenameUser(t *testing.T) {
	channelStore := NewChannelStore()
	channelStore.AddUser("user", "srv", "#chan1")
	channelStore.AddUser("user", "srv", "#chan2")
	channelStore.RenameUser("user", "new", "srv")
	assert.Equal(t, []string{"new"}, channelStore.GetUsers("srv", "#chan1"))
	assert.Equal(t, []string{"new"}, channelStore.GetUsers("srv", "#chan2"))

	channelStore.AddUser("@gotop", "srv", "#chan3")
	channelStore.RenameUser("gotop", "stillgotit", "srv")
	assert.Equal(t, []string{"@stillgotit"}, channelStore.GetUsers("srv", "#chan3"))
}

func TestMode(t *testing.T) {
	channelStore := NewChannelStore()
	channelStore.AddUser("+user", "srv", "#chan")
	channelStore.SetMode("srv", "#chan", "user", "o", "v")
	assert.Equal(t, []string{"@user"}, channelStore.GetUsers("srv", "#chan"))
	channelStore.SetMode("srv", "#chan", "user", "v", "")
	assert.Equal(t, []string{"+user"}, channelStore.GetUsers("srv", "#chan"))
	channelStore.SetMode("srv", "#chan", "user", "", "v")
	assert.Equal(t, []string{"user"}, channelStore.GetUsers("srv", "#chan"))
}

func TestTopic(t *testing.T) {
	channelStore := NewChannelStore()
	assert.Equal(t, "", channelStore.GetTopic("srv", "#chan"))
	channelStore.SetTopic("the topic", "srv", "#chan")
	assert.Equal(t, "the topic", channelStore.GetTopic("srv", "#chan"))
}

func TestFindUserChannels(t *testing.T) {
	channelStore := NewChannelStore()
	channelStore.AddUser("user", "srv", "#chan1")
	channelStore.AddUser("user", "srv", "#chan2")
	channelStore.AddUser("user2", "srv", "#chan3")
	channelStore.AddUser("user", "srv2", "#chan4")
	channelStore.AddUser("@gotop", "srv", "#chan1")

	channels := channelStore.FindUserChannels("user", "srv")
	assert.Len(t, channels, 2)
	assert.Contains(t, channels, "#chan1")
	assert.Contains(t, channels, "#chan2")
	channels = channelStore.FindUserChannels("gotop", "srv")
	assert.Len(t, channels, 1)
	assert.Contains(t, channels, "#chan1")
}
