package server

import (
	"github.com/aabiji/page/epub"
	"github.com/gorilla/mux"
	"net/http"
)

type Storage struct {
	LocalPath string
	NetPath   string
}

// Returns a handler that allows cors when serving files over http
func fileEnableCORS(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

// Serve files from LocalPath on the NetPath http endpoint
func (s *Storage) Mount(router *mux.Router) {
	epub.STORAGE_DIRECTORY = s.LocalPath
	fs := http.FileServer(http.Dir(s.LocalPath))
	router.PathPrefix(s.NetPath).Handler(http.StripPrefix(s.NetPath, fileEnableCORS(fs)))
}
