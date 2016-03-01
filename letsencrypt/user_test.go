package letsencrypt

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xenolf/lego/acme"
)

func tempdir() string {
	f, _ := ioutil.TempDir("", "")
	return f
}

func testUser(t *testing.T, email string) {
	reg := &acme.RegistrationResource{
		URI: "test.com",
		Body: acme.Registration{
			Agreement: "agree?",
		},
	}

	user, err := newUser(email)
	assert.Nil(t, err)
	key := user.GetPrivateKey()
	assert.NotNil(t, key)
	user.Registration = reg

	err = saveUser(user)
	assert.Nil(t, err)

	user, err = getUser(email)
	assert.Nil(t, err)
	assert.Equal(t, email, user.GetEmail())
	assert.Equal(t, key, user.GetPrivateKey())
	assert.Equal(t, reg, user.GetRegistration())
}

func TestUser(t *testing.T) {
	directory = Directory(tempdir())

	testUser(t, "test@test.com")
	testUser(t, "")
}
