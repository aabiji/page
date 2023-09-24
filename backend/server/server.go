package server

import (
	"fmt"
	"github.com/aabiji/page/backend/db"
	"github.com/aabiji/page/backend/epub"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var database db.DB = db.NewDatabase()

// Serve files from localPath on the netPath http endpoint
func serveFiles(router *mux.Router) {
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
	router.PathPrefix(netPath).Handler(http.StripPrefix(netPath, fs)) //FilesAllowCORS(fs)))
}

// Run http server and handle all api endpoints
func Run(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/book/get/{name}", GetBook).Methods("GET")
	router.HandleFunc("/user/auth", AuthAccount).Methods("POST")
	router.HandleFunc("/user/create", CreateAccount).Methods("POST")

	serveFiles(router)
	cors := NewCORS("http://localhost:5173", []string{"GET", "POST"})
	corsRouter := cors.AllowRequests(router)

	fmt.Printf("Running server on http://%s\n", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      corsRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
