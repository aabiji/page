package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/aabiji/page/backend/epub"
	"github.com/gorilla/mux"
)

func setStorageDirectories() {
	extractDir := os.Getenv("EPUB_EXTRACT_DIRECTORY")
	uploadDir := os.Getenv("FILE_UPLOAD_DIRECTORY")
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	if extractDir == "" {
		extractDir = filepath.Join(currentUser.HomeDir, "Page", "BOOKS")
		if err := os.MkdirAll(extractDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	epub.EXTRACT_DIRECTORY = extractDir

	if uploadDir == "" {
		uploadDir = filepath.Join(currentUser.HomeDir, "Page", "FILES")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	FILE_UPLOAD_DIRECTORY = uploadDir
}

func main() {
	setStorageDirectories()

	router := mux.NewRouter()
	router.HandleFunc("/user/login", AuthAccount).Methods("POST")
	router.HandleFunc("/user/create", CreateAccount).Methods("POST")
	router.HandleFunc("/user/book/upload", UserUploadEpub).Methods("POST")
	router.HandleFunc("/user/book/get/{id}", GetUserBookInfo).Methods("GET")

	router.HandleFunc("/book/get/{id}", GetBook).Methods("GET")
	ServeFiles(router)

	addr := "localhost:8080"
	corsRouter := AllowRequests("http://localhost:5173", router)
	fmt.Printf("Running server on http://%s\n", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      corsRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
