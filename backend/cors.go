package main

import (
	"net/http"
	"slices"
	"strings"
)

func isPreflightRequest(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	method := r.Header.Get("Access-Control-Request-Method")
	return r.Method == "OPTIONS" && origin != "" && method != ""
}

// Return a http handler that allows cors when making http requests from a certain origin
func AllowRequests(allowedOrigin string, handler http.Handler) http.Handler {
	allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		allowedOrigin := origin == allowedOrigin
		allowHeader := "Content-Type, withCredentials, Authorization"

		if allowedOrigin {
			w.Header().Add("Origin", "Vary")
			w.Header().Add("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Headers", allowHeader)
		}

		if isPreflightRequest(r) {
			method := r.Header.Get("Access-Control-Request-Method")
			methodsStr := strings.Join(allowedMethods, ", ")
			if allowedOrigin && slices.Contains(allowedMethods, method) {
				w.Header().Set("Access-Control-Allow-Methods", methodsStr)
			}
		}

		handler.ServeHTTP(w, r)
	})
}
