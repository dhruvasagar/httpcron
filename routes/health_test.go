package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	InitHealth(r)

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != "{\"ok\":true}\n" {
		t.Error("Did not get expected HTTP response body, got", w.Body.String())
	}
}
