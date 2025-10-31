package store

import (
	"context"
	"database/sql"
	"time"
)

// PlaneCredential represents a Plane Service Token record
type PlaneCredential struct {
	ID               string
	PlaneWorkspaceID string
	WorkspaceSlug    string
	Kind             string
	TokenEnc         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ListPlaneCredentials returns all plane credentials
func (d *DB) ListPlaneCredentials(ctx context.Context) ([]PlaneCredential, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}

	const q = `
		SELECT id, plane_workspace_id, workspace_slug, kind, token_enc, created_at, updated_at
		FROM plane_credentials
		ORDER BY workspace_slug, updated_at DESC
	`

	rows, err := d.SQL.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PlaneCredential
	for rows.Next() {
		var c PlaneCredential
		if err := rows.Scan(&c.ID, &c.PlaneWorkspaceID, &c.WorkspaceSlug, &c.Kind, &c.TokenEnc, &c.CreatedAt, &c.UpdatedAt); err != nil {
			continue
		}
		items = append(items, c)
	}
	return items, nil
}

// GetPlaneCredential returns a credential by ID
func (d *DB) GetPlaneCredential(ctx context.Context, id string) (*PlaneCredential, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}

	const q = `
		SELECT id, plane_workspace_id, workspace_slug, kind, token_enc, created_at, updated_at
		FROM plane_credentials
		WHERE id = $1::uuid
	`

	var c PlaneCredential
	err := d.SQL.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.PlaneWorkspaceID, &c.WorkspaceSlug, &c.Kind, &c.TokenEnc, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// UpsertPlaneCredential creates or updates a plane credential
// If plane_workspace_id + kind already exists, updates the token
func (d *DB) UpsertPlaneCredential(ctx context.Context, planeWorkspaceID, workspaceSlug, kind, tokenEnc string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}

	const q = `
		INSERT INTO plane_credentials (plane_workspace_id, workspace_slug, kind, token_enc, created_at, updated_at)
		VALUES ($1::uuid, $2, $3, $4, now(), now())
		ON CONFLICT (plane_workspace_id, kind)
		DO UPDATE SET
			workspace_slug = EXCLUDED.workspace_slug,
			token_enc = EXCLUDED.token_enc,
			updated_at = now()
	`

	_, err := d.SQL.ExecContext(ctx, q, planeWorkspaceID, workspaceSlug, kind, tokenEnc)
	return err
}

// DeletePlaneCredential deletes a credential by ID
func (d *DB) DeletePlaneCredential(ctx context.Context, id string) (bool, error) {
	if d == nil || d.SQL == nil {
		return false, sql.ErrConnDone
	}

	const q = `DELETE FROM plane_credentials WHERE id = $1::uuid`
	result, err := d.SQL.ExecContext(ctx, q, id)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}
