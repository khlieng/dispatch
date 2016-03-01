package letsencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/xenolf/lego/acme"
)

const defaultUser = "default"

type User struct {
	Email        string
	Registration *acme.RegistrationResource
	key          *rsa.PrivateKey
}

func (u User) GetEmail() string {
	return u.Email
}

func (u User) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}

func (u User) GetPrivateKey() *rsa.PrivateKey {
	return u.key
}

func newUser(email string) (User, error) {
	var err error
	user := User{Email: email}
	user.key, err = rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return user, err
	}
	return user, nil
}

func getUser(email string) (User, error) {
	var user User

	reg, err := os.Open(directory.UserRegistration(email))
	if err != nil {
		if os.IsNotExist(err) {
			return newUser(email)
		}
		return user, err
	}
	defer reg.Close()

	err = json.NewDecoder(reg).Decode(&user)
	if err != nil {
		return user, err
	}

	user.key, err = loadRSAPrivateKey(directory.UserKey(email))
	if err != nil {
		return user, err
	}

	return user, nil
}

func saveUser(user User) error {
	err := os.MkdirAll(directory.User(user.Email), 0700)
	if err != nil {
		return err
	}

	err = saveRSAPrivateKey(user.key, directory.UserKey(user.Email))
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(&user, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(directory.UserRegistration(user.Email), jsonBytes, 0600)
}

func loadRSAPrivateKey(file string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
}

func saveRSAPrivateKey(key *rsa.PrivateKey, file string) error {
	pemKey := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	keyOut, err := os.Create(file)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	return pem.Encode(keyOut, &pemKey)
}
