package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Println("got / request")
}

func handleGreeting(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "Hello %s\n", vars["name"])
}

func Run(addr string) {
    router := mux.NewRouter()
    router.HandleFunc("/{name}", handleGreeting).Methods("GET")
    router.HandleFunc("/", handleRoot).Methods("GET")

    d := "BOOKS"
    p := "/static/"
    router.PathPrefix(p).Handler(http.StripPrefix(p, http.FileServer(http.Dir(d))))
    http.Handle("/", router)

    server := &http.Server{
        Addr: addr,
        Handler: router,
        ReadTimeout: 30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    fmt.Printf("Running server on http://%s\n", addr)
    log.Fatal(server.ListenAndServe())
}
