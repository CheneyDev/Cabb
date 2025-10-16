package handlers

import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"
    "plane-integration/internal/store"
)

// POST /admin/mappings/repo-project
// Body: { "cnb_repo_id": "group/repo", "plane_workspace_id": "uuid", "plane_project_id": "uuid", "issue_open_state_id": "uuid?", "issue_closed_state_id": "uuid?", "active": true, "sync_direction": "cnb_to_plane|bidirectional", "label_selector": "后端,backend" }
func (h *Handler) AdminRepoProject(c echo.Context) error {
    if !hHasDB(h) {
        return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
    }
    var req struct{
        CNBRepoID string `json:"cnb_repo_id"`
        PlaneWorkspaceID string `json:"plane_workspace_id"`
        PlaneProjectID string `json:"plane_project_id"`
        IssueOpenStateID string `json:"issue_open_state_id"`
        IssueClosedStateID string `json:"issue_closed_state_id"`
        Active *bool `json:"active"`
        SyncDirection string `json:"sync_direction"`
        LabelSelector string `json:"label_selector"`
    }
    if err := c.Bind(&req); err != nil { return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil) }
    if req.CNBRepoID == "" || req.PlaneWorkspaceID == "" || req.PlaneProjectID == "" {
        return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_repo_id/plane_workspace_id/plane_project_id", nil)
    }
    m := store.RepoProjectMapping{
        PlaneProjectID: req.PlaneProjectID,
        PlaneWorkspaceID: req.PlaneWorkspaceID,
        CNBRepoID: req.CNBRepoID,
        IssueOpenStateID: sql.NullString{String: req.IssueOpenStateID, Valid: req.IssueOpenStateID != ""},
        IssueClosedStateID: sql.NullString{String: req.IssueClosedStateID, Valid: req.IssueClosedStateID != ""},
        Active: true,
        SyncDirection: sql.NullString{String: req.SyncDirection, Valid: req.SyncDirection != ""},
        LabelSelector: sql.NullString{String: req.LabelSelector, Valid: req.LabelSelector != ""},
    }
    if req.Active != nil { m.Active = *req.Active }
    if err := h.db.UpsertRepoProjectMapping(c.Request().Context(), m); err != nil {
        return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]any{"result":"ok"})
}

// GET /admin/mappings/repo-project
// Query params (optional): plane_project_id=uuid, cnb_repo_id=org/repo, active=true|false
// Response: { "items": [ { "plane_project_id": "...", "plane_workspace_id": "...", "cnb_repo_id": "...", "issue_open_state_id": "...", "issue_closed_state_id": "...", "active": true, "sync_direction": "...", "label_selector": "..." } ] }
func (h *Handler) AdminRepoProjectList(c echo.Context) error {
    if !hHasDB(h) {
        return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
    }
    planeProjectID := c.QueryParam("plane_project_id")
    cnbRepoID := c.QueryParam("cnb_repo_id")
    activeParam := c.QueryParam("active")

    items, err := h.db.ListRepoProjectMappings(c.Request().Context(), planeProjectID, cnbRepoID, activeParam)
    if err != nil {
        return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
    }
    // Normalize to JSON-friendly map with snake_case keys
    out := make([]map[string]any, 0, len(items))
    for _, m := range items {
        out = append(out, map[string]any{
            "plane_project_id":     m.PlaneProjectID,
            "plane_workspace_id":   m.PlaneWorkspaceID,
            "cnb_repo_id":          m.CNBRepoID,
            "issue_open_state_id":  nullString(m.IssueOpenStateID),
            "issue_closed_state_id": nullString(m.IssueClosedStateID),
            "active":               m.Active,
            "sync_direction":       nullString(m.SyncDirection),
            "label_selector":       nullString(m.LabelSelector),
        })
    }
    return c.JSON(http.StatusOK, map[string]any{"items": out})
}

func nullString(s sql.NullString) any {
    if s.Valid && s.String != "" { return s.String }
    return nil
}

// POST /admin/mappings/pr-states
// Body: { "cnb_repo_id":"group/repo", "plane_project_id":"uuid", ...state ids... }
func (h *Handler) AdminPRStates(c echo.Context) error {
    if !hHasDB(h) { return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil) }
    var req struct{
        CNBRepoID string `json:"cnb_repo_id"`
        PlaneProjectID string `json:"plane_project_id"`
        DraftStateID string `json:"draft_state_id"`
        OpenedStateID string `json:"opened_state_id"`
        ReviewRequestedStateID string `json:"review_requested_state_id"`
        ApprovedStateID string `json:"approved_state_id"`
        MergedStateID string `json:"merged_state_id"`
        ClosedStateID string `json:"closed_state_id"`
    }
    if err := c.Bind(&req); err != nil { return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil) }
    if req.CNBRepoID == "" || req.PlaneProjectID == "" { return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_repo_id/plane_project_id", nil) }
    m := store.PRStateMapping{
        PlaneProjectID: req.PlaneProjectID,
        CNBRepoID: req.CNBRepoID,
        DraftStateID: sql.NullString{String: req.DraftStateID, Valid: req.DraftStateID != ""},
        OpenedStateID: sql.NullString{String: req.OpenedStateID, Valid: req.OpenedStateID != ""},
        ReviewRequestedStateID: sql.NullString{String: req.ReviewRequestedStateID, Valid: req.ReviewRequestedStateID != ""},
        ApprovedStateID: sql.NullString{String: req.ApprovedStateID, Valid: req.ApprovedStateID != ""},
        MergedStateID: sql.NullString{String: req.MergedStateID, Valid: req.MergedStateID != ""},
        ClosedStateID: sql.NullString{String: req.ClosedStateID, Valid: req.ClosedStateID != ""},
    }
    if err := h.db.UpsertPRStateMapping(c.Request().Context(), m); err != nil {
        return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]any{"result":"ok"})
}

func (h *Handler) AdminUsers(c echo.Context) error {
    if !hHasDB(h) { return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil) }
    var req struct{
        Mappings []struct{
            CNBUserID string `json:"cnb_user_id"`
            PlaneUserID string `json:"plane_user_id"`
            DisplayName string `json:"display_name"`
        } `json:"mappings"`
    }
    if err := c.Bind(&req); err != nil { return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil) }
    for _, m := range req.Mappings {
        if m.CNBUserID == "" || m.PlaneUserID == "" { return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_user_id/plane_user_id", nil) }
        if err := h.db.UpsertUserMapping(c.Request().Context(), m.PlaneUserID, m.CNBUserID, m.DisplayName); err != nil {
            return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
        }
    }
    return c.JSON(http.StatusOK, map[string]any{"result":"ok","count":len(req.Mappings)})
}

func (h *Handler) AdminChannelProject(c echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"message": "channel-project mapping API not implemented in scaffold"})
}

// POST /admin/mappings/labels
// Body: { "cnb_repo_id":"group/repo", "plane_project_id":"uuid", "items":[{"cnb_label":"bug","plane_label_id":"uuid"}]}
func (h *Handler) AdminLabels(c echo.Context) error {
    if !hHasDB(h) { return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil) }
    var req struct{
        CNBRepoID string `json:"cnb_repo_id"`
        PlaneProjectID string `json:"plane_project_id"`
        Items []struct{
            CNBLabel string `json:"cnb_label"`
            PlaneLabelID string `json:"plane_label_id"`
        } `json:"items"`
    }
    if err := c.Bind(&req); err != nil { return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil) }
    if req.CNBRepoID == "" || req.PlaneProjectID == "" { return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_repo_id/plane_project_id", nil) }
    for _, it := range req.Items {
        if it.CNBLabel == "" || it.PlaneLabelID == "" { return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_label/plane_label_id", nil) }
        if err := h.db.UpsertLabelMapping(c.Request().Context(), req.PlaneProjectID, req.CNBRepoID, it.CNBLabel, it.PlaneLabelID); err != nil {
            return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
        }
    }
    return c.JSON(http.StatusOK, map[string]any{"result":"ok","count":len(req.Items)})
}
