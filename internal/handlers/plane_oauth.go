package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// OAuthStart redirects the user to Plane's consent page
// Docs: docs/plane-developer-docs/dev-tools/build-plane-app.mdx (authorize-app)
func (h *Handler) PlaneOAuthStart(c echo.Context) error {
	clientID := h.cfg.PlaneClientID
	redirectURI := h.cfg.PlaneRedirectURI
	if clientID == "" || redirectURI == "" {
		return writeError(c, http.StatusBadRequest, "invalid_config", "缺少 PLANE_CLIENT_ID 或 PLANE_REDIRECT_URI", nil)
	}

	state := c.QueryParam("state")
	u, err := url.Parse(h.cfg.PlaneBaseURL)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "invalid_base_url", "PLANE_BASE_URL 无法解析", map[string]any{"base_url": h.cfg.PlaneBaseURL})
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/auth/o/authorize-app/"
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", redirectURI)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

// OAuthCallback exchanges code/app_installation_id for tokens and fetches installation details.
// - If app_installation_id is present: client_credentials → bot token
// - Else if code is present: authorization_code → user token
func (h *Handler) PlaneOAuthCallback(c echo.Context) error {
	if h.cfg.PlaneClientID == "" || h.cfg.PlaneClientSecret == "" || h.cfg.PlaneBaseURL == "" {
		return writeError(c, http.StatusBadRequest, "invalid_config", "缺少 Plane OAuth 配置（client_id/client_secret/base_url）", nil)
	}

	appInstallationID := c.QueryParam("app_installation_id")
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if appInstallationID == "" && code == "" {
		return writeError(c, http.StatusBadRequest, "invalid_request", "缺少 app_installation_id 或 code", nil)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var (
		token *tokenResponse
		err   error
		tType = ""
	)

	if appInstallationID != "" {
		token, err = h.getBotToken(ctx, appInstallationID)
		if err != nil {
			return writeError(c, http.StatusBadGateway, "token_exchange_failed", "获取 Bot Token 失败", map[string]any{"error": err.Error()})
		}
		tType = "bot"
	} else {
		token, err = h.exchangeAuthorizationCode(ctx, code)
		if err != nil {
			return writeError(c, http.StatusBadGateway, "token_exchange_failed", "授权码换取 Token 失败", map[string]any{"error": err.Error()})
		}
		tType = "user"
	}

	// Optionally fetch installation details when we have app_installation_id
	var inst *appInstallation
	if appInstallationID != "" {
		inst, err = h.getAppInstallation(ctx, token.AccessToken, appInstallationID)
		_ = err // 非致命错误，忽略安装信息查询失败（仍返回成功摘要）
	}

	// Compute expires_at (RFC3339 UTC)
	var expiresAt string
	if token != nil && token.ExpiresIn > 0 {
		expiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
	}

	// 持久化 tokens（透明加密存储：待实现）
	if inst != nil && token != nil && hHasDB(h) {
		_ = h.db.UpsertWorkspaceToken(c.Request().Context(), inst.WorkspaceID(), appInstallationID, tType, token.AccessToken, token.RefreshToken, expiresAt, inst.WorkspaceSlug(), inst.AppBot)
	}

	// Build safe response (do not leak tokens)
	resp := map[string]any{
		"result":      "ok",
		"token_type":  tType,
		"state":       state,
		"expires_at":  expiresAt,
		"has_refresh": token.RefreshToken != "",
	}
	if inst != nil {
		resp["workspace"] = map[string]any{
			"id":         inst.WorkspaceID(),
			"slug":       inst.WorkspaceSlug(),
			"app_bot":    inst.AppBot,
			"status":     inst.Status,
			"install_id": inst.ID,
		}
	}

	if strings.EqualFold(c.QueryParam("format"), "json") {
		return c.JSON(http.StatusOK, resp)
	}
	accept := c.Request().Header.Get("Accept")
	ua := c.Request().Header.Get("User-Agent")
	wantsJSON := strings.Contains(accept, "application/json") && !strings.Contains(accept, "text/html")
	isCLI := strings.Contains(strings.ToLower(ua), "curl/") || strings.Contains(strings.ToLower(ua), "httpie/")
	if wantsJSON || isCLI {
		return c.JSON(http.StatusOK, resp)
	}

	wsSlug := ""
	if inst != nil {
		wsSlug = inst.WorkspaceSlug()
	}
	target := h.preferredReturnURL(wsSlug, c.QueryParam("state"), c.QueryParam("return_to"))

	html := h.buildRedirectHTML(target, resp)
	return c.HTML(http.StatusOK, html)
}

// ==== Plane OAuth helpers ====

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type appInstallation struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	Workspace        string `json:"workspace"`
	WorkspaceIDAlt   string `json:"workspace_id"`
	WorkspaceSlugTop string `json:"workspace_slug"`
	AppBot           string `json:"app_bot"`
	WorkspaceDetail  struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"workspace_detail"`
}

func (a *appInstallation) WorkspaceID() string {
	if a.Workspace != "" {
		return a.Workspace
	}
	return a.WorkspaceIDAlt
}
func (a *appInstallation) WorkspaceSlug() string {
	if a.WorkspaceSlugTop != "" {
		return a.WorkspaceSlugTop
	}
	return a.WorkspaceDetail.Slug
}

func (h *Handler) tokenEndpoint() (string, error) {
	u, err := url.Parse(h.cfg.PlaneBaseURL)
	if err != nil {
		return "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/auth/o/token/"
	return u.String(), nil
}

func (h *Handler) appInstallationEndpoint(id string) (string, error) {
	u, err := url.Parse(h.cfg.PlaneBaseURL)
	if err != nil {
		return "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/auth/o/app-installation/"
	q := url.Values{}
	q.Set("id", id)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (h *Handler) getBotToken(ctx context.Context, appInstallationID string) (*tokenResponse, error) {
	endpoint, err := h.tokenEndpoint()
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("app_installation_id", appInstallationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	creds := h.cfg.PlaneClientID + ":" + h.cfg.PlaneClientSecret
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(creds)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token endpoint status=%d body=%s", resp.StatusCode, truncate(string(b), 300))
	}
	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	if tr.AccessToken == "" {
		return nil, errors.New("empty access_token in response")
	}
	return &tr, nil
}

func (h *Handler) exchangeAuthorizationCode(ctx context.Context, code string) (*tokenResponse, error) {
	endpoint, err := h.tokenEndpoint()
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", h.cfg.PlaneClientID)
	form.Set("client_secret", h.cfg.PlaneClientSecret)
	form.Set("redirect_uri", h.cfg.PlaneRedirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token endpoint status=%d body=%s", resp.StatusCode, truncate(string(b), 300))
	}
	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	if tr.AccessToken == "" {
		return nil, errors.New("empty access_token in response")
	}
	return &tr, nil
}

func (h *Handler) getAppInstallation(ctx context.Context, bearerToken, appInstallationID string) (*appInstallation, error) {
	endpoint, err := h.appInstallationEndpoint(appInstallationID)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("installation endpoint status=%d body=%s", resp.StatusCode, truncate(string(b), 300))
	}
	var arr []appInstallation
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, errors.New("empty installation list")
	}
	inst := arr[0]
	return &inst, nil
}

// preferredReturnURL selects a safe redirect destination back to Plane
func (h *Handler) preferredReturnURL(workspaceSlug, state, returnTo string) string {
	allowed := map[string]struct{}{}
	addHost := func(raw string) {
		if raw == "" {
			return
		}
		if u, err := url.Parse(raw); err == nil && u.Host != "" {
			allowed[strings.ToLower(u.Host)] = struct{}{}
		}
	}
	addHost(h.cfg.PlaneAppBaseURL)
	addHost(h.cfg.PlaneBaseURL)

	isAllowed := func(u *url.URL) bool {
		if u == nil || (u.Scheme != "http" && u.Scheme != "https") {
			return false
		}
		host := strings.ToLower(u.Host)
		if _, ok := allowed[host]; ok {
			return true
		}
		if base, err := url.Parse(h.cfg.PlaneBaseURL); err == nil {
			bh := strings.ToLower(base.Host)
			if strings.HasPrefix(bh, "api.") {
				alt := "app." + strings.TrimPrefix(bh, "api.")
				if host == alt {
					return true
				}
			}
		}
		return false
	}

	if returnTo != "" {
		if u, err := url.Parse(returnTo); err == nil && isAllowed(u) {
			return u.String()
		}
	}
	if workspaceSlug != "" {
		if base := h.planeAppBase(); base != nil {
			base.Path = strings.TrimRight(base.Path, "/") + "/" + workspaceSlug + "/settings/integrations/"
			base.RawQuery = ""
			base.Fragment = ""
			return base.String()
		}
	}
	if s := strings.TrimSpace(state); s != "" {
		if u, err := url.Parse(s); err == nil && isAllowed(u) {
			return u.String()
		}
	}
	if base := h.planeAppBase(); base != nil {
		base.Path = "/"
		base.RawQuery = ""
		base.Fragment = ""
		return base.String()
	}
	if u, err := url.Parse(h.cfg.PlaneBaseURL); err == nil {
		u.Path = "/"
		u.RawQuery = ""
		u.Fragment = ""
		return u.String()
	}
	return "/"
}

func (h *Handler) buildRedirectHTML(target string, payload map[string]any) string {
	_ = payload
	esc := func(s string) string {
		r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&#39;")
		return r.Replace(s)
	}
	return "<!DOCTYPE html><html lang=\"zh-CN\"><head>" +
		"<meta charset=\"utf-8\">" +
		"<meta http-equiv=\"refresh\" content=\"2; url=" + esc(target) + "\">" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">" +
		"<meta name=\"color-scheme\" content=\"light dark\">" +
		"<title>返回 Plane</title>" +
		"<style>html,body{height:100%}body{margin:0;font:16px/1.6 -apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;display:grid;place-items:center}.wrap{text-align:center;padding:24px}h1{font-size:20px;margin:0 0 8px}p{margin:6px 0;color:#6b7280}.spinner{width:28px;height:28px;border:3px solid currentColor;border-right-color:transparent;border-radius:50%;margin:14px auto;animation:s .8s linear infinite;opacity:.7}a.btn{display:inline-block;margin-top:10px;padding:8px 12px;border:1px solid currentColor;border-radius:8px;text-decoration:none}@keyframes s{to{transform:rotate(360deg)}}" +
		"</style></head><body><div class=\"wrap\"><div class=\"spinner\" aria-hidden=\"true\"></div><h1>安装完成，正在返回 Plane…</h1><p>若未自动跳转，请点击下方按钮</p><p><a class=\"btn\" href=\"" + esc(target) + "\">返回 Plane</a></p></div>" +
		"<script>(function(){try{var t='" + esc(target) + "'; if(window.opener && !window.opener.closed){try{window.opener.postMessage({type:'plane_installation',status:'ok',target:t}, '*');}catch(e){}} window.location.replace(t);}catch(e){}})</script>" +
		"</body></html>"
}

// planeAppBase returns a parsed URL for PLANE_APP_BASE_URL or derives from PLANE_BASE_URL
func (h *Handler) planeAppBase() *url.URL {
	if h.cfg.PlaneAppBaseURL != "" {
		if u, err := url.Parse(h.cfg.PlaneAppBaseURL); err == nil {
			return u
		}
	}
	if u, err := url.Parse(h.cfg.PlaneBaseURL); err == nil {
		if strings.HasPrefix(strings.ToLower(u.Host), "api.") {
			u.Host = "app." + strings.TrimPrefix(u.Host, "api.")
		}
		u.Path = "/"
		u.RawQuery = ""
		u.Fragment = ""
		return u
	}
	return nil
}
