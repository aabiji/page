package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aabiji/page/backend/epub"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

var database DB = NewDatabase()
const ( // Server error responses
    NOT_FOUND = "Entries not found"
    BAD_CLIENT_REQUEST = "Bad client request"
    SERVER_ERORR = "Internal server error. Please try again"
)
const MAX_UPLOAD_SIZE = 100 << 20 // 100 megabyte limit on uploaded epub files

// Return json containing error value to client signalling an internal server error.
func respondWithError(w http.ResponseWriter, err error) {
    errorCode := http.StatusOK
    if err.Error() == SERVER_ERORR {
        errorCode = http.StatusInternalServerError
    } else if err.Error() == BAD_CLIENT_REQUEST {
        errorCode = http.StatusBadRequest
    }

	w.WriteHeader(errorCode)
	response := map[string]string{"Server error": err.Error()}
	json.NewEncoder(w).Encode(response)
}

// Set cookie header in http response to client.
func setCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: false,
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

// GET /static/* (ex. /static/path/to/file.html)
// Serve files from localPath on the netPath http endpoint
func ServeFiles(router *mux.Router) {
	localPath := os.Getenv("EPUB_STORAGE_DIRECTORY")
	if localPath == "" {
		// Default to directory in current user home if environment variable is not set
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		localPath = filepath.Join(currentUser.HomeDir, "BOOKS")
		os.MkdirAll(localPath, os.ModePerm)
	}

	netPath := "/static/"
	epub.STORAGE_DIRECTORY = localPath

	fs := http.FileServer(http.Dir(localPath))
	router.PathPrefix(netPath).Handler(http.StripPrefix(netPath, fs)) //FilesAllowCORS(fs)))
}

// POST /user/auth
// Validate user login credentials and return cookie containing userId.
// The userId cookie will be used to add user state to other requests made by client.
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
// Validate and create new user account and return cookie containing userId.
// The userId cookie will be used to add user state to other requests made by client.
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := getRequestJson(w, r, &user); err != nil {
		respondWithError(w, errors.New(BAD_CLIENT_REQUEST))
		return
	}

	sql := "SELECT UserId FROM Users WHERE Email=$1"
	_, err := database.Read(sql, []any{user.Email}, []any{&user.Id})
	if err == nil {
		msg := "Account already exists. Create a new one with a different email."
		respondWithError(w, errors.New(msg))
		return
	}

	sql = "INSERT INTO Users (Email, Password) VALUES ($1, $2);"
	err = database.Exec(sql, user.Email, user.Password)
	if err != nil {
		respondWithError(w, errors.New(SERVER_ERORR))
		return
	}

	setCookie(w, r, "userId", user.Id)
}

// Receive file from frontend request and save it locally.
func receiveFile(w http.ResponseWriter, r *http.Request) (string, error) {
    if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
        return "", errors.New(BAD_CLIENT_REQUEST)
    }

    uploadedFile, handler, err := r.FormFile("file")
    if err != nil {
        return "", errors.New(BAD_CLIENT_REQUEST)
    }
    defer uploadedFile.Close()

    // TODO: store files in a preconfigured directory
    filename := filepath.Join(epub.STORAGE_DIRECTORY, "_EPUB", handler.Filename)
    localFile, err := os.Create(filename)
    if err != nil {
        return "", errors.New(SERVER_ERORR)
    }
    defer localFile.Close()

    if _, err := io.Copy(localFile, uploadedFile); err != nil {
        return "", errors.New(SERVER_ERORR)
    }

    return filename, nil
}

// POST /book/upload
// Upload epub file to esrver.
func EpubUpload(w http.ResponseWriter, r *http.Request) {
    filename, err := receiveFile(w, r)
    if err != nil {
        respondWithError(w, err)
        return
    }

    fmt.Println(filename)

    response := map[string]string{"Status": "Epub uploaded successfully"}
    json.NewEncoder(w).Encode(response)
}

// GET /book/get/{name}
// Get book info. NOTE: this function is temporary.
func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := fmt.Sprintf("epub/tests/%s.epub", vars["name"])
	w.Header().Set("Content-Type", "application/json")

	book, err := epub.New(path)
	if err != nil {
		respondWithError(w, err)
		return
	}

	fmt.Println(r.Cookies())
	c, err := r.Cookie("userId")
	if err != nil {
		respondWithError(w, err)
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
