package handlers

import (
	"github.com/labstack/echo/v4"
)

// ServeLogs handles the WebSocket connection for real-time logs.
// It performs manual authentication to support query parameter tokens (useful for WebSockets).
func (h *Handler) ServeLogs(c echo.Context) error {
	// Try standard auth (cookie or query param)
	_, err := h.sessionFromRequest(c)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	// Auth success
	return h.broadcaster.ServeWS(c)
}
