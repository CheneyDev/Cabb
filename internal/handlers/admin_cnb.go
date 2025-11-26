package handlers

import (
	"net/http"
	"sort"
	"strings"

	"cabb/internal/cnb"

	"github.com/labstack/echo/v4"
)

// AdminCNBMembers returns a list of members from the configured CNB repo.
// GET /admin/cnb/members
func (h *Handler) AdminCNBMembers(c echo.Context) error {

	repo := h.cfg.CNBMemberRepo
	if repo == "" {
		repo = "1024hub/plane-test" // Fallback default
	}

	ctx := c.Request().Context()
	client := &cnb.Client{
		BaseURL: h.cfg.CNBBaseURL,
		Token:   h.cfg.CNBAppToken,
	}

	// 1. Fetch direct members (page 1, size 100 - assuming < 100 for now, or we can loop)
	// For better UX, we should probably fetch all or at least a large number.
	direct, err := client.ListRepoMembers(ctx, repo, 1, 100)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "cnb_error", "failed to list direct members", map[string]any{"error": err.Error()})
	}

	// 2. Fetch inherited members
	inheritedGroups, err := client.ListRepoInheritedMembers(ctx, repo, 1, 100)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "cnb_error", "failed to list inherited members", map[string]any{"error": err.Error()})
	}

	// 3. Aggregate and Deduplicate
	// Map by username (or ID) to avoid duplicates. CNB ID is unique.
	uniqueMembers := make(map[string]cnb.Member)

	for _, m := range direct {
		uniqueMembers[m.ID] = m
	}

	for _, group := range inheritedGroups {
		for _, m := range group.Users {
			if _, exists := uniqueMembers[m.ID]; !exists {
				uniqueMembers[m.ID] = m
			}
		}
	}

	// Convert to slice
	var result []cnb.Member
	for _, m := range uniqueMembers {
		result = append(result, m)
	}

	// Sort by username for consistent display
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Username) < strings.ToLower(result[j].Username)
	})

	return c.JSON(http.StatusOK, map[string]any{
		"items": result,
		"total": len(result),
		"repo":  repo,
	})
}
