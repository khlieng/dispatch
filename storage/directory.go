package storage

import (
	"path/filepath"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
)

func DefaultDirectory() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".dispatch")
}

type directory string

func (d directory) Root() string {
	return string(d)
}

func (d directory) LetsEncrypt() string {
	return filepath.Join(d.Root(), "letsencrypt")
}

func (d directory) Users() string {
	return filepath.Join(d.Root(), "users")
}

func (d directory) User(username string) string {
	return filepath.Join(d.Users(), username)
}

func (d directory) Log(username string) string {
	return filepath.Join(d.User(username), "log")
}

func (d directory) Index(username string) string {
	return filepath.Join(d.User(username), "index")
}

func (d directory) Certificate(username string) string {
	return filepath.Join(d.User(username), "cert.pem")
}

func (d directory) Key(username string) string {
	return filepath.Join(d.User(username), "key.pem")
}

func (d directory) Config() string {
	return filepath.Join(d.Root(), "config.toml")
}

func (d directory) Database() string {
	return filepath.Join(d.Root(), "dispatch.db")
}

func (d directory) HMACKey() string {
	return filepath.Join(d.Root(), "hmac.key")
}
