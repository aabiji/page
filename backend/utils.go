package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aabiji/page/backend/epub"
)

// Add json containing an error value to the http response and set the appropriate response code.
func respondWithError(w http.ResponseWriter, err string) {
	errorCode := http.StatusOK
	if err == INTERNAL_ERROR {
		errorCode = http.StatusInternalServerError
	} else if err == BAD_CLIENT_REQUEST {
		errorCode = http.StatusBadRequest
	}

	w.WriteHeader(errorCode)
	response := map[string]string{"Server error": err}
	json.NewEncoder(w).Encode(response)
}

// Add cookie header to the http response sent to the client.
func setCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: false}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(map[string]string{})
}

// Get json payload from the body of a POST request.
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

// Receive a file sent from a frontend request and save it on disk.
func receiveFile(w http.ResponseWriter, r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		return "", errors.New(BAD_CLIENT_REQUEST)
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", errors.New(BAD_CLIENT_REQUEST)
	}
	defer file.Close()

	filename := filepath.Join(FILE_UPLOAD_DIRECTORY, handler.Filename)
	localFile, err := os.Create(filename)
	if err != nil {
		return "", errors.New(INTERNAL_ERROR)
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, file); err != nil {
		return "", errors.New(INTERNAL_ERROR)
	}

	return filename, nil
}

// Receive an epub file sent from a frontend request.
// Insert relavent extracted information from the epub file into the database.
// Return the number of pages in the epub, the id of the newly inserted row and a portential error.
func receiveEpub(w http.ResponseWriter, r *http.Request) (int, int, error) {
	filename, err := receiveFile(w, r)
	if err != nil {
		return 0, 0, err
	}

	fileparts := strings.Split(filename, ".")
	if fileparts[len(fileparts)-1] != "epub" {
		os.Remove(filename)
		return 0, 0, errors.New(BAD_CLIENT_REQUEST)
	}

	e, err := epub.New(filename)
	if err != nil {
		os.Remove(filename)
		return 0, 0, errors.New(INTERNAL_ERROR)
	}

	id, err := insertBook(e)
	if err != nil {
		return 0, 0, errors.New(INTERNAL_ERROR)
	}

	return len(e.Files), id, nil
}

// Get the BookId of a book with a given title
func getBook(title string) (int, error) {
	var id int
	sql := "SELECT BookId FROM Books WHERE Title=$1;"
	if _, err := database.Read(sql, []any{title}, []any{&id}); err != nil {
		return 0, err
	}
	return id, nil
}

// Insert a new entry to the Books table in the database
// and return the id of the newly inserted row.
func insertBook(e epub.Epub) (int, error) {
	info, err := json.Marshal(e.Info)
	if err != nil {
		return 0, errors.New(INTERNAL_ERROR)
	}
	toc, err := json.Marshal(e.TableOfContents)
	if err != nil {
		return 0, errors.New(INTERNAL_ERROR)
	}

	var id int
	id, err = getBook(e.Info.Title)
	if err == nil { // A book with the same title has already been inserted.
		return id, nil
	}

	sql := `
    INSERT INTO Books 
    (Title, CoverImagePath, Files, TableOfContents, Info) 
    VALUES ($1, $2, $3, $4, $5)
    RETURNING BookId;`
	insert := []any{e.Info.Title, e.CoverImagePath, e.Files, toc, info}
	if err := database.ExecScan(sql, insert, &id); err != nil {
		return 0, err
	}

	return id, nil
}
