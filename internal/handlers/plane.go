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
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/labstack/echo/v4"
    "plane-integration/internal/cnb"
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
        token *tokenResponse
        err   error
        tType = ""
    )

    // Prefer app_installation_id flow when present (bot token)
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

    // Decide response format: default to HTML for browsers, JSON for API callers
    // Force JSON if query param format=json
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

    // Build redirect target (prefer explicit return_to, then workspace integrations by slug, then state URL, else fallback app/base)
    wsSlug := ""
    if inst != nil {
        wsSlug = inst.WorkspaceSlug()
    }
    target := h.preferredReturnURL(wsSlug, c.QueryParam("state"), c.QueryParam("return_to"))

    html := h.buildRedirectHTML(target, resp)
    return c.HTML(http.StatusOK, html)
}

func hHasDB(h *Handler) bool { return h != nil && h.db != nil && h.db.SQL != nil }

func (h *Handler) PlaneWebhook(c echo.Context) error {
    // Verify signature if provided
    body, err := io.ReadAll(c.Request().Body)
    if err != nil { return c.NoContent(http.StatusBadRequest) }
    sig := c.Request().Header.Get("X-Plane-Signature")
    if h.cfg.PlaneWebhookSecret != "" {
        if !verifyHMACSHA256(body, h.cfg.PlaneWebhookSecret, sig) {
            return c.NoContent(http.StatusUnauthorized)
        }
    }

    // idempotency
    deliveryID := c.Request().Header.Get("X-Plane-Delivery")
    eventType := c.Request().Header.Get("X-Plane-Event")
    hsum := sha256.Sum256(body)
    sum := hex.EncodeToString(hsum[:])
    if h.dedupe != nil && h.dedupe.CheckAndMark("plane."+eventType, deliveryID, sum) {
        return c.JSON(http.StatusOK, map[string]any{"accepted": true, "delivery_id": deliveryID, "event_type": eventType, "status": "duplicate"})
    }
    if hHasDB(h) && deliveryID != "" {
        dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "plane."+eventType, deliveryID, sum)
        if err == nil && dup {
            return c.JSON(http.StatusOK, map[string]any{"accepted": true, "delivery_id": deliveryID, "event_type": eventType, "status": "duplicate"})
        }
        _ = h.db.UpsertEventDelivery(c.Request().Context(), "plane."+eventType, "incoming", deliveryID, sum, "queued")
    }

    // async process
    go h.processPlaneWebhook(eventType, body, deliveryID)
    return c.JSON(http.StatusOK, map[string]any{"accepted": true, "delivery_id": deliveryID, "event_type": eventType, "status": "queued"})
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
    // Some Plane versions use "workspace" for ID, others may use "workspace_id"
    Workspace        string `json:"workspace"`
    WorkspaceIDAlt   string `json:"workspace_id"`
    // Slug may be directly provided or nested under workspace_detail
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
    // Structured error log (do not include sensitive info)
    // Fields align with AGENTS.md: request_id/source/endpoint/latency_ms/result + error.code
    ep := c.Path()
    if ep == "" { ep = c.Request().URL.Path }
    src := DetectSource(ep)
    // Latency is not directly available here; omit or set -1
    errLog := map[string]any{
        "time":       time.Now().UTC().Format(time.RFC3339),
        "level":      "error",
        "request_id": reqID,
        "endpoint":   ep,
        "source":     src,
        "status":     status,
        "result":     "error",
        "error": map[string]any{
            "code":    code,
            "message": message,
        },
    }
    if b, e := json.Marshal(errLog); e == nil { log.Println(string(b)) }
    return c.JSON(status, map[string]any{
        "error": map[string]any{
            "code":       code,
            "message":    message,
            "details":    details,
            "request_id": reqID,
        },
    })
}

// ==== Plane webhook processing ====
type planeWebhookEnvelope struct {
    Event       string         `json:"event"`
    Action      string         `json:"action"`
    WebhookID   string         `json:"webhook_id"`
    WorkspaceID string         `json:"workspace_id"`
    Data        map[string]any `json:"data"`
}

func (h *Handler) processPlaneWebhook(event string, body []byte, deliveryID string) {
    var env planeWebhookEnvelope
    if err := json.Unmarshal(body, &env); err != nil { return }
    // Trace inbound webhook envelope (decision log, no sensitive data)
    LogStructured("info", map[string]any{
        "event":        "plane.webhook",
        "delivery_id":  deliveryID,
        "x_event":      strings.ToLower(event),
        "action":       strings.ToLower(env.Action),
        "workspace_id": env.WorkspaceID,
    })
    switch strings.ToLower(event) {
    case "issue":
        h.handlePlaneIssueEvent(env, deliveryID)
    case "issue comment", "issue_comment":
        h.handlePlaneIssueComment(env, deliveryID)
    default:
        // ignore others for CNB integration
    }
}

func (h *Handler) handlePlaneIssueEvent(env planeWebhookEnvelope, deliveryID string) {
    if !hHasDB(h) { return }
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    // Extract
    action := strings.ToLower(env.Action)
    data := env.Data
    planeIssueID, _ := dataGetString(data, "id")
    planeProjectID, _ := dataGetString(data, "project")
    name, _ := dataGetString(data, "name")
    descHTML, _ := dataGetString(data, "description_html")

    // Attempt to extract labels for routing
    labels := dataGetLabels(data)
    // Load all mappings for project
    mappings, err := h.db.ListRepoProjectMappingsByPlaneProject(ctx, planeProjectID)

    // Decision log: summary for routing
    LogStructured("info", map[string]any{
        "event":              "plane.issue",
        "delivery_id":        deliveryID,
        "action":             action,
        "plane_issue_id":     planeIssueID,
        "plane_project_id":   planeProjectID,
        "labels":             labels,
        "mappings_count":     len(mappings),
        "outbound_enabled":   h.cfg.CNBOutboundEnabled,
    })

    // CNB outbound fan-out only when enabled and mappings exist
    if err == nil && len(mappings) > 0 && h.cfg.CNBOutboundEnabled {
        cn := &cnb.Client{
            BaseURL:          h.cfg.CNBBaseURL,
            Token:            h.cfg.CNBAppToken,
            IssueCreatePath:  h.cfg.CNBIssueCreatePath,
            IssueUpdatePath:  h.cfg.CNBIssueUpdatePath,
            IssueCommentPath: h.cfg.CNBIssueCommentPath,
        }
        switch action {
        case "create", "created":
            // Fan-out create to repos whose mapping requires bidirectional and label match (selector set)
            for _, m := range mappings {
                dir := strings.ToLower(m.SyncDirection.String)
                if !m.SyncDirection.Valid { dir = "" }
                hit := labelSelectorMatch(m.LabelSelector.String, labels)
                // Existing link check (per repo)
                links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID)
                exists := false
                for _, lk := range links { if lk.Repo == m.CNBRepoID { exists = true; break } }

                decision := "create"
                skip := ""
                if dir != "bidirectional" { decision = "skip"; skip = "direction" }
                if decision != "skip" && !hit { decision = "skip"; skip = "label_miss" }
                if decision != "skip" && exists { decision = "skip"; skip = "already_linked" }

                LogStructured("info", map[string]any{
                    "event":            "plane.issue.route",
                    "delivery_id":      deliveryID,
                    "action":           action,
                    "plane_issue_id":   planeIssueID,
                    "repo":             m.CNBRepoID,
                    "direction":        dir,
                    "label_selector":   m.LabelSelector.String,
                    "label_hit":        hit,
                    "already_linked":   exists,
                    "decision":         decision,
                    "skip_reason":      skip,
                })

                if decision == "skip" { continue }
                iid, err := cn.CreateIssue(ctx, m.CNBRepoID, name, descHTML)
                if err != nil || iid == "" {
                    LogStructured("error", map[string]any{
                        "event":          "plane.issue.cnbrpc",
                        "delivery_id":    deliveryID,
                        "action":         action,
                        "plane_issue_id": planeIssueID,
                        "repo":           m.CNBRepoID,
                        "op":             "create_issue",
                        "error": map[string]any{
                            "code":    "cnb_create_failed",
                            "message": truncate(fmt.Sprintf("%v", err), 200),
                        },
                    })
                    continue
                }
                _ = h.db.CreateIssueLink(ctx, planeIssueID, m.CNBRepoID, iid)
                LogStructured("info", map[string]any{
                    "event":          "plane.issue.cnbrpc",
                    "delivery_id":    deliveryID,
                    "action":         action,
                    "plane_issue_id": planeIssueID,
                    "repo":           m.CNBRepoID,
                    "op":             "create_issue",
                    "result":         "created",
                    "cnb_issue_iid":  iid,
                })
                // Set labels to CNB (ensure labels exist first)
                if len(labels) > 0 {
                    if err := cn.EnsureRepoLabels(ctx, m.CNBRepoID, labels); err != nil {
                        LogStructured("error", map[string]any{"event":"plane.issue.cnbrpc","op":"ensure_labels","repo":m.CNBRepoID,"delivery_id":deliveryID,"error":map[string]any{"code":"cnb_labels_ensure_failed","message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else if err := cn.SetIssueLabels(ctx, m.CNBRepoID, iid, labels); err != nil {
                        LogStructured("error", map[string]any{"event":"plane.issue.cnbrpc","op":"set_issue_labels","repo":m.CNBRepoID,"delivery_id":deliveryID,"error":map[string]any{"code":"cnb_set_labels_failed","message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        LogStructured("info", map[string]any{"event":"plane.issue.cnbrpc","op":"set_issue_labels","repo":m.CNBRepoID,"delivery_id":deliveryID,"result":"set"})
                    }
                }
            }
        case "update", "updated":
            if links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links) > 0 {
                fields := map[string]any{}
                if name != "" { fields["title"] = name }
                if descHTML != "" { fields["body"] = descHTML }
                for _, lk := range links {
                    if err := cn.UpdateIssue(ctx, lk.Repo, lk.Number, fields); err != nil {
                        LogStructured("error", map[string]any{
                            "event":          "plane.issue.cnbrpc",
                            "delivery_id":    deliveryID,
                            "action":         action,
                            "plane_issue_id": planeIssueID,
                            "repo":           lk.Repo,
                            "op":             "update_issue",
                            "error": map[string]any{"code":"cnb_update_failed","message": truncate(fmt.Sprintf("%v", err), 200)},
                        })
                    } else {
                        LogStructured("info", map[string]any{
                            "event":          "plane.issue.cnbrpc",
                            "delivery_id":    deliveryID,
                            "action":         action,
                            "plane_issue_id": planeIssueID,
                            "repo":           lk.Repo,
                            "op":             "update_issue",
                            "result":         "updated",
                        })
                    }
                }
                // Sync labels if provided in webhook
                if len(labels) > 0 {
                    // Ensure labels exist in repo
                    if err := cn.EnsureRepoLabels(ctx, links[0].Repo, labels); err != nil {
                        LogStructured("error", map[string]any{"event":"plane.issue.cnbrpc","op":"ensure_labels","repo":links[0].Repo,"delivery_id":deliveryID,"error":map[string]any{"code":"cnb_labels_ensure_failed","message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        for _, lk := range links {
                            if err := cn.SetIssueLabels(ctx, lk.Repo, lk.Number, labels); err != nil {
                                LogStructured("error", map[string]any{"event":"plane.issue.cnbrpc","op":"set_issue_labels","repo":lk.Repo,"delivery_id":deliveryID,"error":map[string]any{"code":"cnb_set_labels_failed","message": truncate(fmt.Sprintf("%v", err), 200)}})
                            } else {
                                LogStructured("info", map[string]any{"event":"plane.issue.cnbrpc","op":"set_issue_labels","repo":lk.Repo,"delivery_id":deliveryID,"result":"set"})
                            }
                        }
                    }
                }
            }
            // If new labels now match additional mappings, create missing CNB issues
            for _, m := range mappings {
                if !m.SyncDirection.Valid || strings.ToLower(m.SyncDirection.String) != "bidirectional" { continue }
                if !labelSelectorMatch(m.LabelSelector.String, labels) { continue }
                existing, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID)
                found := false
                for _, lk := range existing { if lk.Repo == m.CNBRepoID { found = true; break } }
                if !found {
                    iid, err := cn.CreateIssue(ctx, m.CNBRepoID, name, descHTML)
                    if err != nil || iid == "" {
                        LogStructured("error", map[string]any{
                            "event":          "plane.issue.cnbrpc",
                            "delivery_id":    deliveryID,
                            "action":         action,
                            "plane_issue_id": planeIssueID,
                            "repo":           m.CNBRepoID,
                            "op":             "create_issue",
                            "error": map[string]any{"code":"cnb_create_failed","message": truncate(fmt.Sprintf("%v", err), 200)},
                        })
                    } else {
                        _ = h.db.CreateIssueLink(ctx, planeIssueID, m.CNBRepoID, iid)
                        LogStructured("info", map[string]any{
                            "event":          "plane.issue.cnbrpc",
                            "delivery_id":    deliveryID,
                            "action":         action,
                            "plane_issue_id": planeIssueID,
                            "repo":           m.CNBRepoID,
                            "op":             "create_issue",
                            "result":         "created",
                            "cnb_issue_iid":  iid,
                        })
                    }
                }
            }
        case "delete", "deleted", "close", "closed":
            if links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links) > 0 {
                for _, lk := range links {
                    if err := cn.CloseIssue(ctx, lk.Repo, lk.Number); err != nil {
                        LogStructured("error", map[string]any{
                            "event":          "plane.issue.cnbrpc",
                            "delivery_id":    deliveryID,
                            "action":         action,
                            "plane_issue_id": planeIssueID,
                            "repo":           lk.Repo,
                            "op":             "close_issue",
                            "error": map[string]any{"code":"cnb_close_failed","message": truncate(fmt.Sprintf("%v", err), 200)},
                        })
                    } else {
                        LogStructured("info", map[string]any{
                            "event":          "plane.issue.cnbrpc",
                            "delivery_id":    deliveryID,
                            "action":         action,
                            "plane_issue_id": planeIssueID,
                            "repo":           lk.Repo,
                            "op":             "close_issue",
                            "result":         "closed",
                        })
                    }
                }
            }
        }
    }

    // Notify Feishu thread if bound
    if planeIssueID != "" {
        if tid, err := h.db.FindLarkThreadByPlaneIssue(ctx, planeIssueID); err == nil && tid != "" {
            summary := "Plane 工作项更新: " + truncate(name, 80)
            if action != "" { summary += " (" + action + ")" }
            go h.sendLarkTextToThread("", tid, summary)
        }
    }
    if deliveryID != "" { _ = h.db.UpdateEventDeliveryStatus(ctx, "plane.issue", deliveryID, "succeeded", nil) }
}

func (h *Handler) handlePlaneIssueComment(env planeWebhookEnvelope, deliveryID string) {
    if !hHasDB(h) { return }
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    data := env.Data
    planeIssueID, _ := dataGetString(data, "issue")
    commentHTML, _ := dataGetString(data, "comment_html")

    if planeIssueID == "" || commentHTML == "" { return }
    // CNB outbound when enabled (fan-out to all linked issues)
    if h.cfg.CNBOutboundEnabled {
        if links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links) > 0 {
            cn := &cnb.Client{
                BaseURL:          h.cfg.CNBBaseURL,
                Token:            h.cfg.CNBAppToken,
                IssueCreatePath:  h.cfg.CNBIssueCreatePath,
                IssueUpdatePath:  h.cfg.CNBIssueUpdatePath,
                IssueCommentPath: h.cfg.CNBIssueCommentPath,
            }
            for _, lk := range links {
                if err := cn.AddComment(ctx, lk.Repo, lk.Number, commentHTML); err != nil {
                    LogStructured("error", map[string]any{
                        "event":          "plane.issue_comment.cnbrpc",
                        "delivery_id":    deliveryID,
                        "plane_issue_id": planeIssueID,
                        "repo":           lk.Repo,
                        "op":             "add_comment",
                        "error": map[string]any{"code":"cnb_comment_failed","message": truncate(fmt.Sprintf("%v", err), 200)},
                    })
                } else {
                    LogStructured("info", map[string]any{
                        "event":          "plane.issue_comment.cnbrpc",
                        "delivery_id":    deliveryID,
                        "plane_issue_id": planeIssueID,
                        "repo":           lk.Repo,
                        "op":             "add_comment",
                        "result":         "commented",
                    })
                }
            }
        }
    }
    // Notify Feishu thread if bound (strip tags, keep short)
    if tid, err := h.db.FindLarkThreadByPlaneIssue(ctx, planeIssueID); err == nil && tid != "" {
        txt := commentHTML
        txt = strings.ReplaceAll(txt, "<br>", "\n")
        txt = stripTags(txt)
        if txt == "" { txt = "(空评论)" }
        msg := "Plane 评论: " + truncate(txt, 140)
        go h.sendLarkTextToThread("", tid, msg)
    }
    if deliveryID != "" { _ = h.db.UpdateEventDeliveryStatus(ctx, "plane.issue_comment", deliveryID, "succeeded", nil) }
}

func dataGetString(m map[string]any, key string) (string, bool) {
    if m == nil { return "", false }
    if v, ok := m[key]; ok {
        switch t := v.(type) {
        case string:
            return t, true
        }
    }
    return "", false
}

// dataGetLabels tries to extract label names from Plane webhook payload
func dataGetLabels(m map[string]any) []string {
    names := make([]string, 0, 8)
    if m == nil { return names }
    if v, ok := m["labels"]; ok {
        if arr, ok := v.([]any); ok {
            for _, it := range arr {
                switch t := it.(type) {
                case map[string]any:
                    if n, ok := t["name"].(string); ok && n != "" { names = append(names, n) }
                case string:
                    if t != "" { names = append(names, t) }
                }
            }
        }
    }
    if len(names) == 0 {
        if v, ok := m["label_names"]; ok {
            if arr, ok := v.([]any); ok {
                for _, it := range arr { if s, ok := it.(string); ok && s != "" { names = append(names, s) } }
            }
        }
    }
    return names
}

// labelSelectorMatch returns true if selector contains any token present in labels (case-insensitive)
func labelSelectorMatch(selector string, labels []string) bool {
    selector = strings.TrimSpace(selector)
    if selector == "" { return false }
    tokens := make([]string, 0, 8)
    for _, p := range strings.FieldsFunc(selector, func(r rune) bool { return r == ',' || r == ' ' || r == ';' || r == '|' }) {
        p = strings.TrimSpace(p)
        if p != "" { tokens = append(tokens, strings.ToLower(p)) }
    }
    if len(tokens) == 0 { return false }
    // wildcard support: '*' or 'all' matches any non-empty label set
    for _, t := range tokens { if t == "*" || t == "all" { return len(labels) > 0 } }
    set := make(map[string]struct{}, len(labels))
    for _, l := range labels { if l != "" { set[strings.ToLower(l)] = struct{}{} } }
    for _, tok := range tokens { if _, ok := set[tok]; ok { return true } }
    return false
}

// preferredReturnURL selects a safe redirect destination back to Plane
// Priority:
// 1) return_to query param if absolute http(s) URL and host is allowed
// 2) state if it looks like an absolute http(s) URL and host is allowed
// 3) PLANE_APP_BASE_URL if configured
// 4) Derived from PLANE_BASE_URL (api.* -> app.*) when applicable
// 5) PLANE_BASE_URL as last resort
func (h *Handler) preferredReturnURL(workspaceSlug, state, returnTo string) string {
    // Allowed hosts
    allowed := map[string]struct{}{}
    addHost := func(raw string) {
        if raw == "" { return }
        if u, err := url.Parse(raw); err == nil && u.Host != "" {
            allowed[strings.ToLower(u.Host)] = struct{}{}
        }
    }
    addHost(h.cfg.PlaneAppBaseURL)
    addHost(h.cfg.PlaneBaseURL)

    isAllowed := func(u *url.URL) bool {
        if u == nil || (u.Scheme != "http" && u.Scheme != "https") { return false }
        host := strings.ToLower(u.Host)
        if _, ok := allowed[host]; ok { return true }
        // If PLANE_BASE_URL host starts with api., accept app.<rest> too
        if base, err := url.Parse(h.cfg.PlaneBaseURL); err == nil {
            bh := strings.ToLower(base.Host)
            if strings.HasPrefix(bh, "api.") {
                alt := "app." + strings.TrimPrefix(bh, "api.")
                if host == alt { return true }
            }
        }
        return false
    }

    // 1) return_to
    if returnTo != "" {
        if u, err := url.Parse(returnTo); err == nil && isAllowed(u) {
            return u.String()
        }
    }
    // 2) workspace integrations page (needs workspace slug)
    if workspaceSlug != "" {
        if base := h.planeAppBase(); base != nil {
            base.Path = strings.TrimRight(base.Path, "/") + "/" + workspaceSlug + "/settings/integrations/"
            base.RawQuery = ""
            base.Fragment = ""
            return base.String()
        }
    }
    // 3) state as URL
    if state != "" {
        if u, err := url.Parse(state); err == nil && isAllowed(u) {
            return u.String()
        }
    }
    // 4) PLANE_APP_BASE_URL
    if h.cfg.PlaneAppBaseURL != "" {
        return h.cfg.PlaneAppBaseURL
    }
    // 5) derive from PLANE_BASE_URL (api.* -> app.*)
    if u, err := url.Parse(h.cfg.PlaneBaseURL); err == nil {
        if strings.HasPrefix(strings.ToLower(u.Host), "api.") {
            u.Host = "app." + strings.TrimPrefix(u.Host, "api.")
            u.Path = "/"
            u.RawQuery = ""
            u.Fragment = ""
            return u.String()
        }
        // 6) fallback
        u.Path = "/"
        u.RawQuery = ""
        u.Fragment = ""
        return u.String()
    }
    return "/" // final fallback
}

func (h *Handler) buildRedirectHTML(target string, payload map[string]any) string {
    // Minimal, clean handoff page with auto-redirect
    _ = payload // intentionally unused (no debug output on the page)
    esc := func(s string) string {
        r := strings.NewReplacer(
            "&", "&amp;",
            "<", "&lt;",
            ">", "&gt;",
            "\"", "&quot;",
            "'", "&#39;",
        )
        return r.Replace(s)
    }
    return "<!DOCTYPE html><html lang=\"zh-CN\"><head>" +
        "<meta charset=\"utf-8\">" +
        "<meta http-equiv=\"refresh\" content=\"2; url=" + esc(target) + "\">" +
        "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">" +
        "<meta name=\"color-scheme\" content=\"light dark\">" +
        "<title>返回 Plane</title>" +
        "<style>" +
        "html,body{height:100%}body{margin:0;font:16px/1.6 -apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;display:grid;place-items:center}" +
        ".wrap{text-align:center;padding:24px}h1{font-size:20px;margin:0 0 8px}p{margin:6px 0;color:#6b7280}" +
        ".spinner{width:28px;height:28px;border:3px solid currentColor;border-right-color:transparent;border-radius:50%;margin:14px auto;animation:s .8s linear infinite;opacity:.7}" +
        "a.btn{display:inline-block;margin-top:10px;padding:8px 12px;border:1px solid currentColor;border-radius:8px;text-decoration:none}" +
        "@keyframes s{to{transform:rotate(360deg)}}" +
        "</style></head><body>" +
        "<div class=\"wrap\">" +
        "<div class=\"spinner\" aria-hidden=\"true\"></div>" +
        "<h1>安装完成，正在返回 Plane…</h1>" +
        "<p>若未自动跳转，请点击下方按钮</p>" +
        "<p><a class=\"btn\" href=\"" + esc(target) + "\">返回 Plane</a></p>" +
        "</div>" +
        "<script>(function(){try{var t='" + esc(target) + "'; if(window.opener && !window.opener.closed){try{window.opener.postMessage({type:'plane_installation',status:'ok',target:t}, '*');}catch(e){}} window.location.replace(t);}catch(e){}}</script>" +
        "</body></html>"
}

// planeAppBase returns a parsed URL for PLANE_APP_BASE_URL or derives from PLANE_BASE_URL (api.* -> app.*)
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

// stripTags removes HTML tags (best-effort, not full HTML sanitizer)
func stripTags(s string) string {
    out := make([]rune, 0, len(s))
    inTag := false
    for _, r := range s {
        switch r {
        case '<':
            inTag = true
        case '>':
            inTag = false
        default:
            if !inTag { out = append(out, r) }
        }
    }
    return strings.TrimSpace(string(out))
}
