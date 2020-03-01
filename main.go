package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dhruvasagar/httpcron/db"
	"github.com/dhruvasagar/httpcron/routes"
	"github.com/dhruvasagar/httpcron/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func getBindAddr() string {
	bindAddr := os.Getenv("ADDR")
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}
	return bindAddr
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	return ":" + port
}

func getListenAddr() string {
	return getBindAddr() + getPort()
}

func getWaitTimeout() time.Duration {
	waitStr := os.Getenv("WAIT_TIMEOUT")
	waitInt, _ := strconv.Atoi(waitStr)
	wait := time.Duration(waitInt)
	if wait == 0 {
		wait = 15 * time.Second
	}
	return wait
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := db.New()

	cron := services.NewCron()

	r := mux.NewRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	routes.Init(r, db, cron)

	srv := &http.Server{
		Addr: getListenAddr(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      loggedRouter,
	}

	go func() {
		fmt.Println("Listening on port ", getPort())
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), getWaitTimeout())
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	db.Close()
	os.Exit(0)
}
