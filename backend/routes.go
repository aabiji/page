package main

import (
	"encoding/json"
	"net/http"

	"github.com/aabiji/page/backend/epub"
	"github.com/gorilla/mux"
)

var database DB = NewDatabase()
var FILE_UPLOAD_DIRECTORY string  // Directory where uploaded files will be stored
const MAX_UPLOAD_SIZE = 100 << 20 // 100 megabyte limit on all uploaded files
const (
	USERID             = "userId"
	NOT_FOUND          = "Entries not found"
	BAD_CLIENT_REQUEST = "Bad client request"
	INTERNAL_ERROR     = "Internal server error. Please try again"
	DUPLICATE_ACCOUNT  = "Account already exists. Create a new one with a different email."
	DUPLICATE_BOOK     = "Book is already in the user's collection."
	ACCOUNT_NOT_FOUND  = "Account not found. Forgot your password?"
)

// GET /static/* (ex. /static/path/to/file.html)
//
// Serve requested file from disk to the client.
func ServeFiles(router *mux.Router) {
	route := "/static/"
	fs := http.FileServer(http.Dir(epub.EXTRACT_DIRECTORY))
	router.PathPrefix(route).Handler(http.StripPrefix(route, fs))
}

// POST /user/login
//
// Request payload: {"email": "", "password": "", "confirm": ""}
//
// Response: An empty json response and a cookie containing the user's id.
//
// Validate user login credentials and set a USERID cookie to manage client state.
func AuthAccount(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := getRequestJson(w, r, &user); err != nil {
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}
	if user.Email == "" || user.Password == "" {
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}

	sql := "SELECT UserId FROM Users WHERE Email=$1 AND Password=$2;"
	_, err := database.Read(sql, []any{user.Email, user.Password}, []any{&user.Id})
	if err != nil && err.Error() == NOT_FOUND {
		respondWithError(w, ACCOUNT_NOT_FOUND)
		return
	} else if err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	setCookie(w, r, USERID, user.Id)
}

// POST /user/create
//
// Request payload: {"email": "", "password": "", "confirm":""}
//
// Response: An empty json response and a cookie containing the user's id.
//
// Validate and create new user account and set a USERID cookie to manage client state.
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := getRequestJson(w, r, &user); err != nil {
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}
	if user.Email == "" || user.Password == "" {
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}

	sql := "SELECT UserId FROM Users WHERE Email=$1;"
	_, err := database.Read(sql, []any{user.Email}, []any{&user.Id})
	if err == nil {
		respondWithError(w, DUPLICATE_ACCOUNT)
		return
	}

	sql = "INSERT INTO Users (Email, Password) VALUES ($1,$2) RETURNING UserId;"
	err = database.ExecScan(sql, []any{user.Email, user.Password}, &user.Id)
	if err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	setCookie(w, r, USERID, user.Id)
}

// POST /user/remove
//
// Request payload:
// Cookie with name set to "userId" and value set to the user's id.
//
// Response: Empty json reponse.
//
// Remove all rows in the Users table and the UserBooks table where
// the "UserId" field matches the user id sent in the request cookie.
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(USERID)
	if err != nil {
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}

	usersDelete := "DELETE FROM Users WHERE UserId=$1;"
	if err := database.Exec(usersDelete, c.Value); err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	booksDelete := "DELETE FROM UserBooks WHERE UserId=$1;"
	if err := database.Exec(booksDelete, c.Value); err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{})
}

// POST /user/book/upload
//
// Request payload:
// Multipart form data with field "file".
// Cookie with name set to "userId" and value set to the user's id.
//
// Response: {"BookId": ""}
//
// Upload a user selected epub file to the server. Add it to the user's collection
// of books and return the generated bookId
func UserUploadEpub(w http.ResponseWriter, r *http.Request) {
	pageCount, bookId, err := receiveEpub(w, r)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}

	c, err := r.Cookie(USERID)
	if err != nil { // Cookie not found
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}

	sql := "SELECT UserId FROM UserBooks WHERE BookId=$1;"
	if _, err := database.Read(sql, []any{bookId}, []any{&c.Value}); err == nil {
		respondWithError(w, DUPLICATE_BOOK)
		return
	}

	scrollOffsets := make([]int, pageCount)
	sql = "INSERT INTO UserBooks (UserId, BookId, CurrentPage, ScrollOffsets) VALUES ($1,$2,$3,$4);"
	if err := database.Exec(sql, c.Value, bookId, 0, scrollOffsets); err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	response := map[string]int{"BookId": bookId}
	json.NewEncoder(w).Encode(response)
}

// DELETE /user/book/remove/{id}
//
// Request payload: Cookie with name set to "userId" and value set to the user's id.
//
// Response: Empty json response.
//
// Remove a book by id from the user's collection.
func UserRemoveBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["id"]
	c, err := r.Cookie(USERID)
	if err != nil {
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}

	sql := "DELETE FROM UserBooks WHERE BookId=$1 AND UserId=$2;"
	if err := database.Exec(sql, bookId, c.Value); err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{})
}

// GET /user/book/get/{id}
//
// Request payload: Cookie with name set to "userId" and value set to the user's id.
//
// Response: {"CurrentPage": "", "ScrollOffsets": ""}
//
// Get user specific information related to specific book.
func GetUserBookInfo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(USERID)
	if err != nil { // Cookie not found
		respondWithError(w, BAD_CLIENT_REQUEST)
		return
	}
	bookId := mux.Vars(r)["id"]

	currentPage := 0
	scrollOffsets := []int{}
	sql := "SELECT CurrentPage, ScrollOffsets FROM UserBooks WHERE UserId=$1 AND BookId=$2;"
	_, err = database.Read(sql, []any{c.Value, bookId}, []any{&currentPage, &scrollOffsets})
	if err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	response := map[string]any{
		"CurrentPage":   currentPage,
		"ScrollOffsets": scrollOffsets,
	}
	json.NewEncoder(w).Encode(response)
}

// GET /book/get/{id}
//
// Response:
//
//	{
//	 	"Files": [""],
//	 	"CoverImagePath": "",
//	 	"TableOfContents": [{"Path": "", "Section": ""}],
//			"Info": {
//				"Language": "",
//				"Author": "",
//				"Title": "",
//				"Identifier": "",
//				"Contributor": ""
//				"Rights": "",
//				"Source": "",
//				"Coverage": "",
//				"Relation": "",
//				"Publisher": "",
//				"Description": "",
//				"Date": "",
//				"Subjects": [""],
//			}
//	}
//
// Get detailed information about a book using it's unique id.
func GetBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["id"]

	var imgPath string
	var files []string
	var toc, info []byte
	sql := "SELECT CoverImagePath, Files, TableOfContents, Info FROM Books WHERE BookId=$1;"
	_, err := database.Read(sql, []any{bookId}, []any{&imgPath, &files, &toc, &info})
	if err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	var tocObj []epub.Section
	if err := json.Unmarshal(toc, &tocObj); err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	var infoObj epub.Metadata
	if err := json.Unmarshal(info, &infoObj); err != nil {
		respondWithError(w, INTERNAL_ERROR)
		return
	}

	response := map[string]any{
		"CoverImagePath":  imgPath,
		"Files":           files,
		"TableOfContents": tocObj,
		"Info":            infoObj,
	}
	json.NewEncoder(w).Encode(response)
}
