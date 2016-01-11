package storage

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"
)

var (
	ErrInvalidCert      = errors.New("Invalid certificate")
	ErrCouldNotSaveCert = errors.New("Could not save certificate")
)

func (u *User) GetCertificate() *tls.Certificate {
	u.lock.Lock()
	cert := u.certificate
	u.lock.Unlock()

	return cert
}

func (u *User) SetCertificate(certPEM, keyPEM []byte) error {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return ErrInvalidCert
	}
	u.lock.Lock()
	u.certificate = &cert
	u.lock.Unlock()

	err = os.MkdirAll(Path.User(u.UUID), 0700)
	if err != nil {
		return ErrCouldNotSaveCert
	}

	err = ioutil.WriteFile(Path.Certificate(u.UUID), certPEM, 0600)
	if err != nil {
		return ErrCouldNotSaveCert
	}

	err = ioutil.WriteFile(Path.Key(u.UUID), keyPEM, 0600)
	if err != nil {
		return ErrCouldNotSaveCert
	}

	return nil
}

func (u *User) loadCertificate() error {
	certPEM, err := ioutil.ReadFile(Path.Certificate(u.UUID))
	if err != nil {
		return err
	}

	keyPEM, err := ioutil.ReadFile(Path.Key(u.UUID))
	if err != nil {
		return err
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return err
	}

	u.certificate = &cert
	return nil
}
