package letsencrypt

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/xenolf/lego/acme"
)

const URL = "https://acme-v02.api.letsencrypt.org/directory"
const KeySize = 2048

var directory Directory

func Run(dir, domain, email, port string) (*state, error) {
	directory = Directory(dir)

	user, err := getUser(email)
	if err != nil {
		return nil, err
	}

	client, err := acme.NewClient(URL, &user, acme.RSA2048)
	if err != nil {
		return nil, err
	}
	client.SetHTTPAddress(port)

	if user.Registration == nil {
		user.Registration, err = client.Register(true)
		if err != nil {
			return nil, err
		}

		err = saveUser(user)
		if err != nil {
			return nil, err
		}
	}

	s := &state{
		client: client,
		domain: domain,
	}

	if certExists(domain) {
		if !s.renew() {
			err = s.loadCert()
			if err != nil {
				return nil, err
			}
		}
		s.refreshOCSP()
	} else {
		err = s.obtain()
		if err != nil {
			return nil, err
		}
	}

	go s.maintain()

	return s, nil
}

type state struct {
	client  *acme.Client
	domain  string
	cert    *tls.Certificate
	certPEM []byte
	lock    sync.Mutex
}

func (s *state) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	s.lock.Lock()
	cert := s.cert
	s.lock.Unlock()

	return cert, nil
}

func (s *state) getCertPEM() []byte {
	s.lock.Lock()
	certPEM := s.certPEM
	s.lock.Unlock()

	return certPEM
}

func (s *state) setCert(meta *acme.CertificateResource) {
	cert, err := tls.X509KeyPair(meta.Certificate, meta.PrivateKey)
	if err == nil {
		s.lock.Lock()
		if s.cert != nil {
			cert.OCSPStaple = s.cert.OCSPStaple
		}

		s.cert = &cert
		s.certPEM = meta.Certificate
		s.lock.Unlock()
	}
}

func (s *state) setOCSP(ocsp []byte) {
	cert := tls.Certificate{
		OCSPStaple: ocsp,
	}

	s.lock.Lock()
	if s.cert != nil {
		cert.Certificate = s.cert.Certificate
		cert.PrivateKey = s.cert.PrivateKey
	}
	s.cert = &cert
	s.lock.Unlock()
}

func (s *state) obtain() error {
	cert, err := s.client.ObtainCertificate([]string{s.domain}, true, nil, false)
	if err != nil {
		return err
	}

	s.setCert(cert)
	s.refreshOCSP()

	err = saveCert(cert)
	if err != nil {
		return err
	}

	return nil
}

func (s *state) renew() bool {
	cert, err := ioutil.ReadFile(directory.Cert(s.domain))
	if err != nil {
		return false
	}

	exp, err := acme.GetPEMCertExpiration(cert)
	if err != nil {
		return false
	}

	daysLeft := int(exp.Sub(time.Now().UTC()).Hours() / 24)

	if daysLeft <= 30 {
		metaBytes, err := ioutil.ReadFile(directory.Meta(s.domain))
		if err != nil {
			return false
		}

		key, err := ioutil.ReadFile(directory.Key(s.domain))
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

		newMeta, err := s.client.RenewCertificate(meta, true, false)
		if err != nil {
			return false
		}

		s.setCert(newMeta)

		err = saveCert(newMeta)
		if err != nil {
			return false
		}

		return true
	}

	return false
}

func (s *state) refreshOCSP() {
	ocsp, resp, err := acme.GetOCSPForCert(s.getCertPEM())
	if err == nil && resp.Status == acme.OCSPGood {
		s.setOCSP(ocsp)
	}
}

func (s *state) maintain() {
	renew := time.Tick(24 * time.Hour)
	ocsp := time.Tick(1 * time.Hour)
	for {
		select {
		case <-renew:
			s.renew()

		case <-ocsp:
			s.refreshOCSP()
		}
	}
}

func (s *state) loadCert() error {
	cert, err := ioutil.ReadFile(directory.Cert(s.domain))
	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile(directory.Key(s.domain))
	if err != nil {
		return err
	}

	s.setCert(&acme.CertificateResource{
		Certificate: cert,
		PrivateKey:  key,
	})

	return nil
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

func saveCert(cert *acme.CertificateResource) error {
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
