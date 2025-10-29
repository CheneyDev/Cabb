package handlers

import (
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "cabb/internal/version"
)

func (h *Handler) Healthz(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]any{
        "status":  "ok",
        "version": version.Version,
        "time":    time.Now().Format(time.RFC3339),
        "db":      "not_connected",
    })
}

