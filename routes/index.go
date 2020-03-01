package routes

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "HTTP Cron - Simple Scheduler with an http api")
}

func InitIndex(r *mux.Router) {
	r.HandleFunc("/", index)
}
