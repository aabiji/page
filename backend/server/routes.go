package server

import (
	"encoding/json"
	"fmt"
	"github.com/aabiji/page/backend/db"
	"github.com/aabiji/page/backend/epub"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

// Return json containing error value to client signalling an internal server error.
func errorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	response := map[string]string{
		"Server error": err.Error(),
	}
	json.NewEncoder(w).Encode(response)
}

// Set cookie header in http response to client.
func setCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   r.Host != "localhost:8080",
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(364 * 24 * time.Second),
	}
	http.SetCookie(w, &cookie)

	// Empty json response
	response := map[string]string{}
	json.NewEncoder(w).Encode(response)
}

// Get POST request json payload from request body.
func getRequestJson[T any](w http.ResponseWriter, r *http.Request, data *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

// Validate user login credentials and return cookie containing userId.
// The userId cookie will be used to add user state to other requests made by client.
func AuthAccount(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := getRequestJson(w, r, &user); err != nil {
		errorResponse(w, err)
		return
	}

	user, err := database.ReadUser(user)
	if err != nil {
		errorResponse(w, err)
		return
	}

	setCookie(w, r, "userId", user.Id)
}

// Validate and create new user account and return cookie containing userId.
// The userId cookie will be used to add user state to other requests made by client.
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := getRequestJson(w, r, &user); err != nil {
		errorResponse(w, err)
		return
	}

	// TODO: check for duplicate accounts
	err := database.CreateUser(user)
	if err != nil {
		errorResponse(w, err)
		return
	}

	setCookie(w, r, "userId", user.Id)
}

// Get book info. NOTE: this function is temporary.
func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := fmt.Sprintf("epub/tests/%s.epub", vars["name"])
	w.Header().Set("Content-Type", "application/json")

	book, err := epub.New(path)
	if err != nil {
		errorResponse(w, err)
		return
	}

	c, err := r.Cookie("userId")
	if err != nil {
		errorResponse(w, err)
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
