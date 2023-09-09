package server

import (
    "encoding/json"
    "fmt"
    "net/http"
	"github.com/aabiji/read/epub"
	"github.com/gorilla/mux"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / request")
}

func handleGreeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := fmt.Sprintf("epub/tests/%s.epub", vars["bookName"])

	book, err := epub.New(path)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}
