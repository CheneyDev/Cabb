package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	larkapi "plane-integration/internal/lark"
	planeapi "plane-integration/internal/plane"
	"plane-integration/internal/store"

	"github.com/labstack/echo/v4"
)

// 简化版请求体 - 只需要核心字段
type issueLabelNotifySimplePayload struct {
	RepoSlug    string   `json:"repo_slug"`
	IssueNumber int      `json:"issue_number"`
	Labels      []string `json:"labels"`
}

// IssueLabelNotifySimple 处理简化版 API（只需要 repo_slug + issue_number + labels）
// POST /api/v1/issues/label-sync
func (h *Handler) IssueLabelNotifySimple(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "invalid_token", "鉴权失败（Bearer token 不匹配）", nil)
	}

	body, sum, deliveryID, err := readAndDigestV2(c)
	if err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_body", "读取请求体失败", map[string]any{"error": err.Error()})
	}

	var p issueLabelNotifySimplePayload
	if err := json.Unmarshal(body, &p); err != nil {
		return writeError(c, http.StatusUnprocessableEntity, "invalid_json", "JSON 解析失败", map[string]any{"error": err.Error()})
	}

	// 简单校验
	if strings.TrimSpace(p.RepoSlug) == "" || p.IssueNumber <= 0 || len(p.Labels) == 0 {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少必填字段：repo_slug/issue_number/labels", nil)
	}

	// 内存去重
	if h.dedupe != nil && h.dedupe.CheckAndMark("issue.label.sync", deliveryID, sum) {
		return c.JSON(http.StatusOK, map[string]any{
			"code":    0,
			"message": "success",
			"data": map[string]any{
				"issue_number": p.IssueNumber,
				"processed_at": time.Now().UTC().Format(time.RFC3339),
				"status":       "duplicate",
			},
		})
	}

	// 数据库去重
	if hHasDB(h) && deliveryID != "" {
		dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "issue.label.sync", deliveryID, sum)
		if err == nil && dup {
			return c.JSON(http.StatusOK, map[string]any{
				"code":    0,
				"message": "success",
				"data": map[string]any{
					"issue_number": p.IssueNumber,
					"processed_at": time.Now().UTC().Format(time.RFC3339),
					"status":       "duplicate_db",
				},
			})
		}
	}

	// 记录事件
	if hHasDB(h) && deliveryID != "" {
		_ = h.db.UpsertEventDelivery(c.Request().Context(), "issue.label.sync", "label_sync", deliveryID, sum, "queued")
	}

	// 异步处理（复用完整版逻辑）
	go h.processLabelSyncSimple(p, deliveryID, sum)

	return c.JSON(http.StatusOK, map[string]any{
		"code":    0,
		"message": "success",
		"data": map[string]any{
			"issue_number": p.IssueNumber,
			"processed_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// processLabelSyncSimple 处理简化版标签同步
func (h *Handler) processLabelSyncSimple(p issueLabelNotifySimplePayload, deliveryID, sum string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if !hHasDB(h) {
		return
	}

	// 1. 查询映射
	mapping, err := h.db.GetRepoProjectMapping(ctx, p.RepoSlug)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":     "label.sync.simple",
			"repo_slug": p.RepoSlug,
			"error":     "mapping_not_found",
		})
		return
	}

	// 2. 查找 Plane Issue
	planeIssueID, err := h.db.FindPlaneIssueByCNBIssue(ctx, p.RepoSlug, fmt.Sprintf("%d", p.IssueNumber))
	if err != nil || planeIssueID == "" {
		LogStructured("warn", map[string]any{
			"event":        "label.sync.simple",
			"repo_slug":    p.RepoSlug,
			"issue_number": p.IssueNumber,
			"reason":       "plane_issue_not_found",
		})
		return
	}

	// 3. 过滤 CNB 标签
	cnbLabels := filterCNBLabels(p.Labels)
	if len(cnbLabels) == 0 {
		LogStructured("info", map[string]any{
			"event":  "label.sync.simple",
			"reason": "no_cnb_labels",
		})
		return
	}

	// 4. 映射标签
	planeLabelIDs, err := h.db.MapCNBLabelsToPlane(ctx, mapping.PlaneProjectID, p.RepoSlug, cnbLabels)
	if err != nil || len(planeLabelIDs) == 0 {
		LogStructured("warn", map[string]any{
			"event":  "label.sync.simple",
			"reason": "label_mapping_failed",
		})
		return
	}

	// 5. 获取 token
	token, workspaceSlug, err := h.db.FindBotTokenByWorkspaceID(ctx, mapping.PlaneWorkspaceID)
	if err != nil || token == "" {
		LogStructured("error", map[string]any{
			"event": "label.sync.simple",
			"error": "bot_token_not_found",
		})
		return
	}

	// 6. 更新 Plane
	planeClient := &planeapi.Client{BaseURL: h.cfg.PlaneBaseURL}
	patch := map[string]any{"labels": planeLabelIDs}
	err = planeClient.PatchIssue(ctx, token, workspaceSlug, mapping.PlaneProjectID, planeIssueID, patch)
	if err != nil {
		LogStructured("error", map[string]any{
			"event":   "label.sync.simple",
			"error":   "plane_patch_failed",
			"details": err.Error(),
		})
		return
	}

	LogStructured("info", map[string]any{
		"event":        "label.sync.simple",
		"repo_slug":    p.RepoSlug,
		"issue_number": p.IssueNumber,
		"labels_count": len(planeLabelIDs),
		"result":       "success",
	})

	// 7. 飞书通知
	h.sendLarkNotificationSimple(ctx, mapping, p, cnbLabels)
}

// sendLarkNotificationSimple 简化版飞书通知
func (h *Handler) sendLarkNotificationSimple(ctx context.Context, mapping *store.RepoProjectMapping, p issueLabelNotifySimplePayload, cnbLabels []string) {
	if h.cfg.LarkAppID == "" || h.cfg.LarkAppSecret == "" {
		return
	}

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
		return
	}

	message := fmt.Sprintf("📋 标签更新\n仓库：%s\nIssue：#%d\n标签：%s",
		p.RepoSlug, p.IssueNumber, strings.Join(cnbLabels, ", "))

	for _, link := range links {
		_ = larkClient.SendTextToChat(ctx, token, link.LarkChatID, message)
	}
}
