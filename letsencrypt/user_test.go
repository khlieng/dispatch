package letsencrypt

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tempdir() string {
	f, _ := ioutil.TempDir("", "")
	return f
}

func testUser(t *testing.T, email string) {
	user, err := newUser(email)
	assert.Nil(t, err)
	key := user.GetPrivateKey()
	assert.NotNil(t, key)

	err = saveUser(user)
	assert.Nil(t, err)

	user, err = getUser(email)
	assert.Nil(t, err)
	assert.Equal(t, email, user.GetEmail())
	assert.Equal(t, key, user.GetPrivateKey())
}

func TestUser(t *testing.T) {
	directory = Directory(tempdir())

	testUser(t, "test@test.com")
	testUser(t, "")
}
