package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /
// If FRONTEND_BASE_URL is set, redirect to it; otherwise return 404 JSON
func (h *Handler) Root(c echo.Context) error {
	if h.cfg.FrontendBaseURL != "" {
		return c.Redirect(http.StatusFound, h.cfg.FrontendBaseURL)
	}
	return c.JSON(http.StatusNotFound, map[string]any{"message": "Not Found"})
}
