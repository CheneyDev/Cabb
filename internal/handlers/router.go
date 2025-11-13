package handlers

import (
    "context"
    "sync"
    "time"

    groqp "cabb/internal/ai/providers/groq"
    openaip "cabb/internal/ai/providers/openai"
    "cabb/internal/store"
    "cabb/pkg/config"

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

    // Optional AI namer
    varNamer := initBranchNamer(cfg)

    h := &Handler{
        cfg:                 cfg,
        dedupe:              d,
        db:                  db,
        sessionCookieName:   sessionCookie,
        sessionTTL:          ttl,
        sessionCookieSecure: cfg.AdminSessionSecure,
        aiNamer:             varNamer,
    }

	// Health
	e.GET("/healthz", h.Healthz)
	// Root redirect to frontend (optional)
	e.GET("/", h.Root)

	// Plane Webhook (OAuth fully removed per webhook-only refactor)
	e.POST("/webhooks/plane", h.PlaneWebhook)

	// CNB ingest callbacks via .cnb.yml
	e.POST("/ingest/cnb/issue", h.CNBIngestIssue)
	e.POST("/ingest/cnb/pr", h.CNBIngestPR)
	e.POST("/ingest/cnb/branch", h.CNBIngestBranch)

	// CNB API v1
	api := e.Group("/api/v1")
	api.POST("/issues/label-notify", h.IssueLabelNotify)     // 完整版（11 个字段）
	api.POST("/issues/label-sync", h.IssueLabelNotifySimple) // 简化版（3 个字段）

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
	admin.POST("/mappings", h.AdminMappings)
	admin.GET("/mappings", h.AdminMappingsList)
	admin.POST("/mappings/channel-project", h.AdminChannelProject)
	admin.GET("/links/issues", h.AdminIssueLinksList)
	admin.POST("/links/issues", h.AdminIssueLinksUpsert)
	admin.DELETE("/links/issues", h.AdminIssueLinksDelete)
    admin.GET("/links/lark-threads", h.AdminLarkThreadLinksList)
    admin.POST("/links/lark-threads", h.AdminLarkThreadLinksUpsert)
    admin.DELETE("/links/lark-threads", h.AdminLarkThreadLinksDelete)
    admin.GET("/links/branches", h.AdminBranchIssueLinksList)
	
	// Plane data APIs
	admin.GET("/plane/workspaces", h.AdminPlaneWorkspaces)
	admin.GET("/plane/projects", h.AdminPlaneProjects)

	access := admin.Group("/access")
	access.GET("/users", h.AdminAccessList)
	access.POST("/users", h.AdminAccessCreate)
	access.PATCH("/users/:id", h.AdminAccessUpdate)
	access.POST("/users/:id/reset-password", h.AdminAccessResetPassword)

	// Jobs
	e.POST("/jobs/issue-summary/daily", h.JobIssueSummaryDaily)
	e.POST("/jobs/cleanup/thread-links", h.JobCleanupThreadLinks)
}

type Handler struct {
    cfg                 config.Config
    dedupe              *Deduper
    db                  *store.DB
    sessionCookieName   string
    sessionTTL          time.Duration
    sessionCookieSecure bool
    createLocks         sync.Map // key: repo|planeIssueID -> *sync.Mutex
    aiNamer             interface{ SuggestBranchName(ctx context.Context, title, description string) (string, string, error) }
}

// initBranchNamer wires an AI-based branch namer if enabled and credentials are present.
func initBranchNamer(cfg config.Config) interface{ SuggestBranchName(ctx context.Context, title, description string) (string, string, error) } {
    // Prefer Groq when configured
    if cfg.GroqAPIKey != "" {
        return groqp.New(cfg.GroqModel, cfg.GroqAPIKey, cfg.GroqBaseURL)
    }
    // Fallback to OpenAI if configured
    if cfg.OpenAIAPIKey != "" {
        return openaip.New(cfg.OpenAIModel, cfg.OpenAIAPIKey, cfg.OpenAIBaseURL)
    }
    return nil
}
