package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aabiji/page/backend/epub"
	"github.com/gorilla/mux"
)

var database DB = NewDatabase()
var FILE_UPLOAD_DIRECTORY string // Directory where uploaded files will be stored
const (
	NOT_FOUND          = "Entries not found"
	BAD_CLIENT_REQUEST = "Bad client request"
	SERVER_ERORR       = "Internal server error. Please try again"
	MAX_UPLOAD_SIZE    = 100 << 20 // 100 megabyte limit on all uploaded files
)

// GET /static/* (ex. /static/path/to/file.html)
// Serve requested file from disk to the client.
func ServeFiles(router *mux.Router) {
	route := "/static/"
	fs := http.FileServer(http.Dir(epub.EXTRACT_DIRECTORY))
	router.PathPrefix(route).Handler(http.StripPrefix(route, fs))
}

// POST /user/auth
// Request payload: {"email": "", "password": "", "confirm": ""}
// Response: An empty json response and a cookie containing the user's id.
// Validate user login credentials and set a "userId" cookie to manage client state.
func AuthAccount(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := getRequestJson(w, r, &user); err != nil {
		respondWithError(w, errors.New(BAD_CLIENT_REQUEST))
		return
	}

	sql := "SELECT UserId FROM Users WHERE Email=$1 AND Password=$2;"
	_, err := database.Read(sql, []any{user.Email, user.Password}, []any{&user.Id})
	if err != nil && err.Error() == NOT_FOUND {
		msg := "Account not found. Forgot your password?"
		respondWithError(w, errors.New(msg))
		return
	} else if err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	setCookie(w, r, "userId", user.Id)
}

// POST /user/create
// Request payload: {"email": "", "password": "", "confirm":""}
// Response: An empty json response and a cookie containing the user's id.
// Validate and create new user account and set a "userId" cookie to manage client state.
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := getRequestJson(w, r, &user); err != nil {
		respondWithError(w, errors.New(BAD_CLIENT_REQUEST))
		return
	}

	sql := "SELECT UserId FROM Users WHERE Email=$1;"
	_, err := database.Read(sql, []any{user.Email}, []any{&user.Id})
	if err == nil {
		msg := "Account already exists. Create a new one with a different email."
		respondWithError(w, errors.New(msg))
		return
	}

	sql = "INSERT INTO Users (Email, Password) VALUES ($1,$2) RETURNING UserId;"
	err = database.ExecScan(sql, []any{user.Email, user.Password}, &user.Id)
	if err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	setCookie(w, r, "userId", user.Id)
}

// POST /user/book/upload
// Request payload: Multipart form data with field "file".
// Response: {"BookId": ""}
// Upload a user selected epub file to the server. Add it to the user's collection
// of books and return the generated bookId.
func UserUploadEpub(w http.ResponseWriter, r *http.Request) {
	pageCount, bookId, err := receiveEpub(w, r)
	if err != nil {
		respondWithError(w, err)
		return
	}

	c, err := r.Cookie("userId")
	if err != nil {
		respondWithError(w, errors.New(BAD_CLIENT_REQUEST))
		return
	}

	scrollOffsets := make([]int, pageCount) // TODO: don't allocate please
	sql := "INSERT INTO UserBooks (UserId, BookId, CurrentPage, ScrollOffsets) VALUES ($1,$2,$3,$4);"
	if err := database.Exec(sql, c.Value, bookId, 0, scrollOffsets); err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	response := map[string]int{"BookId": bookId}
	json.NewEncoder(w).Encode(response)
}

// GET /user/book/get/{id}
// Response: {"CurrentPage": "", "ScrollOffsets": ""}
// Get user specific information related to specific book.
func GetUserBookState(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("userId")
	if err != nil {
		respondWithError(w, errors.New(BAD_CLIENT_REQUEST))
		return
	}
	bookId := mux.Vars(r)["id"]

	currentPage := 0
	scrollOffsets := []int{}
	sql := "SELECT CurrentPage, ScrollOffsets FROM UserBooks WHERE UserId=$1 AND BookId=$2;"
	_, err = database.Read(sql, []any{c.Value, bookId}, []any{&currentPage, &scrollOffsets})
	if err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	response := map[string]any{
		"CurrentPage":   currentPage,
		"ScrollOffsets": scrollOffsets,
	}
	json.NewEncoder(w).Encode(response)
}

// GET /book/get/{id}
// Response:
// {
//  	"Files": [""],
//  	"CoverImagePath": "",
//  	"TableOfContents": [{"Path": "", "Section": ""}],
//		"Info": {
//			"Language": "",
//			"Author": "",
//			"Title": "",
//			"Identifier": "",
//			"Contributor": ""
// 			"Rights": "",
//			"Source": "",
//			"Coverage": "",
// 			"Relation": "",
// 			"Publisher": "",
// 			"Description": "",
// 			"Date": "",
// 			"Subjects": [""],
//		}
// }
// Get detailed information about a book using it's unique id.
func GetBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["id"]

	var imgPath string
	var files []string
	var toc, info []byte
	sql := "SELECT CoverImagePath, Files, TableOfContents, Info FROM Books WHERE BookId=$1;"
	_, err := database.Read(sql, []any{bookId}, []any{&imgPath, &files, &toc, &info})
	if err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	var tocObj []epub.Section
	if err := json.Unmarshal(toc, &tocObj); err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	var infoObj epub.Metadata
	if err := json.Unmarshal(info, &infoObj); err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
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
