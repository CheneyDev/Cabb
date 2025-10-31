package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"cabb/internal/handlers"
	"cabb/internal/store"
	"cabb/internal/version"
	"cabb/pkg/config"
)

func main() {
	cfg := config.FromEnv()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Standard middlewares: recover from panics and attach request id (custom, dependency-free)
	e.Use(handlers.Recover())
	e.Use(handlers.RequestID())
	// Structured JSON request logging (per AGENTS.md)
	e.Use(handlers.StructuredLogger())

	// Basic server-level timeouts via stdlib server
	srv := &http.Server{
		Addr:              cfg.Address(),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Optional DB connect
	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		log.Printf("db connect error: %v", err)
	}
	// Auto-run SQL migrations when DB is available
	if db != nil && db.SQL != nil {
		if err := db.RunMigrations(context.Background(), "."); err != nil {
			log.Fatalf("apply migrations failed: %v", err)
		}
		if err := handlers.BootstrapAdminUser(context.Background(), db, cfg); err != nil {
			log.Printf("bootstrap admin user failed: %v", err)
		}
	}

	handlers.RegisterRoutes(e, cfg, db)

	// Start background schedulers
	handlers.StartCleanupScheduler(cfg, db)

	log.Printf("plane-integration %s listening on %s", version.Version, cfg.Address())
	if err := e.StartServer(srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
