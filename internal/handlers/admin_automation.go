package handlers

import (
	"net/http"
	"strings"

	"cabb/internal/store"

	"github.com/labstack/echo/v4"
)

type automationPayload struct {
	TargetRepoURL    string   `json:"target_repo_url"`
	TargetRepoBranch string   `json:"target_repo_branch"`
	PlaneStatuses    []string `json:"plane_statuses"`
	OutputRepoURL    string   `json:"output_repo_url"`
	OutputBranch     string   `json:"output_branch"`
	OutputDir        string   `json:"output_dir"`
	ReportRepoSlug   string   `json:"report_repo_slug"`
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
			OutputDir:        "issue-progress",
			ReportRepoSlug:   "1024hub/plane-test",
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
	cfg := automationPayload{
		TargetRepoURL:    strings.TrimSpace(payload.TargetRepoURL),
		TargetRepoBranch: strings.TrimSpace(payload.TargetRepoBranch),
		PlaneStatuses:    sanitizeStatuses(payload.PlaneStatuses),
		OutputRepoURL:    strings.TrimSpace(payload.OutputRepoURL),
		OutputBranch:     strings.TrimSpace(payload.OutputBranch),
		OutputDir:        strings.TrimSpace(payload.OutputDir),
		ReportRepoSlug:   strings.TrimSpace(payload.ReportRepoSlug),
	}
	if err := h.db.UpsertAutomationConfig(c.Request().Context(), store.AutomationConfig{
		TargetRepoURL:    cfg.TargetRepoURL,
		TargetRepoBranch: cfg.TargetRepoBranch,
		PlaneStatuses:    cfg.PlaneStatuses,
		OutputRepoURL:    cfg.OutputRepoURL,
		OutputBranch:     cfg.OutputBranch,
		OutputDir:        cfg.OutputDir,
		ReportRepoSlug:   cfg.ReportRepoSlug,
	}); err != nil {
		return writeError(c, http.StatusInternalServerError, "db_error", "保存自动化配置失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"saved": true})
}
