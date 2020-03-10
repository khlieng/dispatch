package storage

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
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

	err = ioutil.WriteFile(ConfigPath.Certificate(u.Username), certPEM, 0600)
	if err != nil {
		return ErrCouldNotSaveCert
	}

	err = ioutil.WriteFile(ConfigPath.Key(u.Username), keyPEM, 0600)
	if err != nil {
		return ErrCouldNotSaveCert
	}

	return nil
}

func (u *User) loadCertificate() error {
	certPEM, err := ioutil.ReadFile(ConfigPath.Certificate(u.Username))
	if err != nil {
		return err
	}

	keyPEM, err := ioutil.ReadFile(ConfigPath.Key(u.Username))
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
