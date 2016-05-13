package main

import "net/http"

func headerMiddleware(handler httpHandleFunc) httpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "github.com/tantalic/servemd")
		w.Header().Set("X-Servemd-Version", Version)
		handler(w, r)
	}
}
