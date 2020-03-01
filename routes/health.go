package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func InitHealth(r *mux.Router) {
	r.HandleFunc("/health", healthHandler).Methods("GET")
}
