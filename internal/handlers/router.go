package handlers

import (
	"sync"
	"time"

	"plane-integration/internal/store"
	"plane-integration/pkg/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, cfg config.Config, db *store.DB) {
	// Initialize a lightweight in-memory deduper (5 minutes TTL)
	d := NewDeduper(5 * time.Minute)

	sessionCookie := cfg.AdminSessionCookie
	if sessionCookie == "" {
		sessionCookie = "pi_admin_session"
	}
	ttl := time.Duration(cfg.AdminSessionTTLHours) * time.Hour
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}

	h := &Handler{
		cfg:                 cfg,
		dedupe:              d,
		db:                  db,
		sessionCookieName:   sessionCookie,
		sessionTTL:          ttl,
		sessionCookieSecure: cfg.AdminSessionSecure,
	}

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

	// Admin auth
	auth := e.Group("/admin/auth")
	auth.POST("/login", h.AdminLogin)
	auth.POST("/logout", h.AdminLogout)
	auth.GET("/me", h.AdminMe)

	// Admin protected routes
	admin := e.Group("/admin", h.RequireAdmin)
	admin.POST("/mappings/repo-project", h.AdminRepoProject)
	admin.GET("/mappings/repo-project", h.AdminRepoProjectList)
	admin.POST("/mappings/pr-states", h.AdminPRStates)
	admin.GET("/mappings/users", h.AdminUsersList)
	admin.POST("/mappings/users", h.AdminUsers)
	admin.POST("/mappings/labels", h.AdminLabels)
	admin.POST("/mappings/channel-project", h.AdminChannelProject)

	access := admin.Group("/access")
	access.GET("/users", h.AdminAccessList)
	access.POST("/users", h.AdminAccessCreate)
	access.PATCH("/users/:id", h.AdminAccessUpdate)
	access.POST("/users/:id/reset-password", h.AdminAccessResetPassword)

	// Jobs
	e.POST("/jobs/issue-summary/daily", h.JobIssueSummaryDaily)
}

type Handler struct {
	cfg                 config.Config
	dedupe              *Deduper
	db                  *store.DB
	sessionCookieName   string
	sessionTTL          time.Duration
	sessionCookieSecure bool
	createLocks         sync.Map // key: repo|planeIssueID -> *sync.Mutex
}
