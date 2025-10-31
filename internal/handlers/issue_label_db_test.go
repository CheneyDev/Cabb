package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cabb/internal/store"
	"cabb/pkg/config"
)

// TestWithDatabase 测试数据库集成（需要 DATABASE_URL 环境变量）
func TestWithDatabase(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping database integration tests")
	}

	db, err := store.Open(dbURL)
	require.NoError(t, err, "数据库连接失败")
	defer db.SQL.Close()

	ctx := context.Background()
	err = db.Ping(ctx)
	require.NoError(t, err, "数据库 Ping 失败")

	t.Run("EventDelivery_Idempotency", func(t *testing.T) {
		testEventDeliveryIdempotency(t, db)
	})

	t.Run("RepoProjectMapping_CRUD", func(t *testing.T) {
		testRepoProjectMapping(t, db)
	})

	t.Run("LabelMapping_CRUD", func(t *testing.T) {
		testLabelMapping(t, db)
	})

	t.Run("IssueLabelNotify_WithDB", func(t *testing.T) {
		testIssueLabelNotifyWithDB(t, db)
	})
}

// testEventDeliveryIdempotency 测试事件去重功能
func testEventDeliveryIdempotency(t *testing.T, db *store.DB) {
	ctx := context.Background()
	source := "test"
	eventType := "test.event"
	deliveryID := "test-delivery-" + time.Now().Format("20060102150405")
	payloadSHA := "abc123"

	// 首次插入
	err := db.UpsertEventDelivery(ctx, source, eventType, deliveryID, payloadSHA, "queued")
	require.NoError(t, err, "首次插入事件失败")

	// 检查是否重复
	isDup, err := db.IsDuplicateDelivery(ctx, source, deliveryID, payloadSHA)
	require.NoError(t, err, "检查重复失败")
	assert.True(t, isDup, "应该识别为重复事件")

	// 不同 payloadSHA 不应该重复
	isDup2, err := db.IsDuplicateDelivery(ctx, source, deliveryID, "different-sha")
	require.NoError(t, err)
	assert.False(t, isDup2, "不同 payload 不应识别为重复")
}

// testRepoProjectMapping 测试仓库-项目映射
func testRepoProjectMapping(t *testing.T, db *store.DB) {
	ctx := context.Background()
	cnbRepoID := "test-org/test-repo-" + time.Now().Format("150405")
	workspaceID := "11111111-1111-1111-1111-111111111111"
	projectID := "22222222-2222-2222-2222-222222222222"

	// 插入映射
	mapping := store.RepoProjectMapping{
		CNBRepoID:        cnbRepoID,
		PlaneWorkspaceID: workspaceID,
		PlaneProjectID:   projectID,
		Active:           true,
	}
	err := db.UpsertRepoProjectMapping(ctx, mapping)
	require.NoError(t, err, "插入映射失败")

	// 查询映射
	retrieved, err := db.GetRepoProjectMapping(ctx, cnbRepoID)
	require.NoError(t, err, "查询映射失败")
	require.NotNil(t, retrieved)

	assert.Equal(t, cnbRepoID, retrieved.CNBRepoID)
	assert.Equal(t, workspaceID, retrieved.PlaneWorkspaceID)
	assert.Equal(t, projectID, retrieved.PlaneProjectID)
	assert.True(t, retrieved.Active)

	// 更新映射（改变 workspace_id，保持 project_id 和 cnb_repo_id）
	newWorkspaceID := "33333333-3333-3333-3333-333333333333"
	mapping.PlaneWorkspaceID = newWorkspaceID
	err = db.UpsertRepoProjectMapping(ctx, mapping)
	require.NoError(t, err, "更新映射失败")

	retrieved2, err := db.GetRepoProjectMapping(ctx, cnbRepoID)
	require.NoError(t, err)
	assert.Equal(t, newWorkspaceID, retrieved2.PlaneWorkspaceID, "workspace 应该被更新")
	assert.Equal(t, projectID, retrieved2.PlaneProjectID, "project ID 应该保持不变")
}

// testLabelMapping 测试标签映射
func testLabelMapping(t *testing.T, db *store.DB) {
	ctx := context.Background()
	projectID := "44444444-4444-4444-4444-444444444444"
	repoSlug := "test-org/label-repo-" + time.Now().Format("150405")
	cnbLabel := "bug_CNB"
	planeLabelID := "55555555-5555-5555-5555-555555555555"

	// 插入标签映射
	err := db.UpsertLabelMapping(ctx, projectID, repoSlug, cnbLabel, planeLabelID)
	require.NoError(t, err, "插入标签映射失败")

	// 查询标签映射
	labels := []string{cnbLabel, "feature_CNB", "unknown_CNB"}
	labelIDs, err := db.MapCNBLabelsToPlane(ctx, projectID, repoSlug, labels)
	require.NoError(t, err, "查询标签映射失败")

	assert.Contains(t, labelIDs, planeLabelID, "应该包含映射的标签 ID")
	assert.Len(t, labelIDs, 1, "只有一个标签有映射")

	// 测试 CNB 管理的标签列表
	managedIDs, err := db.GetCNBManagedLabelIDs(ctx, projectID, repoSlug)
	require.NoError(t, err, "获取 CNB 管理标签失败")
	assert.True(t, managedIDs[planeLabelID], "应该标记为 CNB 管理")
}

// testIssueLabelNotifyWithDB 测试带数据库的完整流程
func testIssueLabelNotifyWithDB(t *testing.T, db *store.DB) {
	e := echo.New()
	cfg := config.Config{
		IntegrationToken: "test-token-db",
	}
	h := &Handler{
		cfg:    cfg,
		dedupe: NewDeduper(5 * time.Minute),
		db:     db,
	}

	payload := issueLabelNotifyPayload{
		RepoSlug:     "test/repo-with-db",
		IssueNumber:  999,
		IssueURL:     "https://cnb.cool/test/repo-with-db/-/issues/999",
		Title:        "DB Test Issue",
		State:        "open",
		Author:       authorInfo{Username: "dbtest", Nickname: "DB Tester"},
		Labels:       []string{"test_CNB"},
		LabelTrigger: "test_CNB",
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	deliveryID := "db-test-" + time.Now().Format("20060102150405")

	// 第一次请求
	body, _ := json.Marshal(payload)
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body))
	req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req1.Header.Set("Authorization", "Bearer test-token-db")
	req1.Header.Set("X-Delivery-ID", deliveryID)
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)

	err := h.IssueLabelNotify(c1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec1.Code)

	var resp1 map[string]interface{}
	json.Unmarshal(rec1.Body.Bytes(), &resp1)
	assert.Equal(t, float64(0), resp1["code"])

	// 等待异步处理完成
	time.Sleep(100 * time.Millisecond)

	// 第二次请求（应该在数据库级别去重）
	body2, _ := json.Marshal(payload)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/issues/label-notify", bytes.NewReader(body2))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req2.Header.Set("Authorization", "Bearer test-token-db")
	req2.Header.Set("X-Delivery-ID", deliveryID)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	err = h.IssueLabelNotify(c2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code)

	var resp2 map[string]interface{}
	json.Unmarshal(rec2.Body.Bytes(), &resp2)
	data := resp2["data"].(map[string]interface{})
	
	// 应该是 duplicate_db（数据库级去重）
	status := data["status"].(string)
	assert.Contains(t, []string{"duplicate", "duplicate_db"}, status, "应该识别为重复请求")
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping database connection test")
	}

	db, err := store.Open(dbURL)
	require.NoError(t, err, "数据库连接失败")
	defer db.SQL.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	assert.NoError(t, err, "数据库 Ping 失败")

	// 检查表是否存在
	var count int
	err = db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='event_deliveries'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "event_deliveries 表应该存在")
}
