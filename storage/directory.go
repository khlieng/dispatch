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

func (d directory) Logs() string {
	return filepath.Join(d.Root(), "logs")
}

func (d directory) Log(userID string) string {
	return filepath.Join(d.Logs(), userID+".log")
}

func (d directory) Index(userID string) string {
	return filepath.Join(d.Logs(), userID+".idx")
}

func (d directory) Users() string {
	return filepath.Join(d.Root(), "users")
}

func (d directory) User(userID string) string {
	return filepath.Join(d.Users(), userID)
}

func (d directory) Certificate(userID string) string {
	return filepath.Join(d.User(userID), "cert.pem")
}

func (d directory) Key(userID string) string {
	return filepath.Join(d.User(userID), "key.pem")
}

func (d directory) Config() string {
	return filepath.Join(d.Root(), "config.toml")
}

func (d directory) Database() string {
	return filepath.Join(d.Root(), "dispatch.db")
}
