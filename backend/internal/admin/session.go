package admin

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "admin_session"
const sessionUserKey = "authenticated"

var store *sessions.CookieStore

func initStore(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/admin",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // set true in production (HTTPS)
		SameSite: http.SameSiteStrictMode,
	}
}

func isAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return false
	}
	auth, ok := session.Values[sessionUserKey].(bool)
	return ok && auth
}

func setAuthenticated(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values[sessionUserKey] = true
	return store.Save(r, w, session)
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return
	}
	session.Values[sessionUserKey] = false
	session.Options.MaxAge = -1
	store.Save(r, w, session)
}
