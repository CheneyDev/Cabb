package handlers

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

func (h *Handler) AdminRepoProject(c echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"message": "repo-project mapping API not implemented in scaffold"})
}

func (h *Handler) AdminPRStates(c echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"message": "pr-states mapping API not implemented in scaffold"})
}

func (h *Handler) AdminUsers(c echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"message": "users mapping API not implemented in scaffold"})
}

func (h *Handler) AdminChannelProject(c echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"message": "channel-project mapping API not implemented in scaffold"})
}

