package routes

import (
	"github.com/asdine/storm"
	"github.com/dhruvasagar/httpcron/services"
	"github.com/gorilla/mux"
)

func Init(r *mux.Router, sdb *storm.DB, cron *services.Cron) {
	InitIndex(r)
	InitHealth(r)
	InitAPI(r, sdb, cron)
}
