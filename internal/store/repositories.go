package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// EventDeliveries repo
func (d *DB) UpsertEventDelivery(ctx context.Context, source, eventType, deliveryID, payloadSHA, status string) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	const q = `
INSERT INTO event_deliveries (source, event_type, delivery_id, payload_sha256, status, created_at)
VALUES ($1,$2,$3,$4,$5, now())
ON CONFLICT (source, delivery_id) DO UPDATE SET
  event_type = EXCLUDED.event_type,
  payload_sha256 = EXCLUDED.payload_sha256,
  status = EXCLUDED.status
`
	_, err := d.SQL.ExecContext(ctx, q, source, eventType, deliveryID, payloadSHA, status)
	return err
}

func (d *DB) IsDuplicateDelivery(ctx context.Context, source, deliveryID, payloadSHA string) (bool, error) {
	if d == nil || d.SQL == nil {
		return false, nil
	}
	const q = `SELECT payload_sha256 FROM event_deliveries WHERE source=$1 AND delivery_id=$2`
	var existing sql.NullString
	if err := d.SQL.QueryRowContext(ctx, q, source, deliveryID).Scan(&existing); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if existing.Valid && existing.String == payloadSHA {
		return true, nil
	}
	return false, nil
}

func (d *DB) UpdateEventDeliveryStatus(ctx context.Context, source, deliveryID, status string, nextRetryAt *time.Time) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	const q = `UPDATE event_deliveries SET status=$3, retries=CASE WHEN $3='queued' THEN retries ELSE retries END, next_retry_at=$4 WHERE source=$1 AND delivery_id=$2`
	var nra any
	if nextRetryAt != nil {
		nra = *nextRetryAt
	} else {
		nra = nil
	}
	_, err := d.SQL.ExecContext(ctx, q, source, deliveryID, status, nra)
	return err
}

// RepoProjectMappings repo
type RepoProjectMapping struct {
	PlaneProjectID     string
	PlaneWorkspaceID   string
	CNBRepoID          string
	IssueOpenStateID   sql.NullString
	IssueClosedStateID sql.NullString
	Active             bool
	SyncDirection      sql.NullString
	LabelSelector      sql.NullString
}

func (d *DB) GetRepoProjectMapping(ctx context.Context, cnbRepoID string) (*RepoProjectMapping, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `
SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector
FROM repo_project_mappings
WHERE cnb_repo_id=$1 AND active=true
LIMIT 1`
	var m RepoProjectMapping
	err := d.SQL.QueryRowContext(ctx, q, cnbRepoID).Scan(&m.PlaneProjectID, &m.PlaneWorkspaceID, &m.CNBRepoID, &m.IssueOpenStateID, &m.IssueClosedStateID, &m.Active, &m.SyncDirection, &m.LabelSelector)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *DB) GetRepoProjectMappingByPlaneProject(ctx context.Context, planeProjectID string) (*RepoProjectMapping, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `
SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector
FROM repo_project_mappings
WHERE plane_project_id=$1::uuid AND active=true
LIMIT 1`
	var m RepoProjectMapping
	err := d.SQL.QueryRowContext(ctx, q, planeProjectID).Scan(&m.PlaneProjectID, &m.PlaneWorkspaceID, &m.CNBRepoID, &m.IssueOpenStateID, &m.IssueClosedStateID, &m.Active, &m.SyncDirection, &m.LabelSelector)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// List all mappings for a Plane project (active=true)
func (d *DB) ListRepoProjectMappingsByPlaneProject(ctx context.Context, planeProjectID string) ([]RepoProjectMapping, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `
SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector
FROM repo_project_mappings
WHERE plane_project_id=$1::uuid AND active=true`
	rows, err := d.SQL.QueryContext(ctx, q, planeProjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RepoProjectMapping
	for rows.Next() {
		var m RepoProjectMapping
		if err := rows.Scan(&m.PlaneProjectID, &m.PlaneWorkspaceID, &m.CNBRepoID, &m.IssueOpenStateID, &m.IssueClosedStateID, &m.Active, &m.SyncDirection, &m.LabelSelector); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (d *DB) UpsertRepoProjectMapping(ctx context.Context, m RepoProjectMapping) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE repo_project_mappings SET plane_workspace_id=$2::uuid, issue_open_state_id=$4::uuid, issue_closed_state_id=$5::uuid, active=$6, sync_direction=COALESCE($7::sync_direction, sync_direction), label_selector=COALESCE($8,label_selector), updated_at=now() WHERE plane_project_id=$1::uuid AND cnb_repo_id=$3`
	res, err := d.SQL.ExecContext(ctx, upd, m.PlaneProjectID, m.PlaneWorkspaceID, m.CNBRepoID, nullableUUID(m.IssueOpenStateID), nullableUUID(m.IssueClosedStateID), m.Active, nullableText(m.SyncDirection), nullIfEmpty(m.LabelSelector.String))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO repo_project_mappings (plane_project_id, plane_workspace_id, cnb_repo_id, issue_open_state_id, issue_closed_state_id, active, sync_direction, label_selector, created_at, updated_at) VALUES ($1::uuid,$2::uuid,$3,$4::uuid,$5::uuid,$6,COALESCE($7::sync_direction,'cnb_to_plane'),$8,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, m.PlaneProjectID, m.PlaneWorkspaceID, m.CNBRepoID, nullableUUID(m.IssueOpenStateID), nullableUUID(m.IssueClosedStateID), m.Active, nullableText(m.SyncDirection), nullIfEmpty(m.LabelSelector.String))
	return err
}

// ListRepoProjectMappings lists mappings with optional filters.
// If planeProjectID or cnbRepoID are empty strings, they are ignored.
// activeParam supports: "true" | "false" | "" (ignored).
func (d *DB) ListRepoProjectMappings(ctx context.Context, planeProjectID, cnbRepoID, activeParam string) ([]RepoProjectMapping, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	// Build dynamic WHERE conditions modestly without introducing heavy dependencies
	where := "WHERE 1=1"
	args := []any{}
	idx := 1
	if planeProjectID != "" {
		where += " AND plane_project_id=$" + itoa(idx) + "::uuid"
		args = append(args, planeProjectID)
		idx++
	}
	if cnbRepoID != "" {
		where += " AND cnb_repo_id=$" + itoa(idx)
		args = append(args, cnbRepoID)
		idx++
	}
	if activeParam == "true" || activeParam == "false" {
		where += " AND active=$" + itoa(idx)
		args = append(args, activeParam == "true")
		idx++
	}
	q := "SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector FROM repo_project_mappings " + where
	rows, err := d.SQL.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RepoProjectMapping
	for rows.Next() {
		var m RepoProjectMapping
		if err := rows.Scan(&m.PlaneProjectID, &m.PlaneWorkspaceID, &m.CNBRepoID, &m.IssueOpenStateID, &m.IssueClosedStateID, &m.Active, &m.SyncDirection, &m.LabelSelector); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// tiny helper to avoid strconv import bloat here
func itoa(i int) string { return fmt.Sprintf("%d", i) }

// PR state mapping
type PRStateMapping struct {
	PlaneProjectID         string
	CNBRepoID              string
	DraftStateID           sql.NullString
	OpenedStateID          sql.NullString
	ReviewRequestedStateID sql.NullString
	ApprovedStateID        sql.NullString
	MergedStateID          sql.NullString
	ClosedStateID          sql.NullString
}

// ===== Unified Integration Mappings =====

type IntegrationMappingRec struct {
	ScopeKind     string
	ScopeID       string
	MappingType   string
	LeftSystem    string
	LeftType      string
	LeftKey       string
	RightSystem   string
	RightType     string
	RightKey      string
	Bidirectional bool
	Extras        map[string]any
	Active        bool
}

// UpsertIntegrationMapping inserts or updates a single mapping row.
func (d *DB) UpsertIntegrationMapping(ctx context.Context, m IntegrationMappingRec) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	// Normalize strings
	norm := func(s string) string { return strings.TrimSpace(s) }
	m.ScopeKind = norm(m.ScopeKind)
	m.ScopeID = norm(m.ScopeID)
	m.MappingType = norm(m.MappingType)
	m.LeftSystem = norm(m.LeftSystem)
	m.LeftType = norm(m.LeftType)
	m.LeftKey = norm(m.LeftKey)
	m.RightSystem = norm(m.RightSystem)
	m.RightType = norm(m.RightType)
	m.RightKey = norm(m.RightKey)
	if strings.EqualFold(m.MappingType, "priority") {
		// Priorities: standardize plane priority key to lowercase label
		m.LeftKey = strings.ToLower(m.LeftKey)
	}
	extrasJSON := "{}"
	if m.Extras != nil {
		if b, err := json.Marshal(m.Extras); err == nil {
			extrasJSON = string(b)
		}
	}
	const upd = `UPDATE integration_mappings SET bidirectional=$11, extras=$12::jsonb, active=$13, updated_at=now()
                 WHERE scope_kind=$1 AND scope_id=COALESCE($2, scope_id) AND mapping_type=$3 AND left_system=$4 AND left_type=$5 AND left_key=$6 AND right_system=$7 AND right_type=$8 AND right_key=$9`
	res, err := d.SQL.ExecContext(ctx, upd, m.ScopeKind, nullIfEmpty(m.ScopeID), m.MappingType, m.LeftSystem, m.LeftType, m.LeftKey, m.RightSystem, m.RightType, m.RightKey, m.Bidirectional, extrasJSON, m.Active)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO integration_mappings (scope_kind, scope_id, mapping_type, left_system, left_type, left_key, right_system, right_type, right_key, bidirectional, extras, active, created_at, updated_at)
                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb, $12, now(), now())`
	_, err = d.SQL.ExecContext(ctx, ins, m.ScopeKind, nullIfEmpty(m.ScopeID), m.MappingType, m.LeftSystem, m.LeftType, m.LeftKey, m.RightSystem, m.RightType, m.RightKey, m.Bidirectional, extrasJSON, m.Active)
	return err
}

type IntegrationMappingRow struct {
	ID            int64
	ScopeKind     string
	ScopeID       sql.NullString
	MappingType   string
	LeftSystem    string
	LeftType      string
	LeftKey       string
	RightSystem   string
	RightType     string
	RightKey      string
	Bidirectional bool
	Extras        sql.NullString
	Active        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ListIntegrationMappings lists mappings by scope and type (minimal filters for admin UI).
func (d *DB) ListIntegrationMappings(ctx context.Context, scopeKind, scopeID, mappingType string) ([]IntegrationMappingRow, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	where := "WHERE 1=1"
	args := []any{}
	idx := 1
	if scopeKind != "" {
		where += " AND scope_kind=$" + itoa(idx)
		args = append(args, scopeKind)
		idx++
	}
	if scopeID != "" {
		where += " AND scope_id=$" + itoa(idx)
		args = append(args, scopeID)
		idx++
	}
	if mappingType != "" {
		where += " AND mapping_type=$" + itoa(idx)
		args = append(args, mappingType)
		idx++
	}
	q := "SELECT id, scope_kind, scope_id, mapping_type, left_system, left_type, left_key, right_system, right_type, right_key, bidirectional, extras::text, active, created_at, updated_at FROM integration_mappings " + where + " ORDER BY mapping_type, left_system, left_type, left_key"
	rows, err := d.SQL.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []IntegrationMappingRow
	for rows.Next() {
		var r IntegrationMappingRow
		if err := rows.Scan(&r.ID, &r.ScopeKind, &r.ScopeID, &r.MappingType, &r.LeftSystem, &r.LeftType, &r.LeftKey, &r.RightSystem, &r.RightType, &r.RightKey, &r.Bidirectional, &r.Extras, &r.Active, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		returnErr := rows.Err()
		_ = returnErr
		out = append(out, r)
	}
	return out, rows.Err()
}

// MapPlanePriorityToCNB resolves Plane priority to CNB priority via integration_mappings with scope fallback.
func (d *DB) MapPlanePriorityToCNB(ctx context.Context, planeProjectID, planePriority string) (string, bool, error) {
	if d == nil || d.SQL == nil {
		return "", false, sql.ErrConnDone
	}
	lp := strings.ToLower(strings.TrimSpace(planePriority))
	// First try project scope
	const q = `SELECT right_key FROM integration_mappings
               WHERE active=true AND mapping_type='priority'
                 AND left_system='plane' AND left_type='priority' AND left_key=$1
                 AND right_system='cnb' AND right_type='priority'
                 AND scope_kind='plane_project' AND scope_id=$2
               LIMIT 1`
	var out sql.NullString
	if planeProjectID != "" {
		if err := d.SQL.QueryRowContext(ctx, q, lp, planeProjectID).Scan(&out); err == nil && out.Valid {
			return out.String, true, nil
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", false, err
		}
	}
	// Fallback to global
	const qg = `SELECT right_key FROM integration_mappings
               WHERE active=true AND mapping_type='priority'
                 AND left_system='plane' AND left_type='priority' AND left_key=$1
                 AND right_system='cnb' AND right_type='priority'
                 AND scope_kind='global'
               LIMIT 1`
	if err := d.SQL.QueryRowContext(ctx, qg, lp).Scan(&out); err == nil && out.Valid {
		return out.String, true, nil
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", false, err
	}
	return "", false, nil
}

func (d *DB) UpsertPRStateMapping(ctx context.Context, m PRStateMapping) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE pr_state_mappings SET draft_state_id=$3::uuid, opened_state_id=$4::uuid, review_requested_state_id=$5::uuid, approved_state_id=$6::uuid, merged_state_id=$7::uuid, closed_state_id=$8::uuid, updated_at=now() WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2`
	res, err := d.SQL.ExecContext(ctx, upd, m.PlaneProjectID, m.CNBRepoID, nullableUUID(m.DraftStateID), nullableUUID(m.OpenedStateID), nullableUUID(m.ReviewRequestedStateID), nullableUUID(m.ApprovedStateID), nullableUUID(m.MergedStateID), nullableUUID(m.ClosedStateID))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO pr_state_mappings (plane_project_id, cnb_repo_id, draft_state_id, opened_state_id, review_requested_state_id, approved_state_id, merged_state_id, closed_state_id, created_at, updated_at) VALUES ($1::uuid,$2,$3::uuid,$4::uuid,$5::uuid,$6::uuid,$7::uuid,$8::uuid,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, m.PlaneProjectID, m.CNBRepoID, nullableUUID(m.DraftStateID), nullableUUID(m.OpenedStateID), nullableUUID(m.ReviewRequestedStateID), nullableUUID(m.ApprovedStateID), nullableUUID(m.MergedStateID), nullableUUID(m.ClosedStateID))
	return err
}

func nullableUUID(v sql.NullString) any {
	if v.Valid && v.String != "" {
		return v.String
	}
	return nil
}

func nullableText(v sql.NullString) any {
	if v.Valid && v.String != "" {
		return v.String
	}
	return nil
}

// Workspaces repo (store tokens)
type Workspace struct {
	ID                string
	PlaneWorkspaceID  string
	AppInstallationID sql.NullString
	TokenType         string
	AccessToken       string
	RefreshToken      sql.NullString
	ExpiresAt         sql.NullTime
	WorkspaceSlug     sql.NullString
	AppBot            sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (d *DB) GetWorkspaceBySlug(ctx context.Context, workspaceSlug string) (*Workspace, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `
SELECT id::text, plane_workspace_id::text, app_installation_id, token_type, access_token, refresh_token, expires_at, workspace_slug, app_bot, created_at, updated_at
FROM workspaces
WHERE workspace_slug=$1 AND token_type='bot'
ORDER BY updated_at DESC
LIMIT 1`
	var w Workspace
	if err := d.SQL.QueryRowContext(ctx, q, workspaceSlug).Scan(&w.ID, &w.PlaneWorkspaceID, &w.AppInstallationID, &w.TokenType, &w.AccessToken, &w.RefreshToken, &w.ExpiresAt, &w.WorkspaceSlug, &w.AppBot, &w.CreatedAt, &w.UpdatedAt); err != nil {
		return nil, err
	}
	return &w, nil
}

func (d *DB) UpsertWorkspaceToken(ctx context.Context, planeWorkspaceID, appInstallationID, tokenType, accessToken, refreshToken, expiresAt, workspaceSlug, appBot string) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	// Try update existing row of same plane_workspace_id & token_type, else insert
	const upd = `
UPDATE workspaces SET app_installation_id=$2, access_token=$3, refresh_token=$4, expires_at=$5::timestamptz, workspace_slug=$6, app_bot=$7, updated_at=now()
WHERE plane_workspace_id=$1::uuid AND token_type=$8
`
	res, err := d.SQL.ExecContext(ctx, upd, planeWorkspaceID, appInstallationID, accessToken, nullIfEmpty(refreshToken), nullTime(expiresAt), nullIfEmpty(workspaceSlug), nullIfEmpty(appBot), tokenType)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows > 0 {
		return nil
	}
	const ins = `
INSERT INTO workspaces (plane_workspace_id, app_installation_id, token_type, access_token, refresh_token, expires_at, workspace_slug, app_bot, created_at, updated_at)
VALUES ($1::uuid,$2,$3,$4,$5,$6::timestamptz,$7,$8,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, planeWorkspaceID, appInstallationID, tokenType, accessToken, nullIfEmpty(refreshToken), nullTime(expiresAt), nullIfEmpty(workspaceSlug), nullIfEmpty(appBot))
	return err
}

func (d *DB) FindBotTokenByWorkspaceID(ctx context.Context, planeWorkspaceID string) (accessToken string, workspaceSlug string, err error) {
	if d == nil || d.SQL == nil {
		return "", "", sql.ErrConnDone
	}
	const q = `
SELECT access_token, COALESCE(workspace_slug, '') FROM workspaces
WHERE plane_workspace_id=$1::uuid AND token_type='bot'
ORDER BY updated_at DESC
LIMIT 1`
	err = d.SQL.QueryRowContext(ctx, q, planeWorkspaceID).Scan(&accessToken, &workspaceSlug)
	return
}

// ===== Lark (Feishu) chat-level binding =====

type LarkChatIssueLink struct {
	LarkChatID     string
	LarkThreadID   sql.NullString
	PlaneIssueID   string
	PlaneProjectID sql.NullString
	WorkspaceSlug  sql.NullString
}

// UpsertLarkChatIssueLink binds a Lark chat to a Plane issue (single active binding per chat)
func (d *DB) UpsertLarkChatIssueLink(ctx context.Context, chatID, threadID, planeIssueID, planeProjectID, workspaceSlug string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE chat_issue_links SET lark_thread_id=$2, plane_issue_id=$3::uuid, plane_project_id=$4::uuid, workspace_slug=$5, updated_at=now() WHERE lark_chat_id=$1`
	res, err := d.SQL.ExecContext(ctx, upd, chatID, nullIfEmpty(threadID), planeIssueID, nullIfEmpty(planeProjectID), nullIfEmpty(workspaceSlug))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO chat_issue_links (lark_chat_id, lark_thread_id, plane_issue_id, plane_project_id, workspace_slug, created_at, updated_at) VALUES ($1,$2,$3::uuid,$4::uuid,$5,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, chatID, nullIfEmpty(threadID), planeIssueID, nullIfEmpty(planeProjectID), nullIfEmpty(workspaceSlug))
	return err
}

func (d *DB) GetLarkChatIssueLink(ctx context.Context, chatID string) (*LarkChatIssueLink, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT lark_chat_id, lark_thread_id, plane_issue_id::text, plane_project_id::text, workspace_slug FROM chat_issue_links WHERE lark_chat_id=$1 LIMIT 1`
	var l LarkChatIssueLink
	if err := d.SQL.QueryRowContext(ctx, q, chatID).Scan(&l.LarkChatID, &l.LarkThreadID, &l.PlaneIssueID, &l.PlaneProjectID, &l.WorkspaceSlug); err != nil {
		return nil, err
	}
	return &l, nil
}

// IssueLinks repo
func (d *DB) FindPlaneIssueByCNBIssue(ctx context.Context, cnbRepoID, cnbIssueID string) (planeIssueID string, err error) {
	if d == nil || d.SQL == nil {
		return "", sql.ErrConnDone
	}
	const q = `SELECT plane_issue_id::text FROM issue_links WHERE cnb_repo_id=$1 AND cnb_issue_id=$2 LIMIT 1`
	err = d.SQL.QueryRowContext(ctx, q, cnbRepoID, cnbIssueID).Scan(&planeIssueID)
	return
}

func (d *DB) FindCNBIssueByPlaneIssue(ctx context.Context, planeIssueID string) (cnbRepoID, cnbIssueID string, err error) {
	if d == nil || d.SQL == nil {
		return "", "", sql.ErrConnDone
	}
	const q = `SELECT cnb_repo_id, cnb_issue_id FROM issue_links WHERE plane_issue_id=$1::uuid LIMIT 1`
	err = d.SQL.QueryRowContext(ctx, q, planeIssueID).Scan(&cnbRepoID, &cnbIssueID)
	return
}

// Label mappings
func (d *DB) UpsertLabelMapping(ctx context.Context, planeProjectID, cnbRepoID, cnbLabel, planeLabelID string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE label_mappings SET plane_label_id=$4::uuid, updated_at=now() WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2 AND cnb_label=$3`
	res, err := d.SQL.ExecContext(ctx, upd, planeProjectID, cnbRepoID, cnbLabel, planeLabelID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO label_mappings (plane_project_id, cnb_repo_id, cnb_label, plane_label_id, created_at, updated_at) VALUES ($1::uuid,$2,$3,$4::uuid,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, planeProjectID, cnbRepoID, cnbLabel, planeLabelID)
	return err
}

func (d *DB) MapCNBLabelsToPlane(ctx context.Context, planeProjectID, cnbRepoID string, labels []string) ([]string, error) {
	if d == nil || d.SQL == nil || len(labels) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(labels))
	const q = `SELECT plane_label_id::text FROM label_mappings WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2 AND cnb_label=$3 LIMIT 1`
	for _, lb := range labels {
		var id sql.NullString
		if err := d.SQL.QueryRowContext(ctx, q, planeProjectID, cnbRepoID, lb).Scan(&id); err == nil {
			if id.Valid && id.String != "" {
				out = append(out, id.String)
			}
		}
	}
	return out, nil
}

// GetCNBManagedLabelIDs returns all Plane Label IDs that are managed by CNB
func (d *DB) GetCNBManagedLabelIDs(ctx context.Context, planeProjectID, cnbRepoID string) (map[string]bool, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT plane_label_id::text FROM label_mappings WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2`
	rows, err := d.SQL.QueryContext(ctx, q, planeProjectID, cnbRepoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	managedIDs := make(map[string]bool)
	for rows.Next() {
		var labelID string
		if err := rows.Scan(&labelID); err == nil && labelID != "" {
			managedIDs[labelID] = true
		}
	}
	return managedIDs, rows.Err()
}

// User mappings
type UserMapping struct {
	PlaneUserID string
	CNBUserID   sql.NullString
	LarkUserID  sql.NullString
	DisplayName sql.NullString
	ConnectedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *DB) UpsertUserMapping(ctx context.Context, planeUserID, cnbUserID, displayName string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE user_mappings SET plane_user_id=$1::uuid, display_name=COALESCE($3, display_name), updated_at=now() WHERE cnb_user_id=$2`
	res, err := d.SQL.ExecContext(ctx, upd, planeUserID, cnbUserID, nullIfEmpty(displayName))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO user_mappings (plane_user_id, cnb_user_id, display_name, created_at, updated_at) VALUES ($1::uuid,$2,$3,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, planeUserID, cnbUserID, nullIfEmpty(displayName))
	return err
}

func (d *DB) FindPlaneUserIDsByCNBUsers(ctx context.Context, cnbUserIDs []string) ([]string, error) {
	if d == nil || d.SQL == nil || len(cnbUserIDs) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(cnbUserIDs))
	const q = `SELECT plane_user_id::text FROM user_mappings WHERE cnb_user_id=$1 LIMIT 1`
	for _, u := range cnbUserIDs {
		var id sql.NullString
		if err := d.SQL.QueryRowContext(ctx, q, u).Scan(&id); err == nil {
			if id.Valid && id.String != "" {
				out = append(out, id.String)
			}
		}
	}
	return out, nil
}

// FindCNBUserIDsByPlaneUsers maps Plane user IDs (UUID strings) to CNB user IDs via user_mappings.
// Returns a de-duplicated list of non-empty CNB user IDs.
func (d *DB) FindCNBUserIDsByPlaneUsers(ctx context.Context, planeUserIDs []string) ([]string, error) {
	if d == nil || d.SQL == nil || len(planeUserIDs) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(planeUserIDs))
	seen := make(map[string]struct{}, len(planeUserIDs))
	const q = `SELECT cnb_user_id FROM user_mappings WHERE plane_user_id=$1::uuid LIMIT 1`
	for _, u := range planeUserIDs {
		var id sql.NullString
		if err := d.SQL.QueryRowContext(ctx, q, u).Scan(&id); err == nil {
			if id.Valid && id.String != "" {
				if _, ok := seen[id.String]; !ok {
					out = append(out, id.String)
					seen[id.String] = struct{}{}
				}
			}
		}
	}
	return out, nil
}

// ListUserMappings returns paginated user mapping rows for admin screens.
// limit defaults to 50 and is capped at 200 to avoid overwhelming the UI.
func (d *DB) ListUserMappings(ctx context.Context, planeUserID, cnbUserID, search string, limit int) ([]UserMapping, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	where := "WHERE 1=1"
	args := []any{}
	idx := 1
	if planeUserID != "" {
		where += " AND plane_user_id=$" + itoa(idx) + "::uuid"
		args = append(args, planeUserID)
		idx++
	}
	if cnbUserID != "" {
		where += " AND cnb_user_id=$" + itoa(idx)
		args = append(args, cnbUserID)
		idx++
	}
	if search != "" {
		like := "%" + search + "%"
		where += " AND (display_name ILIKE $" + itoa(idx) + " OR plane_user_id::text ILIKE $" + itoa(idx+1) + " OR cnb_user_id ILIKE $" + itoa(idx+2) + ")"
		args = append(args, like, like, like)
		idx += 3
	}
	q := "SELECT plane_user_id::text, cnb_user_id, lark_user_id, display_name, connected_at, created_at, updated_at FROM user_mappings " + where + " ORDER BY updated_at DESC LIMIT $" + itoa(idx)
	args = append(args, limit)
	rows, err := d.SQL.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserMapping
	for rows.Next() {
		var m UserMapping
		if err := rows.Scan(&m.PlaneUserID, &m.CNBUserID, &m.LarkUserID, &m.DisplayName, &m.ConnectedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (d *DB) CreateIssueLink(ctx context.Context, planeIssueID, cnbRepoID, cnbIssueID string) (bool, error) {
	if d == nil || d.SQL == nil {
		return false, sql.ErrConnDone
	}
	const q = `INSERT INTO issue_links (plane_issue_id, cnb_issue_id, cnb_repo_id, linked_at, created_at, updated_at) VALUES ($1::uuid,$2,$3,now(),now(),now()) ON CONFLICT DO NOTHING`
	res, err := d.SQL.ExecContext(ctx, q, planeIssueID, cnbIssueID, cnbRepoID)
	if err != nil {
		return false, err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return true, nil
	}
	return false, nil
}

type IssueLinkRow struct {
	PlaneIssueID     string
	CNBRepoID        sql.NullString
	CNBIssueID       sql.NullString
	PlaneProjectID   sql.NullString
	PlaneWorkspaceID sql.NullString
	LinkedAt         time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (d *DB) ListIssueLinks(ctx context.Context, planeIssueID, cnbRepoID, cnbIssueID string, limit int) ([]IssueLinkRow, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	if limit <= 0 {
		limit = 50
	} else if limit > 200 {
		limit = 200
	}
	where := "WHERE 1=1"
	args := []any{}
	idx := 1
	if strings.TrimSpace(planeIssueID) != "" {
		where += " AND plane_issue_id=$" + itoa(idx) + "::uuid"
		args = append(args, planeIssueID)
		idx++
	}
	if strings.TrimSpace(cnbRepoID) != "" {
		where += " AND cnb_repo_id=$" + itoa(idx)
		args = append(args, cnbRepoID)
		idx++
	}
	if strings.TrimSpace(cnbIssueID) != "" {
		where += " AND cnb_issue_id=$" + itoa(idx)
		args = append(args, cnbIssueID)
		idx++
	}
	order := " ORDER BY il.updated_at DESC LIMIT $" + itoa(idx)
	args = append(args, limit)
	query := `SELECT il.plane_issue_id::text,
       il.cnb_repo_id,
       il.cnb_issue_id,
       il.linked_at,
       il.created_at,
       il.updated_at,
       rpm.plane_project_id,
       rpm.plane_workspace_id
FROM issue_links il
LEFT JOIN LATERAL (
        SELECT plane_project_id::text AS plane_project_id,
               plane_workspace_id::text AS plane_workspace_id
        FROM repo_project_mappings
        WHERE cnb_repo_id=il.cnb_repo_id AND active=true
        ORDER BY repo_project_mappings.updated_at DESC
        LIMIT 1
) rpm ON true
` + where + order
	rows, err := d.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []IssueLinkRow
	for rows.Next() {
		var row IssueLinkRow
		if err := rows.Scan(&row.PlaneIssueID, &row.CNBRepoID, &row.CNBIssueID, &row.LinkedAt, &row.CreatedAt, &row.UpdatedAt, &row.PlaneProjectID, &row.PlaneWorkspaceID); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (d *DB) DeleteIssueLink(ctx context.Context, planeIssueID, cnbRepoID, cnbIssueID string) (bool, error) {
	if d == nil || d.SQL == nil {
		return false, sql.ErrConnDone
	}
	const q = `DELETE FROM issue_links WHERE plane_issue_id=$1::uuid AND cnb_repo_id=$2 AND cnb_issue_id=$3`
	res, err := d.SQL.ExecContext(ctx, q, planeIssueID, cnbRepoID, cnbIssueID)
	if err != nil {
		return false, err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return true, nil
	}
	return false, nil
}

// ListCNBIssuesByPlaneIssue returns all CNB links for a Plane issue.
func (d *DB) ListCNBIssuesByPlaneIssue(ctx context.Context, planeIssueID string) (links []struct {
	Repo   string
	Number string
}, err error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT cnb_repo_id, cnb_issue_id FROM issue_links WHERE plane_issue_id=$1::uuid`
	rows, err := d.SQL.QueryContext(ctx, q, planeIssueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var repo, iid string
		if err := rows.Scan(&repo, &iid); err != nil {
			return nil, err
		}
		links = append(links, struct {
			Repo   string
			Number string
		}{Repo: repo, Number: iid})
	}
	return links, rows.Err()
}

// Branch links
func (d *DB) UpsertBranchIssueLink(ctx context.Context, planeIssueID, cnbRepoID, branch string, isPrimary bool) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE branch_issue_links SET plane_issue_id=$1::uuid, is_primary=$4, active=true, deleted_at=NULL WHERE cnb_repo_id=$2 AND branch=$3`
	res, err := d.SQL.ExecContext(ctx, upd, planeIssueID, cnbRepoID, branch, isPrimary)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO branch_issue_links (plane_issue_id, cnb_repo_id, branch, is_primary, created_at, active) VALUES ($1::uuid,$2,$3,$4,now(),true)`
	_, err = d.SQL.ExecContext(ctx, ins, planeIssueID, cnbRepoID, branch, isPrimary)
	return err
}

func (d *DB) DeactivateBranchIssueLink(ctx context.Context, cnbRepoID, branch string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const q = `UPDATE branch_issue_links SET active=false, deleted_at=now() WHERE cnb_repo_id=$1 AND branch=$2`
	_, err := d.SQL.ExecContext(ctx, q, cnbRepoID, branch)
	return err
}

// ===== Lark (Feishu) mappings =====

// UpsertLarkThreadLink binds a Lark thread/root message to a Plane issue.
// Optional: planeProjectID/workspaceSlug stored when provided (non-empty).
func (d *DB) UpsertLarkThreadLink(ctx context.Context, larkThreadID, planeIssueID, planeProjectID, workspaceSlug string, syncEnabled bool) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE thread_links SET plane_issue_id=$2::uuid, sync_enabled=$5, workspace_slug=COALESCE($4, workspace_slug), plane_project_id=COALESCE($3::uuid, plane_project_id), updated_at=now() WHERE lark_thread_id=$1`
	res, err := d.SQL.ExecContext(ctx, upd, larkThreadID, planeIssueID, nullIfEmpty(planeProjectID), nullIfEmpty(workspaceSlug), syncEnabled)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO thread_links (lark_thread_id, plane_issue_id, plane_project_id, workspace_slug, sync_enabled, linked_at, created_at, updated_at) VALUES ($1,$2::uuid, $3::uuid, $4, $5, now(),now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, larkThreadID, planeIssueID, nullIfEmpty(planeProjectID), nullIfEmpty(workspaceSlug), syncEnabled)
	return err
}

// FindLarkThreadByPlaneIssue returns a Lark thread id bound to the Plane issue.
func (d *DB) FindLarkThreadByPlaneIssue(ctx context.Context, planeIssueID string) (string, error) {
	if d == nil || d.SQL == nil {
		return "", sql.ErrConnDone
	}
	const q = `SELECT lark_thread_id FROM thread_links WHERE plane_issue_id=$1::uuid LIMIT 1`
	var tid sql.NullString
	if err := d.SQL.QueryRowContext(ctx, q, planeIssueID).Scan(&tid); err != nil {
		return "", err
	}
	if tid.Valid {
		return tid.String, nil
	}
	return "", sql.ErrNoRows
}

// LarkThreadLink holds thread↔issue link with optional metadata
type LarkThreadLink struct {
	LarkThreadID   string
	PlaneIssueID   string
	PlaneProjectID sql.NullString
	WorkspaceSlug  sql.NullString
	LarkChatID     sql.NullString
	SyncEnabled    bool
	LinkedAt       time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (d *DB) GetLarkThreadLink(ctx context.Context, larkThreadID string) (*LarkThreadLink, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT tl.lark_thread_id, tl.plane_issue_id::text, tl.plane_project_id::text, tl.workspace_slug, cil.lark_chat_id, tl.sync_enabled, tl.linked_at, tl.created_at, tl.updated_at FROM thread_links tl LEFT JOIN chat_issue_links cil ON cil.lark_thread_id = tl.lark_thread_id WHERE tl.lark_thread_id=$1 LIMIT 1`
	var tl LarkThreadLink
	err := d.SQL.QueryRowContext(ctx, q, larkThreadID).Scan(&tl.LarkThreadID, &tl.PlaneIssueID, &tl.PlaneProjectID, &tl.WorkspaceSlug, &tl.LarkChatID, &tl.SyncEnabled, &tl.LinkedAt, &tl.CreatedAt, &tl.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tl, nil
}

func (d *DB) ListLarkThreadLinks(ctx context.Context, planeIssueID, larkThreadID string, syncEnabled *bool, limit int) ([]LarkThreadLink, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	if limit <= 0 {
		limit = 50
	} else if limit > 200 {
		limit = 200
	}
	where := "WHERE 1=1"
	args := []any{}
	idx := 1
	if strings.TrimSpace(planeIssueID) != "" {
		where += " AND plane_issue_id=$" + itoa(idx) + "::uuid"
		args = append(args, planeIssueID)
		idx++
	}
	if strings.TrimSpace(larkThreadID) != "" {
		where += " AND lark_thread_id=$" + itoa(idx)
		args = append(args, larkThreadID)
		idx++
	}
	if syncEnabled != nil {
		where += " AND sync_enabled=$" + itoa(idx)
		args = append(args, *syncEnabled)
		idx++
	}
	order := " ORDER BY updated_at DESC LIMIT $" + itoa(idx)
	args = append(args, limit)
	query := "SELECT tl.lark_thread_id, tl.plane_issue_id::text, tl.plane_project_id::text, tl.workspace_slug, cil.lark_chat_id, tl.sync_enabled, tl.linked_at, tl.created_at, tl.updated_at FROM thread_links tl LEFT JOIN chat_issue_links cil ON cil.lark_thread_id = tl.lark_thread_id " + where + order
	rows, err := d.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LarkThreadLink
	for rows.Next() {
		var tl LarkThreadLink
		if err := rows.Scan(&tl.LarkThreadID, &tl.PlaneIssueID, &tl.PlaneProjectID, &tl.WorkspaceSlug, &tl.LarkChatID, &tl.SyncEnabled, &tl.LinkedAt, &tl.CreatedAt, &tl.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, tl)
	}
	return out, rows.Err()
}

func (d *DB) DeleteLarkThreadLink(ctx context.Context, larkThreadID string) (bool, error) {
	if d == nil || d.SQL == nil {
		return false, sql.ErrConnDone
	}
	const q = `DELETE FROM thread_links WHERE lark_thread_id=$1`
	res, err := d.SQL.ExecContext(ctx, q, larkThreadID)
	if err != nil {
		return false, err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return true, nil
	}
	return false, nil
}

// Find bot token by workspace slug
func (d *DB) FindBotTokenByWorkspaceSlug(ctx context.Context, workspaceSlug string) (accessToken string, err error) {
	if d == nil || d.SQL == nil {
		return "", sql.ErrConnDone
	}
	const q = `SELECT access_token FROM workspaces WHERE workspace_slug=$1 AND token_type='bot' ORDER BY updated_at DESC LIMIT 1`
	err = d.SQL.QueryRowContext(ctx, q, workspaceSlug).Scan(&accessToken)
	return
}

// CleanupStaleThreadLinks deletes thread links that are not sync-enabled and not updated since before the cutoff.
func (d *DB) CleanupStaleThreadLinks(ctx context.Context, cutoff time.Time) (int64, error) {
	if d == nil || d.SQL == nil {
		return 0, sql.ErrConnDone
	}
	const q = `DELETE FROM thread_links WHERE sync_enabled=false AND updated_at < $1`
	res, err := d.SQL.ExecContext(ctx, q, cutoff)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// PR links
func (d *DB) UpsertPRLink(ctx context.Context, planeIssueID, cnbRepoID, prIID string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const upd = `UPDATE pr_links SET plane_issue_id=$1::uuid, updated_at=now() WHERE cnb_repo_id=$2 AND pr_iid=$3`
	res, err := d.SQL.ExecContext(ctx, upd, planeIssueID, cnbRepoID, prIID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	const ins = `INSERT INTO pr_links (plane_issue_id, cnb_repo_id, pr_iid, created_at, updated_at) VALUES ($1::uuid,$2,$3,now(),now())`
	_, err = d.SQL.ExecContext(ctx, ins, planeIssueID, cnbRepoID, prIID)
	return err
}

func (d *DB) FindPlaneIssueByCNBPR(ctx context.Context, cnbRepoID, prIID string) (planeIssueID string, err error) {
	if d == nil || d.SQL == nil {
		return "", sql.ErrConnDone
	}
	const q = `SELECT plane_issue_id::text FROM pr_links WHERE cnb_repo_id=$1 AND pr_iid=$2 LIMIT 1`
	err = d.SQL.QueryRowContext(ctx, q, cnbRepoID, prIID).Scan(&planeIssueID)
	return
}

func (d *DB) GetPRStateMapping(ctx context.Context, cnbRepoID string) (*PRStateMapping, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT plane_project_id::text, cnb_repo_id, draft_state_id::text, opened_state_id::text, review_requested_state_id::text, approved_state_id::text, merged_state_id::text, closed_state_id::text FROM pr_state_mappings WHERE cnb_repo_id=$1 LIMIT 1`
	var m PRStateMapping
	err := d.SQL.QueryRowContext(ctx, q, cnbRepoID).Scan(&m.PlaneProjectID, &m.CNBRepoID, &m.DraftStateID, &m.OpenedStateID, &m.ReviewRequestedStateID, &m.ApprovedStateID, &m.MergedStateID, &m.ClosedStateID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// FindDisplayNameByPlaneUserID returns a preferred display name mapped for the plane user, if any.
func (d *DB) FindDisplayNameByPlaneUserID(ctx context.Context, planeUserID string) (string, error) {
	if d == nil || d.SQL == nil {
		return "", sql.ErrConnDone
	}
	const q = `SELECT display_name FROM user_mappings WHERE plane_user_id=$1::uuid LIMIT 1`
	var name sql.NullString
	if err := d.SQL.QueryRowContext(ctx, q, planeUserID).Scan(&name); err != nil {
		return "", err
	}
	if name.Valid {
		return name.String, nil
	}
	return "", sql.ErrNoRows
}

// helpers
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
func nullTime(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// ChannelProjectLink represents a channel-project mapping
type ChannelProjectLink struct {
	LarkChatID     string
	PlaneProjectID string
	NotifyOnCreate bool
}

// GetChannelsByPlaneProject retrieves all Lark chat IDs mapped to a Plane project
func (d *DB) GetChannelsByPlaneProject(ctx context.Context, planeProjectID string) ([]ChannelProjectLink, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `
SELECT lark_chat_id, plane_project_id::text, notify_on_create
FROM channel_project_mappings
WHERE plane_project_id=$1::uuid
ORDER BY created_at DESC`
	rows, err := d.SQL.QueryContext(ctx, q, planeProjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []ChannelProjectLink
	for rows.Next() {
		var link ChannelProjectLink
		if err := rows.Scan(&link.LarkChatID, &link.PlaneProjectID, &link.NotifyOnCreate); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}
