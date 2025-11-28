package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
	"unicode"

	"cabb/internal/lark"

	"github.com/labstack/echo/v4"
)

// MaskedLarkUser represents a Lark user with masked name for login display.
type MaskedLarkUser struct {
	OpenID     string `json:"open_id"`
	Name       string `json:"name"`        // Original name (for internal use)
	MaskedName string `json:"masked_name"` // Masked for display
	Avatar     string `json:"avatar"`
}

// PublicLarkUsers returns a list of Lark users with masked names for login page.
// GET /api/auth/lark-users (no auth required)
func (h *Handler) PublicLarkUsers(c echo.Context) error {
	if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return writeError(c, http.StatusServiceUnavailable, "lark_not_configured", "飞书未配置", nil)
	}

	client := &lark.Client{
		AppID:     h.cfg.LarkAppID,
		AppSecret: h.cfg.LarkAppSecret,
	}

	ctx := c.Request().Context()
	token, _, err := client.TenantAccessToken(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_token_error", "获取飞书凭证失败", nil)
	}

	allUsers, err := client.FindAllUsers(ctx, token)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_api_error", "获取用户列表失败", nil)
	}

	var result []MaskedLarkUser
	for _, u := range allUsers {
		avatar := ""
		if u.Avatar.Avatar72 != "" {
			avatar = u.Avatar.Avatar72
		}
		result = append(result, MaskedLarkUser{
			OpenID:     u.OpenID,
			Name:       u.Name,
			MaskedName: maskName(u.Name),
			Avatar:     avatar,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"items": result,
		"total": len(result),
	})
}

// MagicLinkSend sends a magic link to a Lark user.
// POST /api/auth/magic-link
func (h *Handler) MagicLinkSend(c echo.Context) error {
	if !hHasDB(h) {
		return writeError(c, http.StatusServiceUnavailable, "db_unavailable", "数据库未配置", nil)
	}
	if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return writeError(c, http.StatusServiceUnavailable, "lark_not_configured", "飞书未配置", nil)
	}

	var req struct {
		OpenID string `json:"open_id"`
		Name   string `json:"name"`
		Origin string `json:"origin"`
	}
	if err := c.Bind(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "请求格式错误", nil)
	}
	if strings.TrimSpace(req.OpenID) == "" {
		return writeError(c, http.StatusBadRequest, "missing_open_id", "缺少 open_id", nil)
	}

	ctx := c.Request().Context()

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return writeError(c, http.StatusInternalServerError, "token_error", "生成令牌失败", nil)
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Save token to DB
	_, err := h.db.CreateMagicLinkToken(ctx, req.OpenID, req.Name, token, expiresAt)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "db_error", "保存令牌失败", map[string]any{"error": err.Error()})
	}

	// Build magic link URL - prefer origin from request, then config, then infer
	baseURL := strings.TrimSpace(req.Origin)
	if baseURL == "" {
		baseURL = h.cfg.FrontendBaseURL
	}
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	magicLink := strings.TrimRight(baseURL, "/") + "/api/auth/magic-link/verify?token=" + token

	// Build card message
	card := buildMagicLinkCard(req.Name, magicLink, expiresAt)

	// Send to Lark user
	client := &lark.Client{
		AppID:     h.cfg.LarkAppID,
		AppSecret: h.cfg.LarkAppSecret,
	}
	larkToken, _, err := client.TenantAccessToken(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_token_error", "获取飞书凭证失败", nil)
	}

	if err := client.SendCardToUser(ctx, larkToken, req.OpenID, card); err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_send_error", "发送消息失败", map[string]any{"error": err.Error()})
	}

	// Cleanup old tokens
	_ = h.db.CleanupExpiredMagicLinkTokens(ctx)

	return c.JSON(http.StatusOK, map[string]any{
		"status":  "ok",
		"message": "登录链接已发送到您的飞书",
	})
}

// MagicLinkVerify verifies a magic link token and creates a session.
// GET /api/auth/magic-link/verify?token=xxx
func (h *Handler) MagicLinkVerify(c echo.Context) error {
	if !hHasDB(h) {
		return redirectWithError(c, "数据库未配置")
	}

	token := c.QueryParam("token")
	if strings.TrimSpace(token) == "" {
		return redirectWithError(c, "无效的登录链接")
	}

	ctx := c.Request().Context()

	// Get token from DB
	magicToken, err := h.db.GetMagicLinkToken(ctx, token)
	if err != nil {
		return redirectWithError(c, "登录链接无效或已过期")
	}

	// Check if expired
	if time.Now().After(magicToken.ExpiresAt) {
		return redirectWithError(c, "登录链接已过期，请重新获取")
	}

	// Check if already used
	if magicToken.UsedAt.Valid {
		return redirectWithError(c, "登录链接已使用，请重新获取")
	}

	// Mark as used
	if err := h.db.MarkMagicLinkTokenUsed(ctx, token); err != nil {
		return redirectWithError(c, "验证失败，请重试")
	}

	// Get or create admin user
	displayName := ""
	if magicToken.LarkName.Valid {
		displayName = magicToken.LarkName.String
	}
	if displayName == "" {
		displayName = "飞书用户"
	}

	user, err := h.db.GetOrCreateAdminUserByLarkOpenID(ctx, magicToken.LarkOpenID, displayName, "")
	if err != nil {
		return redirectWithError(c, "创建用户失败")
	}

	if !user.Active {
		return redirectWithError(c, "账号已禁用，请联系管理员")
	}

	// Create session
	sessionToken, err := generateSessionToken(48)
	if err != nil {
		return redirectWithError(c, "创建会话失败")
	}

	expires := time.Now().Add(h.sessionTTL)
	ua := c.Request().UserAgent()
	ip := c.RealIP()

	sess, err := h.db.CreateAdminSession(ctx, user.ID, sessionToken, ua, ip, expires)
	if err != nil {
		return redirectWithError(c, "创建会话失败")
	}

	_ = h.db.RecordAdminLogin(ctx, user.ID, time.Now())

	// Set session cookie
	h.issueSessionCookie(c, sessionToken, sess.ExpiresAt)

	// Redirect to dashboard
	return c.Redirect(http.StatusFound, "/")
}

func redirectWithError(c echo.Context, message string) error {
	return c.Redirect(http.StatusFound, "/login?error="+message)
}

func buildMagicLinkCard(name, link string, expiresAt time.Time) map[string]any {
	expiresStr := expiresAt.Format("15:04")
	return map[string]any{
		"schema": "2.0",
		"config": map[string]any{
			"wide_screen_mode": true,
		},
		"header": map[string]any{
			"title": map[string]any{
				"tag":     "plain_text",
				"content": "🔐 Cabb 登录验证",
			},
			"template": "blue",
		},
		"body": map[string]any{
			"elements": []map[string]any{
				{
					"tag":     "markdown",
					"content": "你好 **" + name + "**，\n\n点击下方按钮登录 Cabb 后台管理系统。",
				},
				{
					"tag":            "column_set",
					"flex_mode":      "none",
					"background_style": "default",
					"columns": []map[string]any{
						{
							"tag":    "column",
							"width":  "weighted",
							"weight": 1,
							"elements": []map[string]any{
								{
									"tag": "button",
									"text": map[string]any{
										"tag":     "plain_text",
										"content": "立即登录",
									},
									"type": "primary",
									"multi_url": map[string]any{
										"url": link,
									},
								},
							},
						},
					},
				},
				{
					"tag":     "markdown",
					"content": "⏰ 链接有效期至 **" + expiresStr + "**，请尽快使用。\n\n如非本人操作，请忽略此消息。",
				},
				{
					"tag": "hr",
				},
				{
					"tag":     "markdown",
					"content": "按钮无法点击？复制以下链接到浏览器打开：\n" + link,
				},
			},
		},
	}
}

// maskName masks a name for privacy.
// Chinese: 王宝哥 -> 王XX
// English: Cheney -> ChXXXX
func maskName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "***"
	}

	runes := []rune(name)
	if len(runes) == 0 {
		return "***"
	}

	// Check if first character is Chinese
	if isChinese(runes[0]) {
		// Chinese name: keep first character
		if len(runes) == 1 {
			return string(runes[0]) + "X"
		}
		return string(runes[0]) + strings.Repeat("X", min(len(runes)-1, 2))
	}

	// English/other name: keep first 2 characters
	if len(runes) <= 2 {
		return string(runes[0]) + "X"
	}
	keepLen := 2
	maskLen := min(len(runes)-keepLen, 4)
	return string(runes[:keepLen]) + strings.Repeat("X", maskLen)
}

func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
