package cookie

import "net/http"

const HostPrefix = "__Host-"

func Harden(r *http.Request, cookie *http.Cookie) *http.Cookie {
	cookie.HttpOnly = true
	cookie.Secure = r.TLS != nil

	if cookie.Path == "" {
		cookie.Path = "/"
	}

	if cookie.Path == "/" && cookie.Secure {
		cookie.Name = HostPrefix + cookie.Name
	}

	if cookie.SameSite == 0 {
		cookie.SameSite = http.SameSiteLaxMode
	}

	return cookie
}

func Set(w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
	http.SetCookie(w, Harden(r, cookie))
}

func Name(r *http.Request, name string) string {
	if r.TLS != nil {
		return HostPrefix + name
	}
	return name
}
