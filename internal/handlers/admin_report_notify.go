package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cabb/internal/cnb"
	"cabb/internal/lark"
	"cabb/internal/store"
	"cabb/pkg/config"

	"github.com/labstack/echo/v4"
)

// AdminReportNotifyConfigGet returns the current report notification configuration
func (h *Handler) AdminReportNotifyConfigGet(c echo.Context) error {
	cfg, err := h.db.GetReportNotifyConfig(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cfg)
}

// AdminReportNotifyConfigSave saves the report notification configuration
func (h *Handler) AdminReportNotifyConfigSave(c echo.Context) error {
	var req store.ReportNotifyConfig
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Validate notify_type
	switch req.NotifyType {
	case "chat", "users", "departments":
		// valid
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid notify_type, must be 'chat', 'users', or 'departments'"})
	}

	if err := h.db.SaveReportNotifyConfig(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// AdminReportNotifyTest sends a test notification based on current config
func (h *Handler) AdminReportNotifyTest(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		OpenID string `json:"open_id"` // optional, for user mode
		ChatID string `json:"chat_id"` // optional, for chat mode
	}
	_ = c.Bind(&req)

	// Get current config to determine test method
	notifyCfg, err := h.db.GetReportNotifyConfig(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get config: " + err.Error()})
	}

	// Build a test card
	testCard := map[string]any{
		"schema": "2.0",
		"config": map[string]any{"wide_screen_mode": true},
		"header": map[string]any{
			"title":    map[string]any{"tag": "plain_text", "content": "📋 测试通知"},
			"template": "blue",
		},
		"body": map[string]any{
			"elements": []map[string]any{
				{"tag": "markdown", "content": "这是一条测试通知，用于验证报告推送功能是否正常工作。"},
				{"tag": "markdown", "content": "**时间**: " + time.Now().Format("2006-01-02 15:04:05")},
				{"tag": "markdown", "content": "**推送方式**: " + notifyCfg.NotifyType},
			},
		},
	}

	larkClient := lark.NewClient(h.cfg.LarkAppID, h.cfg.LarkAppSecret)
	token, _, err := larkClient.TenantAccessToken(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get lark token: " + err.Error()})
	}

	switch notifyCfg.NotifyType {
	case "chat":
		chatID := req.ChatID
		if chatID == "" {
			chatID = notifyCfg.ChatID
		}
		if chatID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "请先配置群聊 Chat ID"})
		}
		if err := larkClient.SendCardToChat(ctx, token, chatID, testCard); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "发送失败: " + err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]any{"status": "ok", "method": "chat", "chat_id": chatID})

	case "users":
		openID := req.OpenID
		if openID == "" && len(notifyCfg.UserIDs) > 0 {
			openID = notifyCfg.UserIDs[0].ID
		}
		if openID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "请先选择至少一个用户"})
		}
		result, err := larkClient.BatchSendCard(ctx, token, []string{openID}, nil, testCard)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "发送失败: " + err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]any{"status": "ok", "method": "users", "message_id": result.MessageID})

	case "departments":
		if len(notifyCfg.DepartmentIDs) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "请先选择至少一个部门"})
		}
		// Only send to first department for test
		deptID := notifyCfg.DepartmentIDs[0].ID
		result, err := larkClient.BatchSendCard(ctx, token, nil, []string{deptID}, testCard)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "发送失败: " + err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]any{"status": "ok", "method": "departments", "message_id": result.MessageID, "department": notifyCfg.DepartmentIDs[0].Name})

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "未知的推送方式"})
	}
}

// AdminReportNotifySend manually triggers a report notification
func (h *Handler) AdminReportNotifySend(c echo.Context) error {
	var req struct {
		ReportType string `json:"report_type"` // daily, weekly, monthly
		Label      string `json:"label"`       // e.g., "2024-01-15" or "2024-01"
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.ReportType == "" || req.Label == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "report_type and label are required"})
	}

	// Trigger async send
	go sendReportWithConfig(h.cfg, h.db, req.ReportType, req.Label)

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "report notification triggered",
	})
}

// AdminLarkDepartments returns all departments from Lark
func (h *Handler) AdminLarkDepartments(c echo.Context) error {
	ctx := c.Request().Context()

	larkClient := lark.NewClient(h.cfg.LarkAppID, h.cfg.LarkAppSecret)
	token, _, err := larkClient.TenantAccessToken(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get lark token: " + err.Error()})
	}

	depts, err := larkClient.ListAllDepartments(ctx, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch departments: " + err.Error()})
	}

	// Convert to response format
	var result []map[string]string
	for _, d := range depts {
		id := d.OpenDepartmentID
		if id == "" {
			id = d.DepartmentID
		}
		result = append(result, map[string]string{
			"id":   id,
			"name": d.Name,
		})
	}

	return c.JSON(http.StatusOK, result)
}

// sendReportWithConfig sends report based on database configuration
func sendReportWithConfig(cfg config.Config, db *store.DB, reportType, label string) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Get notification config from database
	notifyCfg, err := db.GetReportNotifyConfig(ctx)
	if err != nil {
		LogStructured("error", map[string]any{
			"event": "report.notify.config",
			"error": err.Error(),
		})
		return
	}

	// Check if this report type is enabled
	switch reportType {
	case "daily":
		if !notifyCfg.DailyEnabled {
			return
		}
	case "weekly":
		if !notifyCfg.WeeklyEnabled {
			return
		}
	case "monthly":
		if !notifyCfg.MonthlyEnabled {
			return
		}
	}

	// Fetch report JSON from repo
	cnbClient := &cnb.Client{
		BaseURL: cfg.CNBBaseURL,
		Token:   cfg.CNBAppToken,
	}

	var jsonPath string
	switch reportType {
	case "daily":
		jsonPath = "ai-report/daily/report-" + label + ".json"
	case "weekly":
		jsonPath = "ai-report/weekly/report-" + label + ".json"
	case "monthly":
		jsonPath = "ai-report/monthly/report-" + label + ".json"
	default:
		return
	}

	content, err := cnbClient.GetFileContent(ctx, cfg.ReportRepo, cfg.ReportBranch, jsonPath)
	if err != nil {
		LogStructured("error", map[string]any{
			"event": "report.notify.fetch",
			"type":  reportType,
			"label": label,
			"path":  jsonPath,
			"error": err.Error(),
		})
		return
	}

	var report ReportJSON
	if err := json.Unmarshal(content, &report); err != nil {
		LogStructured("error", map[string]any{
			"event": "report.notify.parse",
			"type":  reportType,
			"label": label,
			"error": err.Error(),
		})
		return
	}

	// Build card
	card := buildReportCard(report)

	// Get Lark token
	larkClient := lark.NewClient(cfg.LarkAppID, cfg.LarkAppSecret)
	token, _, err := larkClient.TenantAccessToken(ctx)
	if err != nil {
		LogStructured("error", map[string]any{
			"event": "report.notify.lark_token",
			"error": err.Error(),
		})
		return
	}

	// Send based on notify type
	switch notifyCfg.NotifyType {
	case "chat":
		if notifyCfg.ChatID == "" {
			LogStructured("warn", map[string]any{
				"event":   "report.notify.skip",
				"reason":  "empty chat_id",
			})
			return
		}
		if err := larkClient.SendCardToChat(ctx, token, notifyCfg.ChatID, card); err != nil {
			LogStructured("error", map[string]any{
				"event":   "report.notify.send_chat",
				"chat_id": notifyCfg.ChatID,
				"error":   err.Error(),
			})
			return
		}
		LogStructured("info", map[string]any{
			"event":   "report.notify.sent",
			"type":    reportType,
			"method":  "chat",
			"chat_id": notifyCfg.ChatID,
		})

	case "users":
		if len(notifyCfg.UserIDs) == 0 {
			LogStructured("warn", map[string]any{
				"event":  "report.notify.skip",
				"reason": "empty user_ids",
			})
			return
		}
		// Extract open_ids
		var openIDs []string
		for _, u := range notifyCfg.UserIDs {
			if u.ID != "" {
				openIDs = append(openIDs, u.ID)
			}
		}
		// Batch send (max 200 per request)
		for i := 0; i < len(openIDs); i += 200 {
			end := i + 200
			if end > len(openIDs) {
				end = len(openIDs)
			}
			batch := openIDs[i:end]
			result, err := larkClient.BatchSendCard(ctx, token, batch, nil, card)
			if err != nil {
				LogStructured("error", map[string]any{
					"event":     "report.notify.batch_send",
					"batch":     i / 200,
					"user_count": len(batch),
					"error":     err.Error(),
				})
			} else {
				LogStructured("info", map[string]any{
					"event":      "report.notify.sent",
					"type":       reportType,
					"method":     "users",
					"message_id": result.MessageID,
					"user_count": len(batch),
				})
			}
		}

	case "departments":
		if len(notifyCfg.DepartmentIDs) == 0 {
			LogStructured("warn", map[string]any{
				"event":  "report.notify.skip",
				"reason": "empty department_ids",
			})
			return
		}
		// Extract department ids
		var deptIDs []string
		for _, d := range notifyCfg.DepartmentIDs {
			if d.ID != "" {
				deptIDs = append(deptIDs, d.ID)
			}
		}
		// Batch send (max 200 per request)
		for i := 0; i < len(deptIDs); i += 200 {
			end := i + 200
			if end > len(deptIDs) {
				end = len(deptIDs)
			}
			batch := deptIDs[i:end]
			result, err := larkClient.BatchSendCard(ctx, token, nil, batch, card)
			if err != nil {
				LogStructured("error", map[string]any{
					"event":      "report.notify.batch_send",
					"batch":      i / 200,
					"dept_count": len(batch),
					"error":      err.Error(),
				})
			} else {
				LogStructured("info", map[string]any{
					"event":      "report.notify.sent",
					"type":       reportType,
					"method":     "departments",
					"message_id": result.MessageID,
					"dept_count": len(batch),
				})
			}
		}
	}
}
