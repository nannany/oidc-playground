package middleware

import (
	"net/http"
	"rp/domain"
	"rp/session"
)

// SessionCheck checks if the session is valid
func SessionCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rpSession, _ := session.Store.Get(r, "rp-session")
		userID := rpSession.Values["user"]
		if userID != nil && domain.Users[userID.(string)] != nil {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

	}
}
