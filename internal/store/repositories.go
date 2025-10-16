package store

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"
)

// EventDeliveries repo
func (d *DB) UpsertEventDelivery(ctx context.Context, source, eventType, deliveryID, payloadSHA, status string) error {
    if d == nil || d.SQL == nil { return nil }
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
    if d == nil || d.SQL == nil { return false, nil }
    const q = `SELECT payload_sha256 FROM event_deliveries WHERE source=$1 AND delivery_id=$2`
    var existing sql.NullString
    if err := d.SQL.QueryRowContext(ctx, q, source, deliveryID).Scan(&existing); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return false, nil }
        return false, err
    }
    if existing.Valid && existing.String == payloadSHA { return true, nil }
    return false, nil
}

func (d *DB) UpdateEventDeliveryStatus(ctx context.Context, source, deliveryID, status string, nextRetryAt *time.Time) error {
    if d == nil || d.SQL == nil { return nil }
    const q = `UPDATE event_deliveries SET status=$3, retries=CASE WHEN $3='queued' THEN retries ELSE retries END, next_retry_at=$4 WHERE source=$1 AND delivery_id=$2`
    var nra any
    if nextRetryAt != nil { nra = *nextRetryAt } else { nra = nil }
    _, err := d.SQL.ExecContext(ctx, q, source, deliveryID, status, nra)
    return err
}

// RepoProjectMappings repo
type RepoProjectMapping struct {
    PlaneProjectID   string
    PlaneWorkspaceID string
    CNBRepoID        string
    IssueOpenStateID   sql.NullString
    IssueClosedStateID sql.NullString
    Active bool
    SyncDirection    sql.NullString
    LabelSelector    sql.NullString
}

func (d *DB) GetRepoProjectMapping(ctx context.Context, cnbRepoID string) (*RepoProjectMapping, error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
    const q = `
SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector
FROM repo_project_mappings
WHERE cnb_repo_id=$1 AND active=true
LIMIT 1`
    var m RepoProjectMapping
    err := d.SQL.QueryRowContext(ctx, q, cnbRepoID).Scan(&m.PlaneProjectID, &m.PlaneWorkspaceID, &m.CNBRepoID, &m.IssueOpenStateID, &m.IssueClosedStateID, &m.Active, &m.SyncDirection, &m.LabelSelector)
    if err != nil { return nil, err }
    return &m, nil
}

func (d *DB) GetRepoProjectMappingByPlaneProject(ctx context.Context, planeProjectID string) (*RepoProjectMapping, error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
    const q = `
SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector
FROM repo_project_mappings
WHERE plane_project_id=$1::uuid AND active=true
LIMIT 1`
    var m RepoProjectMapping
    err := d.SQL.QueryRowContext(ctx, q, planeProjectID).Scan(&m.PlaneProjectID, &m.PlaneWorkspaceID, &m.CNBRepoID, &m.IssueOpenStateID, &m.IssueClosedStateID, &m.Active, &m.SyncDirection, &m.LabelSelector)
    if err != nil { return nil, err }
    return &m, nil
}

// List all mappings for a Plane project (active=true)
func (d *DB) ListRepoProjectMappingsByPlaneProject(ctx context.Context, planeProjectID string) ([]RepoProjectMapping, error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
    const q = `
SELECT plane_project_id::text, plane_workspace_id::text, cnb_repo_id, issue_open_state_id::text, issue_closed_state_id::text, active, sync_direction::text, label_selector
FROM repo_project_mappings
WHERE plane_project_id=$1::uuid AND active=true`
    rows, err := d.SQL.QueryContext(ctx, q, planeProjectID)
    if err != nil { return nil, err }
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
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE repo_project_mappings SET plane_workspace_id=$2::uuid, issue_open_state_id=$4::uuid, issue_closed_state_id=$5::uuid, active=$6, sync_direction=COALESCE($7::sync_direction, sync_direction), label_selector=COALESCE($8,label_selector), updated_at=now() WHERE plane_project_id=$1::uuid AND cnb_repo_id=$3`
    res, err := d.SQL.ExecContext(ctx, upd, m.PlaneProjectID, m.PlaneWorkspaceID, m.CNBRepoID, nullableUUID(m.IssueOpenStateID), nullableUUID(m.IssueClosedStateID), m.Active, nullableText(m.SyncDirection), nullIfEmpty(m.LabelSelector.String))
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO repo_project_mappings (plane_project_id, plane_workspace_id, cnb_repo_id, issue_open_state_id, issue_closed_state_id, active, sync_direction, label_selector, created_at, updated_at) VALUES ($1::uuid,$2::uuid,$3,$4::uuid,$5::uuid,$6,COALESCE($7::sync_direction,'cnb_to_plane'),$8,now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, m.PlaneProjectID, m.PlaneWorkspaceID, m.CNBRepoID, nullableUUID(m.IssueOpenStateID), nullableUUID(m.IssueClosedStateID), m.Active, nullableText(m.SyncDirection), nullIfEmpty(m.LabelSelector.String))
    return err
}

// ListRepoProjectMappings lists mappings with optional filters.
// If planeProjectID or cnbRepoID are empty strings, they are ignored.
// activeParam supports: "true" | "false" | "" (ignored).
func (d *DB) ListRepoProjectMappings(ctx context.Context, planeProjectID, cnbRepoID, activeParam string) ([]RepoProjectMapping, error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
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
    if err != nil { return nil, err }
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
    PlaneProjectID string
    CNBRepoID string
    DraftStateID sql.NullString
    OpenedStateID sql.NullString
    ReviewRequestedStateID sql.NullString
    ApprovedStateID sql.NullString
    MergedStateID sql.NullString
    ClosedStateID sql.NullString
}

func (d *DB) UpsertPRStateMapping(ctx context.Context, m PRStateMapping) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE pr_state_mappings SET draft_state_id=$3::uuid, opened_state_id=$4::uuid, review_requested_state_id=$5::uuid, approved_state_id=$6::uuid, merged_state_id=$7::uuid, closed_state_id=$8::uuid, updated_at=now() WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2`
    res, err := d.SQL.ExecContext(ctx, upd, m.PlaneProjectID, m.CNBRepoID, nullableUUID(m.DraftStateID), nullableUUID(m.OpenedStateID), nullableUUID(m.ReviewRequestedStateID), nullableUUID(m.ApprovedStateID), nullableUUID(m.MergedStateID), nullableUUID(m.ClosedStateID))
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO pr_state_mappings (plane_project_id, cnb_repo_id, draft_state_id, opened_state_id, review_requested_state_id, approved_state_id, merged_state_id, closed_state_id, created_at, updated_at) VALUES ($1::uuid,$2,$3::uuid,$4::uuid,$5::uuid,$6::uuid,$7::uuid,$8::uuid,now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, m.PlaneProjectID, m.CNBRepoID, nullableUUID(m.DraftStateID), nullableUUID(m.OpenedStateID), nullableUUID(m.ReviewRequestedStateID), nullableUUID(m.ApprovedStateID), nullableUUID(m.MergedStateID), nullableUUID(m.ClosedStateID))
    return err
}

func nullableUUID(v sql.NullString) any {
    if v.Valid && v.String != "" { return v.String }
    return nil
}

func nullableText(v sql.NullString) any {
    if v.Valid && v.String != "" { return v.String }
    return nil
}

// Workspaces repo (store tokens)
type Workspace struct {
    ID               string
    PlaneWorkspaceID string
    AppInstallationID sql.NullString
    TokenType        string
    AccessToken      string
    RefreshToken     sql.NullString
    ExpiresAt        sql.NullTime
    WorkspaceSlug    sql.NullString
    AppBot           sql.NullString
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (d *DB) UpsertWorkspaceToken(ctx context.Context, planeWorkspaceID, appInstallationID, tokenType, accessToken, refreshToken, expiresAt, workspaceSlug, appBot string) error {
    if d == nil || d.SQL == nil { return nil }
    // Try update existing row of same plane_workspace_id & token_type, else insert
    const upd = `
UPDATE workspaces SET app_installation_id=$2, access_token=$3, refresh_token=$4, expires_at=$5::timestamptz, workspace_slug=$6, app_bot=$7, updated_at=now()
WHERE plane_workspace_id=$1::uuid AND token_type=$8
`
    res, err := d.SQL.ExecContext(ctx, upd, planeWorkspaceID, appInstallationID, accessToken, nullIfEmpty(refreshToken), nullTime(expiresAt), nullIfEmpty(workspaceSlug), nullIfEmpty(appBot), tokenType)
    if err != nil { return err }
    rows, _ := res.RowsAffected()
    if rows > 0 { return nil }
    const ins = `
INSERT INTO workspaces (plane_workspace_id, app_installation_id, token_type, access_token, refresh_token, expires_at, workspace_slug, app_bot, created_at, updated_at)
VALUES ($1::uuid,$2,$3,$4,$5,$6::timestamptz,$7,$8,now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, planeWorkspaceID, appInstallationID, tokenType, accessToken, nullIfEmpty(refreshToken), nullTime(expiresAt), nullIfEmpty(workspaceSlug), nullIfEmpty(appBot))
    return err
}

func (d *DB) FindBotTokenByWorkspaceID(ctx context.Context, planeWorkspaceID string) (accessToken string, workspaceSlug string, err error) {
    if d == nil || d.SQL == nil { return "", "", sql.ErrConnDone }
    const q = `
SELECT access_token, COALESCE(workspace_slug, '') FROM workspaces
WHERE plane_workspace_id=$1::uuid AND token_type='bot'
ORDER BY updated_at DESC
LIMIT 1`
    err = d.SQL.QueryRowContext(ctx, q, planeWorkspaceID).Scan(&accessToken, &workspaceSlug)
    return
}

// IssueLinks repo
func (d *DB) FindPlaneIssueByCNBIssue(ctx context.Context, cnbRepoID, cnbIssueID string) (planeIssueID string, err error) {
    if d == nil || d.SQL == nil { return "", sql.ErrConnDone }
    const q = `SELECT plane_issue_id::text FROM issue_links WHERE cnb_repo_id=$1 AND cnb_issue_id=$2 LIMIT 1`
    err = d.SQL.QueryRowContext(ctx, q, cnbRepoID, cnbIssueID).Scan(&planeIssueID)
    return
}

func (d *DB) FindCNBIssueByPlaneIssue(ctx context.Context, planeIssueID string) (cnbRepoID, cnbIssueID string, err error) {
    if d == nil || d.SQL == nil { return "", "", sql.ErrConnDone }
    const q = `SELECT cnb_repo_id, cnb_issue_id FROM issue_links WHERE plane_issue_id=$1::uuid LIMIT 1`
    err = d.SQL.QueryRowContext(ctx, q, planeIssueID).Scan(&cnbRepoID, &cnbIssueID)
    return
}

// Label mappings
func (d *DB) UpsertLabelMapping(ctx context.Context, planeProjectID, cnbRepoID, cnbLabel, planeLabelID string) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE label_mappings SET plane_label_id=$4::uuid, updated_at=now() WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2 AND cnb_label=$3`
    res, err := d.SQL.ExecContext(ctx, upd, planeProjectID, cnbRepoID, cnbLabel, planeLabelID)
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO label_mappings (plane_project_id, cnb_repo_id, cnb_label, plane_label_id, created_at, updated_at) VALUES ($1::uuid,$2,$3,$4::uuid,now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, planeProjectID, cnbRepoID, cnbLabel, planeLabelID)
    return err
}

func (d *DB) MapCNBLabelsToPlane(ctx context.Context, planeProjectID, cnbRepoID string, labels []string) ([]string, error) {
    if d == nil || d.SQL == nil || len(labels) == 0 { return nil, nil }
    out := make([]string, 0, len(labels))
    const q = `SELECT plane_label_id::text FROM label_mappings WHERE plane_project_id=$1::uuid AND cnb_repo_id=$2 AND cnb_label=$3 LIMIT 1`
    for _, lb := range labels {
        var id sql.NullString
        if err := d.SQL.QueryRowContext(ctx, q, planeProjectID, cnbRepoID, lb).Scan(&id); err == nil {
            if id.Valid && id.String != "" { out = append(out, id.String) }
        }
    }
    return out, nil
}

// User mappings
func (d *DB) UpsertUserMapping(ctx context.Context, planeUserID, cnbUserID, displayName string) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE user_mappings SET plane_user_id=$1::uuid, display_name=COALESCE($3, display_name), updated_at=now() WHERE cnb_user_id=$2`
    res, err := d.SQL.ExecContext(ctx, upd, planeUserID, cnbUserID, nullIfEmpty(displayName))
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO user_mappings (plane_user_id, cnb_user_id, display_name, created_at, updated_at) VALUES ($1::uuid,$2,$3,now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, planeUserID, cnbUserID, nullIfEmpty(displayName))
    return err
}

func (d *DB) FindPlaneUserIDsByCNBUsers(ctx context.Context, cnbUserIDs []string) ([]string, error) {
    if d == nil || d.SQL == nil || len(cnbUserIDs) == 0 { return nil, nil }
    out := make([]string, 0, len(cnbUserIDs))
    const q = `SELECT plane_user_id::text FROM user_mappings WHERE cnb_user_id=$1 LIMIT 1`
    for _, u := range cnbUserIDs {
        var id sql.NullString
        if err := d.SQL.QueryRowContext(ctx, q, u).Scan(&id); err == nil {
            if id.Valid && id.String != "" { out = append(out, id.String) }
        }
    }
    return out, nil
}

func (d *DB) CreateIssueLink(ctx context.Context, planeIssueID, cnbRepoID, cnbIssueID string) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const q = `INSERT INTO issue_links (plane_issue_id, cnb_issue_id, cnb_repo_id, linked_at, created_at, updated_at) VALUES ($1::uuid,$2,$3,now(),now(),now()) ON CONFLICT DO NOTHING`
    _, err := d.SQL.ExecContext(ctx, q, planeIssueID, cnbIssueID, cnbRepoID)
    return err
}

// ListCNBIssuesByPlaneIssue returns all CNB links for a Plane issue.
func (d *DB) ListCNBIssuesByPlaneIssue(ctx context.Context, planeIssueID string) (links []struct{ Repo string; Number string }, err error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
    const q = `SELECT cnb_repo_id, cnb_issue_id FROM issue_links WHERE plane_issue_id=$1::uuid`
    rows, err := d.SQL.QueryContext(ctx, q, planeIssueID)
    if err != nil { return nil, err }
    defer rows.Close()
    for rows.Next() {
        var repo, iid string
        if err := rows.Scan(&repo, &iid); err != nil { return nil, err }
        links = append(links, struct{ Repo string; Number string }{Repo: repo, Number: iid})
    }
    return links, rows.Err()
}

// Branch links
func (d *DB) UpsertBranchIssueLink(ctx context.Context, planeIssueID, cnbRepoID, branch string, isPrimary bool) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE branch_issue_links SET plane_issue_id=$1::uuid, is_primary=$4, active=true, deleted_at=NULL WHERE cnb_repo_id=$2 AND branch=$3`
    res, err := d.SQL.ExecContext(ctx, upd, planeIssueID, cnbRepoID, branch, isPrimary)
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO branch_issue_links (plane_issue_id, cnb_repo_id, branch, is_primary, created_at, active) VALUES ($1::uuid,$2,$3,$4,now(),true)`
    _, err = d.SQL.ExecContext(ctx, ins, planeIssueID, cnbRepoID, branch, isPrimary)
    return err
}

func (d *DB) DeactivateBranchIssueLink(ctx context.Context, cnbRepoID, branch string) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const q = `UPDATE branch_issue_links SET active=false, deleted_at=now() WHERE cnb_repo_id=$1 AND branch=$2`
    _, err := d.SQL.ExecContext(ctx, q, cnbRepoID, branch)
    return err
}

// ===== Lark (Feishu) mappings =====

// UpsertLarkThreadLink binds a Lark thread/root message to a Plane issue.
// Optional: planeProjectID/workspaceSlug stored when provided (non-empty).
func (d *DB) UpsertLarkThreadLink(ctx context.Context, larkThreadID, planeIssueID, planeProjectID, workspaceSlug string, syncEnabled bool) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE thread_links SET plane_issue_id=$2::uuid, sync_enabled=$5, workspace_slug=COALESCE($4, workspace_slug), plane_project_id=COALESCE($3::uuid, plane_project_id), updated_at=now() WHERE lark_thread_id=$1`
    res, err := d.SQL.ExecContext(ctx, upd, larkThreadID, planeIssueID, nullIfEmpty(planeProjectID), nullIfEmpty(workspaceSlug), syncEnabled)
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO thread_links (lark_thread_id, plane_issue_id, plane_project_id, workspace_slug, sync_enabled, linked_at, created_at, updated_at) VALUES ($1,$2::uuid, $3::uuid, $4, $5, now(),now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, larkThreadID, planeIssueID, nullIfEmpty(planeProjectID), nullIfEmpty(workspaceSlug), syncEnabled)
    return err
}

// FindLarkThreadByPlaneIssue returns a Lark thread id bound to the Plane issue.
func (d *DB) FindLarkThreadByPlaneIssue(ctx context.Context, planeIssueID string) (string, error) {
    if d == nil || d.SQL == nil { return "", sql.ErrConnDone }
    const q = `SELECT lark_thread_id FROM thread_links WHERE plane_issue_id=$1::uuid LIMIT 1`
    var tid sql.NullString
    if err := d.SQL.QueryRowContext(ctx, q, planeIssueID).Scan(&tid); err != nil {
        return "", err
    }
    if tid.Valid { return tid.String, nil }
    return "", sql.ErrNoRows
}

// LarkThreadLink holds thread↔issue link with optional metadata
type LarkThreadLink struct {
    LarkThreadID   string
    PlaneIssueID   string
    PlaneProjectID sql.NullString
    WorkspaceSlug  sql.NullString
    SyncEnabled    bool
}

func (d *DB) GetLarkThreadLink(ctx context.Context, larkThreadID string) (*LarkThreadLink, error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
    const q = `SELECT lark_thread_id, plane_issue_id::text, plane_project_id::text, workspace_slug, sync_enabled FROM thread_links WHERE lark_thread_id=$1 LIMIT 1`
    var tl LarkThreadLink
    err := d.SQL.QueryRowContext(ctx, q, larkThreadID).Scan(&tl.LarkThreadID, &tl.PlaneIssueID, &tl.PlaneProjectID, &tl.WorkspaceSlug, &tl.SyncEnabled)
    if err != nil { return nil, err }
    return &tl, nil
}

// Find bot token by workspace slug
func (d *DB) FindBotTokenByWorkspaceSlug(ctx context.Context, workspaceSlug string) (accessToken string, err error) {
    if d == nil || d.SQL == nil { return "", sql.ErrConnDone }
    const q = `SELECT access_token FROM workspaces WHERE workspace_slug=$1 AND token_type='bot' ORDER BY updated_at DESC LIMIT 1`
    err = d.SQL.QueryRowContext(ctx, q, workspaceSlug).Scan(&accessToken)
    return
}

// PR links
func (d *DB) UpsertPRLink(ctx context.Context, planeIssueID, cnbRepoID, prIID string) error {
    if d == nil || d.SQL == nil { return sql.ErrConnDone }
    const upd = `UPDATE pr_links SET plane_issue_id=$1::uuid, updated_at=now() WHERE cnb_repo_id=$2 AND pr_iid=$3`
    res, err := d.SQL.ExecContext(ctx, upd, planeIssueID, cnbRepoID, prIID)
    if err != nil { return err }
    if n, _ := res.RowsAffected(); n > 0 { return nil }
    const ins = `INSERT INTO pr_links (plane_issue_id, cnb_repo_id, pr_iid, created_at, updated_at) VALUES ($1::uuid,$2,$3,now(),now())`
    _, err = d.SQL.ExecContext(ctx, ins, planeIssueID, cnbRepoID, prIID)
    return err
}

func (d *DB) FindPlaneIssueByCNBPR(ctx context.Context, cnbRepoID, prIID string) (planeIssueID string, err error) {
    if d == nil || d.SQL == nil { return "", sql.ErrConnDone }
    const q = `SELECT plane_issue_id::text FROM pr_links WHERE cnb_repo_id=$1 AND pr_iid=$2 LIMIT 1`
    err = d.SQL.QueryRowContext(ctx, q, cnbRepoID, prIID).Scan(&planeIssueID)
    return
}

func (d *DB) GetPRStateMapping(ctx context.Context, cnbRepoID string) (*PRStateMapping, error) {
    if d == nil || d.SQL == nil { return nil, sql.ErrConnDone }
    const q = `SELECT plane_project_id::text, cnb_repo_id, draft_state_id::text, opened_state_id::text, review_requested_state_id::text, approved_state_id::text, merged_state_id::text, closed_state_id::text FROM pr_state_mappings WHERE cnb_repo_id=$1 LIMIT 1`
    var m PRStateMapping
    err := d.SQL.QueryRowContext(ctx, q, cnbRepoID).Scan(&m.PlaneProjectID, &m.CNBRepoID, &m.DraftStateID, &m.OpenedStateID, &m.ReviewRequestedStateID, &m.ApprovedStateID, &m.MergedStateID, &m.ClosedStateID)
    if err != nil { return nil, err }
    return &m, nil
}


// helpers
func nullIfEmpty(s string) any {
    if s == "" { return nil }
    return s
}
func nullTime(s string) any {
    if s == "" { return nil }
    return s
}
