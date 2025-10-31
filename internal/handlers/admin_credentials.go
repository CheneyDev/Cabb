package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"cabb/internal/plane"
)

// GET /admin/plane/credentials
// Returns all Plane Service Token credentials with workspace metadata
func (h *Handler) AdminPlaneCredentialsList(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}

	ctx := c.Request().Context()
	items, err := h.db.ListPlaneCredentials(ctx)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询失败", map[string]any{"error": err.Error()})
	}

	// Fetch workspace metadata for each credential
	planeClient := plane.Client{BaseURL: h.cfg.PlaneBaseURL}
	out := make([]map[string]any, 0, len(items))

	for _, cred := range items {
		workspaceMeta := map[string]any{
			"name": nil,
			"slug": cred.WorkspaceSlug,
		}

		// Try to fetch workspace name from Plane API
		if strings.TrimSpace(cred.TokenEnc) != "" && strings.TrimSpace(cred.WorkspaceSlug) != "" {
			wctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			ws, werr := planeClient.GetWorkspace(wctx, cred.TokenEnc, cred.WorkspaceSlug)
			cancel()

			if werr == nil && ws != nil {
				name := strings.TrimSpace(ws.Name)
				if name == "" {
					name = strings.TrimSpace(ws.Title)
				}
				if name != "" {
					workspaceMeta["name"] = name
				}
			}
		}

		// Mask the token for security (only show last 8 characters)
		maskedToken := ""
		if len(cred.TokenEnc) > 8 {
			maskedToken = "..." + cred.TokenEnc[len(cred.TokenEnc)-8:]
		} else if len(cred.TokenEnc) > 0 {
			maskedToken = "***"
		}

		out = append(out, map[string]any{
			"id":                 cred.ID,
			"plane_workspace_id": cred.PlaneWorkspaceID,
			"workspace_slug":     cred.WorkspaceSlug,
			"workspace_name":     workspaceMeta["name"],
			"kind":               cred.Kind,
			"token_masked":       maskedToken,
			"created_at":         cred.CreatedAt.UTC().Format(time.RFC3339),
			"updated_at":         cred.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{"items": out, "count": len(out)})
}

// POST /admin/plane/credentials
// Body: { "plane_workspace_id": "uuid", "workspace_slug": "my-workspace", "token": "plane_token_xxx" }
// Creates or updates a Service Token credential
func (h *Handler) AdminPlaneCredentialsUpsert(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}

	var req struct {
		PlaneWorkspaceID string `json:"plane_workspace_id"`
		WorkspaceSlug    string `json:"workspace_slug"`
		Token            string `json:"token"`
	}

	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}

	planeWorkspaceID := strings.TrimSpace(req.PlaneWorkspaceID)
	workspaceSlug := strings.TrimSpace(req.WorkspaceSlug)
	token := strings.TrimSpace(req.Token)

	if planeWorkspaceID == "" || workspaceSlug == "" || token == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 plane_workspace_id/workspace_slug/token", nil)
	}

	// Validate token by trying to fetch workspace
	planeClient := plane.Client{BaseURL: h.cfg.PlaneBaseURL}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	ws, err := planeClient.GetWorkspace(ctx, token, workspaceSlug)
	if err != nil {
		return writeError(c, http.StatusUnauthorized, "invalid_token", "Token 校验失败，无法访问 Plane API", map[string]any{"error": err.Error()})
	}

	if ws == nil {
		return writeError(c, http.StatusNotFound, "workspace_not_found", "未找到对应的 Workspace", nil)
	}

	// TODO: Implement transparent encryption when available
	// For now, store token as-is (token_enc column name indicates future encryption)
	tokenEnc := token

	// Upsert credential
	if err := h.db.UpsertPlaneCredential(ctx, planeWorkspaceID, workspaceSlug, "service", tokenEnc); err != nil {
		return writeError(c, http.StatusBadGateway, "save_failed", "保存失败", map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"result": "ok",
		"message": "凭据已保存，Plane 出站功能已启用",
		"workspace": map[string]any{
			"id":   planeWorkspaceID,
			"slug": workspaceSlug,
			"name": optionalString(ws.Name),
		},
	})
}

// DELETE /admin/plane/credentials/:id
// Deletes a credential by ID
func (h *Handler) AdminPlaneCredentialsDelete(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return writeError(c, http.StatusBadRequest, "missing_id", "缺少凭据 ID", nil)
	}

	ctx := c.Request().Context()
	deleted, err := h.db.DeletePlaneCredential(ctx, id)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "delete_failed", "删除失败", map[string]any{"error": err.Error()})
	}

	if !deleted {
		return writeError(c, http.StatusNotFound, "not_found", "未找到对应凭据", nil)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"result": "ok",
		"deleted": true,
		"message": "凭据已删除，Plane 出站功能已禁用",
	})
}
