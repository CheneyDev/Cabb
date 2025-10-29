package handlers

import (
    "context"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/labstack/echo/v4"

    "cabb/internal/store"
    "cabb/pkg/config"
)

func (h *Handler) JobIssueSummaryDaily(c echo.Context) error {
    return c.JSON(http.StatusAccepted, map[string]any{
        "scheduled": true,
        "job":       "issue-summary-daily",
    })
}

// JobCleanupThreadLinks deletes non-sync-enabled thread links older than N days (default 90; min 7)
// Invoke via: POST /jobs/cleanup/thread-links?days=90
func (h *Handler) JobCleanupThreadLinks(c echo.Context) error {
    if !hHasDB(h) {
        return c.JSON(http.StatusOK, map[string]any{"deleted": 0, "cutoff": "", "note": "db unavailable"})
    }
    days := 90
    if v := c.QueryParam("days"); v != "" {
        if n, err := strconv.Atoi(v); err == nil { days = n }
    }
    if days < 7 { days = 7 }
    cutoff := time.Now().AddDate(0, 0, -days)
    ctx, cancel := context.WithTimeout(c.Request().Context(), 12*time.Second)
    defer cancel()
    n, err := h.db.CleanupStaleThreadLinks(ctx, cutoff)
    if err != nil {
        return writeError(c, http.StatusInternalServerError, "cleanup_failed", "清理失败", map[string]any{"error": err.Error()})
    }
    LogStructured("info", map[string]any{"event": "jobs.cleanup.thread_links", "deleted": n, "cutoff": cutoff.UTC().Format(time.RFC3339)})
    return c.JSON(http.StatusOK, map[string]any{
        "deleted": n,
        "cutoff":  cutoff.UTC().Format(time.RFC3339),
        "result":  "ok",
    })
}

// StartCleanupScheduler starts a daily scheduler for cleaning stale thread links.
// It respects cfg.CleanupThreadLinksEnabled, cfg.CleanupThreadLinksDays, cfg.CleanupThreadLinksAt and cfg.Timezone.
func StartCleanupScheduler(cfg config.Config, db *store.DB) {
    if !cfg.CleanupThreadLinksEnabled || db == nil || db.SQL == nil { return }
    go func() {
        // Determine timezone
        loc := time.Local
        if tz := strings.TrimSpace(cfg.Timezone); tz != "" {
            if l, err := time.LoadLocation(tz); err == nil { loc = l }
        }
        // Parse time of day
        hh, mm := parseHHMM(cfg.CleanupThreadLinksAt)
        // Compute first run
        next := nextAt(loc, hh, mm)
        for {
            sleep := time.Until(next)
            if sleep <= 0 { next = next.Add(24 * time.Hour); continue }
            time.Sleep(sleep)
            // Run cleanup
            days := cfg.CleanupThreadLinksDays
            if days < 7 { days = 7 }
            cutoff := time.Now().In(loc).AddDate(0, 0, -days)
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            n, err := db.CleanupStaleThreadLinks(ctx, cutoff)
            cancel()
            if err != nil {
                LogStructured("error", map[string]any{"event": "jobs.cleanup.thread_links.cron", "error": map[string]any{"code": "cleanup_failed", "message": truncate(err.Error(), 200)}, "cutoff": cutoff.UTC().Format(time.RFC3339)})
            } else {
                LogStructured("info", map[string]any{"event": "jobs.cleanup.thread_links.cron", "deleted": n, "cutoff": cutoff.UTC().Format(time.RFC3339)})
            }
            next = next.Add(24 * time.Hour)
        }
    }()
}

func parseHHMM(s string) (int, int) {
    s = strings.TrimSpace(s)
    if s == "" { return 3, 0 }
    parts := strings.SplitN(s, ":", 2)
    if len(parts) != 2 { return 3, 0 }
    h, _ := strconv.Atoi(parts[0])
    m, _ := strconv.Atoi(parts[1])
    if h < 0 || h > 23 { h = 3 }
    if m < 0 || m > 59 { m = 0 }
    return h, m
}

func nextAt(loc *time.Location, hour, min int) time.Time {
    now := time.Now().In(loc)
    run := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, loc)
    if !run.After(now) {
        run = run.Add(24 * time.Hour)
    }
    return run
}
