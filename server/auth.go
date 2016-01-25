package server

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/dgrijalva/jwt-go"

	"github.com/khlieng/dispatch/storage"
)

const (
	cookieName = "dispatch"
)

var (
	hmacKey []byte
)

func initAuth() {
	var err error
	hmacKey, err = getHMACKey()
	if err != nil {
		log.Fatal(err)
	}
}

func getHMACKey() ([]byte, error) {
	key, err := ioutil.ReadFile(storage.Path.HMACKey())
	if err != nil {
		key = make([]byte, 32)
		rand.Read(key)

		err = ioutil.WriteFile(storage.Path.HMACKey(), key, 0600)
		if err != nil {
			return nil, err
		}
	}

	return key, nil
}

func handleAuth(w http.ResponseWriter, r *http.Request) *Session {
	var session *Session

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		authLog(r, "No cookie set")
		session = newUser(w, r)
	} else {
		token, err := parseToken(cookie.Value)

		if err == nil && token.Valid {
			userID := uint64(token.Claims["UserID"].(float64))

			log.Println(r.RemoteAddr, "[Auth] GET", r.URL.Path, "| Valid token | User ID:", userID)

			session = sessions.get(userID)
			if session == nil {
				// A previous anonymous session has been cleaned up, create a new one
				session = newUser(w, r)
			}
		} else {
			if err != nil {
				authLog(r, "Invalid token: "+err.Error())
			} else {
				authLog(r, "Invalid token")
			}

			session = newUser(w, r)
		}
	}

	return session
}

func newUser(w http.ResponseWriter, r *http.Request) *Session {
	user, err := storage.NewUser()
	if err != nil {
		return nil
	}

	log.Println(r.RemoteAddr, "[Auth] Create session | User ID:", user.ID)

	session := NewSession(user)
	sessions.set(user.ID, session)
	go session.run()

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["UserID"] = user.ID
	tokenString, err := token.SignedString(hmacKey)
	if err != nil {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 1, 0),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	return session
}

func parseToken(cookie string) (*jwt.Token, error) {
	return jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacKey, nil
	})

}

func authLog(r *http.Request, s string) {
	log.Println(r.RemoteAddr, "[Auth] GET", r.URL.Path, "|", s)
}
