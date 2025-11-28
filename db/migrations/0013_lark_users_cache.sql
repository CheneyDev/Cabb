-- 0013_lark_users_cache.sql
-- Cache Lark users for faster login page loading

CREATE TABLE IF NOT EXISTS lark_users_cache (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  open_id TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  en_name TEXT,
  avatar_origin TEXT,
  avatar_640 TEXT,
  avatar_240 TEXT,
  avatar_72 TEXT,
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_lark_users_cache_sort ON lark_users_cache(sort_order, name);
