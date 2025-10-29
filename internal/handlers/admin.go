package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"cabb/internal/lark"
	"cabb/internal/plane"
	"cabb/internal/store"
)

// POST /admin/mappings/repo-project
// Body: { "cnb_repo_id": "group/repo", "plane_workspace_id": "uuid", "plane_project_id": "uuid", "issue_open_state_id": "uuid?", "issue_closed_state_id": "uuid?", "active": true, "sync_direction": "cnb_to_plane|bidirectional", "label_selector": "后端,backend" }
func (h *Handler) AdminRepoProject(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		CNBRepoID          string `json:"cnb_repo_id"`
		PlaneWorkspaceID   string `json:"plane_workspace_id"`
		PlaneProjectID     string `json:"plane_project_id"`
		IssueOpenStateID   string `json:"issue_open_state_id"`
		IssueClosedStateID string `json:"issue_closed_state_id"`
		Active             *bool  `json:"active"`
		SyncDirection      string `json:"sync_direction"`
		LabelSelector      string `json:"label_selector"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	if req.CNBRepoID == "" || req.PlaneWorkspaceID == "" || req.PlaneProjectID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_repo_id/plane_workspace_id/plane_project_id", nil)
	}
	m := store.RepoProjectMapping{
		PlaneProjectID:     req.PlaneProjectID,
		PlaneWorkspaceID:   req.PlaneWorkspaceID,
		CNBRepoID:          req.CNBRepoID,
		IssueOpenStateID:   sql.NullString{String: req.IssueOpenStateID, Valid: req.IssueOpenStateID != ""},
		IssueClosedStateID: sql.NullString{String: req.IssueClosedStateID, Valid: req.IssueClosedStateID != ""},
		Active:             true,
		SyncDirection:      sql.NullString{String: req.SyncDirection, Valid: req.SyncDirection != ""},
		LabelSelector:      sql.NullString{String: req.LabelSelector, Valid: req.LabelSelector != ""},
	}
	if req.Active != nil {
		m.Active = *req.Active
	}
	if err := h.db.UpsertRepoProjectMapping(c.Request().Context(), m); err != nil {
		return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok"})
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

	ctx := c.Request().Context()
	items, err := h.db.ListRepoProjectMappings(ctx, planeProjectID, cnbRepoID, activeParam)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
	}
	// Enrich workspace / project metadata when possible
	type workspaceToken struct {
		token string
		slug  string
	}
	type workspaceMeta struct {
		name string
		slug string
	}
	type projectMeta struct {
		name       string
		identifier string
		slug       string
	}

	tokens := make(map[string]workspaceToken, len(items))
	for _, m := range items {
		if _, exists := tokens[m.PlaneWorkspaceID]; exists {
			continue
		}
		accessToken, slug, err := h.db.FindBotTokenByWorkspaceID(ctx, m.PlaneWorkspaceID)
		if err != nil {
			tokens[m.PlaneWorkspaceID] = workspaceToken{}
			continue
		}
		tokens[m.PlaneWorkspaceID] = workspaceToken{token: accessToken, slug: strings.TrimSpace(slug)}
	}

	planeClient := plane.Client{BaseURL: h.cfg.PlaneBaseURL}
	workspaceInfos := make(map[string]workspaceMeta, len(tokens))
	for workspaceID, tk := range tokens {
		info := workspaceMeta{}
		if tk.slug != "" {
			info.slug = tk.slug
		}
		if tk.token != "" && info.slug != "" {
			wctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			ws, err := planeClient.GetWorkspace(wctx, tk.token, info.slug)
			cancel()
			if err == nil && ws != nil {
				if strings.TrimSpace(ws.Name) != "" {
					info.name = strings.TrimSpace(ws.Name)
				} else if strings.TrimSpace(ws.Title) != "" {
					info.name = strings.TrimSpace(ws.Title)
				}
				if strings.TrimSpace(ws.Slug) != "" {
					info.slug = strings.TrimSpace(ws.Slug)
				}
			}
		}
		workspaceInfos[workspaceID] = info
	}

	projectInfos := make(map[string]projectMeta, len(items))
	for _, m := range items {
		key := m.PlaneWorkspaceID + "::" + m.PlaneProjectID
		if _, exists := projectInfos[key]; exists {
			continue
		}
		wm := workspaceInfos[m.PlaneWorkspaceID]
		tk := tokens[m.PlaneWorkspaceID]
		meta := projectMeta{}
		if wm.slug != "" {
			meta.slug = wm.slug
		}
		if tk.token != "" && wm.slug != "" {
			pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			proj, err := planeClient.GetProject(pctx, tk.token, wm.slug, m.PlaneProjectID)
			cancel()
			if err == nil && proj != nil {
				if strings.TrimSpace(proj.Name) != "" {
					meta.name = strings.TrimSpace(proj.Name)
				} else if strings.TrimSpace(proj.Slug) != "" {
					meta.name = strings.TrimSpace(proj.Slug)
				}
				if strings.TrimSpace(proj.Identifier) != "" {
					meta.identifier = strings.TrimSpace(proj.Identifier)
				}
				if strings.TrimSpace(proj.Slug) != "" {
					meta.slug = strings.TrimSpace(proj.Slug)
				}
			}
		}
		projectInfos[key] = meta
	}

	// Normalize to JSON-friendly map with snake_case keys
	out := make([]map[string]any, 0, len(items))
	for _, m := range items {
		wm := workspaceInfos[m.PlaneWorkspaceID]
		pm := projectInfos[m.PlaneWorkspaceID+"::"+m.PlaneProjectID]
		out = append(out, map[string]any{
			"plane_project_id":         m.PlaneProjectID,
			"plane_workspace_id":       m.PlaneWorkspaceID,
			"plane_workspace_slug":     optionalString(wm.slug),
			"plane_workspace_name":     optionalString(wm.name),
			"plane_project_name":       optionalString(pm.name),
			"plane_project_identifier": optionalString(pm.identifier),
			"plane_project_slug":       optionalString(pm.slug),
			"cnb_repo_id":              m.CNBRepoID,
			"issue_open_state_id":      nullString(m.IssueOpenStateID),
			"issue_closed_state_id":    nullString(m.IssueClosedStateID),
			"active":                   m.Active,
			"sync_direction":           nullString(m.SyncDirection),
			"label_selector":           nullString(m.LabelSelector),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out})
}

func nullString(s sql.NullString) any {
	if s.Valid && s.String != "" {
		return s.String
	}
	return nil
}

func optionalString(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullTimeValue(t sql.NullTime) any {
	if t.Valid {
		return t.Time.UTC().Format(time.RFC3339)
	}
	return nil
}

// POST /admin/mappings/pr-states
// Body: { "cnb_repo_id":"group/repo", "plane_project_id":"uuid", ...state ids... }
func (h *Handler) AdminPRStates(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		CNBRepoID              string `json:"cnb_repo_id"`
		PlaneProjectID         string `json:"plane_project_id"`
		DraftStateID           string `json:"draft_state_id"`
		OpenedStateID          string `json:"opened_state_id"`
		ReviewRequestedStateID string `json:"review_requested_state_id"`
		ApprovedStateID        string `json:"approved_state_id"`
		MergedStateID          string `json:"merged_state_id"`
		ClosedStateID          string `json:"closed_state_id"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	if req.CNBRepoID == "" || req.PlaneProjectID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_repo_id/plane_project_id", nil)
	}
	m := store.PRStateMapping{
		PlaneProjectID:         req.PlaneProjectID,
		CNBRepoID:              req.CNBRepoID,
		DraftStateID:           sql.NullString{String: req.DraftStateID, Valid: req.DraftStateID != ""},
		OpenedStateID:          sql.NullString{String: req.OpenedStateID, Valid: req.OpenedStateID != ""},
		ReviewRequestedStateID: sql.NullString{String: req.ReviewRequestedStateID, Valid: req.ReviewRequestedStateID != ""},
		ApprovedStateID:        sql.NullString{String: req.ApprovedStateID, Valid: req.ApprovedStateID != ""},
		MergedStateID:          sql.NullString{String: req.MergedStateID, Valid: req.MergedStateID != ""},
		ClosedStateID:          sql.NullString{String: req.ClosedStateID, Valid: req.ClosedStateID != ""},
	}
	if err := h.db.UpsertPRStateMapping(c.Request().Context(), m); err != nil {
		return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

func (h *Handler) AdminUsersList(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	limit := 50
	if lp := c.QueryParam("limit"); lp != "" {
		if v, err := strconv.Atoi(lp); err == nil {
			limit = v
		}
	}
	planeUserID := c.QueryParam("plane_user_id")
	cnbUserID := c.QueryParam("cnb_user_id")
	search := c.QueryParam("q")
	items, err := h.db.ListUserMappings(c.Request().Context(), planeUserID, cnbUserID, search, limit)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
	}
	out := make([]map[string]any, 0, len(items))
	for _, m := range items {
		out = append(out, map[string]any{
			"plane_user_id": m.PlaneUserID,
			"cnb_user_id":   nullString(m.CNBUserID),
			"lark_user_id":  nullString(m.LarkUserID),
			"display_name":  nullString(m.DisplayName),
			"connected_at":  nullTimeValue(m.ConnectedAt),
			"created_at":    m.CreatedAt.UTC().Format(time.RFC3339),
			"updated_at":    m.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "count": len(out)})
}

func (h *Handler) AdminUsers(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		Mappings []struct {
			CNBUserID   string `json:"cnb_user_id"`
			PlaneUserID string `json:"plane_user_id"`
			DisplayName string `json:"display_name"`
		} `json:"mappings"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	for _, m := range req.Mappings {
		if m.CNBUserID == "" || m.PlaneUserID == "" {
			return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_user_id/plane_user_id", nil)
		}
		if err := h.db.UpsertUserMapping(c.Request().Context(), m.PlaneUserID, m.CNBUserID, m.DisplayName); err != nil {
			return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok", "count": len(req.Mappings)})
}

func (h *Handler) AdminChannelProject(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"message": "channel-project mapping API not implemented in scaffold"})
}

// POST /admin/mappings/labels
// Body: { "cnb_repo_id":"group/repo", "plane_project_id":"uuid", "items":[{"cnb_label":"bug","plane_label_id":"uuid"}]}
func (h *Handler) AdminLabels(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		CNBRepoID      string `json:"cnb_repo_id"`
		PlaneProjectID string `json:"plane_project_id"`
		Items          []struct {
			CNBLabel     string `json:"cnb_label"`
			PlaneLabelID string `json:"plane_label_id"`
		} `json:"items"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	if req.CNBRepoID == "" || req.PlaneProjectID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_repo_id/plane_project_id", nil)
	}
	for _, it := range req.Items {
		if it.CNBLabel == "" || it.PlaneLabelID == "" {
			return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 cnb_label/plane_label_id", nil)
		}
		if err := h.db.UpsertLabelMapping(c.Request().Context(), req.PlaneProjectID, req.CNBRepoID, it.CNBLabel, it.PlaneLabelID); err != nil {
			return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok", "count": len(req.Items)})
}

// POST /admin/mappings
//
//	Body: {
//	  "scope_kind":"plane_project","scope_id":"<uuid>","mapping_type":"priority",
//	  "items":[{"left":{"system":"plane","type":"priority","key":"urgent"},"right":{"system":"cnb","type":"priority","key":"P0"},"bidirectional":true,"extras":{},"active":true}]
//	}
func (h *Handler) AdminMappings(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		ScopeKind   string `json:"scope_kind"`
		ScopeID     string `json:"scope_id"`
		MappingType string `json:"mapping_type"`
		Items       []struct {
			Left struct {
				System string `json:"system"`
				Type   string `json:"type"`
				Key    string `json:"key"`
			} `json:"left"`
			Right struct {
				System string `json:"system"`
				Type   string `json:"type"`
				Key    string `json:"key"`
			} `json:"right"`
			Bidirectional bool            `json:"bidirectional"`
			Extras        json.RawMessage `json:"extras"`
			Active        *bool           `json:"active"`
		} `json:"items"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	if strings.TrimSpace(req.ScopeKind) == "" || strings.TrimSpace(req.MappingType) == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 scope_kind/mapping_type", nil)
	}
	for _, it := range req.Items {
		if strings.TrimSpace(it.Left.System) == "" || strings.TrimSpace(it.Left.Type) == "" || strings.TrimSpace(it.Left.Key) == "" {
			return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 left.system/type/key", nil)
		}
		if strings.TrimSpace(it.Right.System) == "" || strings.TrimSpace(it.Right.Type) == "" || strings.TrimSpace(it.Right.Key) == "" {
			return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 right.system/type/key", nil)
		}
		// Parse extras JSON if provided
		var extras map[string]any
		if len(it.Extras) > 0 {
			_ = json.Unmarshal(it.Extras, &extras)
		}
		active := true
		if it.Active != nil {
			active = *it.Active
		}
		m := store.IntegrationMappingRec{
			ScopeKind:     req.ScopeKind,
			ScopeID:       req.ScopeID,
			MappingType:   req.MappingType,
			LeftSystem:    it.Left.System,
			LeftType:      it.Left.Type,
			LeftKey:       it.Left.Key,
			RightSystem:   it.Right.System,
			RightType:     it.Right.Type,
			RightKey:      it.Right.Key,
			Bidirectional: it.Bidirectional,
			Extras:        extras,
			Active:        active,
		}
		if err := h.db.UpsertIntegrationMapping(c.Request().Context(), m); err != nil {
			return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok", "count": len(req.Items)})
}

// GET /admin/mappings?scope_kind=plane_project&scope_id=<id>&mapping_type=priority
func (h *Handler) AdminMappingsList(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	scopeKind := c.QueryParam("scope_kind")
	scopeID := c.QueryParam("scope_id")
	mappingType := c.QueryParam("mapping_type")
	items, err := h.db.ListIntegrationMappings(c.Request().Context(), scopeKind, scopeID, mappingType)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
	}
	out := make([]map[string]any, 0, len(items))
	for _, m := range items {
		var extras any
		if m.Extras.Valid && m.Extras.String != "" {
			_ = json.Unmarshal([]byte(m.Extras.String), &extras)
		}
		out = append(out, map[string]any{
			"id":            m.ID,
			"scope_kind":    m.ScopeKind,
			"scope_id":      nullString(m.ScopeID),
			"mapping_type":  m.MappingType,
			"left":          map[string]any{"system": m.LeftSystem, "type": m.LeftType, "key": m.LeftKey},
			"right":         map[string]any{"system": m.RightSystem, "type": m.RightType, "key": m.RightKey},
			"bidirectional": m.Bidirectional,
			"extras":        extras,
			"active":        m.Active,
			"created_at":    m.CreatedAt.UTC().Format(time.RFC3339),
			"updated_at":    m.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "count": len(out)})
}

// GET /admin/links/issues
func (h *Handler) AdminIssueLinksList(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	limit := 50
	if l := strings.TrimSpace(c.QueryParam("limit")); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	planeIssueID := strings.TrimSpace(c.QueryParam("plane_issue_id"))
	cnbRepoID := strings.TrimSpace(c.QueryParam("cnb_repo_id"))
	cnbIssueID := strings.TrimSpace(c.QueryParam("cnb_issue_id"))
	ctx := c.Request().Context()
	items, err := h.db.ListIssueLinks(ctx, planeIssueID, cnbRepoID, cnbIssueID, limit)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
	}
	trimNull := func(ns sql.NullString) string {
		if ns.Valid {
			return strings.TrimSpace(ns.String)
		}
		return ""
	}
	planeClient := plane.Client{BaseURL: h.cfg.PlaneBaseURL}
	type workspaceToken struct {
		token string
		slug  string
	}
	tokens := make(map[string]workspaceToken, len(items))
	for _, it := range items {
		wsID := trimNull(it.PlaneWorkspaceID)
		if wsID == "" {
			continue
		}
		if _, exists := tokens[wsID]; exists {
			continue
		}
		accessToken, slug, terr := h.db.FindBotTokenByWorkspaceID(ctx, wsID)
		if terr != nil {
			tokens[wsID] = workspaceToken{}
			continue
		}
		tokens[wsID] = workspaceToken{token: strings.TrimSpace(accessToken), slug: strings.TrimSpace(slug)}
	}
	type workspaceMeta struct {
		name string
		slug string
	}
	workspaceInfos := make(map[string]workspaceMeta, len(tokens))
	for wsID, tk := range tokens {
		meta := workspaceMeta{slug: tk.slug}
		if tk.token != "" && meta.slug != "" {
			wctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			ws, werr := planeClient.GetWorkspace(wctx, tk.token, meta.slug)
			cancel()
			if werr == nil && ws != nil {
				name := strings.TrimSpace(ws.Name)
				if name == "" {
					name = strings.TrimSpace(ws.Title)
				}
				if name != "" {
					meta.name = name
				}
				if slug := strings.TrimSpace(ws.Slug); slug != "" {
					meta.slug = slug
				}
			}
		}
		workspaceInfos[wsID] = meta
	}
	type projectMeta struct {
		name       string
		identifier string
		slug       string
	}
	projectInfos := make(map[string]projectMeta, len(items))
	for _, it := range items {
		wsID := trimNull(it.PlaneWorkspaceID)
		projectID := trimNull(it.PlaneProjectID)
		if wsID == "" || projectID == "" {
			continue
		}
		key := wsID + "::" + projectID
		if _, exists := projectInfos[key]; exists {
			continue
		}
		meta := projectMeta{}
		tk := tokens[wsID]
		slug := workspaceInfos[wsID].slug
		if slug == "" {
			slug = tk.slug
		}
		if tk.token != "" && slug != "" {
			pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			proj, perr := planeClient.GetProject(pctx, tk.token, slug, projectID)
			cancel()
			if perr == nil && proj != nil {
				if name := strings.TrimSpace(proj.Name); name != "" {
					meta.name = name
				} else if alt := strings.TrimSpace(proj.Slug); alt != "" {
					meta.name = alt
				}
				if identifier := strings.TrimSpace(proj.Identifier); identifier != "" {
					meta.identifier = identifier
				}
				if pslug := strings.TrimSpace(proj.Slug); pslug != "" {
					meta.slug = pslug
				}
			}
		}
		projectInfos[key] = meta
	}
	issueNames := make(map[string]string, len(items))
	for _, it := range items {
		wsID := trimNull(it.PlaneWorkspaceID)
		projectID := trimNull(it.PlaneProjectID)
		issueID := strings.TrimSpace(it.PlaneIssueID)
		if wsID == "" || projectID == "" || issueID == "" {
			continue
		}
		key := wsID + "::" + projectID + "::" + issueID
		if _, exists := issueNames[key]; exists {
			continue
		}
		tk := tokens[wsID]
		slug := workspaceInfos[wsID].slug
		if slug == "" {
			slug = tk.slug
		}
		if tk.token != "" && slug != "" {
			ictx, cancel := context.WithTimeout(ctx, 5*time.Second)
			name, ierr := planeClient.GetIssueName(ictx, tk.token, slug, projectID, issueID)
			cancel()
			if ierr == nil {
				if trimmed := strings.TrimSpace(name); trimmed != "" {
					issueNames[key] = trimmed
					continue
				}
			}
		}
		issueNames[key] = ""
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		wsID := trimNull(it.PlaneWorkspaceID)
		projectID := trimNull(it.PlaneProjectID)
		issueID := strings.TrimSpace(it.PlaneIssueID)
		wm := workspaceInfos[wsID]
		pm := projectInfos[wsID+"::"+projectID]
		issueName := ""
		if wsID != "" && projectID != "" && issueID != "" {
			issueName = issueNames[wsID+"::"+projectID+"::"+issueID]
		}
		out = append(out, map[string]any{
			"plane_issue_id":           it.PlaneIssueID,
			"plane_issue_name":         optionalString(issueName),
			"cnb_repo_id":              nullString(it.CNBRepoID),
			"cnb_issue_id":             nullString(it.CNBIssueID),
			"plane_project_id":         nullString(it.PlaneProjectID),
			"plane_project_name":       optionalString(pm.name),
			"plane_project_identifier": optionalString(pm.identifier),
			"plane_project_slug":       optionalString(pm.slug),
			"plane_workspace_id":       nullString(it.PlaneWorkspaceID),
			"plane_workspace_name":     optionalString(wm.name),
			"plane_workspace_slug":     optionalString(wm.slug),
			"linked_at":                it.LinkedAt.UTC().Format(time.RFC3339),
			"created_at":               it.CreatedAt.UTC().Format(time.RFC3339),
			"updated_at":               it.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "count": len(out)})
}

// POST /admin/links/issues
func (h *Handler) AdminIssueLinksUpsert(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		PlaneIssueID string `json:"plane_issue_id"`
		CNBRepoID    string `json:"cnb_repo_id"`
		CNBIssueID   string `json:"cnb_issue_id"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	planeIssueID := strings.TrimSpace(req.PlaneIssueID)
	cnbRepoID := strings.TrimSpace(req.CNBRepoID)
	cnbIssueID := strings.TrimSpace(req.CNBIssueID)
	if planeIssueID == "" || cnbRepoID == "" || cnbIssueID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 plane_issue_id/cnb_repo_id/cnb_issue_id", nil)
	}
	inserted, err := h.db.CreateIssueLink(c.Request().Context(), planeIssueID, cnbRepoID, cnbIssueID)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
	}
	status := map[string]any{"result": "ok", "inserted": inserted}
	if !inserted {
		status["message"] = "记录已存在"
	}
	return c.JSON(http.StatusOK, status)
}

// DELETE /admin/links/issues
func (h *Handler) AdminIssueLinksDelete(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		PlaneIssueID string `json:"plane_issue_id"`
		CNBRepoID    string `json:"cnb_repo_id"`
		CNBIssueID   string `json:"cnb_issue_id"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	planeIssueID := strings.TrimSpace(req.PlaneIssueID)
	cnbRepoID := strings.TrimSpace(req.CNBRepoID)
	cnbIssueID := strings.TrimSpace(req.CNBIssueID)
	if planeIssueID == "" || cnbRepoID == "" || cnbIssueID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 plane_issue_id/cnb_repo_id/cnb_issue_id", nil)
	}
	deleted, err := h.db.DeleteIssueLink(c.Request().Context(), planeIssueID, cnbRepoID, cnbIssueID)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "delete_failed", "删除失败", map[string]any{"error": err.Error()})
	}
	if !deleted {
		return writeError(c, http.StatusNotFound, "not_found", "未找到对应映射", nil)
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok", "deleted": true})
}

// GET /admin/links/lark-threads
func (h *Handler) AdminLarkThreadLinksList(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	limit := 50
	if l := strings.TrimSpace(c.QueryParam("limit")); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	planeIssueID := strings.TrimSpace(c.QueryParam("plane_issue_id"))
	larkThreadID := strings.TrimSpace(c.QueryParam("lark_thread_id"))
	var syncEnabled *bool
	if sv := strings.TrimSpace(c.QueryParam("sync_enabled")); sv != "" {
		if b, err := strconv.ParseBool(sv); err == nil {
			syncEnabled = &b
		}
	}
	ctx := c.Request().Context()
	items, err := h.db.ListLarkThreadLinks(ctx, planeIssueID, larkThreadID, syncEnabled, limit)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
	}
	trimNull := func(ns sql.NullString) string {
		if ns.Valid {
			return strings.TrimSpace(ns.String)
		}
		return ""
	}
	planeClient := plane.Client{BaseURL: h.cfg.PlaneBaseURL}
	type workspaceMeta struct {
		token string
		name  string
	}
	workspaceInfos := make(map[string]workspaceMeta, len(items))
	for _, it := range items {
		slug := trimNull(it.WorkspaceSlug)
		if slug == "" {
			continue
		}
		if _, exists := workspaceInfos[slug]; exists {
			continue
		}
		meta := workspaceMeta{}
		if token, err := h.db.FindBotTokenByWorkspaceSlug(ctx, slug); err == nil {
			meta.token = strings.TrimSpace(token)
		}
		if meta.token != "" {
			wctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			ws, err := planeClient.GetWorkspace(wctx, meta.token, slug)
			cancel()
			if err == nil && ws != nil {
				name := strings.TrimSpace(ws.Name)
				if name == "" {
					name = strings.TrimSpace(ws.Title)
				}
				meta.name = name
			}
		}
		workspaceInfos[slug] = meta
	}
	type projectMeta struct {
		name       string
		identifier string
		slug       string
	}
	projectInfos := make(map[string]projectMeta, len(items))
	for _, it := range items {
		slug := trimNull(it.WorkspaceSlug)
		projectID := trimNull(it.PlaneProjectID)
		if slug == "" || projectID == "" {
			continue
		}
		key := slug + "::" + projectID
		if _, exists := projectInfos[key]; exists {
			continue
		}
		meta := projectMeta{}
		if wm, ok := workspaceInfos[slug]; ok && wm.token != "" {
			pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			proj, err := planeClient.GetProject(pctx, wm.token, slug, projectID)
			cancel()
			if err == nil && proj != nil {
				if name := strings.TrimSpace(proj.Name); name != "" {
					meta.name = name
				}
				if identifier := strings.TrimSpace(proj.Identifier); identifier != "" {
					meta.identifier = identifier
				}
				if pslug := strings.TrimSpace(proj.Slug); pslug != "" {
					meta.slug = pslug
				}
			}
		}
		projectInfos[key] = meta
	}
	issueNames := make(map[string]string, len(items))
	for _, it := range items {
		slug := trimNull(it.WorkspaceSlug)
		projectID := trimNull(it.PlaneProjectID)
		issueID := strings.TrimSpace(it.PlaneIssueID)
		if slug == "" || projectID == "" || issueID == "" {
			continue
		}
		key := slug + "::" + projectID + "::" + issueID
		if _, exists := issueNames[key]; exists {
			continue
		}
		if wm, ok := workspaceInfos[slug]; ok && wm.token != "" {
			ictx, cancel := context.WithTimeout(ctx, 5*time.Second)
			name, ierr := planeClient.GetIssueName(ictx, wm.token, slug, projectID, issueID)
			cancel()
			if ierr == nil {
				trimmed := strings.TrimSpace(name)
				if trimmed != "" {
					issueNames[key] = trimmed
					continue
				}
			}
		}
		issueNames[key] = ""
	}
	chatNames := make(map[string]string)
	chatIDs := make([]string, 0)
	for _, it := range items {
		if chatID := trimNull(it.LarkChatID); chatID != "" {
			if _, exists := chatNames[chatID]; !exists {
				chatNames[chatID] = ""
				chatIDs = append(chatIDs, chatID)
			}
		}
	}
	if len(chatIDs) > 0 && strings.TrimSpace(h.cfg.LarkAppID) != "" && strings.TrimSpace(h.cfg.LarkAppSecret) != "" {
		larkClient := lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
		tctx, cancel := context.WithTimeout(ctx, 8*time.Second)
		tenantToken, _, tokenErr := larkClient.TenantAccessToken(tctx)
		cancel()
		tenantToken = strings.TrimSpace(tenantToken)
		if tokenErr == nil && tenantToken != "" {
			for _, chatID := range chatIDs {
				cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				chat, cerr := larkClient.GetChat(cctx, tenantToken, chatID)
				cancel()
				if cerr == nil && chat != nil {
					if name := strings.TrimSpace(chat.Name); name != "" {
						chatNames[chatID] = name
					}
				}
			}
		}
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		slug := trimNull(it.WorkspaceSlug)
		projectID := trimNull(it.PlaneProjectID)
		issueID := strings.TrimSpace(it.PlaneIssueID)
		chatID := trimNull(it.LarkChatID)
		wm := workspaceInfos[slug]
		pm := projectInfos[slug+"::"+projectID]
		issueName := ""
		if slug != "" && projectID != "" && issueID != "" {
			issueName = issueNames[slug+"::"+projectID+"::"+issueID]
		}
		chatName := ""
		if chatID != "" {
			chatName = chatNames[chatID]
		}
		out = append(out, map[string]any{
			"lark_thread_id":           it.LarkThreadID,
			"lark_chat_id":             nullString(it.LarkChatID),
			"lark_chat_name":           optionalString(chatName),
			"plane_issue_id":           it.PlaneIssueID,
			"plane_issue_name":         optionalString(issueName),
			"plane_project_id":         nullString(it.PlaneProjectID),
			"plane_project_name":       optionalString(pm.name),
			"plane_project_identifier": optionalString(pm.identifier),
			"plane_project_slug":       optionalString(pm.slug),
			"workspace_slug":           nullString(it.WorkspaceSlug),
			"workspace_name":           optionalString(wm.name),
			"sync_enabled":             it.SyncEnabled,
			"linked_at":                it.LinkedAt.UTC().Format(time.RFC3339),
			"created_at":               it.CreatedAt.UTC().Format(time.RFC3339),
			"updated_at":               it.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "count": len(out)})
}

// POST /admin/links/lark-threads
func (h *Handler) AdminLarkThreadLinksUpsert(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		LarkThreadID   string `json:"lark_thread_id"`
		PlaneIssueID   string `json:"plane_issue_id"`
		PlaneProjectID string `json:"plane_project_id"`
		WorkspaceSlug  string `json:"workspace_slug"`
		SyncEnabled    *bool  `json:"sync_enabled"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	larkThreadID := strings.TrimSpace(req.LarkThreadID)
	planeIssueID := strings.TrimSpace(req.PlaneIssueID)
	if larkThreadID == "" || planeIssueID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 lark_thread_id/plane_issue_id", nil)
	}
	planeProjectID := strings.TrimSpace(req.PlaneProjectID)
	workspaceSlug := strings.TrimSpace(req.WorkspaceSlug)
	syncEnabled := false
	if req.SyncEnabled != nil {
		syncEnabled = *req.SyncEnabled
	}
	if err := h.db.UpsertLarkThreadLink(c.Request().Context(), larkThreadID, planeIssueID, planeProjectID, workspaceSlug, syncEnabled); err != nil {
		return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

// DELETE /admin/links/lark-threads
func (h *Handler) AdminLarkThreadLinksDelete(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		LarkThreadID string `json:"lark_thread_id"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	larkThreadID := strings.TrimSpace(req.LarkThreadID)
	if larkThreadID == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 lark_thread_id", nil)
	}
	deleted, err := h.db.DeleteLarkThreadLink(c.Request().Context(), larkThreadID)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "delete_failed", "删除失败", map[string]any{"error": err.Error()})
	}
	if !deleted {
		return writeError(c, http.StatusNotFound, "not_found", "未找到对应映射", nil)
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok", "deleted": true})
}
