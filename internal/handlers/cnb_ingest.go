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

	planeapi "cabb/internal/plane"

	"github.com/labstack/echo/v4"
)

type cnbIssuePayload struct {
	Event       string   `json:"event"`
	Repo        string   `json:"repo"`
	IssueIID    string   `json:"issue_iid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CommentHTML string   `json:"comment_html"`
	Comment     string   `json:"comment"`
	Labels      []string `json:"labels"`
	Assignees   []string `json:"assignees"`
	Priority    string   `json:"priority"`
}

type cnbPRPayload struct {
	Event    string `json:"event"`
	Action   string `json:"action"`
	Repo     string `json:"repo"`
	PRIID    string `json:"pr_iid"`
	IssueIID string `json:"issue_iid"`
}

type cnbBranchPayload struct {
	Event    string `json:"event"`
	Action   string `json:"action"`
	Repo     string `json:"repo"`
	Branch   string `json:"branch"`
	IssueIID string `json:"issue_iid"`
}

func (h *Handler) CNBIngestIssue(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "invalid_token", "鉴权失败（Bearer token 不匹配）", nil)
	}
	body, sum, deliveryID, err := readAndDigest(c)
	if err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_body", "读取请求体失败", map[string]any{"error": err.Error()})
	}
	var p cnbIssuePayload
	if err := json.Unmarshal(body, &p); err != nil {
		return writeError(c, http.StatusUnprocessableEntity, "invalid_json", "JSON 解析失败", map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(p.Event) == "" || strings.TrimSpace(p.Repo) == "" || strings.TrimSpace(p.IssueIID) == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少必填字段：event/repo/issue_iid", nil)
	}
	// Idempotency (temporary in-memory). DB-backed available below.
	if h.dedupe != nil && h.dedupe.CheckAndMark("cnb.issue", deliveryID, sum) {
		return c.JSON(http.StatusOK, map[string]any{
			"accepted":       true,
			"source":         "cnb.issue",
			"event_type":     p.Event,
			"delivery_id":    deliveryID,
			"payload_sha256": sum,
			"status":         "duplicate",
		})
	}
	// DB-level idempotency
	if hHasDB(h) && deliveryID != "" {
		dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "cnb.issue", deliveryID, sum)
		if err == nil && dup {
			return c.JSON(http.StatusOK, map[string]any{
				"accepted":       true,
				"source":         "cnb.issue",
				"event_type":     p.Event,
				"delivery_id":    deliveryID,
				"payload_sha256": sum,
				"status":         "duplicate",
			})
		}
		_ = h.db.UpsertEventDelivery(c.Request().Context(), "cnb.issue", p.Event, deliveryID, sum, "queued")
	}

	go h.processCNBIssue(p, deliveryID, sum)
	return c.JSON(http.StatusAccepted, map[string]any{
		"accepted":       true,
		"source":         "cnb.issue",
		"event_type":     p.Event,
		"delivery_id":    deliveryID,
		"payload_sha256": sum,
		"status":         "queued",
	})
}

func (h *Handler) CNBIngestPR(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "invalid_token", "鉴权失败（Bearer token 不匹配）", nil)
	}
	body, sum, deliveryID, err := readAndDigest(c)
	if err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_body", "读取请求体失败", map[string]any{"error": err.Error()})
	}
	var p cnbPRPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return writeError(c, http.StatusUnprocessableEntity, "invalid_json", "JSON 解析失败", map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(p.Event) == "" || strings.TrimSpace(p.Repo) == "" || strings.TrimSpace(p.PRIID) == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少必填字段：event/repo/pr_iid", nil)
	}
	// Normalize event_type: prefer action if present
	evt := p.Event
	if strings.TrimSpace(p.Action) != "" {
		evt = p.Action
	}
	if h.dedupe != nil && h.dedupe.CheckAndMark("cnb.pr", deliveryID, sum) {
		return c.JSON(http.StatusOK, map[string]any{
			"accepted":       true,
			"source":         "cnb.pr",
			"event_type":     evt,
			"delivery_id":    deliveryID,
			"payload_sha256": sum,
			"status":         "duplicate",
		})
	}
	if hHasDB(h) && deliveryID != "" {
		dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "cnb.pr", deliveryID, sum)
		if err == nil && dup {
			return c.JSON(http.StatusOK, map[string]any{
				"accepted":       true,
				"source":         "cnb.pr",
				"event_type":     evt,
				"delivery_id":    deliveryID,
				"payload_sha256": sum,
				"status":         "duplicate",
			})
		}
		_ = h.db.UpsertEventDelivery(c.Request().Context(), "cnb.pr", evt, deliveryID, sum, "queued")
	}
	go h.processCNBPR(p, evt, deliveryID, sum)
	return c.JSON(http.StatusAccepted, map[string]any{
		"accepted":       true,
		"source":         "cnb.pr",
		"event_type":     evt,
		"delivery_id":    deliveryID,
		"payload_sha256": sum,
		"status":         "queued",
	})
}

func (h *Handler) CNBIngestBranch(c echo.Context) error {
	if !h.authorizeIntegration(c) {
		return writeError(c, http.StatusUnauthorized, "invalid_token", "鉴权失败（Bearer token 不匹配）", nil)
	}
	body, sum, deliveryID, err := readAndDigest(c)
	if err != nil {
		return writeError(c, http.StatusBadRequest, "invalid_body", "读取请求体失败", map[string]any{"error": err.Error()})
	}
	var p cnbBranchPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return writeError(c, http.StatusUnprocessableEntity, "invalid_json", "JSON 解析失败", map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(p.Event) == "" || strings.TrimSpace(p.Repo) == "" || strings.TrimSpace(p.Branch) == "" {
		return writeError(c, http.StatusBadRequest, "missing_fields", "缺少必填字段：event/repo/branch", nil)
	}
	evt := p.Event
	if strings.TrimSpace(p.Action) != "" {
		evt = p.Action
	}
	if h.dedupe != nil && h.dedupe.CheckAndMark("cnb.branch", deliveryID, sum) {
		return c.JSON(http.StatusOK, map[string]any{
			"accepted":       true,
			"source":         "cnb.branch",
			"event_type":     evt,
			"delivery_id":    deliveryID,
			"payload_sha256": sum,
			"status":         "duplicate",
		})
	}
	if hHasDB(h) && deliveryID != "" {
		dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "cnb.branch", deliveryID, sum)
		if err == nil && dup {
			return c.JSON(http.StatusOK, map[string]any{
				"accepted":       true,
				"source":         "cnb.branch",
				"event_type":     evt,
				"delivery_id":    deliveryID,
				"payload_sha256": sum,
				"status":         "duplicate",
			})
		}
		_ = h.db.UpsertEventDelivery(c.Request().Context(), "cnb.branch", evt, deliveryID, sum, "queued")
	}
	go h.processCNBBranch(p, evt, deliveryID, sum)
	return c.JSON(http.StatusAccepted, map[string]any{
		"accepted":       true,
		"source":         "cnb.branch",
		"event_type":     evt,
		"delivery_id":    deliveryID,
		"payload_sha256": sum,
		"status":         "queued",
	})
}

func (h *Handler) authorizeIntegration(c echo.Context) bool {
	if h.cfg.IntegrationToken == "" {
		return true // allow in scaffold when not configured
	}
	auth := c.Request().Header.Get("Authorization")
	want := "Bearer " + h.cfg.IntegrationToken
	return auth == want
}

// readAndDigest reads the raw body for JSON decoding and returns sha256(payload)
// Also extracts a best-effort delivery id from header "X-CNB-Delivery"
func readAndDigest(c echo.Context) (body []byte, sum string, deliveryID string, err error) {
	body, err = io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, "", "", err
	}
	// restore body for any later binds (not used now, but safe)
	c.Request().Body = io.NopCloser(strings.NewReader(string(body)))
	h := sha256.Sum256(body)
	sum = hex.EncodeToString(h[:])
	deliveryID = c.Request().Header.Get("X-CNB-Delivery")
	return body, sum, deliveryID, nil
}

// === processors ===
func (h *Handler) processCNBIssue(p cnbIssuePayload, deliveryID, sum string) {
	// Only handle a few core events for now: issue.open, issue.close
	if !h.cfg.PlaneOutboundEnabled {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if !hHasDB(h) {
		return
	}
	mapping, err := h.db.GetRepoProjectMapping(ctx, p.Repo)
	if err != nil {
		return
	}
	token, slug, err := h.db.FindBotTokenByWorkspaceID(ctx, mapping.PlaneWorkspaceID)
	if err != nil || token == "" || slug == "" {
		return
	}
	cl := &planeapi.Client{BaseURL: h.cfg.PlaneBaseURL}

	switch p.Event {
	case "issue.open":
		// check existing link
		if id, _ := h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID); id != "" {
			return
		}
		// create a minimal issue
		name := fmt.Sprintf("%s#%s", p.Repo, p.IssueIID)
		payload := map[string]any{"name": name}
		// optional description
		if strings.TrimSpace(p.Description) != "" {
			payload["description_html"] = fmt.Sprintf("<p>%s</p>", p.Description)
		}
		// optional labels
		if len(p.Labels) > 0 {
			if ids, _ := h.db.MapCNBLabelsToPlane(ctx, mapping.PlaneProjectID, p.Repo, p.Labels); len(ids) > 0 {
				payload["labels"] = ids
			}
		}
		// optional assignees
		if len(p.Assignees) > 0 {
			if uids, _ := h.db.FindPlaneUserIDsByCNBUsers(ctx, p.Assignees); len(uids) > 0 {
				payload["assignees"] = uids
			}
		}
		// optional priority (CNB -> Plane mapping)
		if pr := strings.TrimSpace(p.Priority); pr != "" {
			if planePri, ok := mapCNBPriorityToPlane(pr); ok {
				payload["priority"] = planePri
			}
		}
		if mapping.IssueOpenStateID.Valid {
			payload["state"] = mapping.IssueOpenStateID.String
		}
		LogStructured("info", map[string]any{
			"event":         "cnb.issue.planerpc",
			"delivery_id":   deliveryID,
			"repo":          p.Repo,
			"cnb_issue_iid": p.IssueIID,
			"op":            "create_issue",
			"fields_keys":   keysOf(payload),
		})
		issueID, err := cl.CreateIssue(ctx, token, slug, mapping.PlaneProjectID, payload)
		if err != nil {
			return
		}
                _, _ = h.db.CreateIssueLink(ctx, issueID, p.Repo, p.IssueIID)
		_ = cl.AddComment(ctx, token, slug, mapping.PlaneProjectID, issueID, fmt.Sprintf("<p>Linked from CNB issue <code>%s#%s</code></p>", p.Repo, p.IssueIID))
		LogStructured("info", map[string]any{
			"event":          "cnb.issue.planerpc",
			"delivery_id":    deliveryID,
			"repo":           p.Repo,
			"cnb_issue_iid":  p.IssueIID,
			"plane_issue_id": issueID,
			"op":             "create_issue",
			"result":         "created",
		})
		if deliveryID != "" {
			_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.issue", deliveryID, "succeeded", nil)
		}
	case "issue.close":
		if id, _ := h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID); id != "" {
			// set closed state if configured
			if mapping.IssueClosedStateID.Valid {
				_ = cl.PatchIssue(ctx, token, slug, mapping.PlaneProjectID, id, map[string]any{"state": mapping.IssueClosedStateID.String})
			}
		}
		if deliveryID != "" {
			_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.issue", deliveryID, "succeeded", nil)
		}
	case "issue.reopen":
		if id, _ := h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID); id != "" {
			if mapping.IssueOpenStateID.Valid {
				_ = cl.PatchIssue(ctx, token, slug, mapping.PlaneProjectID, id, map[string]any{"state": mapping.IssueOpenStateID.String})
			}
		}
		if deliveryID != "" {
			_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.issue", deliveryID, "succeeded", nil)
		}
	case "issue.comment", "comment.create":
		if id, _ := h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID); id != "" {
			content := p.CommentHTML
			if strings.TrimSpace(content) == "" && strings.TrimSpace(p.Comment) != "" {
				content = fmt.Sprintf("<p>%s</p>", p.Comment)
			}
			if strings.TrimSpace(content) != "" {
				_ = cl.AddComment(ctx, token, slug, mapping.PlaneProjectID, id, content)
				if deliveryID != "" {
					_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.issue", deliveryID, "succeeded", nil)
				}
			}
		}
	case "issue.update":
		if id, _ := h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID); id != "" {
			patch := map[string]any{}
			if strings.TrimSpace(p.Title) != "" {
				patch["name"] = p.Title
			}
			if len(p.Labels) > 0 {
				if ids, _ := h.db.MapCNBLabelsToPlane(ctx, mapping.PlaneProjectID, p.Repo, p.Labels); len(ids) > 0 {
					patch["labels"] = ids
				}
			}
			if len(p.Assignees) > 0 {
				if uids, _ := h.db.FindPlaneUserIDsByCNBUsers(ctx, p.Assignees); len(uids) > 0 {
					patch["assignees"] = uids
				}
			}
			if pr := strings.TrimSpace(p.Priority); pr != "" {
				if planePri, ok := mapCNBPriorityToPlane(pr); ok {
					patch["priority"] = planePri
				}
			}
			if len(patch) == 0 {
				LogStructured("info", map[string]any{
					"event":          "cnb.issue.planerpc",
					"delivery_id":    deliveryID,
					"repo":           p.Repo,
					"plane_issue_id": id,
					"op":             "update_issue",
					"decision":       "skip",
					"skip_reason":    "no_supported_fields",
				})
			} else {
				LogStructured("info", map[string]any{
					"event":          "cnb.issue.planerpc",
					"delivery_id":    deliveryID,
					"repo":           p.Repo,
					"plane_issue_id": id,
					"op":             "update_issue",
					"fields_keys":    keysOf(patch),
				})
				if err := cl.PatchIssue(ctx, token, slug, mapping.PlaneProjectID, id, patch); err == nil {
					LogStructured("info", map[string]any{
						"event":          "cnb.issue.planerpc",
						"delivery_id":    deliveryID,
						"repo":           p.Repo,
						"plane_issue_id": id,
						"op":             "update_issue",
						"result":         "updated",
					})
				}
			}
			if deliveryID != "" {
				_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.issue", deliveryID, "succeeded", nil)
			}
		}
	default:
		// no-op for other events
	}
}

// mapCNBPriorityToPlane converts CNB priority string to Plane priority label.
// CNB: P0/P1/P2/P3/""/-1P/-2P; Plane: urgent/high/medium/low/none
func mapCNBPriorityToPlane(cnb string) (string, bool) {
	s := strings.ToUpper(strings.TrimSpace(cnb))
	switch s {
	case "P0":
		return "urgent", true
	case "P1":
		return "high", true
	case "P2":
		return "medium", true
	case "P3":
		return "low", true
	case "", "NONE":
		return "none", true
	case "-1P", "-2P":
		// 待确认：CNB 负优先级在 Plane 中无直接等价，暂映射为 low
		return "low", true
	}
	return "", false
}

func (h *Handler) processCNBPR(p cnbPRPayload, evt, deliveryID, sum string) {
	if !h.cfg.PlaneOutboundEnabled {
		return
	}
	if !hHasDB(h) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	mapping, err := h.db.GetRepoProjectMapping(ctx, p.Repo)
	if err != nil {
		return
	}
	token, slug, err := h.db.FindBotTokenByWorkspaceID(ctx, mapping.PlaneWorkspaceID)
	if err != nil || token == "" || slug == "" {
		return
	}
	cl := &planeapi.Client{BaseURL: h.cfg.PlaneBaseURL}

	// resolve issue link: priority issue_iid; else pr_links
	var planeIssueID string
	if strings.TrimSpace(p.IssueIID) != "" {
		planeIssueID, _ = h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID)
		if planeIssueID != "" {
			_ = h.db.UpsertPRLink(ctx, planeIssueID, p.Repo, p.PRIID)
		}
	}
	if planeIssueID == "" {
		planeIssueID, _ = h.db.FindPlaneIssueByCNBPR(ctx, p.Repo, p.PRIID)
	}
	if planeIssueID == "" {
		return
	}

	// map PR event to Plane state via pr_state_mappings (if configured)
	prMap, _ := h.db.GetPRStateMapping(ctx, p.Repo)
	var targetState string
	switch strings.ToLower(evt) {
	case "opened", "open", "ready_for_review":
		if prMap != nil && prMap.OpenedStateID.Valid {
			targetState = prMap.OpenedStateID.String
		}
	case "review_requested":
		if prMap != nil && prMap.ReviewRequestedStateID.Valid {
			targetState = prMap.ReviewRequestedStateID.String
		}
	case "approved":
		if prMap != nil && prMap.ApprovedStateID.Valid {
			targetState = prMap.ApprovedStateID.String
		}
	case "merged", "merge":
		if prMap != nil && prMap.MergedStateID.Valid {
			targetState = prMap.MergedStateID.String
		}
	case "closed", "close":
		if prMap != nil && prMap.ClosedStateID.Valid {
			targetState = prMap.ClosedStateID.String
		}
	}
	if targetState != "" {
		if err := cl.PatchIssue(ctx, token, slug, mapping.PlaneProjectID, planeIssueID, map[string]any{"state": targetState}); err != nil {
			return
		}
	}
	// Add a marker comment on significant events
	switch strings.ToLower(evt) {
	case "opened", "merged", "closed":
		_ = cl.AddComment(ctx, token, slug, mapping.PlaneProjectID, planeIssueID, fmt.Sprintf("<p>PR <code>%s#%s</code> %s</p>", p.Repo, p.PRIID, strings.ToLower(evt)))
	}
	if deliveryID != "" {
		_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.pr", deliveryID, "succeeded", nil)
	}
}

func (h *Handler) processCNBBranch(p cnbBranchPayload, evt, deliveryID, sum string) {
	if !h.cfg.PlaneOutboundEnabled {
		return
	}
	if !hHasDB(h) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	mapping, err := h.db.GetRepoProjectMapping(ctx, p.Repo)
	if err != nil {
		return
	}
	token, slug, err := h.db.FindBotTokenByWorkspaceID(ctx, mapping.PlaneWorkspaceID)
	if err != nil || token == "" || slug == "" {
		return
	}
	cl := &planeapi.Client{BaseURL: h.cfg.PlaneBaseURL}

	// If payload specifies issue_iid, we can resolve the Plane issue
	if strings.TrimSpace(p.IssueIID) == "" {
		return
	}
	planeIssueID, err := h.db.FindPlaneIssueByCNBIssue(ctx, p.Repo, p.IssueIID)
	if err != nil || planeIssueID == "" {
		return
	}

	switch strings.ToLower(evt) {
	case "create", "created":
		_ = h.db.UpsertBranchIssueLink(ctx, planeIssueID, p.Repo, p.Branch, true)
		// move to open/in-progress state if configured
		if mapping.IssueOpenStateID.Valid {
			_ = cl.PatchIssue(ctx, token, slug, mapping.PlaneProjectID, planeIssueID, map[string]any{"state": mapping.IssueOpenStateID.String})
		}
		if deliveryID != "" {
			_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.branch", deliveryID, "succeeded", nil)
		}
	case "delete", "deleted":
		// TBD: 状态回滚策略（待确认）。目前仅标记 branch link 失效。
		_ = h.db.DeactivateBranchIssueLink(ctx, p.Repo, p.Branch)
		if deliveryID != "" {
			_ = h.db.UpdateEventDeliveryStatus(ctx, "cnb.branch", deliveryID, "succeeded", nil)
		}
	default:
		// ignore
	}
}
