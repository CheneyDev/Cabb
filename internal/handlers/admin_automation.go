package handlers

import (
	"net/http"
	"path"
	"strings"

	"cabb/internal/store"

	"github.com/labstack/echo/v4"
)

type automationPayload struct {
	TargetRepoURL    string                   `json:"target_repo_url"`
	TargetRepoBranch string                   `json:"target_repo_branch"`
	PlaneStatuses    []string                 `json:"plane_statuses"`
	OutputRepoURL    string                   `json:"output_repo_url"`
	OutputBranch     string                   `json:"output_branch"`
	OutputDir        string                   `json:"output_dir"`
	ReportRepoSlug   string                   `json:"report_repo_slug"`
	ReportRepos      []store.ReportRepoConfig `json:"report_repos"`
}

func sanitizeStatuses(arr []string) []string {
	out := []string{}
	for _, s := range arr {
		v := strings.TrimSpace(s)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func sanitizeReportRepos(arr []store.ReportRepoConfig) []store.ReportRepoConfig {
	out := []store.ReportRepoConfig{}
	for _, r := range arr {
		slug := strings.TrimSpace(r.Slug)
		if slug == "" && r.RepoURL != "" {
			slug = path.Base(strings.TrimSuffix(strings.TrimSpace(r.RepoURL), "/"))
		}
		if slug == "" {
			continue
		}
		out = append(out, store.ReportRepoConfig{
			RepoURL:     strings.TrimSpace(r.RepoURL),
			Branch:      strings.TrimSpace(r.Branch),
			Slug:        slug,
			DisplayName: strings.TrimSpace(r.DisplayName),
		})
	}
	return out
}

// AdminAutomationGet returns the global automation configuration.
func (h *Handler) AdminAutomationGet(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库不可用", nil)
	}
	ctx := c.Request().Context()
	cfg, err := h.db.GetAutomationConfig(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "db_error", "查询自动化配置失败", map[string]any{"error": err.Error()})
	}
	if cfg == nil {
		// Defaults
		return c.JSON(http.StatusOK, automationPayload{
			TargetRepoURL:    "",
			TargetRepoBranch: "main",
			PlaneStatuses:    []string{},
			OutputRepoURL:    "https://cnb.cool/1024hub/plane-test",
			OutputBranch:     "main",
			OutputDir:        "ai-report",
			ReportRepoSlug:   "1024hub/plane-test",
			ReportRepos:      []store.ReportRepoConfig{},
		})
	}
	return c.JSON(http.StatusOK, automationPayload{
		TargetRepoURL:    cfg.TargetRepoURL,
		TargetRepoBranch: cfg.TargetRepoBranch,
		PlaneStatuses:    cfg.PlaneStatuses,
		OutputRepoURL:    cfg.OutputRepoURL,
		OutputBranch:     cfg.OutputBranch,
		OutputDir:        cfg.OutputDir,
		ReportRepoSlug:   cfg.ReportRepoSlug,
		ReportRepos:      cfg.ReportRepos,
	})
}

// AdminAutomationSave upserts the global automation configuration.
func (h *Handler) AdminAutomationSave(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库不可用", nil)
	}
	var payload automationPayload
	if err := c.Bind(&payload); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "请求体无法解析", nil)
	}
	LogStructured("info", map[string]any{
		"event":          "admin.automation.save.begin",
		"report_repos":   len(payload.ReportRepos),
		"target_repo":    strings.TrimSpace(payload.TargetRepoURL),
		"output_repo":    strings.TrimSpace(payload.OutputRepoURL),
		"report_repo":    strings.TrimSpace(payload.ReportRepoSlug),
		"target_branch":  strings.TrimSpace(payload.TargetRepoBranch),
		"output_branch":  strings.TrimSpace(payload.OutputBranch),
		"plane_statuses": len(payload.PlaneStatuses),
	})
	cfg := automationPayload{
		TargetRepoURL:    strings.TrimSpace(payload.TargetRepoURL),
		TargetRepoBranch: strings.TrimSpace(payload.TargetRepoBranch),
		PlaneStatuses:    sanitizeStatuses(payload.PlaneStatuses),
		OutputRepoURL:    strings.TrimSpace(payload.OutputRepoURL),
		OutputBranch:     strings.TrimSpace(payload.OutputBranch),
		OutputDir:        strings.TrimSpace(payload.OutputDir),
		ReportRepoSlug:   strings.TrimSpace(payload.ReportRepoSlug),
		ReportRepos:      sanitizeReportRepos(payload.ReportRepos),
	}
	if err := h.db.UpsertAutomationConfig(c.Request().Context(), store.AutomationConfig{
		TargetRepoURL:    cfg.TargetRepoURL,
		TargetRepoBranch: cfg.TargetRepoBranch,
		PlaneStatuses:    cfg.PlaneStatuses,
		OutputRepoURL:    cfg.OutputRepoURL,
		OutputBranch:     cfg.OutputBranch,
		OutputDir:        cfg.OutputDir,
		ReportRepoSlug:   cfg.ReportRepoSlug,
		ReportRepos:      cfg.ReportRepos,
	}); err != nil {
		LogStructured("error", map[string]any{
			"event": "admin.automation.save.failed",
			"error": err.Error(),
		})
		return writeError(c, http.StatusInternalServerError, "db_error", "保存自动化配置失败", map[string]any{"error": err.Error()})
	}
	LogStructured("info", map[string]any{
		"event":        "admin.automation.save.success",
		"report_repos": len(cfg.ReportRepos),
		"target_repo":  cfg.TargetRepoURL,
		"output_repo":  cfg.OutputRepoURL,
	})
	return c.JSON(http.StatusOK, map[string]any{"saved": true})
}
