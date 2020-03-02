package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/dhruvasagar/httpcron/db"
	"github.com/dhruvasagar/httpcron/services"
	"github.com/gorilla/mux"
)

func renderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func index(sdb *storm.DB, cron *services.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cronEntries []db.CronEntry
		err := sdb.All(&cronEntries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		renderJSON(w, cronEntries)
	}
}

func createOrUpdate(
	r *http.Request,
	sdb *storm.DB,
	cron *services.Cron,
) (*db.CronEntry, error) {
	var cronEntry db.CronEntry
	err := json.NewDecoder(r.Body).Decode(&cronEntry)
	if err != nil {
		return nil, err
	}
	entryID, err := cron.Add(cronEntry.Spec, cronEntry.Command)
	if err != nil {
		return nil, err
	}
	cronEntry.EntryID = entryID
	return &cronEntry, sdb.Save(&cronEntry)
}

func create(sdb *storm.DB, cron *services.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cronEntry, err := createOrUpdate(r, sdb, cron)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		renderJSON(w, cronEntry)
	}
}

func get(sdb *storm.DB, cron *services.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		var cronEntry db.CronEntry
		err := sdb.One("ID", id, &cronEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		renderJSON(w, cronEntry)
	}
}

func update(sdb *storm.DB, cron *services.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cronEntry, err := createOrUpdate(r, sdb, cron)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		renderJSON(w, cronEntry)
	}
}

func del(sdb *storm.DB, cron *services.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.Atoi(idStr)

		var cronEntry db.CronEntry
		err := sdb.One("ID", id, &cronEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = sdb.DeleteStruct(&cronEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cron.Remove(cronEntry.EntryID)

		w.WriteHeader(http.StatusOK)
		renderJSON(w, map[string]bool{"ok": true})
	}
}

func InitCronsAPI(r *mux.Router, sdb *storm.DB, cron *services.Cron) {
	s := r.PathPrefix(
		"/crons/",
	).Headers(
		"Content-Type", "application/json",
	).Subrouter()

	s.HandleFunc("/", index(sdb, cron)).Methods("GET")
	s.HandleFunc("/", create(sdb, cron)).Methods("POST")
	s.HandleFunc("/{id}", get(sdb, cron)).Methods("GET")
	s.HandleFunc("/{id}", update(sdb, cron)).Methods("PUT")
	s.HandleFunc("/{id}", del(sdb, cron)).Methods("DELETE")
}
