package main

import (
	"encoding/json"
	"errors"
	"github.com/aabiji/page/backend/epub"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

// Receive file from frontend request and save it locally.
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
		return "", errors.New(SERVER_ERORR)
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, file); err != nil {
		return "", errors.New(SERVER_ERORR)
	}

	return filename, nil
}

// Insert row into the Books table in the database.
// Return the id of the inserted row.
func insertBook(e epub.Epub) (int, error) {
	info, err := json.Marshal(e.Info)
	if err != nil {
		return 0, errors.New(SERVER_ERORR)
	}
	toc, err := json.Marshal(e.TableOfContents)
	if err != nil {
		return 0, errors.New(SERVER_ERORR)
	}

	var id int
	sql := `
    INSERT INTO Books 
    (CoverImagePath, Files, TableOfContents, Info) 
    VALUES ($1,$2,$3,$4)
    RETURNING BookId;`
	insert := []any{e.CoverImagePath, e.Files, toc, info}
	if err := database.ExecScan(sql, insert, &id); err != nil {
		return 0, err
	}

	return id, nil
}

// Receive epub file from frontend and insert it into the database.
// Return the number of pages in the epub and the id of the inserted row.
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
		return 0, 0, errors.New(SERVER_ERORR)
	}

	id, err := insertBook(e)
	if err != nil {
		return 0, 0, errors.New(SERVER_ERORR)
	}

	return len(e.Files), id, nil
}
