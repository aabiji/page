package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user/auth", AuthAccount).Methods("POST")
	router.HandleFunc("/user/create", CreateAccount).Methods("POST")

	router.HandleFunc("/book/get/{name}", GetBook).Methods("GET")
    router.HandleFunc("/book/upload", EpubUpload).Methods("POST")
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
