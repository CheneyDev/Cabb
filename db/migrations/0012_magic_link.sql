-- 0012_magic_link.sql
-- MagicLink authentication via Feishu (Lark)

-- Create magic_link_tokens table for passwordless login
CREATE TABLE IF NOT EXISTS magic_link_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lark_open_id TEXT NOT NULL,
  lark_name TEXT,
  token TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_magic_link_tokens_token ON magic_link_tokens(token);
CREATE INDEX IF NOT EXISTS idx_magic_link_tokens_lark_open_id ON magic_link_tokens(lark_open_id);

-- Add lark_open_id to admin_users for linking
ALTER TABLE admin_users ADD COLUMN IF NOT EXISTS lark_open_id TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_users_lark_open_id ON admin_users(lark_open_id) WHERE lark_open_id IS NOT NULL;

-- Make password_hash nullable (for MagicLink-only users)
ALTER TABLE admin_users ALTER COLUMN password_hash DROP NOT NULL;

-- Cleanup: drop unused tables from initial migration
DROP TABLE IF EXISTS lark_accounts;
DROP TABLE IF EXISTS commit_summaries;
