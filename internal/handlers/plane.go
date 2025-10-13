package handlers

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
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
    // Basic config validation
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
        token  *tokenResponse
        err    error
        tType  = ""
        result = map[string]any{}
    )

    // Prefer app_installation_id flow when present (bot token)
    if appInstallationID != "" {
        token, err = h.getBotToken(ctx, appInstallationID)
        if err != nil {
            return writeError(c, http.StatusBadGateway, "token_exchange_failed", "获取 Bot Token 失败", map[string]any{"error": err.Error()})
        }
        tType = "bot"
        result["app_installation_id"] = appInstallationID
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
        if err != nil {
            // Non-fatal: return token metadata even if installation lookup fails
            result["installation_lookup_error"] = err.Error()
        }
    }

    // Compute expires_at (RFC3339 UTC)
    var expiresAt string
    if token != nil && token.ExpiresIn > 0 {
        expiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
    }

    // TODO: 持久化 tokens（透明加密存储），落库 workspaces（待接入 DB）

    // Build safe response (do not leak tokens)
    resp := map[string]any{
        "result":       "ok",
        "token_type":   tType,
        "state":        state,
        "expires_at":   expiresAt,
        "has_refresh":  token.RefreshToken != "",
    }
    // Include minimal workspace details if available
    if inst != nil {
        resp["workspace"] = map[string]any{
            "id":        inst.WorkspaceID(),
            "slug":      inst.WorkspaceSlug(),
            "app_bot":   inst.AppBot,
            "status":    inst.Status,
            "install_id": inst.ID,
        }
    }

    return c.JSON(http.StatusOK, resp)
}

func (h *Handler) PlaneWebhook(c echo.Context) error {
    // Verify signature if provided
    body, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    sig := c.Request().Header.Get("X-Plane-Signature")
    if h.cfg.PlaneWebhookSecret != "" {
        if !verifyHMACSHA256(body, h.cfg.PlaneWebhookSecret, sig) {
            return c.NoContent(http.StatusUnauthorized)
        }
    }

    // Acknowledge; actual routing will be implemented later.
    deliveryID := c.Request().Header.Get("X-Plane-Delivery")
    eventType := c.Request().Header.Get("X-Plane-Event")
    return c.JSON(http.StatusOK, map[string]any{
        "accepted":    true,
        "delivery_id": deliveryID,
        "event_type":  eventType,
    })
}

func verifyHMACSHA256(body []byte, secret, got string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    want := hex.EncodeToString(mac.Sum(nil))
    // Some platforms prefix with "sha256=..."; accept either raw or prefixed.
    got = strings.TrimSpace(got)
    if strings.HasPrefix(strings.ToLower(got), "sha256=") {
        got = got[len("sha256="):]
    }
    // Constant-time compare
    wantB, err1 := hex.DecodeString(want)
    gotB, err2 := hex.DecodeString(got)
    if err1 != nil || err2 != nil {
        return false
    }
    return hmac.Equal(wantB, gotB)
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
    AppBot           string `json:"app_bot"`
    WorkspaceDetail  struct {
        ID   string `json:"id"`
        Slug string `json:"slug"`
        Name string `json:"name"`
    } `json:"workspace_detail"`
}

func (a *appInstallation) WorkspaceID() string  { return a.Workspace }
func (a *appInstallation) WorkspaceSlug() string { return a.WorkspaceDetail.Slug }

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
    // Basic auth: client_id:client_secret
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
        return nil, fmt.Errorf("app-installation status=%d body=%s", resp.StatusCode, truncate(string(b), 300))
    }
    // The API returns an array; take first element when present
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

// ==== helpers ====

func truncate(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return s[:n]
}

func writeError(c echo.Context, status int, code, message string, details map[string]any) error {
    if details == nil {
        details = map[string]any{}
    }
    reqID := c.Response().Header().Get("X-Request-ID")
    if reqID == "" {
        reqID = c.Request().Header.Get("X-Request-ID")
    }
    return c.JSON(status, map[string]any{
        "error": map[string]any{
            "code":       code,
            "message":    message,
            "details":    details,
            "request_id": reqID,
        },
    })
}
