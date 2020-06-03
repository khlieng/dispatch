package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateGetSetUsers(t *testing.T) {
	state := newState(NewClient(&Config{}))
	users := []string{"a", "b"}
	state.setUsers(users, "#chan")
	assert.Equal(t, users, state.getUsers("#chan"))
	state.setUsers(users, "#chan")
	assert.Equal(t, users, state.getUsers("#chan"))
}

func TestStateAddRemoveUser(t *testing.T) {
	state := newState(NewClient(&Config{}))
	state.addUser("user", "#chan")
	state.addUser("user", "#chan")
	assert.Len(t, state.getUsers("#chan"), 1)
	state.addUser("user2", "#chan")
	assert.Equal(t, []string{"user", "user2"}, state.getUsers("#chan"))
	state.removeUser("user", "#chan")
	assert.Equal(t, []string{"user2"}, state.getUsers("#chan"))
}

func TestStateRemoveUserAll(t *testing.T) {
	state := newState(NewClient(&Config{}))
	state.addUser("user", "#chan1")
	state.addUser("user", "#chan2")
	state.removeUserAll("user")
	assert.Empty(t, state.getUsers("#chan1"))
	assert.Empty(t, state.getUsers("#chan2"))
}

func TestStateRenameUser(t *testing.T) {
	state := newState(NewClient(&Config{}))
	state.addUser("user", "#chan1")
	state.addUser("user", "#chan2")
	state.renameUser("user", "new")
	assert.Equal(t, []string{"new"}, state.getUsers("#chan1"))
	assert.Equal(t, []string{"new"}, state.getUsers("#chan2"))

	state.addUser("@gotop", "#chan3")
	state.renameUser("gotop", "stillgotit")
	assert.Equal(t, []string{"@stillgotit"}, state.getUsers("#chan3"))
}

func TestStateMode(t *testing.T) {
	state := newState(NewClient(&Config{}))
	state.addUser("+user", "#chan")
	state.setMode("#chan", "user", "o", "v")
	assert.Equal(t, []string{"@user"}, state.getUsers("#chan"))
	state.setMode("#chan", "user", "v", "")
	assert.Equal(t, []string{"@user"}, state.getUsers("#chan"))
	state.setMode("#chan", "user", "", "o")
	assert.Equal(t, []string{"+user"}, state.getUsers("#chan"))
	state.setMode("#chan", "user", "q", "")
	assert.Equal(t, []string{"~user"}, state.getUsers("#chan"))
}

func TestStateTopic(t *testing.T) {
	state := newState(NewClient(&Config{}))
	assert.Equal(t, "", state.getTopic("#chan"))
	state.setTopic("the topic", "#chan")
	assert.Equal(t, "the topic", state.getTopic("#chan"))
}

func TestStateChannelUserMode(t *testing.T) {
	user := NewUser("&test")
	assert.Equal(t, "test", user.nick)
	assert.Equal(t, "a", string(user.modes[0]))
	assert.Equal(t, "&test", user.String())

	user.RemoveModes("a")
	assert.Equal(t, "test", user.String())
	user.AddModes("o")
	assert.Equal(t, "@test", user.String())
	user.AddModes("q")
	assert.Equal(t, "~test", user.String())
	user.AddModes("v")
	assert.Equal(t, "~test", user.String())
	user.RemoveModes("qo")
	assert.Equal(t, "+test", user.String())
	user.RemoveModes("v")
	assert.Equal(t, "test", user.String())
}
