package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"plane-integration/pkg/config"
)

// TestIssueLabelNotify_Success 测试成功案例
func TestIssueLabelNotify_Success(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-123",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     nil,
	}

	payload := issueLabelNotifyPayload{
		RepoSlug:     "test/repo",
		IssueNumber:  123,
		IssueURL:     "https://cnb.cool/test/repo/-/issues/123",
		Title:        "Test Issue",
		State:        "open",
		Author:       authorInfo{Username: "test", Nickname: "Test User"},
		Description:  "Test description",
		Labels:       []string{"bug_CNB", "feature"},
		LabelTrigger: "bug_CNB",
		UpdatedAt:    "2025-10-30T00:00:00Z",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer test-token-123")
	req.Header.Set("X-Delivery-ID", "test-delivery-1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.IssueLabelNotify(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
}

// TestIssueLabelNotify_Unauthorized 测试鉴权失败
func TestIssueLabelNotify_Unauthorized(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-123",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     nil,
	}

	payload := issueLabelNotifyPayload{
		RepoSlug:    "test/repo",
		IssueNumber: 123,
		Labels:      []string{"bug_CNB"},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer wrong-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.IssueLabelNotify(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "invalid_token", errObj["code"])
}

// TestIssueLabelNotify_MissingFields 测试缺少必填字段
func TestIssueLabelNotify_MissingFields(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-123",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     nil,
	}

	payload := map[string]interface{}{
		"repo_slug": "test/repo",
		// 缺少 issue_number 和 labels
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer test-token-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.IssueLabelNotify(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "missing_fields", errObj["code"])
}

// TestIssueLabelNotify_InvalidJSON 测试无效 JSON
func TestIssueLabelNotify_InvalidJSON(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-123",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     nil,
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer test-token-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.IssueLabelNotify(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "invalid_json", errObj["code"])
}

// TestIssueLabelNotify_Idempotency 测试幂等性
func TestIssueLabelNotify_Idempotency(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-123",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     nil,
	}

	payload := issueLabelNotifyPayload{
		RepoSlug:     "test/repo",
		IssueNumber:  456,
		IssueURL:     "https://cnb.cool/test/repo/-/issues/456",
		Title:        "Test Idempotency",
		State:        "open",
		Author:       authorInfo{Username: "test", Nickname: "Test User"},
		Labels:       []string{"test_CNB"},
		LabelTrigger: "test_CNB",
		UpdatedAt:    "2025-10-30T00:00:00Z",
	}

	deliveryID := "test-idempotent-123"

	// 第一次请求
	body, _ := json.Marshal(payload)
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body))
	req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req1.Header.Set("Authorization", "Bearer test-token-123")
	req1.Header.Set("X-Delivery-ID", deliveryID)
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)

	err := h.IssueLabelNotify(c1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// 第二次请求（应该被标记为 duplicate）
	body2, _ := json.Marshal(payload)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body2))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req2.Header.Set("Authorization", "Bearer test-token-123")
	req2.Header.Set("X-Delivery-ID", deliveryID)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	err = h.IssueLabelNotify(c2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code)

	var resp2 map[string]interface{}
	json.Unmarshal(rec2.Body.Bytes(), &resp2)
	data := resp2["data"].(map[string]interface{})
	assert.Equal(t, "duplicate", data["status"])
}

// TestIssueLabelSync_Success 测试简化版 API 成功案例
func TestIssueLabelSync_Success(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-123",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     nil,
	}

	payload := map[string]any{
		"repo_slug":    "test/repo",
		"issue_number": 789,
		"labels":       []string{"bug_CNB", "feature_CNB"},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-sync", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer test-token-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.IssueLabelNotifySimple(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
}

// TestFilterCNBLabels 测试标签过滤函数
func TestFilterCNBLabels(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "只有 CNB 标签",
			input:    []string{"bug_CNB", "feature_CNB"},
			expected: []string{"bug_CNB", "feature_CNB"},
		},
		{
			name:     "混合标签",
			input:    []string{"bug_CNB", "priority", "feature_CNB", "urgent"},
			expected: []string{"bug_CNB", "feature_CNB"},
		},
		{
			name:     "没有 CNB 标签",
			input:    []string{"bug", "priority", "urgent"},
			expected: nil,
		},
		{
			name:     "空标签列表",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "包含 emoji 的 CNB 标签",
			input:    []string{"🚧 处理中_CNB", "bug", "✅ 已完成_CNB"},
			expected: []string{"🚧 处理中_CNB", "✅ 已完成_CNB"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterCNBLabels(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
