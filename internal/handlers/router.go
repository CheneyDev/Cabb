package handlers

import (
    "time"

    "github.com/labstack/echo/v4"
    "plane-integration/internal/store"
    "plane-integration/pkg/config"
)

func RegisterRoutes(e *echo.Echo, cfg config.Config, db *store.DB) {
    // Initialize a lightweight in-memory deduper (5 minutes TTL)
    d := NewDeduper(5 * time.Minute)

    h := &Handler{cfg: cfg, dedupe: d, db: db}

    // Health
    e.GET("/healthz", h.Healthz)
    // Root redirect to frontend (optional)
    e.GET("/", h.Root)

    // Plane OAuth + Webhook
    e.GET("/plane/oauth/start", h.PlaneOAuthStart)
    e.GET("/plane/oauth/callback", h.PlaneOAuthCallback)
    e.POST("/webhooks/plane", h.PlaneWebhook)

    // CNB ingest callbacks via .cnb.yml
    e.POST("/ingest/cnb/issue", h.CNBIngestIssue)
    e.POST("/ingest/cnb/pr", h.CNBIngestPR)
    e.POST("/ingest/cnb/branch", h.CNBIngestBranch)

    // Feishu (Lark)
    e.POST("/webhooks/lark/events", h.LarkEvents)
    e.POST("/webhooks/lark/interactivity", h.LarkInteractivity)
    e.POST("/webhooks/lark/commands", h.LarkCommands)

    // Admin mappings
    e.POST("/admin/mappings/repo-project", h.AdminRepoProject)
    e.GET("/admin/mappings/repo-project", h.AdminRepoProjectList)
    e.POST("/admin/mappings/pr-states", h.AdminPRStates)
    e.POST("/admin/mappings/users", h.AdminUsers)
    e.POST("/admin/mappings/labels", h.AdminLabels)
    e.POST("/admin/mappings/channel-project", h.AdminChannelProject)

    // Jobs
    e.POST("/jobs/issue-summary/daily", h.JobIssueSummaryDaily)
}

type Handler struct {
    cfg config.Config
    dedupe *Deduper
    db *store.DB
}
