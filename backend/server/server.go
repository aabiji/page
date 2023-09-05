package server

import (
	"fmt"
    "github.com/aabiji/read/epub"
	"github.com/gorilla/mux"
    "github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / request")
}

func handleGreeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    path := fmt.Sprintf("epub/tests/%s.epub", vars["bookName"])

    e, err := epub.New(path)
	if err != nil {
        fmt.Fprintf(w, "%s\n", err.Error())
        return
	}

    fmt.Fprintf(w, "%v", e.Files)
}

func Run(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/{bookName}", handleGreeting).Methods("GET")
	router.HandleFunc("/", handleRoot).Methods("GET")
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
