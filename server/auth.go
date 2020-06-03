package server

import (
	"log"
	"net/http"

	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
)

func (d *Dispatch) handleAuth(w http.ResponseWriter, r *http.Request, createUser, refresh bool) *State {
	var state *State

	cookie, err := r.Cookie(session.CookieName)
	if err != nil {
		if createUser {
			state, err = d.newUser(w, r)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		session := d.states.getSession(cookie.Value)
		if session != nil {
			key := session.Key()

			if !session.Expired() {
				state = d.states.get(session.UserID)

				if refresh {
					newKey, err := session.Refresh()
					if err != nil {
						log.Println(err)
					}

					if newKey != "" {
						d.states.setSession(session)
						d.states.deleteSession(key)
						session.SetCookie(w, r)
					}
				}
			} else {
				d.states.deleteSession(key)
			}
		}

		if state != nil {
			log.Println(r.RemoteAddr, "[Auth] GET", r.URL.Path, "| Valid token | User ID:", state.user.ID)
		} else if createUser {
			state, err = d.newUser(w, r)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return state
}

func (d *Dispatch) newUser(w http.ResponseWriter, r *http.Request) (*State, error) {
	user, err := storage.NewUser(d.Store)
	if err != nil {
		return nil, err
	}

	log.Println(r.RemoteAddr, "[Auth] New anonymous user | ID:", user.ID)

	session, err := session.New(user.ID)
	if err != nil {
		return nil, err
	}
	d.states.setSession(session)

	state := NewState(user, d)
	d.states.set(state)
	go state.run()

	session.SetCookie(w, r)

	return state, nil
}
