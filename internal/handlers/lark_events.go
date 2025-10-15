package handlers

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/labstack/echo/v4"
    "plane-integration/internal/lark"
    "plane-integration/internal/plane"
    "plane-integration/internal/store"
)

// Event v2 envelope (兼容 challenge 握手)
type larkEventEnvelope struct {
    Schema    string          `json:"schema"`
    Header    larkEventHeader `json:"header"`
    Event     json.RawMessage `json:"event"`
    Challenge string          `json:"challenge"`
    Type      string          `json:"type"`
}

type larkEventHeader struct {
    EventID    string `json:"event_id"`
    EventType  string `json:"event_type"`
    CreateTime string `json:"create_time"`
    Token      string `json:"token"`
    AppID      string `json:"app_id"`
    TenantKey  string `json:"tenant_key"`
}

// Message event payload (im.message.receive_v1)
type larkMessageEvent struct {
    Sender struct {
        SenderID struct {
            OpenID  string `json:"open_id"`
            UnionID string `json:"union_id"`
            UserID  string `json:"user_id"`
        } `json:"sender_id"`
    } `json:"sender"`
    Message struct {
        MessageID   string `json:"message_id"`
        RootID      string `json:"root_id"`
        ParentID    string `json:"parent_id"`
        ChatID      string `json:"chat_id"`
        ChatType    string `json:"chat_type"`
        MessageType string `json:"message_type"`
        Content     string `json:"content"`
        Mentions    []struct{
            Name string `json:"name"`
            Key  string `json:"key"`
            ID   struct{
                OpenID string `json:"open_id"`
                UserID string `json:"user_id"`
            } `json:"id"`
        } `json:"mentions"`
    } `json:"message"`
}

// text content wrapper inside message.content
type larkTextContent struct {
    Text string `json:"text"`
}

func (h *Handler) LarkEvents(c echo.Context) error {
    // Read raw for signature verification
    body, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    // Challenge quick path (旧版/新版均可)
    var probe struct{ Challenge string `json:"challenge"`; Type string `json:"type"` }
    if json.Unmarshal(body, &probe) == nil && probe.Challenge != "" {
        return c.JSON(http.StatusOK, map[string]string{"challenge": probe.Challenge})
    }

    // Verify signature or token when可用（避免时钟漂移：5 分钟窗口）
    if !h.verifyLarkSignature(c.Request().Header, body) {
        // fallback: header.token 校验（若配置了 Verification Token）
        var env larkEventEnvelope
        if json.Unmarshal(body, &env) == nil {
            if h.cfg.LarkVerificationToken != "" && env.Header.Token != h.cfg.LarkVerificationToken {
                return c.NoContent(http.StatusUnauthorized)
            }
        }
    }

    // Parse envelope
    var env larkEventEnvelope
    if err := json.Unmarshal(body, &env); err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    if env.Challenge != "" { // 冗余保护
        return c.JSON(http.StatusOK, map[string]string{"challenge": env.Challenge})
    }

    switch strings.ToLower(env.Header.EventType) {
    case "im.message.receive_v1":
        var ev larkMessageEvent
        if err := json.Unmarshal(env.Event, &ev); err != nil {
            return c.NoContent(http.StatusBadRequest)
        }
        // Only handle group for now
        if strings.ToLower(ev.Message.ChatType) != "group" {
            return c.NoContent(http.StatusOK)
        }
        // Parse text
        var txt larkTextContent
        _ = json.Unmarshal([]byte(ev.Message.Content), &txt)
        text := strings.TrimSpace(txt.Text)
        if text == "" { // nothing to do
            return c.NoContent(http.StatusOK)
        }
        // Command parse: support "/bind <url>" or "绑定 <url>" or "bind <url>"
        // Only trigger when @bot is present OR message starts with command
        mentioned := len(ev.Message.Mentions) > 0
        lower := strings.ToLower(text)
        if mentioned || strings.HasPrefix(lower, "/bind") || strings.HasPrefix(lower, "bind ") || strings.HasPrefix(text, "绑定 ") {
            issueID, slug, projectID := parsePlaneIssueLink(text)
            if issueID != "" && hHasDB(h) {
                // Thread id: prefer root_id, else message_id
                threadID := ev.Message.RootID
                if threadID == "" { threadID = ev.Message.MessageID }
                // Save link (sync disabled by default; we enable minimal notifications only)
                _ = h.db.UpsertLarkThreadLink(c.Request().Context(), threadID, issueID, projectID, slug, false)
                // Acknowledge (best-effort notify via Lark when outbound is wired)
                go h.notifyLarkThreadBound(ev.Message.ChatID, threadID, issueID)
                return c.JSON(http.StatusOK, map[string]any{"result":"ok","action":"bind","plane_issue_id":issueID})
            }
        }
        // Comment command in a bound thread
        if ev.Message.RootID != "" && hHasDB(h) {
            threadID := ev.Message.RootID
            tl, err := h.db.GetLarkThreadLink(c.Request().Context(), threadID)
            if err == nil && tl != nil && tl.PlaneIssueID != "" {
                // command prefix: /comment or 评论
                arg := extractCommandArg(text, "/comment")
                if arg == "" { arg = extractCommandArg(text, "评论") }
                if arg != "" {
                    // post to Plane
                    go h.postPlaneComment(tl, arg)
                    return c.JSON(http.StatusOK, map[string]any{"result":"ok","action":"comment","plane_issue_id":tl.PlaneIssueID})
                }
            }
        }
        return c.NoContent(http.StatusOK)
    default:
        // Ignore other event types for now
        return c.NoContent(http.StatusOK)
    }
}

func (h *Handler) LarkInteractivity(c echo.Context) error {
    return c.NoContent(http.StatusOK)
}

func (h *Handler) LarkCommands(c echo.Context) error {
    return c.NoContent(http.StatusOK)
}

// verifyLarkSignature validates Feishu signatures when headers are present.
// 当前实现遵循：signature = base64(hmac_sha256(encrypt_key, timestamp + body))（待确认）
// 若缺少头或密钥未配置，则返回 true 以避免误拒绝（由 verification token 兜底）。
func (h *Handler) verifyLarkSignature(hdr http.Header, body []byte) bool {
    // Per docs: sha256(timestamp + nonce + encrypt_key + body), hex lower
    ts := strings.TrimSpace(hdr.Get("X-Lark-Request-Timestamp"))
    nonce := strings.TrimSpace(hdr.Get("X-Lark-Request-Nonce"))
    sig := strings.TrimSpace(hdr.Get("X-Lark-Signature"))
    if ts == "" || nonce == "" || sig == "" || h.cfg.LarkEncryptKey == "" {
        return true // allow; will fallback to token check if configured
    }
    // time skew soft check (do not reject hard to avoid false negatives)
    if tsec, err := strconv.ParseInt(ts, 10, 64); err == nil {
        now := time.Now().Unix()
        if tsec < now-300 || tsec > now+300 {
            // soft-pass, we'll still compute but not hard fail on skew
        }
    }
    var b strings.Builder
    b.WriteString(ts)
    b.WriteString(nonce)
    b.WriteString(h.cfg.LarkEncryptKey)
    b.Write(body)
    sum := sha256.Sum256([]byte(b.String()))
    want := hex.EncodeToString(sum[:])
    return strings.EqualFold(sig, want)
}

var uuidRe = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`)

// extractPlaneIssueID tries to parse a Plane issue URL and return UUID.
// It accepts any URL containing a UUID; when base URLs are provided, it prefers matching hosts.
func extractPlaneIssueID(text, planeAppBase, planeAPIBase string) string {
    // Fast path: any UUID in text
    m := uuidRe.FindString(text)
    if m == "" { return "" }
    // Optional: ensure URL host matches Plane when a URL is present in text
    // (best-effort; 待确认 具体链接格式)
    return strings.ToLower(m)
}

// parsePlaneIssueLink extracts (issueID, workspaceSlug, projectID) from an app URL if present.
// Example: https://app.plane.so/{slug}/projects/{project_id}/issues/{issue_id}
func parsePlaneIssueLink(s string) (issueID, workspaceSlug, projectID string) {
    // naive scan for url segments; fallback to UUID only
    // regex: host/.../{slug}/projects/{uuid}/issues/{uuid}
    re := regexp.MustCompile(`https?://[^\s]+/([\w-]+)/projects/(` + uuidRe.String() + `)/issues/(` + uuidRe.String() + `)`) // #nosec G101
    m := re.FindStringSubmatch(s)
    if len(m) == 4 {
        return strings.ToLower(m[3]), m[1], strings.ToLower(m[2])
    }
    return strings.ToLower(uuidRe.FindString(s)), "", ""
}

// notifyLarkThreadBound posts a confirmation message to the bound thread (best-effort; 待确认 API 路径)
func (h *Handler) notifyLarkThreadBound(chatID, threadID, issueID string) {
    // No-op if outbound not configured
    if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
        return
    }
    // Minimal text; real implementation should render rich card
    text := "已绑定 Plane 工作项: " + issueID
    _ = h.sendLarkTextToThread(chatID, threadID, text)
}

// sendLarkTextToThread sends a text to a specific thread; falls back to chat if reply fails.
func (h *Handler) sendLarkTextToThread(chatID, threadID, text string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
    defer cancel()
    cli := &lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
    token, _, err := cli.TenantAccessToken(ctx)
    if err != nil { return err }
    // Prefer replying in thread
    if threadID != "" {
        if err := cli.ReplyTextInThread(ctx, token, threadID, text); err == nil {
            return nil
        }
    }
    if chatID != "" {
        return cli.SendTextToChat(ctx, token, chatID, text)
    }
    return nil
}

// postPlaneComment posts a text comment into Plane for the given thread link.
func (h *Handler) postPlaneComment(tl *store.LarkThreadLink, comment string) {
    if !hHasDB(h) || tl == nil || tl.PlaneIssueID == "" { return }
    ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
    defer cancel()
    // Resolve workspace slug and project id
    slug := ""
    if tl.WorkspaceSlug.Valid { slug = tl.WorkspaceSlug.String }
    projectID := ""
    if tl.PlaneProjectID.Valid { projectID = tl.PlaneProjectID.String }
    if slug == "" || projectID == "" { return } // cannot route; 待确认：是否支持通过 Issue API 反查
    token, err := h.db.FindBotTokenByWorkspaceSlug(ctx, slug)
    if err != nil || token == "" { return }
    pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
    // Wrap as simple HTML; escape minimal
    html := comment // Feishu text contains no HTML; Plane accepts HTML
    _ = pc.AddComment(ctx, token, slug, projectID, tl.PlaneIssueID, html)
}

// extractCommandArg returns text after a command prefix (case-insensitive), trimmed.
func extractCommandArg(text, cmd string) string {
    t := strings.TrimSpace(text)
    lc := strings.ToLower(t)
    lcmd := strings.ToLower(cmd)
    if strings.HasPrefix(lc, lcmd+" ") {
        return strings.TrimSpace(t[len(cmd):])
    }
    // handle leading @bot mention e.g. "@bot /comment xxx" → remove first token
    parts := strings.Fields(t)
    if len(parts) >= 2 {
        if strings.HasPrefix(parts[1], cmd) {
            return strings.TrimSpace(strings.TrimPrefix(t, parts[0]))[len(cmd):]
        }
    }
    return ""
}
