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
		Brief       string `json:"brief"`
		CommitCount int    `json:"commit_count"`
		Impact      string `json:"impact"`
		Members     []struct {
			Name         string   `json:"name"`
			LarkUserID   string   `json:"lark_user_id,omitempty"`
			Commits      int      `json:"commits"`
			Achievements []string `json:"achievements"`
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
	// Determine title based on type
	var titlePrefix string
	switch r.Meta.Type {
	case "daily":
		titlePrefix = "项目开发日报"
	case "weekly":
		titlePrefix = "项目开发周报"
	case "monthly":
		titlePrefix = "项目开发月报"
	default:
		titlePrefix = "项目开发汇报"
	}

	// Parse time range for header tags
	startDate := r.Meta.TimeRange.Start
	endDate := r.Meta.TimeRange.End
	if len(startDate) > 10 {
		startDate = startDate[:10]
	}
	if len(endDate) > 10 {
		endDate = endDate[:10]
	}

	// Build brief summary for each repo
	var briefSummary strings.Builder
	for _, repo := range r.Repos {
		briefSummary.WriteString("• **")
		briefSummary.WriteString(repo.DisplayName)
		briefSummary.WriteString("** ：")
		if repo.Brief != "" {
			briefSummary.WriteString(repo.Brief)
		} else if len(repo.Members) == 0 {
			briefSummary.WriteString("无提交记录")
		} else {
			briefSummary.WriteString("有 ")
			briefSummary.WriteString(strings.TrimSpace(strings.Split(strings.TrimPrefix(repo.Slug, "/"), "/")[len(strings.Split(repo.Slug, "/"))-1]))
			briefSummary.WriteString(" 次提交")
		}
		briefSummary.WriteString("\n")
	}

	// Build highlights content
	var highlightsContent strings.Builder
	for _, h := range r.Highlights {
		highlightsContent.WriteString("- ")
		highlightsContent.WriteString(h)
		highlightsContent.WriteString("\n")
	}
	highlightsStr := strings.TrimSpace(highlightsContent.String())
	if highlightsStr == "" {
		highlightsStr = "无"
	}

	// Build body elements
	elements := []map[string]any{
		{
			"tag":     "markdown",
			"content": "### <font color='violet'>开发进展</font>",
		},
		{
			"tag":       "markdown",
			"content":   strings.TrimSpace(briefSummary.String()),
			"text_size": "normal_v2",
		},
		{
			"tag":        "column_set",
			"flex_mode":  "stretch",
			"horizontal_spacing": "8px",
			"columns": []map[string]any{
				{
					"tag":              "column",
					"width":            "weighted",
					"weight":           1,
					"background_style": "blue-50",
					"padding":          "12px 12px 12px 12px",
					"vertical_spacing": "4px",
					"elements": []map[string]any{
						{
							"tag":       "markdown",
							"content":   "**<font color='blue'>关键进展</font>**",
							"text_size": "normal_v2",
							"icon": map[string]any{
								"tag":   "standard_icon",
								"token": "hot_outlined",
								"color": "grey",
							},
						},
						{
							"tag":     "markdown",
							"content": highlightsStr,
						},
					},
				},
			},
		},
		{"tag": "hr", "margin": "12px 0px 0px 0px"},
	}

	// Add repo sections
	for _, repo := range r.Repos {
		// Repo title
		elements = append(elements, map[string]any{
			"tag":     "markdown",
			"content": "### <font color='violet'>" + repo.DisplayName + "</font>",
		})

		// Build member elements with person_list
		memberElements := []map[string]any{
			{
				"tag":       "markdown",
				"content":   "**<font color='violet'>主要贡献者</font>**",
				"text_size": "normal_v2",
				"margin":    "0px 0px 8px 0px",
				"icon": map[string]any{
					"tag":   "standard_icon",
					"token": "group_outlined",
					"color": "grey",
				},
			},
		}

		// Add each member
		for _, m := range repo.Members {
			// Build achievements string
			var achievements strings.Builder
			for _, a := range m.Achievements {
				achievements.WriteString("- ")
				achievements.WriteString(a)
				achievements.WriteString("\n")
			}
			achievementsStr := strings.TrimSpace(achievements.String())
			if achievementsStr == "" {
				achievementsStr = "- 无详细记录"
			}

			// Person column (with or without person_list)
			personColumn := map[string]any{
				"tag":            "column",
				"width":          "weighted",
				"weight":         1,
				"vertical_align": "top",
				"elements":       []map[string]any{},
			}

			if m.LarkUserID != "" {
				// Use person_list component if we have Lark user ID
				personColumn["elements"] = []map[string]any{
					{
						"tag":  "person_list",
						"size": "small",
						"persons": []map[string]any{
							{"id": m.LarkUserID},
						},
						"show_avatar": true,
						"show_name":   true,
					},
				}
			} else {
				// Fallback to markdown with name
				personColumn["elements"] = []map[string]any{
					{
						"tag":       "markdown",
						"content":   "**" + m.Name + "**",
						"text_size": "normal_v2",
					},
				}
			}

			// Achievements column
			achievementsColumn := map[string]any{
				"tag":            "column",
				"width":          "weighted",
				"weight":         5,
				"vertical_align": "top",
				"elements": []map[string]any{
					{
						"tag":       "markdown",
						"content":   achievementsStr,
						"text_size": "normal_v2",
					},
					{"tag": "hr", "margin": "0px 0px 0px 0px"},
				},
			}

			memberElements = append(memberElements, map[string]any{
				"tag":                "column_set",
				"horizontal_spacing": "8px",
				"columns":            []map[string]any{personColumn, achievementsColumn},
			})
		}

		// Add impact if present
		if repo.Impact != "" {
			memberElements = append(memberElements,
				map[string]any{
					"tag":       "markdown",
					"content":   "**<font color='violet'>影响与价值</font>**",
					"text_size": "normal_v2",
					"margin":    "8px 0px 0px 0px",
					"icon": map[string]any{
						"tag":   "standard_icon",
						"token": "safe-settings_outlined",
						"color": "grey",
					},
				},
				map[string]any{
					"tag":       "markdown",
					"content":   repo.Impact,
					"text_size": "normal_v2",
				},
			)
		}

		// Wrap in column_set
		elements = append(elements, map[string]any{
			"tag":        "column_set",
			"flex_mode":  "stretch",
			"horizontal_spacing": "12px",
			"columns": []map[string]any{
				{
					"tag":              "column",
					"width":            "weighted",
					"weight":           2,
					"vertical_spacing": "4px",
					"elements":         memberElements,
				},
			},
		})

		elements = append(elements, map[string]any{"tag": "hr", "margin": "12px 0px 0px 0px"})
	}

	return map[string]any{
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
		"header": map[string]any{
			"title": map[string]any{
				"tag":     "plain_text",
				"content": titlePrefix,
			},
			"text_tag_list": []map[string]any{
				{
					"tag":   "text_tag",
					"text":  map[string]any{"tag": "plain_text", "content": startDate},
					"color": "blue",
				},
				{
					"tag":   "text_tag",
					"text":  map[string]any{"tag": "plain_text", "content": "-"},
					"color": "blue",
				},
				{
					"tag":   "text_tag",
					"text":  map[string]any{"tag": "plain_text", "content": endDate},
					"color": "blue",
				},
			},
			"template": "blue",
			"icon": map[string]any{
				"tag":   "standard_icon",
				"token": "code_outlined",
			},
			"padding": "12px 12px 12px 12px",
		},
		"body": map[string]any{
			"direction":          "vertical",
			"horizontal_spacing": "8px",
			"vertical_spacing":   "8px",
			"padding":            "8px 12px 8px 12px",
			"elements":           elements,
		},
	}
}
