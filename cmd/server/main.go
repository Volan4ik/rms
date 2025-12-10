package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/rms/internal/config"
	"github.com/example/rms/internal/db"
	api "github.com/example/rms/internal/http"
	"github.com/example/rms/internal/repository"
)

// @title Restaurant Management System API
// @version 1.0
// @description REST API for restaurant hall, orders, and warehouse management
// @BasePath /api
func main() {
	cfg := config.Load()
	database := db.MustConnect(cfg)
	defer database.Close()

	repo := repository.New(database)
	router := api.NewRouter(repo)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
