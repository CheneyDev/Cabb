package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	larkapi "cabb/internal/lark"
	planeapi "cabb/internal/plane"
	"cabb/internal/store"

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

	// 4. 获取 Plane Service Token
	token := strings.TrimSpace(h.cfg.PlaneServiceToken)
	if token == "" {
		LogStructured("error", map[string]any{
			"event": "label.sync.simple",
			"error": "plane_service_token_not_configured",
		})
		return
	}

	// 4.1 获取 workspace_slug
	workspaceSlug := strings.TrimSpace(mapping.WorkspaceSlug.String)
	if !mapping.WorkspaceSlug.Valid || workspaceSlug == "" {
		LogStructured("error", map[string]any{
			"event": "label.sync.simple",
			"error": "workspace_slug_not_configured",
		})
		return
	}

	// 5. 映射标签（自动创建）
	planeLabelIDs := []string{}
	for _, cnbLabel := range cnbLabels {
		ids, _ := h.db.MapCNBLabelsToPlane(ctx, mapping.PlaneProjectID, p.RepoSlug, []string{cnbLabel})
		if len(ids) > 0 {
			planeLabelIDs = append(planeLabelIDs, ids[0])
			continue
		}
		labelID, err := h.findOrCreatePlaneLabel(ctx, token, workspaceSlug, mapping.PlaneProjectID, p.RepoSlug, cnbLabel)
		if err != nil {
			continue
		}
		planeLabelIDs = append(planeLabelIDs, labelID)
	}

	if len(planeLabelIDs) == 0 {
		LogStructured("warn", map[string]any{
			"event":  "label.sync.simple",
			"reason": "no_valid_label_mappings",
		})
		return
	}

	// 6. 更新 Plane（增量更新）
	planeClient := &planeapi.Client{BaseURL: h.cfg.PlaneBaseURL}

	// 6.1 获取当前标签
	currentLabelIDs, err := planeClient.GetIssueLabels(ctx, token, workspaceSlug, mapping.PlaneProjectID, planeIssueID)
	if err != nil {
		LogStructured("error", map[string]any{
			"event": "label.sync.simple",
			"error": "get_current_labels_failed",
		})
		return
	}

	// 6.2 获取 CNB 管理的标签 ID
	cnbManagedIDs, err := h.db.GetCNBManagedLabelIDs(ctx, mapping.PlaneProjectID, p.RepoSlug)
	if err != nil {
		LogStructured("error", map[string]any{
			"event": "label.sync.simple",
			"error": "get_cnb_managed_labels_failed",
		})
		return
	}

	// 6.3 保留非 CNB 管理的标签
	preservedLabelIDs := make([]string, 0)
	for _, labelID := range currentLabelIDs {
		if !cnbManagedIDs[labelID] {
			preservedLabelIDs = append(preservedLabelIDs, labelID)
		}
	}

	// 6.4 合并并去重
	finalLabelIDs := uniqueStrings(append(preservedLabelIDs, planeLabelIDs...))

	// 6.5 更新
	patch := map[string]any{"labels": finalLabelIDs}
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
		"cnb_labels":   len(planeLabelIDs),
		"preserved":    len(preservedLabelIDs),
		"total":        len(finalLabelIDs),
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
