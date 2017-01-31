package main

import (
	"net/http"
	"strings"
)

func headerMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "github.com/tantalic/servemd")
		w.Header().Set("X-Servemd-Version", Version)
		handler(w, r)
	}
}

func basicAuthMiddleware(handler http.HandlerFunc, users []string) http.HandlerFunc {
	// If no users are defined don't add the middleware
	if len(users) == 0 {
		return handler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			goto unauthorized
		}

		for _, u := range users {
			userpass := strings.SplitN(u, ":", 2)
			if len(userpass) == 2 && username == userpass[0] && password == userpass[1] {
				handler(w, r)
				return
			}
		}

	unauthorized:
		w.Header().Set("WWW-Authenticate", `Basic realm="Protected Content"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
