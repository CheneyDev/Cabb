package main

import (
    "log"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"

    "plane-integration/internal/handlers"
    "plane-integration/internal/store"
    "plane-integration/internal/version"
    "plane-integration/pkg/config"
)

func main() {
    cfg := config.FromEnv()

    e := echo.New()
    e.HideBanner = true
    e.HidePort = true

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

    handlers.RegisterRoutes(e, cfg, db)

    log.Printf("plane-integration %s listening on %s", version.Version, cfg.Address())
    if err := e.StartServer(srv); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
