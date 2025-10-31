package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	larkapi "plane-integration/internal/lark"
	planeapi "plane-integration/internal/plane"
	"plane-integration/internal/store"

	"github.com/labstack/echo/v4"
)

type issueLabelNotifyPayload struct {
	RepoSlug     string       `json:"repo_slug"`
	IssueNumber  int          `json:"issue_number"`
	IssueURL     string       `json:"issue_url"`
	Title        string       `json:"title"`
	State        string       `json:"state"`
	Author       authorInfo   `json:"author"`
	Description  string       `json:"description"`
	Labels       []string     `json:"labels"`
	LabelTrigger string       `json:"label_trigger"`
	UpdatedAt    string       `json:"updated_at"`
	EventContext eventContext `json:"event_context"`
}

type authorInfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type eventContext struct {
	EventType string `json:"event_type"`
	Branch    string `json:"branch"`
}

// IssueLabelNotify handles POST /api/v1/issues/label-notify
// Receives issue label change notifications from CNB job-get-issues-info
func (h *Handler) IssueLabelNotify(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "invalid_token", "鉴权失败（Bearer token 不匹配）", nil)
	}

	body, sum, deliveryID, err := readAndDigestV2(c)
	if err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_body", "读取请求体失败", map[string]any{"error": err.Error()})
	}

	var p issueLabelNotifyPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return writeError(c, http.StatusUnprocessableEntity, "invalid_json", "JSON 解析失败", map[string]any{"error": err.Error()})
	}

	if err := validateIssueLabelPayload(p); err != nil {
		return writeError(c, http.StatusBadRequest, "missing_fields", err.Error(), nil)
	}

	// In-memory idempotency
	if h.dedupe != nil && h.dedupe.CheckAndMark("issue.label.notify", deliveryID, sum) {
		return c.JSON(http.StatusOK, map[string]any{
			"code":    0,
			"message": "success",
			"data": map[string]any{
				"issue_number":   p.IssueNumber,
				"processed_at":   time.Now().UTC().Format(time.RFC3339),
				"status":         "duplicate",
				"delivery_id":    deliveryID,
				"payload_sha256": sum,
			},
		})
	}

	// DB-level idempotency
	if hHasDB(h) && deliveryID != "" {
		dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "issue.label.notify", deliveryID, sum)
		if err == nil && dup {
			return c.JSON(http.StatusOK, map[string]any{
				"code":    0,
				"message": "success",
				"data": map[string]any{
					"issue_number":   p.IssueNumber,
					"processed_at":   time.Now().UTC().Format(time.RFC3339),
					"status":         "duplicate_db",
					"delivery_id":    deliveryID,
					"payload_sha256": sum,
				},
			})
		}
	}

	// Record delivery
	if hHasDB(h) && deliveryID != "" {
		_ = h.db.UpsertEventDelivery(c.Request().Context(), "issue.label.notify", "label_notify", deliveryID, sum, "queued")
	}

	// Process asynchronously
	go h.processIssueLabelNotify(p, deliveryID, sum)

	return c.JSON(http.StatusOK, map[string]any{
		"code":    0,
		"message": "success",
		"data": map[string]any{
			"issue_number": p.IssueNumber,
			"processed_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func validateIssueLabelPayload(p issueLabelNotifyPayload) error {
	if strings.TrimSpace(p.RepoSlug) == "" {
		return errMissingField("repo_slug")
	}
	if p.IssueNumber <= 0 {
		return errMissingField("issue_number")
	}
	if strings.TrimSpace(p.IssueURL) == "" {
		return errMissingField("issue_url")
	}
	if strings.TrimSpace(p.Title) == "" {
		return errMissingField("title")
	}
	if strings.TrimSpace(p.State) == "" {
		return errMissingField("state")
	}
	if strings.TrimSpace(p.Author.Username) == "" {
		return errMissingField("author.username")
	}
	if strings.TrimSpace(p.Author.Nickname) == "" {
		return errMissingField("author.nickname")
	}
	if len(p.Labels) == 0 {
		return errMissingField("labels")
	}
	if strings.TrimSpace(p.LabelTrigger) == "" {
		return errMissingField("label_trigger")
	}
	if strings.TrimSpace(p.UpdatedAt) == "" {
		return errMissingField("updated_at")
	}
	return nil
}

func errMissingField(field string) error {
	return echo.NewHTTPError(http.StatusBadRequest, "缺少必填字段："+field)
}

func readAndDigestV2(c echo.Context) (body []byte, sum string, deliveryID string, err error) {
	body, err = io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, "", "", err
	}
	c.Request().Body = io.NopCloser(strings.NewReader(string(body)))
	h := sha256.Sum256(body)
	sum = hex.EncodeToString(h[:])
	deliveryID = c.Request().Header.Get("X-Delivery-ID")
	if deliveryID == "" {
		deliveryID = c.Request().Header.Get("X-Request-ID")
	}
	return body, sum, deliveryID, nil
}

// processIssueLabelNotify handles the business logic asynchronously
func (h *Handler) processIssueLabelNotify(p issueLabelNotifyPayload, deliveryID, sum string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if !hHasDB(h) {
		LogStructured("warn", map[string]any{
			"event":       "issue.label.notify.process",
			"delivery_id": deliveryID,
			"decision":    "skip",
			"reason":      "no_database",
		})
		return
	}

	// 1. 查询 repo-project 映射
	mapping, err := h.db.GetRepoProjectMapping(ctx, p.RepoSlug)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":       "issue.label.notify.process",
			"delivery_id": deliveryID,
			"repo_slug":   p.RepoSlug,
			"error":       "mapping_not_found",
		})
		return
	}

	// 2. 查找对应的 Plane Issue
	planeIssueID, err := h.db.FindPlaneIssueByCNBIssue(ctx, p.RepoSlug, fmt.Sprintf("%d", p.IssueNumber))
	if err != nil || planeIssueID == "" {
		LogStructured("warn", map[string]any{
			"event":        "issue.label.notify.process",
			"delivery_id":  deliveryID,
			"repo_slug":    p.RepoSlug,
			"issue_number": p.IssueNumber,
			"decision":     "skip",
			"reason":       "plane_issue_not_found",
		})
		return
	}

	// 3. 过滤 _CNB 后缀的标签
	cnbLabels := filterCNBLabels(p.Labels)
	if len(cnbLabels) == 0 {
		LogStructured("info", map[string]any{
			"event":        "issue.label.notify.process",
			"delivery_id":  deliveryID,
			"repo_slug":    p.RepoSlug,
			"issue_number": p.IssueNumber,
			"decision":     "skip",
			"reason":       "no_cnb_labels",
		})
		return
	}

	// 4. 映射 CNB 标签到 Plane 标签 ID
	planeLabelIDs, err := h.db.MapCNBLabelsToPlane(ctx, mapping.PlaneProjectID, p.RepoSlug, cnbLabels)
	if err != nil || len(planeLabelIDs) == 0 {
		LogStructured("warn", map[string]any{
			"event":       "issue.label.notify.process",
			"delivery_id": deliveryID,
			"cnb_labels":  cnbLabels,
			"decision":    "skip",
			"reason":      "label_mapping_failed",
		})
		return
	}

	// 5. 获取 bot token
	token, workspaceSlug, err := h.db.FindBotTokenByWorkspaceID(ctx, mapping.PlaneWorkspaceID)
	if err != nil || token == "" {
		LogStructured("error", map[string]any{
			"event":       "issue.label.notify.process",
			"delivery_id": deliveryID,
			"error":       "bot_token_not_found",
		})
		return
	}

	// 6. 更新 Plane Issue 标签（增量更新，只替换 CNB 管理的标签）
	planeClient := &planeapi.Client{BaseURL: h.cfg.PlaneBaseURL}

	// 6.1 获取当前 Issue 的所有标签
	currentLabelIDs, err := planeClient.GetIssueLabels(ctx, token, workspaceSlug, mapping.PlaneProjectID, planeIssueID)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":          "issue.label.notify.process",
			"delivery_id":    deliveryID,
			"plane_issue_id": planeIssueID,
			"error":          "get_current_labels_failed",
			"details":        err.Error(),
		})
		return
	}

	// 6.2 获取 CNB 管理的标签 ID 列表（用于识别哪些标签可以被替换）
	cnbManagedIDs, err := h.db.GetCNBManagedLabelIDs(ctx, mapping.PlaneProjectID, p.RepoSlug)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":       "issue.label.notify.process",
			"delivery_id": deliveryID,
			"error":       "get_cnb_managed_labels_failed",
			"details":     err.Error(),
		})
		return
	}

	// 6.3 过滤出非 CNB 管理的标签（需要保留）
	preservedLabelIDs := make([]string, 0)
	for _, labelID := range currentLabelIDs {
		if !cnbManagedIDs[labelID] {
			// 不是 CNB 管理的标签，需要保留
			preservedLabelIDs = append(preservedLabelIDs, labelID)
		}
	}

	// 6.4 合并：保留的标签 + 新的 CNB 标签
	finalLabelIDs := append(preservedLabelIDs, planeLabelIDs...)

	// 6.5 去重
	uniqueLabelIDs := uniqueStrings(finalLabelIDs)

	// 6.6 更新到 Plane
	patch := map[string]any{"labels": uniqueLabelIDs}
	err = planeClient.PatchIssue(ctx, token, workspaceSlug, mapping.PlaneProjectID, planeIssueID, patch)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":          "issue.label.notify.process",
			"delivery_id":    deliveryID,
			"plane_issue_id": planeIssueID,
			"error":          "plane_patch_failed",
			"details":        err.Error(),
		})
		return
	}

	LogStructured("info", map[string]any{
		"event":            "issue.label.notify.process",
		"delivery_id":      deliveryID,
		"repo_slug":        p.RepoSlug,
		"issue_number":     p.IssueNumber,
		"plane_issue_id":   planeIssueID,
		"cnb_labels_count": len(planeLabelIDs),
		"preserved_count":  len(preservedLabelIDs),
		"total_count":      len(uniqueLabelIDs),
		"result":           "success",
	})

	// 7. 发送飞书通知（如果配置了 channel-project 映射）
	h.sendLarkNotificationForLabelChange(ctx, mapping, p, planeIssueID, cnbLabels)
}

// filterCNBLabels 提取以 _CNB 结尾的标签
func filterCNBLabels(labels []string) []string {
	var cnbLabels []string
	for _, label := range labels {
		if strings.HasSuffix(label, "_CNB") {
			cnbLabels = append(cnbLabels, label)
		}
	}
	return cnbLabels
}

// sendLarkNotificationForLabelChange 发送飞书标签变更通知
func (h *Handler) sendLarkNotificationForLabelChange(ctx context.Context, mapping *store.RepoProjectMapping, p issueLabelNotifyPayload, planeIssueID string, cnbLabels []string) {
	if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return
	}

	// 查询 channel-project 映射
	links, err := h.db.GetChannelsByPlaneProject(ctx, mapping.PlaneProjectID)
	if err != nil || len(links) == 0 {
		return
	}

	larkClient := &larkapi.Client{
		AppID:     h.cfg.LarkAppID,
		AppSecret: h.cfg.LarkAppSecret,
		BaseURL:   "https://open.feishu.cn",
	}

	token, _, err := larkClient.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{
			"event": "lark.notify.label_change",
			"error": "get_tenant_token_failed",
		})
		return
	}

	// 构建通知消息
	message := buildLabelChangeMessage(p, cnbLabels)

	// 向所有绑定的飞书群组发送通知
	for _, link := range links {
		if err := larkClient.SendTextToChat(ctx, token, link.LarkChatID, message); err != nil {
			LogStructured("error", map[string]any{
				"event":   "lark.notify.label_change",
				"chat_id": link.LarkChatID,
				"error":   err.Error(),
			})
		} else {
			LogStructured("info", map[string]any{
				"event":        "lark.notify.label_change",
				"chat_id":      link.LarkChatID,
				"issue_number": p.IssueNumber,
				"result":       "success",
			})
		}
	}
}

// buildLabelChangeMessage 构建标签变更通知消息
func buildLabelChangeMessage(p issueLabelNotifyPayload, cnbLabels []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 Issue 标签更新\n\n"))
	sb.WriteString(fmt.Sprintf("仓库：%s\n", p.RepoSlug))
	sb.WriteString(fmt.Sprintf("Issue：#%d - %s\n", p.IssueNumber, p.Title))
	sb.WriteString(fmt.Sprintf("状态：%s\n", p.State))
	sb.WriteString(fmt.Sprintf("标签：%s\n", strings.Join(cnbLabels, ", ")))
	sb.WriteString(fmt.Sprintf("触发标签：%s\n", p.LabelTrigger))
	if p.IssueURL != "" {
		sb.WriteString(fmt.Sprintf("\n🔗 查看详情：%s", p.IssueURL))
	}
	return sb.String()
}

// uniqueStrings removes duplicates from a string slice
func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != "" && !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
