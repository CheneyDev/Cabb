package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"cabb/internal/version"
)

func (h *Handler) Healthz(c echo.Context) error {
	dbStatus := "not_configured"
	if h.db != nil && h.db.SQL != nil {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := h.db.Ping(ctx); err == nil {
			dbStatus = "connected"
		} else {
			dbStatus = "error: " + err.Error()
		}
	}
	
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "ok",
		"version": version.Version,
		"time":    time.Now().Format(time.RFC3339),
		"db":      dbStatus,
	})
}
