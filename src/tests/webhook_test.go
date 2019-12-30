package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lungria/spendshelf-backend/src/api"
)

func TestWebHookGet(t *testing.T) {
	server := api.WebHookAPI{}
	server.InitRouter(":8080")

	req, err := http.NewRequest("GET", "/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	server.HTTPServer.Handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code: %v,got: %v", http.StatusOK, recorder.Code)
	}
}

func TestWebHookPostBadRequest(t *testing.T) {

	server := api.WebHookAPI{}
	server.InitRouter(":8080")

	req, err := http.NewRequest("POST", "/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	server.HTTPServer.Handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Wrong status code. Expected %v, got %v", http.StatusBadRequest, recorder.Code)
	}
}
