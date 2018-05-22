package server

import (
	"log"
	"net/http"
	"time"

	"github.com/khlieng/dispatch/storage"
)

const (
	cookieName = "dispatch"
)

func handleAuth(w http.ResponseWriter, r *http.Request, createUser bool) *Session {
	var session *Session

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if createUser {
			session = newUser(w, r)
		}
	} else {
		session = sessions.get(cookie.Value)
		if session != nil {
			log.Println(r.RemoteAddr, "[Auth] GET", r.URL.Path, "| Valid token | User ID:", session.user.ID)
		} else if createUser {
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

	session, err := NewSession(user)
	if err != nil {
		return nil
	}
	sessions.set(session)
	go session.run()

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    session.id,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 1, 0),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	return session
}
