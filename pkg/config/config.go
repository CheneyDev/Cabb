package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        int
	DatabaseURL string
	Timezone    string

	// Plane (Webhook-only)
	PlaneBaseURL       string
	PlaneWebhookSecret string
	PlaneServiceToken  string // Global Service Token for outbound API calls

    // Feishu (Lark)
    LarkAppID             string
    LarkAppSecret         string
    LarkEncryptKey        string
    LarkVerificationToken string

	// CNB
	CNBAppToken        string
	IntegrationToken   string
	CNBBaseURL         string
	CNBOutboundEnabled bool

	// Optional CNB path overrides
	CNBIssueCreatePath  string
	CNBIssueUpdatePath  string
	CNBIssueCommentPath string

	// Crypto
	EncryptionKey string

	// Optional: redirect backend root to frontend
	FrontendBaseURL string

	// Admin console auth
	AdminSessionCookie     string
	AdminSessionTTLHours   int
	AdminSessionSecure     bool
	AdminBootstrapEmail    string
	AdminBootstrapPassword string
	AdminBootstrapName     string

	// Jobs / Cleanup
	CleanupThreadLinksEnabled bool
	CleanupThreadLinksDays    int
	CleanupThreadLinksAt      string // HH:MM in cfg.Timezone
}

func FromEnv() Config {
    cfg := Config{
		Port:        intFromEnv("PORT", 8080),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Timezone:    strFromEnv("TIMEZONE", "Local"),

		PlaneBaseURL:       strFromEnv("PLANE_BASE_URL", "https://api.plane.so"),
		PlaneWebhookSecret: os.Getenv("PLANE_WEBHOOK_SECRET"),
		PlaneServiceToken:  os.Getenv("PLANE_SERVICE_TOKEN"),

        LarkAppID:             os.Getenv("LARK_APP_ID"),
        LarkAppSecret:         os.Getenv("LARK_APP_SECRET"),
        LarkEncryptKey:        os.Getenv("LARK_ENCRYPT_KEY"),
        LarkVerificationToken: os.Getenv("LARK_VERIFICATION_TOKEN"),

		CNBAppToken:        os.Getenv("CNB_APP_TOKEN"),
		IntegrationToken:   os.Getenv("INTEGRATION_TOKEN"),
		CNBBaseURL:         os.Getenv("CNB_BASE_URL"),
		CNBOutboundEnabled: boolFromEnv("CNB_OUTBOUND_ENABLED", true),

		CNBIssueCreatePath:  os.Getenv("CNB_ISSUE_CREATE_PATH"),
		CNBIssueUpdatePath:  os.Getenv("CNB_ISSUE_UPDATE_PATH"),
		CNBIssueCommentPath: os.Getenv("CNB_ISSUE_COMMENT_PATH"),

		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),

		FrontendBaseURL: os.Getenv("FRONTEND_BASE_URL"),

		AdminSessionCookie:     strFromEnv("ADMIN_SESSION_COOKIE", "pi_admin_session"),
		AdminSessionTTLHours:   intFromEnv("ADMIN_SESSION_TTL_HOURS", 12),
		AdminSessionSecure:     boolFromEnv("ADMIN_SESSION_SECURE", false),
		AdminBootstrapEmail:    os.Getenv("ADMIN_BOOTSTRAP_EMAIL"),
		AdminBootstrapPassword: os.Getenv("ADMIN_BOOTSTRAP_PASSWORD"),
		AdminBootstrapName:     strFromEnv("ADMIN_BOOTSTRAP_NAME", "Plane Admin"),

		CleanupThreadLinksEnabled: boolFromEnv("CLEANUP_THREAD_LINKS_ENABLED", true),
		CleanupThreadLinksDays:    intFromEnv("CLEANUP_THREAD_LINKS_DAYS", 90),
		CleanupThreadLinksAt:      strFromEnv("CLEANUP_THREAD_LINKS_AT", "03:00"),
	}
	return cfg
}

func (c Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func intFromEnv(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return def
}

func strFromEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func boolFromEnv(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "t", "true", "yes", "y", "on":
		return true
	case "0", "f", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}
