package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestIndex(t *testing.T) {
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	InitIndex(r)

	req := httptest.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != "HTTP Cron - Simple Scheduler with an http api" {
		t.Error("Did not get expected HTTP response body, got", w.Body.String())
	}
}
