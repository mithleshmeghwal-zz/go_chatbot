package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"os"
	"net/url"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.CreateRoutes()

	code := m.Run()

	os.Exit(code)
}

func TestVerificationEndpoint(t *testing.T) {

	req, _ := http.NewRequest("GET", "localhost/webhook", nil)
	q := url.Values{}
	q.Add("hub.challenge", "123456789")
	q.Add("hub.mode", "subscribe")
	q.Add("hub.verify_token", os.Getenv("VERIFY_TOKEN"))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// func TestMessagesEndpoint(t *testing.T) {
	
// 	req, _ := http.NewRequest("POST", "localhost/webhook", nil)

// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// }

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.R.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
