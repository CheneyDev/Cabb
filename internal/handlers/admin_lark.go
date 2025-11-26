package handlers

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"cabb/internal/lark"

	"github.com/labstack/echo/v4"
)

// AdminLarkUsers returns a list of Lark users.
// GET /admin/lark/users
func (h *Handler) AdminLarkUsers(c echo.Context) error {
	// 1. Init Client
	client := &lark.Client{
		AppID:     h.cfg.LarkAppID,
		AppSecret: h.cfg.LarkAppSecret,
	}

	// 2. Get Tenant Token
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	token, _, err := client.TenantAccessToken(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_token_error", "failed to get tenant token", map[string]any{"error": err.Error()})
	}

	// 3. Fetch All Users (Recursive)
	allUsers, err := client.FindAllUsers(ctx, token)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_api_error", "failed to fetch users", map[string]any{"error": err.Error()})
	}

	// 4. Sort by Name
	sort.Slice(allUsers, func(i, j int) bool {
		return strings.ToLower(allUsers[i].Name) < strings.ToLower(allUsers[j].Name)
	})

	return c.JSON(http.StatusOK, map[string]any{
		"items": allUsers,
		"total": len(allUsers),
	})
}
