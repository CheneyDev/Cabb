package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) AdminAccessList(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	users, err := h.db.ListAdminUsers(c.Request().Context())
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询系统用户失败", map[string]any{"error": err.Error()})
	}
	out := make([]map[string]any, 0, len(users))
	for i := range users {
		user := users[i]
		out = append(out, serializeAdminUser(&user))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "count": len(out)})
}

func (h *Handler) AdminAccessCreate(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
		Role        string `json:"role"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	displayName := strings.TrimSpace(req.DisplayName)
	password := strings.TrimSpace(req.Password)
	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = "admin"
	}
	if email == "" || displayName == "" || password == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "请填写 email/display_name/password", nil)
	}
	if role != "admin" {
		return writeError(c, http.StatusBadRequest, "invalid_role", "当前仅支持 admin 角色", nil)
	}
	if len(password) < 8 {
		return writeError(c, http.StatusBadRequest, "weak_password", "密码至少 8 位", nil)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "hash_failed", "加密密码失败", nil)
	}
	user, err := h.db.CreateAdminUser(c.Request().Context(), email, displayName, string(hash), role, true)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return writeError(c, http.StatusConflict, "duplicate_email", "邮箱已存在", nil)
		}
		return writeError(c, http.StatusBadGateway, "create_failed", "创建系统用户失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]any{"user": serializeAdminUser(user)})
}

func (h *Handler) AdminAccessUpdate(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return writeError(c, http.StatusBadRequest, "missing_id", "缺少用户 ID", nil)
	}
	var req struct {
		DisplayName *string `json:"display_name"`
		Role        *string `json:"role"`
		Active      *bool   `json:"active"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	ctx := c.Request().Context()
	user, err := h.db.GetAdminUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return writeError(c, http.StatusNotFound, "not_found", "用户不存在", nil)
		}
		return writeError(c, http.StatusBadGateway, "query_failed", "查询用户失败", map[string]any{"error": err.Error()})
	}
	displayName := user.DisplayName
	role := user.Role
	active := user.Active
	if req.DisplayName != nil {
		trimmed := strings.TrimSpace(*req.DisplayName)
		if trimmed == "" {
			return writeError(c, http.StatusBadRequest, "invalid_display_name", "显示名不能为空", nil)
		}
		displayName = trimmed
	}
	if req.Role != nil {
		candidate := strings.TrimSpace(*req.Role)
		if candidate == "" {
			candidate = "admin"
		}
		if candidate != "admin" {
			return writeError(c, http.StatusBadRequest, "invalid_role", "当前仅支持 admin 角色", nil)
		}
		role = candidate
	}
	if req.Active != nil {
		active = *req.Active
	}
	if err := h.db.UpdateAdminUser(ctx, id, displayName, role, active); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return writeError(c, http.StatusNotFound, "not_found", "用户不存在", nil)
		}
		return writeError(c, http.StatusBadGateway, "update_failed", "更新用户失败", map[string]any{"error": err.Error()})
	}
	updated, err := h.db.GetAdminUserByID(ctx, id)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "query_failed", "查询用户失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"user": serializeAdminUser(updated)})
}

func (h *Handler) AdminAccessResetPassword(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return writeError(c, http.StatusBadRequest, "missing_id", "缺少用户 ID", nil)
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	password := strings.TrimSpace(req.Password)
	if len(password) < 8 {
		return writeError(c, http.StatusBadRequest, "weak_password", "密码至少 8 位", nil)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "hash_failed", "加密密码失败", nil)
	}
	if err := h.db.UpdateAdminUserPassword(c.Request().Context(), id, string(hash)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return writeError(c, http.StatusNotFound, "not_found", "用户不存在", nil)
		}
		return writeError(c, http.StatusBadGateway, "update_failed", "重置密码失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"result": "ok"})
}
