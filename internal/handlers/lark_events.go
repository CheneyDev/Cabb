package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/labstack/echo/v4"
)

// Minimal event struct to handle challenge handshake
type larkEventEnvelope struct {
    Challenge string `json:"challenge"`
    Type      string `json:"type"`
}

func (h *Handler) LarkEvents(c echo.Context) error {
    var env larkEventEnvelope
    if err := json.NewDecoder(c.Request().Body).Decode(&env); err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    if env.Challenge != "" {
        return c.JSON(http.StatusOK, map[string]string{"challenge": env.Challenge})
    }
    return c.NoContent(http.StatusOK)
}

func (h *Handler) LarkInteractivity(c echo.Context) error {
    return c.NoContent(http.StatusOK)
}

func (h *Handler) LarkCommands(c echo.Context) error {
    return c.NoContent(http.StatusOK)
}

