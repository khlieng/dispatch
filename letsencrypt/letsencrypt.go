package letsencrypt

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/xenolf/lego/acme"
)

const URL = "https://acme-v01.api.letsencrypt.org/directory"
const KeySize = 2048

var directory Directory

func Run(dir, domain, email, port string, onChange func()) (string, string, error) {
	directory = Directory(dir)

	user, err := getUser(email)
	if err != nil {
		return "", "", nil
	}

	client, err := acme.NewClient(URL, &user, KeySize)
	client.ExcludeChallenges([]string{"tls-sni-01"})
	client.SetHTTPPort(port)

	if user.Registration == nil {
		user.Registration, err = client.Register()
		if err != nil {
			return "", "", err
		}

		err = client.AgreeToTOS()
		if err != nil {
			return "", "", err
		}

		err = saveUser(user)
		if err != nil {
			return "", "", err
		}
	}

	if certExists(domain) {
		renew(client, domain)
	} else {
		err = obtain(client, domain)
		if err != nil {
			return "", "", err
		}
	}

	go keepRenewed(client, domain, onChange)

	return directory.Cert(domain), directory.Key(domain), nil
}

func obtain(client *acme.Client, domain string) error {
	cert, errors := client.ObtainCertificate([]string{domain}, false)
	if err := errors[domain]; err != nil {
		if _, ok := err.(acme.TOSError); ok {
			err := client.AgreeToTOS()
			if err != nil {
				return err
			}
			return obtain(client, domain)
		}

		return err
	}

	err := saveCert(cert)
	if err != nil {
		return err
	}

	return nil
}

func renew(client *acme.Client, domain string) bool {
	cert, err := ioutil.ReadFile(directory.Cert(domain))
	if err != nil {
		return false
	}

	exp, err := acme.GetPEMCertExpiration(cert)
	if err != nil {
		return false
	}

	daysLeft := int(exp.Sub(time.Now().UTC()).Hours() / 24)

	if daysLeft <= 30 {
		metaBytes, err := ioutil.ReadFile(directory.Meta(domain))
		if err != nil {
			return false
		}

		key, err := ioutil.ReadFile(directory.Key(domain))
		if err != nil {
			return false
		}

		var meta acme.CertificateResource
		err = json.Unmarshal(metaBytes, &meta)
		if err != nil {
			return false
		}
		meta.Certificate = cert
		meta.PrivateKey = key

	Renew:
		newMeta, err := client.RenewCertificate(meta, false)
		if err != nil {
			if _, ok := err.(acme.TOSError); ok {
				err := client.AgreeToTOS()
				if err != nil {
					return false
				}
				goto Renew
			}
			return false
		}

		err = saveCert(newMeta)
		if err != nil {
			return false
		}

		return true
	}

	return false
}

func keepRenewed(client *acme.Client, domain string, onChange func()) {
	for {
		time.Sleep(24 * time.Hour)
		if renew(client, domain) {
			onChange()
		}
	}
}

func certExists(domain string) bool {
	if _, err := os.Stat(directory.Cert(domain)); err != nil {
		return false
	}
	if _, err := os.Stat(directory.Key(domain)); err != nil {
		return false
	}
	return true
}

func saveCert(cert acme.CertificateResource) error {
	err := os.MkdirAll(directory.Domain(cert.Domain), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(directory.Cert(cert.Domain), cert.Certificate, 0600)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(directory.Key(cert.Domain), cert.PrivateKey, 0600)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(&cert, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(directory.Meta(cert.Domain), jsonBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
