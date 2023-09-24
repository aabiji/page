package server

import (
	"net/http"
	"slices"
	"strings"
)

type CORS struct {
	methods        string
	allowedOrigin  string
	allowedMethods []string
}

func NewCORS(origin string, methods []string) CORS {
	c := CORS{allowedMethods: methods, allowedOrigin: origin}
	c.methods = strings.Join(methods, ", ")
	return c
}

// Return a handler that allows cors when accessing files over http
func FilesAllowCORS(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

func isPreflightRequest(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	method := r.Header.Get("Acess-Control-Request-Method")
	return r.Method == "OPTIONS" && origin != "" && method != ""
}

// Return a http handler that allows cors when making http requests from a certain origin
func (c *CORS) AllowRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		allowedOrigin := origin == c.allowedOrigin

		if allowedOrigin {
			w.Header().Add("Origin", "Vary")
			w.Header().Add("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if isPreflightRequest(r) {
			method := r.Header.Get("Access-Control-Request-Method")
			if allowedOrigin && slices.Contains(c.allowedMethods, method) {
				w.Header().Set("Acess-Control-Allow-Methods", c.methods)
			}
		}

		handler.ServeHTTP(w, r)
	})
}
