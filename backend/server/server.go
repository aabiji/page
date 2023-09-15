package server

import (
	"fmt"
	"github.com/aabiji/page/epub"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// Returns a handler that allows cors when serving files over http
func fileEnableCORS(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

// Serve files from localPath on the netPath http endpoint
func ServeFiles(router *mux.Router) {
	localPath := os.Getenv("EPUB_STORAGE_DIRECTORY")
	if localPath == "" {
		// Default to directory in current user home if environment variable is not set
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		localPath = filepath.Join(currentUser.HomeDir, "BOOKS")
		os.MkdirAll(localPath, os.ModePerm)
	}

	netPath := "/static/"
	epub.STORAGE_DIRECTORY = localPath

	fs := http.FileServer(http.Dir(localPath))
	router.PathPrefix(netPath).Handler(http.StripPrefix(netPath, fileEnableCORS(fs)))
}

func Run(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/book/get/{name}", getBookInfo).Methods("GET")
	corsRouter := cors.Default().Handler(router)
	ServeFiles(router)

	fmt.Printf("Running server on http://%s\n", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      corsRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
