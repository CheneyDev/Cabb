package store

import (
	"context"
	"database/sql"
	"time"
)

type MagicLinkToken struct {
	ID         string
	LarkOpenID string
	LarkName   sql.NullString
	Token      string
	ExpiresAt  time.Time
	UsedAt     sql.NullTime
	CreatedAt  time.Time
}

// CreateMagicLinkToken creates a new magic link token for a Lark user.
func (d *DB) CreateMagicLinkToken(ctx context.Context, larkOpenID, larkName, token string, expiresAt time.Time) (*MagicLinkToken, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `INSERT INTO magic_link_tokens (lark_open_id, lark_name, token, expires_at, created_at) 
		VALUES ($1, $2, $3, $4, now()) 
		RETURNING id::text, lark_open_id, lark_name, token, expires_at, used_at, created_at`
	var t MagicLinkToken
	if err := d.SQL.QueryRowContext(ctx, q, larkOpenID, nullIfEmpty(larkName), token, expiresAt).Scan(
		&t.ID, &t.LarkOpenID, &t.LarkName, &t.Token, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetMagicLinkToken retrieves a magic link token by token string.
func (d *DB) GetMagicLinkToken(ctx context.Context, token string) (*MagicLinkToken, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT id::text, lark_open_id, lark_name, token, expires_at, used_at, created_at 
		FROM magic_link_tokens WHERE token = $1`
	var t MagicLinkToken
	if err := d.SQL.QueryRowContext(ctx, q, token).Scan(
		&t.ID, &t.LarkOpenID, &t.LarkName, &t.Token, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// MarkMagicLinkTokenUsed marks a token as used.
func (d *DB) MarkMagicLinkTokenUsed(ctx context.Context, token string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const q = `UPDATE magic_link_tokens SET used_at = now() WHERE token = $1`
	_, err := d.SQL.ExecContext(ctx, q, token)
	return err
}

// CleanupExpiredMagicLinkTokens removes expired tokens older than 24 hours.
func (d *DB) CleanupExpiredMagicLinkTokens(ctx context.Context) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	const q = `DELETE FROM magic_link_tokens WHERE expires_at < now() - interval '24 hours'`
	_, err := d.SQL.ExecContext(ctx, q)
	return err
}

// GetOrCreateAdminUserByLarkOpenID finds or creates an admin user by Lark open_id.
func (d *DB) GetOrCreateAdminUserByLarkOpenID(ctx context.Context, larkOpenID, displayName, email string) (*AdminUser, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	
	// Try to find existing user
	const findQ = `SELECT id::text, email, display_name, password_hash, role, active, last_login_at, created_at, updated_at 
		FROM admin_users WHERE lark_open_id = $1 LIMIT 1`
	var u AdminUser
	err := d.SQL.QueryRowContext(ctx, findQ, larkOpenID).Scan(
		&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Role, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == nil {
		return &u, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	
	// Create new user (no password for MagicLink users)
	if email == "" {
		email = larkOpenID + "@lark.local"
	}
	const createQ = `INSERT INTO admin_users (email, display_name, password_hash, role, active, lark_open_id, created_at, updated_at) 
		VALUES ($1, $2, '', 'admin', true, $3, now(), now()) 
		RETURNING id::text, email, display_name, password_hash, role, active, last_login_at, created_at, updated_at`
	if err := d.SQL.QueryRowContext(ctx, createQ, email, displayName, larkOpenID).Scan(
		&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Role, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}
