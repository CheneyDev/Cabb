package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type AdminUser struct {
	ID           string
	Email        string
	DisplayName  string
	PasswordHash string
	Role         string
	Active       bool
	LastLoginAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AdminSession struct {
	ID          string
	AdminUserID string
	Token       string
	UserAgent   sql.NullString
	IPAddress   sql.NullString
	ExpiresAt   time.Time
	RevokedAt   sql.NullTime
	CreatedAt   time.Time
}

type AdminSessionWithUser struct {
	Session AdminSession
	User    AdminUser
}

func (d *DB) GetAdminUserByEmail(ctx context.Context, email string) (*AdminUser, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT id::text, email, display_name, password_hash, role, active, last_login_at, created_at, updated_at FROM admin_users WHERE email=$1 LIMIT 1`
	var u AdminUser
	if err := d.SQL.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Role, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) GetAdminUserByID(ctx context.Context, id string) (*AdminUser, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT id::text, email, display_name, password_hash, role, active, last_login_at, created_at, updated_at FROM admin_users WHERE id=$1::uuid`
	var u AdminUser
	if err := d.SQL.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Role, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) CreateAdminUser(ctx context.Context, email, displayName, passwordHash, role string, active bool) (*AdminUser, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `INSERT INTO admin_users (email, display_name, password_hash, role, active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5, now(), now()) RETURNING id::text, email, display_name, password_hash, role, active, last_login_at, created_at, updated_at`
	var u AdminUser
	if err := d.SQL.QueryRowContext(ctx, q, email, displayName, passwordHash, role, active).Scan(&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Role, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) UpdateAdminUser(ctx context.Context, id, displayName, role string, active bool) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const q = `UPDATE admin_users SET display_name=$2, role=$3, active=$4, updated_at=now() WHERE id=$1::uuid`
	res, err := d.SQL.ExecContext(ctx, q, id, displayName, role, active)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *DB) UpdateAdminUserPassword(ctx context.Context, id, passwordHash string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const q = `UPDATE admin_users SET password_hash=$2, updated_at=now() WHERE id=$1::uuid`
	res, err := d.SQL.ExecContext(ctx, q, id, passwordHash)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *DB) ListAdminUsers(ctx context.Context) ([]AdminUser, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT id::text, email, display_name, password_hash, role, active, last_login_at, created_at, updated_at FROM admin_users ORDER BY created_at ASC`
	rows, err := d.SQL.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AdminUser
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Role, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (d *DB) RecordAdminLogin(ctx context.Context, id string, ts time.Time) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const q = `UPDATE admin_users SET last_login_at=$2, updated_at=now() WHERE id=$1::uuid`
	_, err := d.SQL.ExecContext(ctx, q, id, ts)
	return err
}

func (d *DB) CreateAdminSession(ctx context.Context, adminUserID, token, userAgent, ip string, expiresAt time.Time) (*AdminSession, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `INSERT INTO admin_sessions (admin_user_id, session_token, user_agent, ip_address, expires_at, created_at) VALUES ($1::uuid,$2,$3,$4,$5, now()) RETURNING id::text, admin_user_id::text, session_token, user_agent, ip_address, expires_at, revoked_at, created_at`
	var s AdminSession
	if err := d.SQL.QueryRowContext(ctx, q, adminUserID, token, nullIfEmpty(userAgent), nullIfEmpty(ip), expiresAt).Scan(&s.ID, &s.AdminUserID, &s.Token, &s.UserAgent, &s.IPAddress, &s.ExpiresAt, &s.RevokedAt, &s.CreatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *DB) GetAdminSessionWithUser(ctx context.Context, token string) (*AdminSessionWithUser, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT s.id::text, s.admin_user_id::text, s.session_token, s.user_agent, s.ip_address, s.expires_at, s.revoked_at, s.created_at, u.id::text, u.email, u.display_name, u.password_hash, u.role, u.active, u.last_login_at, u.created_at, u.updated_at FROM admin_sessions s JOIN admin_users u ON u.id = s.admin_user_id WHERE s.session_token=$1`
	var res AdminSessionWithUser
	if err := d.SQL.QueryRowContext(ctx, q, token).Scan(&res.Session.ID, &res.Session.AdminUserID, &res.Session.Token, &res.Session.UserAgent, &res.Session.IPAddress, &res.Session.ExpiresAt, &res.Session.RevokedAt, &res.Session.CreatedAt, &res.User.ID, &res.User.Email, &res.User.DisplayName, &res.User.PasswordHash, &res.User.Role, &res.User.Active, &res.User.LastLoginAt, &res.User.CreatedAt, &res.User.UpdatedAt); err != nil {
		return nil, err
	}
	return &res, nil
}

func (d *DB) RevokeAdminSession(ctx context.Context, token string) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}
	const q = `UPDATE admin_sessions SET revoked_at=now() WHERE session_token=$1`
	res, err := d.SQL.ExecContext(ctx, q, token)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *DB) CleanupExpiredAdminSessions(ctx context.Context) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	const q = `DELETE FROM admin_sessions WHERE expires_at < now() OR revoked_at IS NOT NULL AND revoked_at < now() - interval '7 days'`
	_, err := d.SQL.ExecContext(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrConnDone) {
		return err
	}
	return nil
}
