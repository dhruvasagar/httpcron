package handlers

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestAuthoriziationHandler(t *testing.T) {
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/", dummyHandler).Methods("GET")

	r.Use(AuthorizationHandler)

	req := httptest.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != "Unauthorized Access\n" {
		t.Error("Did not get expected HTTP response body, got", w.Body.String())
	}
}

func TestAuthoriziationHandlerSuccess(t *testing.T) {
	os.Setenv("USER_TOKEN", "test")

	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/", dummyHandler).Methods("GET")

	r.Use(AuthorizationHandler)

	req := httptest.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("test")))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
}
