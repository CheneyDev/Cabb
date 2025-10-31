-- Migration 0011: Plane Credentials (webhook-only refactor)
-- Store Plane Service Token / PAT for optional outbound calls
-- Decoupled from legacy OAuth columns in workspaces table

CREATE TABLE IF NOT EXISTS plane_credentials (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_workspace_id uuid NOT NULL,
  workspace_slug text NOT NULL,
  kind text NOT NULL DEFAULT 'service', -- 'service' for Service Token / PAT
  token_enc text NOT NULL, -- encrypted token (透明加密待实现)
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_workspace_id, kind)
);

CREATE INDEX IF NOT EXISTS idx_plane_credentials_workspace ON plane_credentials(plane_workspace_id);
CREATE INDEX IF NOT EXISTS idx_plane_credentials_slug ON plane_credentials(workspace_slug);

COMMENT ON TABLE plane_credentials IS 'Stores Plane Service Token for PLANE_OUTBOUND_ENABLED mode; webhook-only refactor';
COMMENT ON COLUMN plane_credentials.token_enc IS 'Encrypted token; use transparent encryption when available';
