package storage

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func DefaultDirectory() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".dispatch")
}

type directory struct {
	dataRoot   string
	configRoot string
}

func (d directory) DataRoot() string {
	return d.dataRoot
}

func (d directory) ConfigRoot() string {
	return d.configRoot
}

func (d directory) LetsEncrypt() string {
	return filepath.Join(d.ConfigRoot(), "letsencrypt")
}

func (d directory) Users() string {
	return filepath.Join(d.DataRoot(), "users")
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

func (d directory) Downloads(username string) string {
	return filepath.Join(d.User(username), "downloads")
}

func (d directory) Config() string {
	return filepath.Join(d.ConfigRoot(), "config.toml")
}

func (d directory) Database() string {
	return filepath.Join(d.DataRoot(), "dispatch.db")
}
