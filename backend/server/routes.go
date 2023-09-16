package server

import (
	"encoding/json"
	"fmt"
	"github.com/aabiji/page/epub"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func handleError(w http.ResponseWriter, err error) {
	j, _ := json.Marshal(map[string]string{
		"Server error": fmt.Sprintf("%s", err.Error()),
	})
	w.Write(j)
}

func setExampleCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "userId",
		Value:    "user123",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(364 * 24 * time.Second),
	}
	http.SetCookie(w, &cookie)

	message := map[string]string{
		"Status": "Cookie set",
	}
	json.NewEncoder(w).Encode(message)
}

func getBookInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := fmt.Sprintf("epub/tests/%s.epub", vars["name"])
	w.Header().Set("Content-Type", "application/json")

	book, err := epub.New(path)
	if err != nil {
		handleError(w, err)
		return
	}

	c, err := r.Cookie("userId")
	if err != nil {
		handleError(w, err)
		return
	}
	fmt.Println(c.Name, c.Value)

	scrollOffsets := []int{}
	for i := 0; i < len(book.Files); i++ {
		scrollOffsets = append(scrollOffsets, 0)
	}
	userBookInfo := map[string]any{
		"Epub":              book,
		"CurrentPage":       0,
		"FileScrollOffsets": scrollOffsets,
	}
	json.NewEncoder(w).Encode(userBookInfo)
}
