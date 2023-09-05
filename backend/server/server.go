package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func fileEnableCORS(fs http.Handler) http.HandlerFunc {
	// returns a handler that allows cors when serving files over http
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / request")
}

func handleGreeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s\n", vars["name"])
}

func Run(addr string) {
	// Handling routes
	router := mux.NewRouter()
	router.HandleFunc("/{name}", handleGreeting).Methods("GET")
	router.HandleFunc("/", handleRoot).Methods("GET")
	corsRouter := cors.Default().Handler(router)

	// Serving static files
	d := "BOOKS"
	p := "/static/"
	fs := http.FileServer(http.Dir(d))
	router.PathPrefix(p).Handler(http.StripPrefix(p, fileEnableCORS(fs)))

	fmt.Printf("Running server on http://%s\n", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      corsRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
