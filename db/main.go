package db

import (
	"os"

	"github.com/asdine/storm"
	log "github.com/sirupsen/logrus"
)

func getDBPath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "httpcron.db"
	}
	return dbPath
}

func New() *storm.DB {
	db, err := storm.Open(getDBPath())
	if err != nil {
		log.Fatal("Couldn't open database", err)
	}
	return db
}
