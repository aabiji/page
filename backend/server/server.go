package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func Run(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/book/get/{name}", getBookInfo).Methods("GET")
	corsRouter := cors.Default().Handler(router)

	storage := Storage{LocalPath: "BOOKS", NetPath: "/static/"}
	storage.Mount(router)

	fmt.Printf("Running server on http://%s\n", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      corsRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
