package handlers

import (
	"strings"
	"time"

	"cabb/internal/store"
	"cabb/pkg/config"
)

// ReportJSON represents the AI-generated report JSON structure
type ReportJSON struct {
	Meta struct {
		Type      string `json:"type"`
		Label     string `json:"label"`
		TimeRange struct {
			Start    string `json:"start"`
			End      string `json:"end"`
			Timezone string `json:"timezone"`
		} `json:"time_range"`
		GeneratedAt string `json:"generated_at"`
	} `json:"meta"`
	Summary    string   `json:"summary"`
	Highlights []string `json:"highlights"`
	Repos      []struct {
		Slug        string `json:"slug"`
		DisplayName string `json:"display_name"`
		CommitCount int    `json:"commit_count"`
		Members     []struct {
			Name         string   `json:"name"`
			Role         string   `json:"role"`
			Commits      int      `json:"commits"`
			Achievements []string `json:"achievements"`
			Impact       string   `json:"impact"`
			Risks        []string `json:"risks"`
		} `json:"members"`
	} `json:"repos"`
}

// StartReportScheduler starts the scheduled report notification tasks.
// Schedule:
//   - Daily: Mon-Thu 09:55 (previous day), Fri 17:30 (same day)
//   - Weekly: Mon 09:55
//   - Monthly: 1st of month 11:00
func StartReportScheduler(cfg config.Config, db *store.DB) {
	// Check if Lark credentials are configured
	if cfg.LarkAppID == "" || cfg.LarkAppSecret == "" {
		return
	}

	go func() {
		loc := time.Local
		if tz := strings.TrimSpace(cfg.Timezone); tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now().In(loc)
			hhmm := now.Format("15:04")
			weekday := now.Weekday()
			day := now.Day()

			// Daily report
			switch weekday {
			case time.Monday, time.Tuesday, time.Wednesday, time.Thursday:
				// Send previous day's report at 09:55
				if hhmm == "09:55" {
					yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
					go sendReportWithConfig(cfg, db, "daily", yesterday)
				}
			case time.Friday:
				// Generate at 17:00, send at 17:30 (same day)
				if hhmm == "17:30" {
					today := now.Format("2006-01-02")
					go sendReportWithConfig(cfg, db, "daily", today)
				}
			}

			// Weekly report: Monday 09:55
			if weekday == time.Monday && hhmm == "09:55" {
				// Last week's report - use Monday's date as label or week range
				lastMonday := now.AddDate(0, 0, -7)
				lastSunday := now.AddDate(0, 0, -1)
				weekLabel := lastMonday.Format("2006-01-02") + "_to_" + lastSunday.Format("2006-01-02")
				go sendReportWithConfig(cfg, db, "weekly", weekLabel)
			}

			// Monthly report: 1st of month at 11:00
			if day == 1 && hhmm == "11:00" {
				// Last month's report
				lastMonth := now.AddDate(0, -1, 0)
				monthLabel := lastMonth.Format("2006-01")
				go sendReportWithConfig(cfg, db, "monthly", monthLabel)
			}
		}
	}()
}



func buildReportCard(r ReportJSON) map[string]any {
	// Determine title and color based on type
	var titlePrefix, headerColor string
	switch r.Meta.Type {
	case "daily":
		titlePrefix = "📋 日报"
		headerColor = "blue"
	case "weekly":
		titlePrefix = "📊 周报"
		headerColor = "green"
	case "monthly":
		titlePrefix = "📈 月报"
		headerColor = "purple"
	default:
		titlePrefix = "📋 汇报"
		headerColor = "blue"
	}

	// Build highlights content
	var highlights strings.Builder
	for _, h := range r.Highlights {
		highlights.WriteString("• ")
		highlights.WriteString(h)
		highlights.WriteString("\n")
	}
	highlightsStr := strings.TrimSpace(highlights.String())
	if highlightsStr == "" {
		highlightsStr = "无"
	}

	// Build members content (limit to 2000 chars)
	var members strings.Builder
	for _, repo := range r.Repos {
		members.WriteString("**[")
		members.WriteString(repo.Slug)
		members.WriteString("]** ")
		members.WriteString(repo.DisplayName)
		members.WriteString("\n")
		for _, m := range repo.Members {
			members.WriteString("• ")
			members.WriteString(m.Name)
			members.WriteString("：")
			members.WriteString(strings.Join(m.Achievements, "；"))
			members.WriteString("\n")
		}
		if members.Len() > 1800 {
			members.WriteString("...")
			break
		}
	}
	membersStr := strings.TrimSpace(members.String())
	if membersStr == "" {
		membersStr = "无提交记录"
	}

	return map[string]any{
		"schema": "2.0",
		"config": map[string]any{
			"wide_screen_mode": true,
		},
		"header": map[string]any{
			"title": map[string]any{
				"tag":     "plain_text",
				"content": titlePrefix + "（" + r.Meta.Label + "）",
			},
			"template": headerColor,
		},
		"body": map[string]any{
			"elements": []map[string]any{
				{
					"tag":     "markdown",
					"content": "**时间范围**：" + r.Meta.TimeRange.Start + " 至 " + r.Meta.TimeRange.End,
				},
				{
					"tag":     "markdown",
					"content": "**概览**\n" + r.Summary,
				},
				{
					"tag": "hr",
				},
				{
					"tag":     "markdown",
					"content": "**亮点**\n" + highlightsStr,
				},
				{
					"tag": "hr",
				},
				{
					"tag":     "markdown",
					"content": "**工作汇总**\n" + membersStr,
				},
			},
		},
	}
}
