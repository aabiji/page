package main

// TODO --> test ServeFiles, UserUploadBook, GetUserBookInfo, GetBook

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type object map[string]any
type httpHandler func(w http.ResponseWriter, r *http.Request)

func randomString(size int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		arr[i] = chars[rand.Intn(len(chars))]
	}
	return string(arr)
}

func findCookie(res *http.Response, name string) *http.Cookie {
	for _, c := range res.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func attachFileToRequest(filename string) (*bytes.Buffer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContents := &bytes.Buffer{}
	writer := multipart.NewWriter(fileContents)
	defer writer.Close()

	filepart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(filepart, file)
	if err != nil {
		return nil, err
	}
	return fileContents, nil
}

func callApi(endpoint, method string, payload object, handler httpHandler) (*http.Response, object, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, object{}, err
	}
	reader := bytes.NewReader(encoded)
	requqest := httptest.NewRequest(method, endpoint, reader)
	w := httptest.NewRecorder()
	handler(w, requqest)
	response := w.Result()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, object{}, err
	}
	defer response.Body.Close()

	responseJson := object{}
	err = json.Unmarshal(body, &responseJson)
	if err != nil {
		return nil, object{}, err
	}

	return response, responseJson, nil
}

func TestAuth(t *testing.T) {
	payload := object{}
	// test for incorrect payload
	response, responseJson, err := callApi("/user/login", http.MethodPost, payload, AuthAccount)
	if err != nil {
		t.Errorf("Exepecting nil error, got %v", err)
	}
	if val, _ := responseJson["Server error"]; val != BAD_CLIENT_REQUEST {
		t.Errorf("Expecting %s, got %s", BAD_CLIENT_REQUEST, val)
	}

	// Test for login where the account is not found
	payload["password"] = "supersecretpassword"
	payload["email"] = randomString(5) + "@email.com"
	response, responseJson, err = callApi("/user/login", http.MethodPost, payload, AuthAccount)
	if err != nil {
		t.Errorf("Exepecting nil error, got %v", err)
	}
	if val, _ := responseJson["Server error"]; val != ACCOUNT_NOT_FOUND {
		t.Errorf("Expecting %s, got %s", ACCOUNT_NOT_FOUND, val)
	}

	// Test for normal login
	_, _, _ = callApi("/user/create", http.MethodPost, payload, CreateAccount)
	response, responseJson, err = callApi("/user/login", http.MethodPost, payload, AuthAccount)
	if len(responseJson) != 0 {
		t.Errorf("Expecting emoty json response, got %v", responseJson)
	}
	c := findCookie(response, USERID)
	if c == nil {
		t.Errorf("Expected %s cookie value to be non nil, found %v", USERID, c)
	}
}

func TestCreateAccount(t *testing.T) {
	payload := object{}
	// Test for incorrect payload
	response, responseJson, err := callApi("/usr/create", http.MethodPost, payload, CreateAccount)
	if err != nil {
		t.Errorf("Expecting nil error, got %v", err)
	}
	if val, _ := responseJson["Server error"]; val != BAD_CLIENT_REQUEST {
		t.Errorf("Expecting %s, got %s", BAD_CLIENT_REQUEST, val)
	}

	// Test for regular signup
	payload["password"] = "supersecretpassword"
	payload["email"] = randomString(5) + "@email.com"
	response, responseJson, err = callApi("/user/create", http.MethodPost, payload, CreateAccount)
	if err != nil {
		t.Errorf("Expecting nil error, got %v", err)
	}
	if len(responseJson) != 0 {
		t.Errorf("Expected response to be an empty json object, found %v", responseJson)
	}
	c := findCookie(response, USERID)
	if c == nil {
		t.Errorf("Expected %s cookie value to be non nil, found %v", USERID, c)
	}

	// Test for duplicate accounts (same email used by more than 1 account)
	response, responseJson, err = callApi("/user/create", http.MethodPost, payload, CreateAccount)
	if err != nil {
		t.Errorf("Expecting nil error, got %v", err)
	}

	if val, _ := responseJson["Server error"]; val != DUPLICATE_ACCOUNT {
		t.Errorf("Expecting %s, found %s", DUPLICATE_ACCOUNT, val)
	}
}
