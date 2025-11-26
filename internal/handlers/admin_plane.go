package handlers

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"cabb/internal/plane"

	"github.com/labstack/echo/v4"
)

// AdminPlaneMembers returns a list of members from a Plane workspace.
// GET /admin/plane/members?workspace_slug=xxx
func (h *Handler) AdminPlaneMembers(c echo.Context) error {
	// 1. Get Workspace Slug
	slug := c.QueryParam("workspace_slug")
	if slug == "" {
		// Try to find a default workspace from DB
		if hHasDB(h) {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()
			var dbSlug string
			err := h.db.SQL.QueryRowContext(ctx, `
				SELECT workspace_slug FROM repo_project_mappings 
				WHERE workspace_slug IS NOT NULL AND workspace_slug != '' 
				LIMIT 1
			`).Scan(&dbSlug)
			if err == nil && dbSlug != "" {
				slug = dbSlug
			}
		}
	}
	if slug == "" {
		return writeError(c, http.StatusBadRequest, "missing_param", "workspace_slug is required", nil)
	}

	// 2. Init Client
	client := &plane.Client{
		BaseURL: h.cfg.PlaneBaseURL,
	}
	token := h.cfg.PlaneServiceToken

	// 3. Fetch Members
	ctx := c.Request().Context()
	members, err := client.ListWorkspaceMembers(ctx, token, slug)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "plane_error", "failed to list members", map[string]any{"error": err.Error()})
	}

	// 4. Sort by Display Name
	sort.Slice(members, func(i, j int) bool {
		return strings.ToLower(members[i].DisplayName) < strings.ToLower(members[j].DisplayName)
	})

	return c.JSON(http.StatusOK, map[string]any{
		"items": members,
		"total": len(members),
		"slug":  slug,
	})
}
