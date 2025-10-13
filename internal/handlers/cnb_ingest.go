package handlers

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

type cnbIssuePayload struct {
    Event   string `json:"event"`
    Repo    string `json:"repo"`
    IssueIID string `json:"issue_iid"`
}

type cnbPRPayload struct {
    Event string `json:"event"`
    Action string `json:"action"`
    Repo string `json:"repo"`
    PRIid string `json:"pr_iid"`
}

type cnbBranchPayload struct {
    Event  string `json:"event"`
    Action string `json:"action"`
    Repo   string `json:"repo"`
    Branch string `json:"branch"`
}

func (h *Handler) CNBIngestIssue(c echo.Context) error {
    if !h.authorizeIntegration(c) {
        return c.NoContent(http.StatusUnauthorized)
    }
    var p cnbIssuePayload
    if err := c.Bind(&p); err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    return c.JSON(http.StatusAccepted, map[string]any{
        "accepted": true,
        "source":   "cnb.issue",
        "payload":  p,
    })
}

func (h *Handler) CNBIngestPR(c echo.Context) error {
    if !h.authorizeIntegration(c) {
        return c.NoContent(http.StatusUnauthorized)
    }
    var p cnbPRPayload
    if err := c.Bind(&p); err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    return c.JSON(http.StatusAccepted, map[string]any{
        "accepted": true,
        "source":   "cnb.pr",
        "payload":  p,
    })
}

func (h *Handler) CNBIngestBranch(c echo.Context) error {
    if !h.authorizeIntegration(c) {
        return c.NoContent(http.StatusUnauthorized)
    }
    var p cnbBranchPayload
    if err := c.Bind(&p); err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    return c.JSON(http.StatusAccepted, map[string]any{
        "accepted": true,
        "source":   "cnb.branch",
        "payload":  p,
    })
}

func (h *Handler) authorizeIntegration(c echo.Context) bool {
    if h.cfg.IntegrationToken == "" {
        return true // allow in scaffold when not configured
    }
    auth := c.Request().Header.Get("Authorization")
    want := "Bearer " + h.cfg.IntegrationToken
    return auth == want
}

