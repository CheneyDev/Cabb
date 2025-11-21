package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cabb/internal/lark"
	"cabb/internal/plane"
	"cabb/internal/store"

	"github.com/labstack/echo/v4"
)

// Max acceptable age for incoming Lark events (seconds)
const larkEventMaxAgeSeconds = 60

// retryUpdateInteractiveCard updates a card via callback token with bounded retries and persists states into DB.
// Backoff: 0s, 5s, 15s, 30s (4 attempts total). Writes event_deliveries with source=lark.card.update.
func (h *Handler) retryUpdateInteractiveCard(callbackToken string, card map[string]any) {
	if strings.TrimSpace(callbackToken) == "" || card == nil {
		return
	}
	// Prepare identifiers
	b, _ := json.Marshal(card)
	sum := sha256.Sum256(b)
	payloadSHA := hex.EncodeToString(sum[:])
	// Record/initialize delivery row
	if hHasDB(h) {
		_ = h.db.UpsertEventDelivery(context.Background(), "lark.card.update", "card.update", callbackToken, payloadSHA, "queued")
	}
	backoffs := []time.Duration{0, 5 * time.Second, 15 * time.Second, 30 * time.Second}
	cli := &lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
	for i, d := range backoffs {
		if d > 0 {
			time.Sleep(d)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		token, _, err := cli.TenantAccessToken(ctx)
		if err == nil {
			err = cli.UpdateInteractiveCard(ctx, token, callbackToken, card)
		}
		cancel()
		if err == nil {
			LogStructured("info", map[string]any{"event": "lark.card.update.ok", "attempt": i + 1})
			if hHasDB(h) {
				_ = h.db.UpdateEventDeliveryStatus(context.Background(), "lark.card.update", callbackToken, "succeeded", nil)
			}
			return
		}
		// Schedule next retry if any
		LogStructured("error", map[string]any{"event": "lark.card.update.fail", "attempt": i + 1, "error": err.Error()})
		if hHasDB(h) {
			if i < len(backoffs)-1 {
				next := time.Now().Add(backoffs[i+1])
				_ = h.db.UpdateEventDeliveryRetry(context.Background(), "lark.card.update", callbackToken, next)
			} else {
				_ = h.db.UpdateEventDeliveryStatus(context.Background(), "lark.card.update", callbackToken, "failed", nil)
			}
		}
	}
}

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
		Mentions    []struct {
			Name string `json:"name"`
			Key  string `json:"key"`
			ID   struct {
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
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Challenge quick path (旧版/新版均可)
	var probe struct {
		Challenge string `json:"challenge"`
		Type      string `json:"type"`
	}
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

	var env larkEventEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if env.Challenge != "" {
		return c.JSON(http.StatusOK, map[string]string{"challenge": env.Challenge})
	}

	evtID := strings.TrimSpace(env.Header.EventID)
	sum := sha256.Sum256(body)
	payloadSHA := hex.EncodeToString(sum[:])
	if evtID != "" && h.dedupe != nil && h.dedupe.CheckAndMark("lark.events", evtID, payloadSHA) {
		LogStructured("info", map[string]any{"event": "lark.event.duplicate.mem", "event_id": evtID})
		return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "duplicate"})
	}
	if hHasDB(h) && evtID != "" {
		if dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "lark.events", evtID, payloadSHA); err == nil && dup {
			LogStructured("info", map[string]any{"event": "lark.event.duplicate.db", "event_id": evtID})
			return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "duplicate"})
		}
		_ = h.db.UpsertEventDelivery(c.Request().Context(), "lark.events", strings.ToLower(env.Header.EventType), evtID, payloadSHA, "seen")
	}
	if larkEventMaxAgeSeconds > 0 {
		maxAge := larkEventMaxAgeSeconds
		ct := strings.TrimSpace(env.Header.CreateTime)
		if ct != "" {
			if ts, err := strconv.ParseInt(ct, 10, 64); err == nil {
				if ts > 1_000_000_000_000 {
					ts = ts / 1000
				}
				evtT := time.Unix(ts, 0)
				age := time.Since(evtT)
				if age > time.Duration(maxAge)*time.Second {
					LogStructured("info", map[string]any{"event": "lark.event.expired", "event_id": evtID, "age_ms": age.Milliseconds(), "max_age_s": maxAge})
					if hHasDB(h) && evtID != "" {
						_ = h.db.UpdateEventDeliveryStatus(c.Request().Context(), "lark.events", evtID, "expired", nil)
					}
					return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "expired"})
				}
			}
		}
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
		if text == "" {
			return c.NoContent(http.StatusOK)
		}
		lower := strings.ToLower(text)
		isBind := strings.HasPrefix(lower, "/bind") || strings.HasPrefix(lower, "bind ") || strings.HasPrefix(text, "绑定 ")
		if !isBind && len(ev.Message.Mentions) > 0 {
			parts := strings.Fields(text)
			if len(parts) >= 2 {
				sec := strings.ToLower(parts[1])
				if strings.HasPrefix(sec, "/bind") || strings.HasPrefix(sec, "bind") || strings.HasPrefix(sec, "绑定") {
					isBind = true
				}
			}
		}

		LogStructured("info", map[string]any{
			"event":          "lark.message.debug",
			"text":           text,
			"lower":          lower,
			"is_bind":        isBind,
			"mentions_count": len(ev.Message.Mentions),
			"chat_id":        ev.Message.ChatID,
			"message_type":   ev.Message.MessageType,
		})
		if isBind {
			// Determine thread target
			threadID := ev.Message.RootID
			if threadID == "" {
				threadID = ev.Message.MessageID
			}

			issueID, slug, projectID := parsePlaneIssueLink(text)
			LogStructured("info", map[string]any{
				"event":        "bind.parse",
				"chat_id":      ev.Message.ChatID,
				"thread_id":    threadID,
				"has_issue_id": issueID != "",
				"has_slug":     slug != "",
				"has_project":  projectID != "",
			})

			// Resolve short link KEY-N via Plane API when possible
			if issueID == "" {
				// Try resolving via workspace browse short link: /{slug}/browse/{KEY-N}
				if slug != "" {
					if seq := extractBrowseSequence(text); seq != "" && hHasDB(h) {
						LogStructured("info", map[string]any{"event": "bind.resolve_shortlink.start", "chat_id": ev.Message.ChatID, "thread_id": threadID, "slug": slug, "seq": seq})
						rctx, cancel := context.WithTimeout(c.Request().Context(), 8*time.Second)
						defer cancel()
						token, _ := h.ensurePlaneBotToken(rctx, slug)
						if token != "" {
							pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
							if iid, pid, err := pc.GetIssueBySequence(rctx, token, slug, seq); err == nil {
								issueID, projectID = iid, pid
								LogStructured("info", map[string]any{"event": "bind.resolve_shortlink.ok", "chat_id": ev.Message.ChatID, "thread_id": threadID, "issue_id": issueID, "project_id": projectID})
							}
						} else {
							LogStructured("warn", map[string]any{"event": "bind.resolve_shortlink.no_token", "chat_id": ev.Message.ChatID, "thread_id": threadID, "slug": slug})
						}
					}
				}
			}
			if issueID == "" {
				if seq := extractBrowseSequence(text); seq != "" {
					go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "绑定失败：无法解析短链接 "+seq+"。请先在本服务完成 Plane 应用安装（获取 bot token），或改用完整链接：/bind https://app.plane.so/{slug}/projects/{project}/issues/{issue}")
					LogStructured("warn", map[string]any{"event": "bind.error.shortlink_unresolved", "chat_id": ev.Message.ChatID, "thread_id": threadID, "slug": slug, "seq": seq})
				} else {
					go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "绑定失败：未检测到 Plane 工作项链接或 UUID。示例：/bind https://app.plane.so/{slug}/projects/{project}/issues/{issue}")
					LogStructured("warn", map[string]any{"event": "bind.error.missing_issue_id", "chat_id": ev.Message.ChatID, "thread_id": threadID})
				}
				return c.JSON(http.StatusOK, map[string]any{"result": "error", "action": "bind", "error": "missing_issue_id"})
			}
			if !hHasDB(h) {
				go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "绑定失败：服务未连接数据库，请稍后重试或联系管理员。")
				LogStructured("error", map[string]any{"event": "bind.error.db_unavailable", "chat_id": ev.Message.ChatID, "thread_id": threadID})
				return c.JSON(http.StatusOK, map[string]any{"result": "error", "action": "bind", "plane_issue_id": issueID, "error": "db_unavailable"})
			}
			// Prevent duplicate binding within the same chat
			if hHasDB(h) && ev.Message.ChatID != "" {
				if cl, err := h.db.GetLarkChatIssueLink(c.Request().Context(), ev.Message.ChatID); err == nil && cl != nil && cl.PlaneIssueID != "" {
					// Same issue already bound for this chat
					if strings.EqualFold(cl.PlaneIssueID, issueID) {
						// Reply "ISSUE {title} 已绑定" and do not create duplicate thread link
						// Prefer existing mapping's project/slug for title/url
						slugEff := nsToString(cl.WorkspaceSlug)
						projEff := nsToString(cl.PlaneProjectID)
						go func() { _ = h.postBoundAlready(ev.Message.ChatID, threadID, slugEff, projEff, issueID) }()
						LogStructured("info", map[string]any{"event": "bind.duplicate.same_issue", "chat_id": ev.Message.ChatID, "thread_id": threadID, "issue_id": issueID})
						return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "bind", "status": "duplicate", "plane_issue_id": issueID})
					}
					// Different issue already bound → overwrite directly (no confirmation needed)
					LogStructured("info", map[string]any{"event": "bind.different.overwrite", "chat_id": ev.Message.ChatID, "thread_id": threadID, "current_issue_id": cl.PlaneIssueID, "new_issue_id": issueID})
					// Fall through to persist new binding
				}
			}

			// Persist link
			if err := h.db.UpsertLarkThreadLink(c.Request().Context(), threadID, issueID, projectID, slug, false); err != nil {
				LogStructured("error", map[string]any{"event": "bind.persist.thread_link.error", "chat_id": ev.Message.ChatID, "thread_id": threadID, "issue_id": issueID, "error": err.Error()})
				go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "绑定失败：内部错误，请稍后重试。")
				return c.JSON(http.StatusOK, map[string]any{"result": "error", "action": "bind", "plane_issue_id": issueID, "error": "upsert_failed"})
			}
			// Also bind chat -> issue for out-of-thread commands
			_ = h.db.UpsertLarkChatIssueLink(c.Request().Context(), ev.Message.ChatID, threadID, issueID, projectID, slug)
			// Success ack with details; prefer rich post with anchor title when possible
			go h.postBindAck(ev.Message.ChatID, threadID, slug, projectID, issueID)
			// Post share link to Plane issue as comment
			go h.postChatShareLinkToPlane(ev.Message.ChatID, slug, projectID, issueID)
			LogStructured("info", map[string]any{"event": "bind.persist.ok", "chat_id": ev.Message.ChatID, "thread_id": threadID, "issue_id": issueID, "project_id": projectID, "slug": slug})
			return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "bind", "plane_issue_id": issueID})
		}
		// Comment/sync commands and auto-sync in a bound thread
		if ev.Message.RootID != "" && hHasDB(h) {
			threadID := ev.Message.RootID
			tl, err := h.db.GetLarkThreadLink(c.Request().Context(), threadID)
			if err == nil && tl != nil && tl.PlaneIssueID != "" {
				arg := extractCommandArg(text, "/comment")
				if arg == "" {
					arg = extractCommandArg(text, "评论")
				}
				if arg != "" {
					trimmed := strings.TrimSpace(arg)
					go h.postPlaneCommentWithRetry(tl, trimmed, ev.Message.MessageID)
					// Best-effort user feedback in thread
					go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "评论已同步至 Plane")
					return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "comment", "plane_issue_id": tl.PlaneIssueID})
				}

				lct := strings.ToLower(strings.TrimSpace(text))
				if strings.HasPrefix(lct, "/sync") || strings.HasPrefix(text, "开启同步") || strings.HasPrefix(text, "关闭同步") {
					enable := false
					if strings.HasPrefix(lct, "/sync on") || strings.HasPrefix(text, "开启同步") {
						enable = true
					} else if strings.HasPrefix(lct, "/sync off") || strings.HasPrefix(text, "关闭同步") {
						enable = false
					} else {
						go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "用法：/sync on 开启线程自动同步；/sync off 关闭自动同步。也可发送‘开启同步’或‘关闭同步’。")
						return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "sync_help"})
					}
					slug := ""
					if tl.WorkspaceSlug.Valid {
						slug = tl.WorkspaceSlug.String
					}
					projectID := ""
					if tl.PlaneProjectID.Valid {
						projectID = tl.PlaneProjectID.String
					}
					if err := h.db.UpsertLarkThreadLink(c.Request().Context(), threadID, tl.PlaneIssueID, projectID, slug, enable); err != nil {
						go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "设置失败：内部错误，请稍后重试。")
						return c.JSON(http.StatusOK, map[string]any{"result": "error", "action": "sync_toggle", "error": "upsert_failed"})
					}
					if enable {
						go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "已开启线程自动同步：该线程新消息将自动同步为 Plane 评论。")
					} else {
						go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "已关闭线程自动同步：该线程新消息将不再同步至 Plane。")
					}
					return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "sync_toggle", "enabled": enable})
				}

				// 3) auto-sync: non-command text in a bound thread when enabled
				if tl.SyncEnabled && ev.Message.MessageType == "text" {
					// skip slash-command messages
					if !strings.HasPrefix(lct, "/") && strings.TrimSpace(text) != "" {
						go h.postPlaneCommentWithRetry(tl, text, ev.Message.MessageID)
						return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "sync_auto", "plane_issue_id": tl.PlaneIssueID})
					}
				}
			}
		}
		if hHasDB(h) && ev.Message.RootID == "" && ev.Message.ChatID != "" {
			lct := strings.ToLower(strings.TrimSpace(text))
			if strings.HasPrefix(lct, "/sync") || strings.HasPrefix(text, "开启同步") || strings.HasPrefix(text, "关闭同步") {
				cl, err := h.db.GetLarkChatIssueLink(c.Request().Context(), ev.Message.ChatID)
				if err == nil && cl != nil && cl.PlaneIssueID != "" && cl.LarkThreadID.Valid {
					enable := false
					if strings.HasPrefix(lct, "/sync on") || strings.HasPrefix(text, "开启同步") {
						enable = true
					}
					if strings.HasPrefix(lct, "/sync off") || strings.HasPrefix(text, "关闭同步") { /* keep false */
					}
					slug := ""
					if cl.WorkspaceSlug.Valid {
						slug = cl.WorkspaceSlug.String
					}
					pid := ""
					if cl.PlaneProjectID.Valid {
						pid = cl.PlaneProjectID.String
					}
					if err := h.db.UpsertLarkThreadLink(c.Request().Context(), cl.LarkThreadID.String, cl.PlaneIssueID, pid, slug, enable); err == nil {
						if enable {
							go h.sendLarkTextToThread(ev.Message.ChatID, cl.LarkThreadID.String, "已开启线程自动同步")
						} else {
							go h.sendLarkTextToThread(ev.Message.ChatID, cl.LarkThreadID.String, "已关闭线程自动同步")
						}
						return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "sync_toggle", "enabled": enable, "scope": "chat_fallback"})
					}
				}
			}
		}
		if arg := extractCommandArg(text, "/comment"); arg != "" {
			if hHasDB(h) && ev.Message.ChatID != "" {
				if cl, err := h.db.GetLarkChatIssueLink(c.Request().Context(), ev.Message.ChatID); err == nil && cl != nil && cl.PlaneIssueID != "" {
					tl := &store.LarkThreadLink{
						LarkThreadID:   nsToString(cl.LarkThreadID),
						PlaneIssueID:   cl.PlaneIssueID,
						PlaneProjectID: cl.PlaneProjectID,
						WorkspaceSlug:  cl.WorkspaceSlug,
						SyncEnabled:    false,
					}
					go h.postPlaneCommentWithRetry(tl, strings.TrimSpace(arg), ev.Message.MessageID)
					t := nsToString(cl.LarkThreadID)
					go h.sendLarkTextToThread(ev.Message.ChatID, t, "评论已同步至 Plane")
					return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "comment", "plane_issue_id": cl.PlaneIssueID, "scope": "chat_fallback"})
				}
			}
			// fallback guidance
			threadID := ev.Message.RootID
			if threadID == "" {
				threadID = ev.Message.MessageID
			}
			go h.sendLarkTextToThread(ev.Message.ChatID, threadID, "未绑定目标工作项。请先使用 /bind 绑定，或在绑定的话题中回复 /comment。")
			return c.JSON(http.StatusOK, map[string]any{"result": "error", "action": "comment", "error": "no_binding"})
		}
		if len(ev.Message.Mentions) > 0 && isLikelyPureBotMention(text) {
			threadID := ev.Message.RootID
			if threadID == "" {
				threadID = ev.Message.MessageID
			}
			go func() { _ = h.postBotHelpCard(ev.Message.ChatID, threadID) }()
			return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "help_card"})
		}
		return c.NoContent(http.StatusOK)
	case "im.chat.member.bot.added_v1":
		var ev struct {
			ChatID string `json:"chat_id"`
		}
		if err := json.Unmarshal(env.Event, &ev); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		chatID := strings.TrimSpace(ev.ChatID)
		if chatID == "" {
			return c.NoContent(http.StatusOK)
		}
		go func() { _ = h.postBotHelpCard(chatID, "") }()
		return c.JSON(http.StatusOK, map[string]any{"result": "ok", "action": "bot_added"})
	case "card.action.trigger", "card.action.trigger_v1":
		var act struct {
			Token  string `json:"token"`
			Action struct {
				Value map[string]any `json:"value"`
			} `json:"action"`
		}
		_ = json.Unmarshal(env.Event, &act)
		LogStructured("info", map[string]any{"event": "lark.events.card_action.receive"})
		return h.handleLarkCardAction(c, act.Action.Value, act.Token)
	default:
		// Ignore other event types for now
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) LarkInteractivity(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	LogStructured("info", map[string]any{"event": "lark.interactivity.receive", "content_length": len(body)})
	// Challenge support
	var probe struct {
		Challenge string `json:"challenge"`
		Type      string `json:"type"`
	}
	if json.Unmarshal(body, &probe) == nil && probe.Challenge != "" {
		LogStructured("info", map[string]any{"event": "lark.interactivity.challenge"})
		return c.JSON(http.StatusOK, map[string]string{"challenge": probe.Challenge})
	}
	if !h.verifyLarkSignature(c.Request().Header, body) {
		var env larkEventEnvelope
		if json.Unmarshal(body, &env) == nil {
			if h.cfg.LarkVerificationToken != "" && env.Header.Token != h.cfg.LarkVerificationToken {
				LogStructured("warn", map[string]any{"event": "lark.interactivity.unauthorized", "reason": "token_mismatch"})
				return c.NoContent(http.StatusUnauthorized)
			}
		}
	}
	// Dedupe + lateness guard when we can parse envelope
	var baseEnv larkEventEnvelope
	_ = json.Unmarshal(body, &baseEnv)
	if baseEnv.Header.EventID != "" {
		evtID := strings.TrimSpace(baseEnv.Header.EventID)
		sum := sha256.Sum256(body)
		payloadSHA := hex.EncodeToString(sum[:])
		if h.dedupe != nil && h.dedupe.CheckAndMark("lark.interactivity", evtID, payloadSHA) {
			LogStructured("info", map[string]any{"event": "lark.interactivity.duplicate.mem", "event_id": evtID})
			return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "duplicate"})
		}
		if hHasDB(h) {
			if dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "lark.interactivity", evtID, payloadSHA); err == nil && dup {
				LogStructured("info", map[string]any{"event": "lark.interactivity.duplicate.db", "event_id": evtID})
				return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "duplicate"})
			}
			_ = h.db.UpsertEventDelivery(c.Request().Context(), "lark.interactivity", strings.ToLower(baseEnv.Header.EventType), evtID, payloadSHA, "seen")
		}
		if larkEventMaxAgeSeconds > 0 {
			ct := strings.TrimSpace(baseEnv.Header.CreateTime)
			if ct != "" {
				if ts, err := strconv.ParseInt(ct, 10, 64); err == nil {
					if ts > 1_000_000_000_000 {
						ts = ts / 1000
					}
					evtT := time.Unix(ts, 0)
					age := time.Since(evtT)
					if age > time.Duration(larkEventMaxAgeSeconds)*time.Second {
						LogStructured("info", map[string]any{"event": "lark.interactivity.expired", "event_id": evtID, "age_ms": age.Milliseconds()})
						if hHasDB(h) {
							_ = h.db.UpdateEventDeliveryStatus(c.Request().Context(), "lark.interactivity", evtID, "expired", nil)
						}
						return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "expired"})
					}
				}
			}
		}
	}
	// Try envelope form first (im.message.card.action.trigger)
	var env larkEventEnvelope
	if json.Unmarshal(body, &env) == nil && len(env.Event) > 0 {
		LogStructured("info", map[string]any{"event": "lark.interactivity.parsed", "event_type": env.Header.EventType})
		// Minimal action struct with callback token/context
		var act struct {
			Token  string `json:"token"`
			Action struct {
				Value         map[string]any `json:"value"`
				OpenMessageID string         `json:"open_message_id"`
			} `json:"action"`
			Context struct {
				OpenChatID    string `json:"open_chat_id"`
				OpenMessageID string `json:"open_message_id"`
			} `json:"context"`
		}
		_ = json.Unmarshal(env.Event, &act)
		return h.handleLarkCardAction(c, act.Action.Value, act.Token)
	}
	// Fallback: direct action payload
	var payload struct {
		Action struct {
			Value map[string]any `json:"value"`
		} `json:"action"`
	}
	if json.Unmarshal(body, &payload) == nil && payload.Action.Value != nil {
		LogStructured("info", map[string]any{"event": "lark.interactivity.value_only"})
		return h.handleLarkCardAction(c, payload.Action.Value, "")
	}
	LogStructured("warn", map[string]any{"event": "lark.interactivity.ignored"})
	return c.NoContent(http.StatusOK)
}

func (h *Handler) LarkCommands(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if !h.verifyLarkSignature(c.Request().Header, body) {
		var env larkEventEnvelope
		if json.Unmarshal(body, &env) == nil {
			if h.cfg.LarkVerificationToken != "" && env.Header.Token != h.cfg.LarkVerificationToken {
				return c.NoContent(http.StatusUnauthorized)
			}
		}
	}
	var env larkEventEnvelope
	_ = json.Unmarshal(body, &env)
	evtID := strings.TrimSpace(env.Header.EventID)
	if evtID != "" {
		sum := sha256.Sum256(body)
		payloadSHA := hex.EncodeToString(sum[:])
		if h.dedupe != nil && h.dedupe.CheckAndMark("lark.commands", evtID, payloadSHA) {
			LogStructured("info", map[string]any{"event": "lark.commands.duplicate.mem", "event_id": evtID})
			return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "duplicate"})
		}
		if hHasDB(h) {
			if dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "lark.commands", evtID, payloadSHA); err == nil && dup {
				LogStructured("info", map[string]any{"event": "lark.commands.duplicate.db", "event_id": evtID})
				return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "duplicate"})
			}
			_ = h.db.UpsertEventDelivery(c.Request().Context(), "lark.commands", strings.ToLower(env.Header.EventType), evtID, payloadSHA, "seen")
		}
		if larkEventMaxAgeSeconds > 0 {
			ct := strings.TrimSpace(env.Header.CreateTime)
			if ct != "" {
				if ts, err := strconv.ParseInt(ct, 10, 64); err == nil {
					if ts > 1_000_000_000_000 {
						ts = ts / 1000
					}
					evtT := time.Unix(ts, 0)
					age := time.Since(evtT)
					if age > time.Duration(larkEventMaxAgeSeconds)*time.Second {
						LogStructured("info", map[string]any{"event": "lark.commands.expired", "event_id": evtID, "age_ms": age.Milliseconds()})
						if hHasDB(h) {
							_ = h.db.UpdateEventDeliveryStatus(c.Request().Context(), "lark.commands", evtID, "expired", nil)
						}
						return c.JSON(http.StatusOK, map[string]any{"result": "ok", "status": "expired"})
					}
				}
			}
		}
	}
	return c.NoContent(http.StatusOK)
}

// handleLarkCardAction processes a minimal value map with our custom fields.
func (h *Handler) handleLarkCardAction(c echo.Context, val map[string]any, callbackToken string) error {
	if val == nil {
		return c.NoContent(http.StatusOK)
	}
	op, _ := val["op"].(string)
	chatID, _ := val["chat_id"].(string)
	threadID, _ := val["thread_id"].(string)
	LogStructured("info", map[string]any{"event": "lark.card.action", "op": op, "chat_id": chatID, "thread_id": threadID})
	switch op {
	case "rebind_confirm":
		newIssue, _ := val["new_issue_id"].(string)
		newProj, _ := val["new_project_id"].(string)
		newSlug, _ := val["new_slug"].(string)
		if newIssue == "" || threadID == "" || chatID == "" {
			LogStructured("error", map[string]any{"event": "lark.card.action.error", "op": op, "reason": "missing_fields"})
			return c.JSON(http.StatusOK, map[string]any{
				"toast": map[string]any{"type": "error", "content": "参数缺失，操作失败"},
			})
		}
		if !hHasDB(h) {
			LogStructured("error", map[string]any{"event": "lark.card.action.error", "op": op, "reason": "db_unavailable"})
			return c.JSON(http.StatusOK, map[string]any{
				"toast": map[string]any{"type": "error", "content": "服务暂不可用，请稍后重试"},
			})
		}
		// Persist new binding (idempotent)
		if err := h.db.UpsertLarkThreadLink(c.Request().Context(), threadID, newIssue, newProj, newSlug, false); err != nil {
			LogStructured("error", map[string]any{"event": "lark.card.action.persist.error", "op": op, "error": err.Error()})
			return c.JSON(http.StatusOK, map[string]any{
				"toast": map[string]any{"type": "error", "content": "保存失败，请稍后重试"},
			})
		}
		_ = h.db.UpsertLarkChatIssueLink(c.Request().Context(), chatID, threadID, newIssue, newProj, newSlug)
		// Build new card
		newURL := h.planeIssueURL(newSlug, newProj, newIssue)
		display := newURL
		if display == "" {
			display = newIssue
		} else {
			display = "[打开 Issue](" + display + ")"
		}
		card := map[string]any{
			"schema": "2.0",
			"config": map[string]any{"update_multi": true, "summary": map[string]any{"content": "已改绑"}},
			"header": map[string]any{
				"title":    map[string]any{"tag": "plain_text", "content": "绑定已更新"},
				"template": "green",
				"icon":     map[string]any{"tag": "standard_icon", "token": "check_outlined"},
			},
			"body": map[string]any{
				"direction":        "vertical",
				"vertical_spacing": "small",
				"elements": []any{
					map[string]any{"tag": "markdown", "content": "已改绑为：\n" + display},
				},
			},
		}
		// Prefer delayed update with callback token, return toast immediately
		if callbackToken != "" && h.cfg.LarkAppID != "" && h.cfg.LarkAppSecret != "" {
			go h.retryUpdateInteractiveCard(callbackToken, card)
		}
		LogStructured("info", map[string]any{"event": "lark.card.action.persist.ok", "op": op, "chat_id": chatID, "thread_id": threadID, "new_issue_id": newIssue})
		return c.JSON(http.StatusOK, map[string]any{"toast": map[string]any{"type": "success", "content": "已改绑"}})
	case "rebind_cancel":
		currIssue, _ := val["curr_issue_id"].(string)
		currProj, _ := val["curr_project_id"].(string)
		currSlug, _ := val["curr_slug"].(string)
		currURL := h.planeIssueURL(currSlug, currProj, currIssue)
		display := currURL
		if display == "" {
			display = currIssue
		} else {
			display = "[打开 Issue](" + display + ")"
		}
		card := map[string]any{
			"schema": "2.0",
			"config": map[string]any{"update_multi": true, "summary": map[string]any{"content": "已取消改绑"}},
			"header": map[string]any{
				"title":    map[string]any{"tag": "plain_text", "content": "已取消改绑"},
				"template": "yellow",
				"icon":     map[string]any{"tag": "standard_icon", "token": "info_outlined"},
			},
			"body": map[string]any{
				"direction":        "vertical",
				"vertical_spacing": "small",
				"elements": []any{
					map[string]any{"tag": "markdown", "content": "已保留当前绑定：\n" + display},
				},
			},
		}
		if callbackToken != "" && h.cfg.LarkAppID != "" && h.cfg.LarkAppSecret != "" {
			go h.retryUpdateInteractiveCard(callbackToken, card)
		}
		LogStructured("info", map[string]any{"event": "lark.card.action.ok", "op": op, "chat_id": chatID, "thread_id": threadID})
		return c.JSON(http.StatusOK, map[string]any{"toast": map[string]any{"type": "info", "content": "已取消"}})
	default:
		LogStructured("info", map[string]any{"event": "lark.card.action.ignored", "op": op})
		return c.JSON(http.StatusOK, map[string]any{"toast": map[string]any{"type": "info", "content": "无操作"}})
	}
}

// verifyLarkSignature validates Feishu request signatures when headers are present.
// 按官方文档（docs/feishu/003-服务端API/事件与回调）：
// signature = sha256(timestamp + nonce + encrypt_key + raw_body) 的十六进制字符串（大小写不敏感）。
// 若缺少头或未配置 Encrypt Key，则返回 true（不拒绝），由 Verification Token 兜底校验。
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
	if m == "" {
		return ""
	}
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
	// fallback UUID only
	if u := strings.ToLower(uuidRe.FindString(s)); u != "" {
		return u, "", ""
	}
	// support browse short links: /{slug}/browse/{KEY-N}
	reb := regexp.MustCompile(`https?://[^\s]+/([\w-]+)/browse/([A-Za-z0-9]+-[0-9]+)`) // #nosec G101
	mb := reb.FindStringSubmatch(s)
	if len(mb) == 3 {
		// we only return slug here; the caller may resolve sequence -> (issueID, projectID)
		return "", mb[1], ""
	}
	return "", "", ""
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
	LogStructured("info", map[string]any{"event": "lark.send.text.start", "chat_id": chatID, "thread_id": threadID})
	token, _, err := cli.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{"event": "lark.token.error", "error": err.Error()})
		return err
	}
	if threadID != "" {
		if err := cli.ReplyTextInThread(ctx, token, threadID, text); err == nil {
			LogStructured("info", map[string]any{"event": "lark.send.text.ok", "way": "thread", "chat_id": chatID, "thread_id": threadID})
			return nil
		} else {
			LogStructured("warn", map[string]any{"event": "lark.send.text.fail", "way": "thread", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
		}
	}
	if chatID != "" {
		if err := cli.SendTextToChat(ctx, token, chatID, text); err != nil {
			LogStructured("error", map[string]any{"event": "lark.send.text.fail", "way": "chat", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
			return err
		}
		LogStructured("info", map[string]any{"event": "lark.send.text.ok", "way": "chat", "chat_id": chatID, "thread_id": threadID})
		return nil
	}
	return nil
}

// postPlaneComment posts a text comment into Plane for the given thread link.
// postPlaneCommentWithRetry posts a Plane comment with idempotency and bounded retries.
// delivery_id: prefer Lark message_id; fallback to sha(issue_id|comment)
func (h *Handler) postPlaneCommentWithRetry(tl *store.LarkThreadLink, comment, larkMessageID string) {
	if !hHasDB(h) || tl == nil || tl.PlaneIssueID == "" {
		return
	}
	// Resolve routing
	slug := ""
	if tl.WorkspaceSlug.Valid {
		slug = tl.WorkspaceSlug.String
	}
	projectID := ""
	if tl.PlaneProjectID.Valid {
		projectID = tl.PlaneProjectID.String
	}
	if slug == "" || projectID == "" {
		return
	}
	// Build delivery identifiers
	trimmed := strings.TrimSpace(comment)
	if trimmed == "" {
		return
	}
	var deliveryID string
	if strings.TrimSpace(larkMessageID) != "" {
		deliveryID = strings.TrimSpace(larkMessageID)
	} else {
		sum := sha256.Sum256([]byte(tl.PlaneIssueID + "|" + trimmed))
		deliveryID = hex.EncodeToString(sum[:])
	}
	payloadSHA := deliveryID // stable enough; could also hash full payload
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if dup, err := h.db.IsDuplicateDelivery(ctx, "lark.comment", deliveryID, payloadSHA); err == nil && dup {
		cancel()
		LogStructured("info", map[string]any{"event": "lark.comment.duplicate.db", "delivery_id": deliveryID})
		return
	}
	_ = h.db.UpsertEventDelivery(ctx, "lark.comment", "comment", deliveryID, payloadSHA, "queued")
	cancel()

	// Retry plan
	backoffs := []time.Duration{0, 5 * time.Second, 15 * time.Second}
	for i, d := range backoffs {
		if d > 0 {
			time.Sleep(d)
		}
		rctx, rcancel := context.WithTimeout(context.Background(), 12*time.Second)
		token, _ := h.ensurePlaneBotToken(rctx, slug)
		var err error
		if token == "" {
			err = fmt.Errorf("missing_plane_token")
		} else {
			pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
			html := trimmed
			err = pc.AddComment(rctx, token, slug, projectID, tl.PlaneIssueID, html)
		}
		rcancel()
		if err == nil {
			LogStructured("info", map[string]any{"event": "lark.comment.post.ok", "attempt": i + 1, "issue_id": tl.PlaneIssueID})
			_ = h.db.UpdateEventDeliveryStatus(context.Background(), "lark.comment", deliveryID, "succeeded", nil)
			return
		}
		LogStructured("error", map[string]any{"event": "lark.comment.post.fail", "attempt": i + 1, "issue_id": tl.PlaneIssueID, "error": truncate(err.Error(), 200)})
		if i < len(backoffs)-1 {
			_ = h.db.UpdateEventDeliveryRetry(context.Background(), "lark.comment", deliveryID, time.Now().Add(backoffs[i+1]))
		} else {
			_ = h.db.UpdateEventDeliveryStatus(context.Background(), "lark.comment", deliveryID, "failed", nil)
		}
	}
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

// planeIssueURL builds an app URL for the issue when enough context is available.
func (h *Handler) planeIssueURL(slug, projectID, issueID string) string {
	if slug == "" || projectID == "" || issueID == "" {
		return ""
	}
	// Derive app URL from base URL (api.plane.so -> app.plane.so)
	base := h.cfg.PlaneBaseURL
	if strings.Contains(base, "api.plane.so") {
		base = strings.Replace(base, "api.plane.so", "app.plane.so", 1)
	} else if base != "" {
		base = strings.Replace(base, "//api.", "//app.", 1)
	}
	if base == "" {
		base = "https://app.plane.so"
	}
	base = strings.TrimRight(base, "/")
	var b strings.Builder
	b.WriteString(base)
	b.WriteString("/")
	b.WriteString(slug)
	b.WriteString("/projects/")
	b.WriteString(projectID)
	b.WriteString("/issues/")
	b.WriteString(issueID)
	return b.String()
}

// extractBrowseSequence extracts KEY-N from a /{slug}/browse/KEY-N style link
func extractBrowseSequence(s string) string {
	re := regexp.MustCompile(`https?://[^\s]+/[\w-]+/browse/([A-Za-z0-9]+-[0-9]+)`) // #nosec G101
	m := re.FindStringSubmatch(s)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

func isLikelyPureBotMention(text string) bool {
	t := strings.TrimSpace(text)
	if t == "" {
		return false
	}
	parts := strings.Fields(t)
	if len(parts) == 0 {
		return false
	}
	if !strings.HasPrefix(parts[0], "@") {
		return false
	}
	if len(parts) == 1 {
		return true
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, parts[0]))
	if rest == "" {
		return true
	}
	l := strings.ToLower(rest)
	if l == "help" || l == "/help" || l == "帮助" {
		return true
	}
	return false
}

func buildLarkBotHelpCard(boundIssueMarkdown string, repoBranchMarkdown string) map[string]any {
	commands := "- **/bind {issue_link}** - 绑定 Plane Issue 到当前群聊\n" +
		"- **/comment {text}** - 添加评论到 issue"
	bodyElements := []any{
		map[string]any{
			"tag":        "markdown",
			"content":    "机器人支持的常用命令：",
			"margin":     "0px 0px 0px 0px",
			"element_id": "intro",
		},
		map[string]any{
			"tag":                "column_set",
			"flex_mode":          "stretch",
			"horizontal_spacing": "8px",
			"horizontal_align":   "left",
			"columns": []any{
				map[string]any{
					"tag":              "column",
					"width":            "weighted",
					"background_style": "blue-50",
					"elements": []any{
						map[string]any{
							"tag":        "markdown",
							"content":    "**<font color='blue'>📌 工作项命令</font>**",
							"element_id": "cmd_group_title",
						},
						map[string]any{
							"tag":        "markdown",
							"content":    commands,
							"element_id": "cmd_group_body",
						},
					},
					"padding":          "12px 12px 12px 12px",
					"vertical_spacing": "8px",
					"horizontal_align": "left",
					"vertical_align":   "top",
					"weight":           1,
				},
			},
			"margin": "0px 0px 0px 0px",
		},
		map[string]any{
			"tag":        "markdown",
			"content":    "**使用示例**：`/bind https://work.1024hub.org:4430/chen/browse/CHENE-1`",
			"text_align": "left",
			"text_size":  "normal_v2",
			"margin":     "0px 0px 0px 0px",
			"element_id": "example",
		},
	}
	if strings.TrimSpace(boundIssueMarkdown) != "" {
		bodyElements = append(bodyElements, map[string]any{
			"tag":                "column_set",
			"flex_mode":          "stretch",
			"horizontal_spacing": "8px",
			"horizontal_align":   "left",
			"columns": []any{
				map[string]any{
					"tag":              "column",
					"width":            "weighted",
					"background_style": "blue-50",
					"elements": []any{
						map[string]any{
							"tag":     "markdown",
							"content": "**<font color='blue'>🔗 已绑定 Issue</font>**",
						},
						map[string]any{
							"tag":        "markdown",
							"content":    boundIssueMarkdown,
							"text_align": "left",
							"text_size":  "normal_v2",
						},
					},
					"padding":          "12px 12px 12px 12px",
					"vertical_spacing": "4px",
					"horizontal_align": "left",
					"vertical_align":   "top",
					"weight":           1,
				},
			},
			"margin": "0px 0px 0px 0px",
		})
	}
	if strings.TrimSpace(repoBranchMarkdown) != "" {
		bodyElements = append(bodyElements, map[string]any{
			"tag":                "column_set",
			"flex_mode":          "stretch",
			"horizontal_spacing": "8px",
			"horizontal_align":   "left",
			"columns": []any{
				map[string]any{
					"tag":              "column",
					"width":            "weighted",
					"background_style": "blue-50",
					"elements": []any{
						map[string]any{
							"tag":     "markdown",
							"content": "**<font color='blue'>🌿 关联代码仓库与分支</font>**",
						},
						map[string]any{
							"tag":        "markdown",
							"content":    repoBranchMarkdown,
							"text_align": "left",
							"text_size":  "normal_v2",
						},
					},
					"padding":          "12px 12px 12px 12px",
					"vertical_spacing": "4px",
					"horizontal_align": "left",
					"vertical_align":   "top",
					"weight":           1,
				},
			},
			"margin": "0px 0px 0px 0px",
		})
	}
	helpURL := "https://github.com/CheneyDev/Cabb#%E9%A3%9E%E4%B9%A6%E9%9B%86%E6%88%90"
	bodyElements = append(bodyElements, map[string]any{
		"tag": "button",
		"text": map[string]any{
			"tag":     "plain_text",
			"content": "查看教程",
		},
		"type":  "primary_filled",
		"width": "fill",
		"size":  "large",
		"behaviors": []any{
			map[string]any{
				"type":        "open_url",
				"default_url": helpURL,
				"pc_url":      "",
				"ios_url":     "",
				"android_url": "",
			},
		},
		"margin":     "4px 0px 4px 0px",
		"element_id": "btn_help",
	})
	card := map[string]any{
		"schema": "2.0",
		"config": map[string]any{
			"update_multi": true,
			"style": map[string]any{
				"text_size": map[string]any{
					"normal_v2": map[string]any{
						"default": "normal",
						"pc":      "normal",
						"mobile":  "heading",
					},
				},
			},
		},
		"body": map[string]any{
			"direction": "vertical",
			"elements":  bodyElements,
		},
		"header": map[string]any{
			"title": map[string]any{
				"tag":     "plain_text",
				"content": "机器人使用指南",
			},
			"subtitle": map[string]any{
				"tag":     "plain_text",
				"content": "",
			},
			"template": "blue",
			"icon": map[string]any{
				"tag":   "standard_icon",
				"token": "robot_outlined",
			},
			"padding": "12px 8px 12px 8px",
		},
	}
	return card
}

func (h *Handler) postBotHelpCard(chatID, threadID string) error {
	if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return nil
	}
	boundIssueMarkdown := ""
	repoBranchMarkdown := ""
	if hHasDB(h) {
		ctxDB, cancelDB := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancelDB()
		slug := ""
		projectID := ""
		issueID := ""
		if threadID != "" {
			if tl, err := h.db.GetLarkThreadLink(ctxDB, threadID); err == nil && tl != nil && tl.PlaneIssueID != "" {
				issueID = tl.PlaneIssueID
				if tl.WorkspaceSlug.Valid {
					slug = tl.WorkspaceSlug.String
				}
				if tl.PlaneProjectID.Valid {
					projectID = tl.PlaneProjectID.String
				}
			}
		}
		if issueID == "" && chatID != "" {
			if cl, err := h.db.GetLarkChatIssueLink(ctxDB, chatID); err == nil && cl != nil && cl.PlaneIssueID != "" {
				issueID = cl.PlaneIssueID
				if cl.WorkspaceSlug.Valid {
					slug = cl.WorkspaceSlug.String
				}
				if cl.PlaneProjectID.Valid {
					projectID = cl.PlaneProjectID.String
				}
			}
		}
		if issueID != "" {
			url := h.planeIssueURL(slug, projectID, issueID)
			display := issueID
			if slug != "" && projectID != "" && h.cfg.PlaneBaseURL != "" {
				if token, _ := h.ensurePlaneBotToken(ctxDB, slug); token != "" {
					pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
					if name, err := pc.GetIssueName(ctxDB, token, slug, projectID, issueID); err == nil && strings.TrimSpace(name) != "" {
						if url != "" {
							display = "[" + escapeMD(strings.TrimSpace(name)) + "](" + url + ")"
						} else {
							display = strings.TrimSpace(name)
						}
					} else if url != "" {
						display = "[" + issueID + "](" + url + ")"
					}
				} else if url != "" {
					display = "[" + issueID + "](" + url + ")"
				}
			} else if url != "" {
				display = "[" + issueID + "](" + url + ")"
			}
			boundIssueMarkdown = "- " + display

			active := true
			if rows, err := h.db.ListBranchIssueLinks(ctxDB, issueID, "", "", &active, 5); err == nil {
				var repoBranches []string
				for _, r := range rows {
					if !r.CNBRepoID.Valid || !r.Branch.Valid {
						continue
					}
					repo := strings.TrimSpace(r.CNBRepoID.String)
					branch := strings.TrimSpace(r.Branch.String)
					if repo == "" || branch == "" {
						continue
					}
					linkText := repo + " / " + branch
					item := escapeMD(linkText)
					if base := strings.TrimSpace(h.cfg.CNBBaseURL); base != "" {
						if strings.Contains(base, "api.cnb.cool") {
							base = strings.Replace(base, "api.cnb.cool", "cnb.cool", 1)
						} else {
							base = strings.Replace(base, "//api.", "//", 1)
						}
						base = strings.TrimRight(base, "/")
						path := "/" + repo + "/-/tree/" + branch
						item = "[" + escapeMD(linkText) + "](" + base + path + ")"
					}
					repoBranches = append(repoBranches, item)
				}
				if len(repoBranches) > 0 {
					var b strings.Builder
					for _, line := range repoBranches {
						b.WriteString("- ")
						b.WriteString(line)
						b.WriteString("\n")
					}
					repoBranchMarkdown = strings.TrimSpace(b.String())
				}
			}
		}
	}
	card := buildLarkBotHelpCard(boundIssueMarkdown, repoBranchMarkdown)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	cli := &lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
	LogStructured("info", map[string]any{"event": "lark.send.help_card.start", "chat_id": chatID, "thread_id": threadID})
	token, _, err := cli.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{"event": "lark.token.error", "error": err.Error()})
		return err
	}
	if threadID != "" {
		if err := cli.ReplyCardInThread(ctx, token, threadID, card); err == nil {
			LogStructured("info", map[string]any{"event": "lark.send.help_card.ok", "way": "thread", "chat_id": chatID, "thread_id": threadID})
			return nil
		}
		LogStructured("warn", map[string]any{"event": "lark.send.help_card.fail", "way": "thread", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
	}
	if chatID != "" {
		if err := cli.SendCardToChat(ctx, token, chatID, card); err != nil {
			LogStructured("error", map[string]any{"event": "lark.send.help_card.fail", "way": "chat", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
			return err
		}
		LogStructured("info", map[string]any{"event": "lark.send.help_card.ok", "way": "chat", "chat_id": chatID, "thread_id": threadID})
		return nil
	}
	return nil
}

// nsToString returns the string value of a sql.NullString or empty when invalid
func nsToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// ensurePlaneBotToken returns the global Plane service token from config.
// Returns empty if PLANE_SERVICE_TOKEN is not configured.
func (h *Handler) ensurePlaneBotToken(ctx context.Context, workspaceSlug string) (string, error) {
	return strings.TrimSpace(h.cfg.PlaneServiceToken), nil
}

// postBoundAlready posts an idempotent notice that the chat is already bound to the issue.
func (h *Handler) postBoundAlready(chatID, threadID, slug, projectID, issueID string) error {
	url := h.planeIssueURL(slug, projectID, issueID)
	title := ""
	if slug != "" && projectID != "" && issueID != "" && hHasDB(h) {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		if token, _ := h.ensurePlaneBotToken(ctx, slug); token != "" {
			pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
			if name, err := pc.GetIssueName(ctx, token, slug, projectID, issueID); err == nil {
				title = name
			}
		}
	}
	if title != "" && h.cfg.LarkAppID != "" && h.cfg.LarkAppSecret != "" {
		return h.sendLarkPostToThread(chatID, threadID, title, url, "已绑定（无需重复绑定）")
	}
	msg := "已绑定："
	if url != "" {
		msg += url
	} else {
		msg += issueID
	}
	msg += "（无需重复绑定）"
	return h.sendLarkTextToThread(chatID, threadID, msg)
}

// postBindAck best-effort sends a rich anchor link with issue title; fallback to plain text
func (h *Handler) postBindAck(chatID, threadID, slug, projectID, issueID string) {
	url := h.planeIssueURL(slug, projectID, issueID)
	// Try to fetch issue title when we can route
	title := ""
	if slug != "" && projectID != "" && issueID != "" && hHasDB(h) {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		if token, _ := h.ensurePlaneBotToken(ctx, slug); token != "" {
			pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
			if name, err := pc.GetIssueName(ctx, token, slug, projectID, issueID); err == nil {
				title = name
			}
		}
	}
	// If we have a title, send as post anchor; else plain text
	if title != "" && h.cfg.LarkAppID != "" && h.cfg.LarkAppSecret != "" {
		if err := h.sendLarkPostToThread(chatID, threadID, title, url, "绑定成功"); err == nil {
			return
		}
		// fall through to text fallback on error
	}
	// Fallback text
	msg := "绑定成功："
	if url != "" {
		msg += url
	} else {
		msg += issueID
	}
	if slug == "" || projectID == "" {
		msg += "\n提示：未包含工作区/项目，暂不支持线程 /comment 同步；请使用完整链接重新绑定。"
	} else {
		msg += "\n提示：在该线程使用 /comment 添加评论可同步至 Plane。"
	}
	_ = h.sendLarkTextToThread(chatID, threadID, msg)
}

// sendLarkPostToThread composes a minimal zh_cn post: "ISSUE {anchor(title)} {suffix}", and replies in thread
func (h *Handler) sendLarkPostToThread(chatID, threadID, anchorText, href, suffix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	cli := &lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
	LogStructured("info", map[string]any{"event": "lark.send.post.start", "chat_id": chatID, "thread_id": threadID})
	token, _, err := cli.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{"event": "lark.token.error", "error": err.Error()})
		return err
	}
	line := []map[string]any{
		{"tag": "text", "text": "ISSUE "},
		{"tag": "a", "text": anchorText, "href": href},
		{"tag": "text", "text": " " + suffix},
	}
	post := map[string]any{
		"zh_cn": map[string]any{
			"title":   "",
			"content": []any{line},
		},
	}
	if threadID != "" {
		if err := cli.ReplyPostInThread(ctx, token, threadID, post); err == nil {
			LogStructured("info", map[string]any{"event": "lark.send.post.ok", "way": "thread", "chat_id": chatID, "thread_id": threadID})
			return nil
		} else {
			LogStructured("warn", map[string]any{"event": "lark.send.post.fail", "way": "thread", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
		}
	}
	if chatID != "" {
		if err := cli.SendPostToChat(ctx, token, chatID, post); err != nil {
			LogStructured("error", map[string]any{"event": "lark.send.post.fail", "way": "chat", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
			return err
		}
		LogStructured("info", map[string]any{"event": "lark.send.post.ok", "way": "chat", "chat_id": chatID, "thread_id": threadID})
		return nil
	}
	return nil
}

// postRebindConfirmCard sends an interactive card to confirm switching binding to a new issue.
func (h *Handler) postRebindConfirmCard(chatID, threadID, currSlug, currProjectID, currIssueID, newSlug, newProjectID, newIssueID string) error {
	if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return nil
	}
	currURL := h.planeIssueURL(currSlug, currProjectID, currIssueID)
	newURL := h.planeIssueURL(newSlug, newProjectID, newIssueID)
	currTitle := ""
	newTitle := ""
	if currSlug != "" && currProjectID != "" && currIssueID != "" && hHasDB(h) {
		rctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		if token, _ := h.ensurePlaneBotToken(rctx, currSlug); token != "" {
			pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
			if name, err := pc.GetIssueName(rctx, token, currSlug, currProjectID, currIssueID); err == nil {
				currTitle = name
			}
		}
	}
	if newSlug != "" && newProjectID != "" && newIssueID != "" && hHasDB(h) {
		rctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		if token, _ := h.ensurePlaneBotToken(rctx, newSlug); token != "" {
			pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
			if name, err := pc.GetIssueName(rctx, token, newSlug, newProjectID, newIssueID); err == nil {
				newTitle = name
			}
		}
	}
	currDisplay := currIssueID
	if currTitle != "" && currURL != "" {
		currDisplay = "[" + escapeMD(currTitle) + "](" + currURL + ")"
	} else if currURL != "" {
		currDisplay = currURL
	}
	newDisplay := newIssueID
	if newTitle != "" && newURL != "" {
		newDisplay = "[" + escapeMD(newTitle) + "](" + newURL + ")"
	} else if newURL != "" {
		newDisplay = newURL
	}
	summary := "确认Issue绑定请求"
	if currTitle != "" && newTitle != "" {
		summary = "确认：" + currTitle + " → " + newTitle
	}
	currDetail := "**" + currDisplay + "**"
	newDetail := "**" + newDisplay + "**"
	card := map[string]any{
		"schema": "2.0",
		"config": map[string]any{
			"update_multi": true,
			"style": map[string]any{
				"text_size": map[string]any{
					"normal_v2": map[string]any{"default": "normal", "pc": "normal", "mobile": "heading"},
				},
			},
			"summary": map[string]any{"content": summary},
		},
		"body": map[string]any{
			"direction": "vertical",
			"elements": []any{
				map[string]any{"tag": "markdown", "content": "您发起了新的 Issue 绑定请求，是否确认更换绑定关系？", "text_align": "left", "text_size": "normal_v2", "margin": "0px 0px 0px 0px", "element_id": "intro"},
				map[string]any{
					"tag":                "column_set",
					"flex_mode":          "stretch",
					"horizontal_spacing": "12px",
					"horizontal_align":   "left",
					"columns": []any{
						map[string]any{
							"tag":              "column",
							"width":            "weighted",
							"background_style": "blue-50",
							"padding":          "12px 12px 12px 12px",
							"vertical_spacing": "4px",
							"horizontal_align": "left",
							"vertical_align":   "top",
							"weight":           1,
							"elements": []any{
								map[string]any{"tag": "div", "text": map[string]any{"tag": "plain_text", "content": "当前绑定", "text_align": "left", "text_size": "normal_v2", "text_color": "blue"}, "icon": map[string]any{"tag": "standard_icon", "token": "info_outlined", "color": "grey"}},
								map[string]any{"tag": "markdown", "content": currDetail, "text_align": "left", "text_size": "normal_v2", "element_id": "curr_detail"},
							},
						},
						map[string]any{
							"tag":              "column",
							"width":            "weighted",
							"background_style": "violet-50",
							"padding":          "12px 12px 12px 12px",
							"vertical_spacing": "4px",
							"horizontal_align": "left",
							"vertical_align":   "top",
							"weight":           1,
							"elements": []any{
								map[string]any{"tag": "div", "text": map[string]any{"tag": "plain_text", "content": "新的请求", "text_align": "left", "text_size": "normal_v2", "text_color": "violet"}, "icon": map[string]any{"tag": "standard_icon", "token": "more-add_outlined", "color": "grey"}},
								map[string]any{"tag": "markdown", "content": newDetail, "text_align": "left", "text_size": "normal_v2", "element_id": "new_detail"},
							},
						},
					},
					"margin": "0px 0px 0px 0px",
				},
				map[string]any{"tag": "hr", "margin": "0px 0px 0px 0px", "element_id": "divider"},
				map[string]any{
					"tag":                "column_set",
					"flex_mode":          "stretch",
					"horizontal_spacing": "8px",
					"horizontal_align":   "left",
					"columns": []any{
						map[string]any{
							"tag":              "column",
							"width":            "auto",
							"vertical_spacing": "8px",
							"horizontal_align": "left",
							"vertical_align":   "top",
							"elements": []any{
								map[string]any{
									"tag": "button",
									"text": map[string]any{
										"tag":     "plain_text",
										"content": "确认换绑",
									},
									"type":  "primary_filled",
									"width": "default",
									"size":  "large",
									"behaviors": []any{
										map[string]any{
											"type": "callback",
											"value": map[string]any{
												"op":              "rebind_confirm",
												"chat_id":         chatID,
												"thread_id":       threadID,
												"curr_issue_id":   currIssueID,
												"curr_project_id": currProjectID,
												"curr_slug":       currSlug,
												"new_issue_id":    newIssueID,
												"new_project_id":  newProjectID,
												"new_slug":        newSlug,
											},
										},
									},
									"margin":     "4px 0px 4px 0px",
									"element_id": "btn_confirm",
								},
							},
						},
						map[string]any{
							"tag":              "column",
							"width":            "auto",
							"vertical_spacing": "8px",
							"horizontal_align": "left",
							"vertical_align":   "top",
							"elements": []any{
								map[string]any{
									"tag": "button",
									"text": map[string]any{
										"tag":     "plain_text",
										"content": "保持当前绑定",
									},
									"type":  "default",
									"width": "default",
									"size":  "large",
									"behaviors": []any{
										map[string]any{
											"type": "callback",
											"value": map[string]any{
												"op":        "rebind_cancel",
												"chat_id":   chatID,
												"thread_id": threadID,
											},
										},
									},
									"margin":     "4px 0px 4px 0px",
									"element_id": "btn_cancel",
								},
							},
						},
					},
					"margin": "0px 0px 0px 0px",
				},
			},
		},
		"header": map[string]any{
			"title": map[string]any{"tag": "plain_text", "content": "确认Issue绑定请求"},
			"subtitle": map[string]any{
				"tag":     "plain_text",
				"content": "",
			},
			"template": "blue",
			"icon":     map[string]any{"tag": "standard_icon", "token": "link-copy_outlined"},
			"padding":  "12px 12px 12px 12px",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	cli := &lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
	LogStructured("info", map[string]any{"event": "lark.send.card.start", "chat_id": chatID, "thread_id": threadID})
	token, _, err := cli.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{"event": "lark.token.error", "error": err.Error()})
		return err
	}
	if err := cli.ReplyCardInThread(ctx, token, threadID, card); err != nil {
		LogStructured("error", map[string]any{"event": "lark.send.card.fail", "way": "thread", "chat_id": chatID, "thread_id": threadID, "error": err.Error()})
		if chatID != "" {
			LogStructured("info", map[string]any{"event": "lark.send.card.fallback_to_chat.start", "chat_id": chatID, "thread_id": threadID})
			if e2 := cli.SendCardToChat(ctx, token, chatID, card); e2 != nil {
				LogStructured("error", map[string]any{"event": "lark.send.card.fallback_to_chat.fail", "chat_id": chatID, "thread_id": threadID, "error": e2.Error()})
				return err
			}
			LogStructured("info", map[string]any{"event": "lark.send.card.fallback_to_chat.ok", "chat_id": chatID, "thread_id": threadID})
			return nil
		}
		return err
	}
	LogStructured("info", map[string]any{"event": "lark.send.card.ok", "way": "thread", "chat_id": chatID, "thread_id": threadID})
	return nil
}

// escapeMD escapes square brackets and parentheses in markdown link text
func escapeMD(s string) string {
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}

// postChatShareLinkToPlane fetches the chat share link and posts it as a comment to the Plane issue.
func (h *Handler) postChatShareLinkToPlane(chatID, slug, projectID, issueID string) {
	if !hHasDB(h) || h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return
	}
	if chatID == "" || slug == "" || projectID == "" || issueID == "" {
		LogStructured("warn", map[string]any{
			"event":      "post_share_link.skip",
			"reason":     "missing_params",
			"chat_id":    chatID,
			"slug":       slug,
			"project_id": projectID,
			"issue_id":   issueID,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get Plane token
	token, err := h.ensurePlaneBotToken(ctx, slug)
	if err != nil || token == "" {
		LogStructured("error", map[string]any{
			"event":   "post_share_link.token_error",
			"chat_id": chatID,
			"slug":    slug,
			"error":   err,
		})
		return
	}

	// Get Lark tenant token
	larkCli := &lark.Client{AppID: h.cfg.LarkAppID, AppSecret: h.cfg.LarkAppSecret}
	larkToken, _, err := larkCli.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":   "post_share_link.lark_token_error",
			"chat_id": chatID,
			"error":   err.Error(),
		})
		return
	}

	// Get chat info for display name
	chatInfo, err := larkCli.GetChat(ctx, larkToken, chatID)
	chatName := "飞书群聊"
	if err == nil && chatInfo != nil && chatInfo.Name != "" {
		chatName = chatInfo.Name
	}

	// Get chat share link
	shareLink, err := larkCli.GetChatShareLink(ctx, larkToken, chatID, "year")
	if err != nil {
		LogStructured("error", map[string]any{
			"event":   "post_share_link.get_link_error",
			"chat_id": chatID,
			"error":   err.Error(),
		})
		return
	}

	if shareLink == nil || shareLink.ShareLink == "" {
		LogStructured("warn", map[string]any{
			"event":   "post_share_link.empty_link",
			"chat_id": chatID,
		})
		return
	}

	// Build HTML comment with share link
	commentHTML := fmt.Sprintf(
		`<p>🔗 <strong>飞书群聊已绑定</strong></p><p>群聊名称：%s</p><p><a href="%s" target="_blank">点击加入群聊</a></p>`,
		chatName,
		shareLink.ShareLink,
	)

	// Post comment to Plane
	pc := &plane.Client{BaseURL: h.cfg.PlaneBaseURL}
	if err := pc.AddComment(ctx, token, slug, projectID, issueID, commentHTML); err != nil {
		LogStructured("error", map[string]any{
			"event":    "post_share_link.add_comment_error",
			"chat_id":  chatID,
			"issue_id": issueID,
			"error":    err.Error(),
		})
		return
	}

	LogStructured("info", map[string]any{
		"event":      "post_share_link.ok",
		"chat_id":    chatID,
		"chat_name":  chatName,
		"issue_id":   issueID,
		"share_link": shareLink.ShareLink,
	})
}
