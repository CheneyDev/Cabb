package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"cabb/internal/store"
)

const (
	adminSessionContextKey = "admin_session"
)

var (
	errNoSession      = errors.New("no session")
	errSessionExpired = errors.New("session expired")
	errSessionRevoked = errors.New("session revoked")
	errUserInactive   = errors.New("user inactive")
)

func (h *Handler) AdminLogin(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "解析失败", nil)
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少 email/password", nil)
	}
	ctx := c.Request().Context()
	user, err := h.db.GetAdminUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return writeError(c, http.StatusUnauthorized, "invalid_credentials", "账号或密码错误", nil)
		}
		return writeError(c, http.StatusBadGateway, "query_failed", "查询用户失败", map[string]any{"error": err.Error()})
	}
	if !user.Active {
		return writeError(c, http.StatusForbidden, "user_inactive", "账号已禁用，请联系管理员", nil)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return writeError(c, http.StatusUnauthorized, "invalid_credentials", "账号或密码错误", nil)
	}
	token, err := generateSessionToken(48)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "token_error", "生成会话失败", nil)
	}
	expires := time.Now().Add(h.sessionTTL)
	ua := c.Request().UserAgent()
	ip := c.RealIP()
	sess, err := h.db.CreateAdminSession(ctx, user.ID, token, ua, ip, expires)
	if err != nil {
		return writeError(c, http.StatusBadGateway, "session_create_failed", "创建会话失败", map[string]any{"error": err.Error()})
	}
	_ = h.db.RecordAdminLogin(ctx, user.ID, time.Now())
	_ = h.db.CleanupExpiredAdminSessions(ctx)
	h.issueSessionCookie(c, token, sess.ExpiresAt)
	return c.JSON(http.StatusOK, map[string]any{
		"user":    serializeAdminUser(user),
		"session": map[string]any{"expires_at": sess.ExpiresAt.UTC().Format(time.RFC3339)},
	})
}

func (h *Handler) AdminLogout(c echo.Context) error {
	if !hHasDB(h) {
		return c.JSON(http.StatusOK, map[string]any{"result": "ok"})
	}
	cookie, err := c.Cookie(h.sessionCookieName)
	if err == nil && cookie != nil && cookie.Value != "" {
		token := cookie.Value
		if err := h.db.RevokeAdminSession(c.Request().Context(), token); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return writeError(c, http.StatusBadGateway, "session_revoke_failed", "注销失败", map[string]any{"error": err.Error()})
		}
	}
	h.clearSessionCookie(c)
	return c.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

func (h *Handler) AdminMe(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	if sess, ok := c.Get(adminSessionContextKey).(*store.AdminSessionWithUser); ok && sess != nil {
		return c.JSON(http.StatusOK, map[string]any{
			"user":    serializeAdminUser(&sess.User),
			"session": map[string]any{"expires_at": sess.Session.ExpiresAt.UTC().Format(time.RFC3339)},
		})
	}
	sess, err := h.sessionFromRequest(c)
	if err != nil {
		return h.handleAuthError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"user":    serializeAdminUser(&sess.User),
		"session": map[string]any{"expires_at": sess.Session.ExpiresAt.UTC().Format(time.RFC3339)},
	})
}

func (h *Handler) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := h.sessionFromRequest(c)
		if err != nil {
			return h.handleAuthError(c, err)
		}
		c.Set(adminSessionContextKey, sess)
		return next(c)
	}
}

func (h *Handler) sessionFromRequest(c echo.Context) (*store.AdminSessionWithUser, error) {
	if !hHasDB(h) {
		return nil, errNoSession
	}
	cookie, err := c.Cookie(h.sessionCookieName)
	if err != nil || cookie == nil || strings.TrimSpace(cookie.Value) == "" {
		return nil, errNoSession
	}
	token := strings.TrimSpace(cookie.Value)
	sess, err := h.db.GetAdminSessionWithUser(c.Request().Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNoSession
		}
		return nil, err
	}
	if sess.Session.RevokedAt.Valid {
		return nil, errSessionRevoked
	}
	if time.Now().After(sess.Session.ExpiresAt) {
		return nil, errSessionExpired
	}
	if !sess.User.Active {
		return nil, errUserInactive
	}
	return sess, nil
}

func (h *Handler) handleAuthError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, errNoSession):
		return writeError(c, http.StatusUnauthorized, "unauthorized", "请先登录", nil)
	case errors.Is(err, errSessionExpired):
		h.clearSessionCookie(c)
		return writeError(c, http.StatusUnauthorized, "session_expired", "登录已过期，请重新登录", nil)
	case errors.Is(err, errSessionRevoked):
		h.clearSessionCookie(c)
		return writeError(c, http.StatusUnauthorized, "session_revoked", "会话已失效，请重新登录", nil)
	case errors.Is(err, errUserInactive):
		h.clearSessionCookie(c)
		return writeError(c, http.StatusForbidden, "user_inactive", "账号已禁用，请联系管理员", nil)
	default:
		return writeError(c, http.StatusInternalServerError, "auth_failed", "鉴权失败", map[string]any{"error": err.Error()})
	}
}

func (h *Handler) issueSessionCookie(c echo.Context, token string, expires time.Time) {
	cookie := &http.Cookie{
		Name:     h.sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.sessionCookieSecure,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	}
	if h.sessionTTL > 0 {
		cookie.MaxAge = int(h.sessionTTL.Seconds())
	}
	c.SetCookie(cookie)
}

func (h *Handler) clearSessionCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     h.sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.sessionCookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	c.SetCookie(cookie)
}

func serializeAdminUser(u *store.AdminUser) map[string]any {
	if u == nil {
		return map[string]any{}
	}
	lastLogin := ""
	if u.LastLoginAt.Valid {
		lastLogin = u.LastLoginAt.Time.UTC().Format(time.RFC3339)
	}
	return map[string]any{
		"id":            u.ID,
		"email":         u.Email,
		"display_name":  u.DisplayName,
		"role":          u.Role,
		"active":        u.Active,
		"last_login_at": nullableString(lastLogin),
		"created_at":    u.CreatedAt.UTC().Format(time.RFC3339),
		"updated_at":    u.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func nullableString(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func generateSessionToken(n int) (string, error) {
	if n <= 0 {
		n = 32
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
