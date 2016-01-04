package letsencrypt

import (
	"path/filepath"
)

type Directory string

func (d Directory) Domain(domain string) string {
	return filepath.Join(string(d), "certs", domain)
}

func (d Directory) Cert(domain string) string {
	return filepath.Join(d.Domain(domain), "cert.pem")
}

func (d Directory) Key(domain string) string {
	return filepath.Join(d.Domain(domain), "key.pem")
}

func (d Directory) Meta(domain string) string {
	return filepath.Join(d.Domain(domain), "metadata.json")
}

func (d Directory) User(email string) string {
	if email == "" {
		email = defaultUser
	}
	return filepath.Join(string(d), "users", email)
}

func (d Directory) UserRegistration(email string) string {
	return filepath.Join(d.User(email), "registration.json")
}

func (d Directory) UserKey(email string) string {
	return filepath.Join(d.User(email), "key.pem")
}
