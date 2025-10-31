package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations scans db/migrations for .sql files and executes them in lexicographic order.
// It is safe to call multiple times if migrations are idempotent (recommended: use IF NOT EXISTS).
func (d *DB) RunMigrations(ctx context.Context, root string) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	// root is typically the repo root ("."), we resolve migrations under db/migrations
	dir := filepath.Join(root, "db", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		// If migrations directory is unavailable, skip silently (entrypoint may handle it)
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// Ensure schema_migrations exists
	const createMeta = `CREATE TABLE IF NOT EXISTS schema_migrations (filename text PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())`
	if _, err := d.SQL.ExecContext(ctx, createMeta); err != nil {
		return err
	}
	// Load applied filenames
	applied := map[string]struct{}{}
	rows, err := d.SQL.QueryContext(ctx, `SELECT filename FROM schema_migrations`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var f string
			if err := rows.Scan(&f); err == nil {
				applied[f] = struct{}{}
			}
		}
		_ = rows.Err()
	}
	// collect *.sql
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(strings.ToLower(name), ".sql") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	sort.Strings(files)
	for _, f := range files {
		base := filepath.Base(f)
		if _, ok := applied[base]; ok {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		sqlText := string(b)
		if strings.TrimSpace(sqlText) == "" {
			continue
		}
		tx, err := d.SQL.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, sqlText); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", base, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1) ON CONFLICT (filename) DO NOTHING`, base); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s failed: %w", base, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// RunEmbeddedMigrations allows wiring an embedded FS if needed (not used currently).
func (d *DB) RunEmbeddedMigrations(ctx context.Context, fsys fs.FS, dir string) error {
	if d == nil || d.SQL == nil {
		return nil
	}
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		b, err := fs.ReadFile(fsys, filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if _, err := d.SQL.ExecContext(ctx, string(b)); err != nil {
			return fmt.Errorf("migration %s failed: %w", name, err)
		}
	}
	return nil
}
