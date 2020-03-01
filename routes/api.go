package routes

import (
	"github.com/asdine/storm"
	"github.com/dhruvasagar/httpcron/handlers"
	"github.com/dhruvasagar/httpcron/routes/api"
	"github.com/dhruvasagar/httpcron/services"
	"github.com/gorilla/mux"
)

func InitAPI(r *mux.Router, sdb *storm.DB, cron *services.Cron) {
	s := r.PathPrefix(
		"/api",
	).Headers(
		"Content-Type", "application/json",
	).Subrouter()

	s.Use(handlers.AuthorizationHandler)

	api.InitCronsAPI(r, sdb, cron)
}
