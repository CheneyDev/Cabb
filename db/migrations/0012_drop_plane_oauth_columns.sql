-- Migration 0012: Drop Plane OAuth Columns (webhook-only refactor, Phase C)
-- Remove deprecated OAuth-related columns from workspaces table
-- Only execute after confirming no dependencies remain

ALTER TABLE workspaces DROP COLUMN IF EXISTS token_type;
ALTER TABLE workspaces DROP COLUMN IF EXISTS access_token;
ALTER TABLE workspaces DROP COLUMN IF EXISTS refresh_token;
ALTER TABLE workspaces DROP COLUMN IF EXISTS expires_at;
ALTER TABLE workspaces DROP COLUMN IF EXISTS app_installation_id;
ALTER TABLE workspaces DROP COLUMN IF EXISTS app_bot;

-- Keep plane_workspace_id and workspace_slug for webhook source identification
