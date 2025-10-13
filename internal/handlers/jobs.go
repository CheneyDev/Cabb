package handlers

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

func (h *Handler) JobIssueSummaryDaily(c echo.Context) error {
    return c.JSON(http.StatusAccepted, map[string]any{
        "scheduled": true,
        "job":       "issue-summary-daily",
    })
}

