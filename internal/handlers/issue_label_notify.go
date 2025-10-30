package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

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
		return
	}

	// TODO: Implement business logic here based on requirements:
	// 1. Check repo-project mapping
	// 2. Find or create corresponding Plane issue
	// 3. Sync labels to Plane
	// 4. Send Lark notification if channel-project mapping exists
	// 5. Update issue links table

	// Example: Log the received notification
	_ = ctx
	_ = deliveryID
	_ = sum
	// Placeholder for actual implementation
}
