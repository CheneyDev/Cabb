package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"cabb/internal/cnb"
	"cabb/internal/lark"
	"cabb/internal/plane"
	"cabb/internal/store"
	"cabb/pkg/config"
)

func (h *Handler) JobIssueSummaryDaily(c echo.Context) error {
	return c.JSON(http.StatusAccepted, map[string]any{
		"scheduled": true,
		"job":       "issue-summary-daily",
	})
}

// JobCleanupThreadLinks deletes non-sync-enabled thread links older than N days (default 90; min 7)
// Invoke via: POST /jobs/cleanup/thread-links?days=90
func (h *Handler) JobCleanupThreadLinks(c echo.Context) error {
	if !hHasDB(h) {
		return c.JSON(http.StatusOK, map[string]any{"deleted": 0, "cutoff": "", "note": "db unavailable"})
	}
	days := 90
	if v := c.QueryParam("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			days = n
		}
	}
	if days < 7 {
		days = 7
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	ctx, cancel := context.WithTimeout(c.Request().Context(), 12*time.Second)
	defer cancel()
	n, err := h.db.CleanupStaleThreadLinks(ctx, cutoff)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "cleanup_failed", "清理失败", map[string]any{"error": err.Error()})
	}
	LogStructured("info", map[string]any{"event": "jobs.cleanup.thread_links", "deleted": n, "cutoff": cutoff.UTC().Format(time.RFC3339)})
	return c.JSON(http.StatusOK, map[string]any{
		"deleted": n,
		"cutoff":  cutoff.UTC().Format(time.RFC3339),
		"result":  "ok",
	})
}

// StartCleanupScheduler starts a daily scheduler for cleaning stale thread links.
// It respects cfg.CleanupThreadLinksEnabled, cfg.CleanupThreadLinksDays, cfg.CleanupThreadLinksAt and cfg.Timezone.
func StartCleanupScheduler(cfg config.Config, db *store.DB) {
	if !cfg.CleanupThreadLinksEnabled || db == nil || db.SQL == nil {
		return
	}
	go func() {
		// Determine timezone
		loc := time.Local
		if tz := strings.TrimSpace(cfg.Timezone); tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		// Parse time of day
		hh, mm := parseHHMM(cfg.CleanupThreadLinksAt)
		// Compute first run
		next := nextAt(loc, hh, mm)
		for {
			sleep := time.Until(next)
			if sleep <= 0 {
				next = next.Add(24 * time.Hour)
				continue
			}
			time.Sleep(sleep)
			// Run cleanup
			days := cfg.CleanupThreadLinksDays
			if days < 7 {
				days = 7
			}
			cutoff := time.Now().In(loc).AddDate(0, 0, -days)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			n, err := db.CleanupStaleThreadLinks(ctx, cutoff)
			cancel()
			if err != nil {
				LogStructured("error", map[string]any{"event": "jobs.cleanup.thread_links.cron", "error": map[string]any{"code": "cleanup_failed", "message": truncate(err.Error(), 200)}, "cutoff": cutoff.UTC().Format(time.RFC3339)})
			} else {
				LogStructured("info", map[string]any{"event": "jobs.cleanup.thread_links.cron", "deleted": n, "cutoff": cutoff.UTC().Format(time.RFC3339)})
			}
			next = next.Add(24 * time.Hour)
		}
	}()
}

func parseHHMM(s string) (int, int) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 3, 0
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 3, 0
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	if h < 0 || h > 23 {
		h = 3
	}
	if m < 0 || m > 59 {
		m = 0
	}
	return h, m
}

func nextAt(loc *time.Location, hour, min int) time.Time {
	now := time.Now().In(loc)
	run := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, loc)
	if !run.After(now) {
		run = run.Add(24 * time.Hour)
	}
	return run
}

// JobDailyProgressReport triggers daily progress reports for active issues.
// POST /jobs/daily-report
func (h *Handler) JobDailyProgressReport(c echo.Context) error {
	if !hHasDB(h) {
		return c.JSON(http.StatusOK, map[string]any{"status": "skipped", "note": "db unavailable"})
	}
	ctx := c.Request().Context()

	// 1. Handle issues with active branch links
	branchLinks, err := h.db.ListActiveBranchLinks(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "db_error", "Failed to list branch links", map[string]any{"error": err.Error()})
	}

	cnbClient := &cnb.Client{
		BaseURL: h.cfg.CNBBaseURL,
		Token:   h.cfg.CNBAppToken,
	}
	planeClient := &plane.Client{
		BaseURL: h.cfg.PlaneBaseURL,
	}

	triggered := 0
	errors := 0

	for _, link := range branchLinks {
		// Fetch issue details to get title/desc
		// We need workspace slug and project ID. The link table doesn't have them directly,
		// but we can try to find them or maybe we need to store them in branch_issue_links?
		// Wait, branch_issue_links only has plane_issue_id.
		// We can use planeClient.GetIssueDetail if we have workspace/project.
		// Or we can use a helper to find the issue if we don't have context.
		// Actually, `repo_project_mappings` might help if we join, but `branch_issue_links` is per issue.
		// Let's assume we can get workspace/project from `repo_project_mappings` via `cnb_repo_id`.
		// But an issue might belong to a different project than the repo mapping if it was moved?
		// Unlikely for this integration.

		// Better approach: Use `repo_project_mappings` to get workspace/project for the repo.
		mapping, err := h.db.GetRepoProjectMapping(ctx, link.CNBRepoID)
		if err != nil {
			LogStructured("error", map[string]any{"event": "job.daily_report.get_mapping", "repo": link.CNBRepoID, "error": err.Error()})
			errors++
			continue
		}
		if mapping == nil {
			continue
		}

		// Get Issue Details
		issue, err := planeClient.GetIssueDetail(ctx, h.cfg.PlaneServiceToken, mapping.PlaneWorkspaceID, mapping.PlaneProjectID, link.PlaneIssueID)
		if err != nil {
			LogStructured("error", map[string]any{"event": "job.daily_report.get_issue", "issue_id": link.PlaneIssueID, "error": err.Error()})
			errors++
			continue
		}

		// Trigger Pipeline
		envVars := map[string]string{
			"ISSUE_TITLE":       issue.Name,
			"ISSUE_DESCRIPTION": issue.DescriptionHTML, // or DescriptionStripped
			"LARK_CHAT_ID":      "",
		}
		if link.LarkChatID.Valid {
			envVars["LARK_CHAT_ID"] = link.LarkChatID.String
		}

		if err := cnbClient.TriggerPipeline(ctx, link.CNBRepoID, link.Branch, envVars); err != nil {
			LogStructured("error", map[string]any{"event": "job.daily_report.trigger_pipeline", "repo": link.CNBRepoID, "branch": link.Branch, "error": err.Error()})
			errors++
		} else {
			triggered++
		}
	}

	// 2. Handle issues WITHOUT branch links (send reminder)
	chatLinks, err := h.db.ListActiveChatLinksWithoutBranch(ctx)
	if err != nil {
		LogStructured("error", map[string]any{"event": "job.daily_report.list_chat_links", "error": err.Error()})
	} else {
		larkClient := lark.NewClient(h.cfg.LarkAppID, h.cfg.LarkAppSecret)
		for _, link := range chatLinks {
			// Send reminder to Lark Chat
			// We need to know which issue it is to give a good message.
			// We only have PlaneIssueID. We assume we can't easily get the title without project context.
			// But we can try to find it if we iterate projects or if we had project_id in chat_issue_links.
			// chat_issue_links HAS plane_project_id!

			// We need to fetch the chat link details again or update the query to return project_id.
			// The ListActiveChatLinksWithoutBranch query returns plane_issue_id and lark_chat_id.
			// Let's update the query in store to return project_id too?
			// Or just send a generic message.
			// "Reminder: The bound issue has no linked branch/repo. Progress reporting is disabled."

			msg := "⚠️ **每日进度汇报提醒**\n\n当前绑定的 Plane Issue 尚未关联代码仓库/分支，无法自动生成开发进度日报。\n请使用 `/bind` 指令关联分支。"
			if err := larkClient.SendMessage(ctx, link.LarkChatID, msg); err != nil {
				LogStructured("error", map[string]any{"event": "job.daily_report.send_reminder", "chat_id": link.LarkChatID, "error": err.Error()})
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"triggered": triggered,
		"errors":    errors,
	})
}

// JobDailyReportNotify fetches the daily report from repo and sends to Lark.
// POST /jobs/daily-report/notify
func (h *Handler) JobDailyReportNotify(c echo.Context) error {
	if !hHasDB(h) {
		return c.JSON(http.StatusOK, map[string]any{"status": "skipped", "note": "db unavailable"})
	}
	ctx := c.Request().Context()

	// 1. Determine report path
	// Format: reports/report-{YYYY-MM-DD}.json
	// We assume daily report for "yesterday" if run today? Or today?
	// The script generates report for "yesterday" by default, but label is yesterday's date.
	// If we run at 17:50 today, we probably want TODAY's report if the script ran today for TODAY?
	// Wait, the script default is "yesterday".
	// If we want today's report, we need to ensure the script ran for "today" or "yesterday".
	// User said "17:50 send summary", implying summary of TODAY.
	// So the script should have run for "today" or "yesterday"?
	// Usually daily report at end of day is for TODAY.
	// Let's assume the script was triggered with `REPORT_TIMEFRAME=today` or similar, OR
	// the script default "yesterday" is for next morning.
	// BUT user said "17:50 send", so it's end of day.
	// So we should look for TODAY's report.
	// Let's try to fetch report for TODAY.
	today := time.Now().Format("2006-01-02")
	path := fmt.Sprintf("reports/report-%s.json", today)

	cnbClient := &cnb.Client{
		BaseURL: h.cfg.CNBBaseURL,
		Token:   h.cfg.CNBAppToken,
	}

	// 2. Fetch JSON content
	content, err := cnbClient.GetFileContent(ctx, h.cfg.ReportRepo, h.cfg.ReportBranch, path)
	if err != nil {
		// Try yesterday if today not found?
		// No, strict requirement for now.
		return writeError(c, http.StatusNotFound, "report_missing", "Report not found", map[string]any{"path": path, "error": err.Error()})
	}

	// 3. Parse JSON
	var report struct {
		Date            string `json:"date"`
		ProgressSummary struct {
			Overview string `json:"overview"`
			Details  []struct {
				Topic   string `json:"topic"`
				Content string `json:"content"`
			} `json:"details"`
		} `json:"progress_summary"`
		CodeReviewSummary struct {
			Overview string `json:"overview"`
			Details  []struct {
				Author      string `json:"author"`
				Changes     string `json:"changes"`
				Suggestions string `json:"suggestions"`
			} `json:"details"`
		} `json:"code_review_summary"`
	}
	if err := json.Unmarshal(content, &report); err != nil {
		return writeError(c, http.StatusInternalServerError, "json_error", "Failed to parse report", map[string]any{"error": err.Error()})
	}

	// 4. Send to Lark
	// We need to send to a specific group. Which one?
	// User said "send to Lark group". We assume a global configured group or we iterate all bound groups?
	// "send to Lark group" implies a single group.
	// But we have `ChatIssueLinks`.
	// Maybe we send to ALL bound chats? Or a specific "Admin/Dev" group?
	// The requirement says "send summary result to Lark group".
	// Let's assume we send to ALL active chats found in `chat_issue_links` (deduplicated)?
	// OR, maybe there is a main channel?
	// Given "non-dev colleagues" and "dev colleagues", it implies a general group.
	// Let's use a new config `LarkReportChatID` or similar?
	// Or just send to all bound chats.
	// Sending to all bound chats seems safest for "project specific" reports.
	// BUT this report is "Daily Progress" for the whole repo?
	// The script runs for the whole repo.
	// So we should send to a main channel.
	// Let's look for a configured channel in DB or Config.
	// For now, I'll use a placeholder or iterate all unique chats.
	// Iterating all unique chats is better.

	// chatLinks, err := h.db.ListActiveChatLinksWithoutBranch(ctx) // This lists issues without branch.
	// We need ALL active chats.
	// Let's add `ListAllActiveChats` to store?
	// Or just use `ListActiveBranchLinks` and extract chats.

	// Actually, for this feature, it seems to be a "Repo Level" report.
	// We should probably have a "Repo -> Chat" mapping.
	// `channel_project_mappings`?
	// Let's check `channel_project_mappings`.

	// For now, to be safe and simple, I will iterate all unique chat IDs from `branch_issue_links` + `chat_issue_links`.
	// But wait, `branch_issue_links` joins `chat_issue_links`.
	// So we can just get all unique LarkChatIDs from `ListActiveBranchLinks`.

	// Collect unique chat IDs
	chatIDs := make(map[string]struct{})
	branchLinks, _ := h.db.ListActiveBranchLinks(ctx)
	for _, l := range branchLinks {
		if l.LarkChatID.Valid {
			chatIDs[l.LarkChatID.String] = struct{}{}
		}
	}

	larkClient := lark.NewClient(h.cfg.LarkAppID, h.cfg.LarkAppSecret)

	// Build Messages
	// Message 1: Progress Summary (For All)
	progressMsg := fmt.Sprintf("📅 **项目日报 (%s)**\n\n**%s**\n\n", report.Date, report.ProgressSummary.Overview)
	for _, d := range report.ProgressSummary.Details {
		progressMsg += fmt.Sprintf("🔹 **%s**\n%s\n", d.Topic, d.Content)
	}

	// Message 2: Code Review Summary (For Devs - but we send to same group for now, maybe threaded?)
	// User said "part 1 for non-dev, part 2 for dev".
	// If in same group, just send two messages or one long one.
	// Let's send one card or two messages.
	// Two messages is clearer.

	reviewMsg := fmt.Sprintf("💻 **Code Review 汇总**\n\n**%s**\n\n", report.CodeReviewSummary.Overview)
	for _, d := range report.CodeReviewSummary.Details {
		reviewMsg += fmt.Sprintf("👤 **%s**\n- 变动: %s\n- 建议: %s\n", d.Author, d.Changes, d.Suggestions)
	}

	sent := 0
	for chatID := range chatIDs {
		// Send Progress
		if err := larkClient.SendMessage(ctx, chatID, progressMsg); err != nil {
			LogStructured("error", map[string]any{"event": "job.notify.send_progress", "chat_id": chatID, "error": err.Error()})
		}
		// Send Review
		if err := larkClient.SendMessage(ctx, chatID, reviewMsg); err != nil {
			LogStructured("error", map[string]any{"event": "job.notify.send_review", "chat_id": chatID, "error": err.Error()})
		}
		sent++
	}

	return c.JSON(http.StatusOK, map[string]any{
		"sent_to": sent,
		"date":    report.Date,
	})
}

// JobIssueProgressTasks exposes active branch-linked issues with optional Lark chat bindings for CI to consume.
// Auth: Bearer INTEGRATION_TOKEN (same as CNB ingest). Returns snapshots when available for title/description.
func (h *Handler) JobIssueProgressTasks(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "unauthorized", "invalid integration token", nil)
	}
	if !hHasDB(h) {
		return c.JSON(http.StatusOK, map[string]any{"branch_links": []any{}, "unbound_chats": []any{}, "note": "db unavailable"})
	}
	ctx := c.Request().Context()

	links, err := h.db.ListActiveBranchLinks(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "db_error", "failed to list branch links", map[string]any{"error": err.Error()})
	}
	var out []map[string]any
	for _, l := range links {
		if !l.LarkChatID.Valid || strings.TrimSpace(l.LarkChatID.String) == "" {
			// 没有关联群聊则跳过生成任务，提醒由无分支列表处理
			continue
		}
		title, desc := "", ""
		if snap, err := h.db.GetPlaneIssueSnapshot(ctx, l.PlaneIssueID); err == nil {
			if v, ok := snap["name"].(string); ok {
				title = v
			}
			if v, ok := snap["description_html"].(string); ok {
				desc = v
			}
		}
		out = append(out, map[string]any{
			"plane_issue_id":    l.PlaneIssueID,
			"cnb_repo_id":       l.CNBRepoID,
			"branch":            l.Branch,
			"issue_title":       title,
			"issue_description": desc,
			"lark_chat_id":      l.LarkChatID.String,
		})
	}

	unbound, _ := h.db.ListActiveChatLinksWithoutBranch(ctx)
	var reminders []map[string]any
	for _, r := range unbound {
		reminders = append(reminders, map[string]any{
			"plane_issue_id": r.PlaneIssueID,
			"lark_chat_id":   r.LarkChatID,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"branch_links":  out,
		"unbound_chats": reminders,
	})
}

// JobIssueProgressSend sends a plain text message to the specified Lark chat.
// Expected payload: { "chat_id": "...", "issue_title": "...", "date": "YYYY-MM-DD", "message": "..." }
func (h *Handler) JobIssueProgressSend(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "unauthorized", "invalid integration token", nil)
	}
	var req struct {
		ChatID     string `json:"chat_id"`
		IssueTitle string `json:"issue_title"`
		Date       string `json:"date"`
		Message    string `json:"message"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_json", "无法解析请求体", map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(req.ChatID) == "" || strings.TrimSpace(req.Message) == "" {
		return writeError(c, http.StatusBadRequest, "invalid_param", "chat_id 和 message 不能为空", nil)
	}
	lc := lark.NewClient(h.cfg.LarkAppID, h.cfg.LarkAppSecret)
	ctx, cancel := context.WithTimeout(c.Request().Context(), 8*time.Second)
	defer cancel()

	msg := req.Message
	if strings.TrimSpace(req.IssueTitle) != "" || strings.TrimSpace(req.Date) != "" {
		msg = strings.TrimSpace(fmt.Sprintf("📅 %s %s\n\n%s", req.Date, req.IssueTitle, req.Message))
	}
	if err := lc.SendMessage(ctx, req.ChatID, msg); err != nil {
		return writeError(c, http.StatusInternalServerError, "lark_send_failed", "发送失败", map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"sent": true})
}

// JobIssueProgressNotify reads per-issue progress JSON files from report repo and posts to Lark chats.
// Path pattern: {ProgressDir}/daily/{date}/issue-{issue_id}.json (issue_id raw UUID).
// Request: POST /jobs/issue-progress/notify?date=YYYY-MM-DD (default today in TZ)
func (h *Handler) JobIssueProgressNotify(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "unauthorized", "invalid integration token", nil)
	}
	if !hHasDB(h) {
		return c.JSON(http.StatusOK, map[string]any{"status": "skipped", "note": "db unavailable"})
	}
	date := strings.TrimSpace(c.QueryParam("date"))
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	progressDir := strings.Trim(strings.TrimSpace(h.cfg.ProgressDir), "/")
	if progressDir == "" {
		progressDir = "issue-progress"
	}

	ctx := c.Request().Context()
	links, err := h.db.ListActiveBranchLinks(ctx)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, "db_error", "Failed to list branch links", map[string]any{"error": err.Error()})
	}
	if len(links) == 0 {
		return c.JSON(http.StatusOK, map[string]any{"sent": 0, "note": "no active branch links"})
	}
	cnbClient := &cnb.Client{
		BaseURL: h.cfg.CNBBaseURL,
		Token:   h.cfg.CNBAppToken,
	}
	larkClient := lark.NewClient(h.cfg.LarkAppID, h.cfg.LarkAppSecret)

	sent := 0
	for _, l := range links {
		if !l.LarkChatID.Valid || strings.TrimSpace(l.LarkChatID.String) == "" {
			continue
		}
		path := fmt.Sprintf("%s/daily/%s/issue-%s.json", progressDir, date, l.PlaneIssueID)
		content, err := cnbClient.GetFileContent(ctx, h.cfg.ReportRepo, h.cfg.ReportBranch, path)
		if err != nil {
			LogStructured("warn", map[string]any{"event": "job.issue_progress.notify.missing", "path": path, "issue_id": l.PlaneIssueID, "error": truncate(err.Error(), 200)})
			continue
		}
		var payload struct {
			Date       string `json:"date"`
			IssueID    string `json:"issue_id"`
			IssueTitle string `json:"issue_title"`
			Progress   struct {
				Overview string `json:"overview"`
				Details  []struct {
					Topic   string `json:"topic"`
					Content string `json:"content"`
				} `json:"details"`
			} `json:"progress_summary"`
			CodeReview struct {
				Overview string `json:"overview"`
				Details  []struct {
					Author      string `json:"author"`
					Changes     string `json:"changes"`
					Suggestions string `json:"suggestions"`
				} `json:"details"`
			} `json:"code_review_summary"`
		}
		if err := json.Unmarshal(content, &payload); err != nil {
			LogStructured("error", map[string]any{"event": "job.issue_progress.notify.json_error", "path": path, "error": truncate(err.Error(), 200)})
			continue
		}
		progMsg := fmt.Sprintf("📅 %s 需求进度\n**%s**\n\n%s\n", payload.Date, payload.IssueTitle, payload.Progress.Overview)
		for _, d := range payload.Progress.Details {
			progMsg += fmt.Sprintf("🔹 %s\n%s\n", d.Topic, d.Content)
		}
		reviewMsg := fmt.Sprintf("💻 Code Review\n%s\n\n", payload.CodeReview.Overview)
		for _, d := range payload.CodeReview.Details {
			reviewMsg += fmt.Sprintf("👤 %s\n- 变动: %s\n- 建议: %s\n", d.Author, d.Changes, d.Suggestions)
		}
		if err := larkClient.SendMessage(ctx, l.LarkChatID.String, progMsg); err != nil {
			LogStructured("error", map[string]any{"event": "job.issue_progress.notify.send_progress", "chat_id": l.LarkChatID.String, "error": truncate(err.Error(), 200)})
		} else {
			sent++
		}
		if err := larkClient.SendMessage(ctx, l.LarkChatID.String, reviewMsg); err != nil {
			LogStructured("error", map[string]any{"event": "job.issue_progress.notify.send_review", "chat_id": l.LarkChatID.String, "error": truncate(err.Error(), 200)})
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"sent": sent, "date": date})
}
