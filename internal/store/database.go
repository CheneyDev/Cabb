package store

import (
    "context"
    "database/sql"
)

type DB struct {
    SQL *sql.DB
}

func (d *DB) Ping(ctx context.Context) error {
    if d == nil || d.SQL == nil {
        return nil
    }
    return d.SQL.PingContext(ctx)
}

