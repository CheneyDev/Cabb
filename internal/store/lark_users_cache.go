package store

import (
	"context"
	"database/sql"
	"time"
)

type LarkUserCache struct {
	ID           string
	OpenID       string
	Name         string
	EnName       sql.NullString
	AvatarOrigin sql.NullString
	Avatar640    sql.NullString
	Avatar240    sql.NullString
	Avatar72     sql.NullString
	SortOrder    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ListLarkUsersCache returns all cached Lark users ordered by sort_order and name.
func (d *DB) ListLarkUsersCache(ctx context.Context) ([]LarkUserCache, error) {
	if d == nil || d.SQL == nil {
		return nil, sql.ErrConnDone
	}
	const q = `SELECT id::text, open_id, name, en_name, avatar_origin, avatar_640, avatar_240, avatar_72, sort_order, created_at, updated_at 
		FROM lark_users_cache ORDER BY sort_order, name`
	rows, err := d.SQL.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LarkUserCache
	for rows.Next() {
		var u LarkUserCache
		if err := rows.Scan(&u.ID, &u.OpenID, &u.Name, &u.EnName, &u.AvatarOrigin, &u.Avatar640, &u.Avatar240, &u.Avatar72, &u.SortOrder, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// RefreshLarkUsersCache replaces all cached users with new data.
func (d *DB) RefreshLarkUsersCache(ctx context.Context, users []LarkUserCache) error {
	if d == nil || d.SQL == nil {
		return sql.ErrConnDone
	}

	tx, err := d.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing cache
	if _, err := tx.ExecContext(ctx, `DELETE FROM lark_users_cache`); err != nil {
		return err
	}

	// Insert new users
	const insertQ = `INSERT INTO lark_users_cache (open_id, name, en_name, avatar_origin, avatar_640, avatar_240, avatar_72, sort_order, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())`
	for _, u := range users {
		if _, err := tx.ExecContext(ctx, insertQ, u.OpenID, u.Name, nullIfEmpty(u.EnName.String), nullIfEmpty(u.AvatarOrigin.String), nullIfEmpty(u.Avatar640.String), nullIfEmpty(u.Avatar240.String), nullIfEmpty(u.Avatar72.String), u.SortOrder); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetLarkUsersCacheCount returns the number of cached users.
func (d *DB) GetLarkUsersCacheCount(ctx context.Context) (int, error) {
	if d == nil || d.SQL == nil {
		return 0, sql.ErrConnDone
	}
	var count int
	if err := d.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM lark_users_cache`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// GetLarkUsersCacheLastUpdated returns the last updated time of the cache.
func (d *DB) GetLarkUsersCacheLastUpdated(ctx context.Context) (time.Time, error) {
	if d == nil || d.SQL == nil {
		return time.Time{}, sql.ErrConnDone
	}
	var t sql.NullTime
	if err := d.SQL.QueryRowContext(ctx, `SELECT MAX(updated_at) FROM lark_users_cache`).Scan(&t); err != nil {
		return time.Time{}, err
	}
	if !t.Valid {
		return time.Time{}, nil
	}
	return t.Time, nil
}
