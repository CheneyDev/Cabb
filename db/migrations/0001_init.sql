-- Schema initialization for plane-integration
-- Requires Postgres 16+

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
DO $$ BEGIN
  CREATE TYPE token_type AS ENUM ('bot', 'user');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE sync_direction AS ENUM ('cnb_to_plane', 'bidirectional');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE summary_status AS ENUM ('pending', 'succeeded', 'failed');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Workspaces: Plane installations & tokens
CREATE TABLE IF NOT EXISTS workspaces (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_workspace_id uuid NOT NULL,
  app_installation_id uuid,
  token_type token_type NOT NULL,
  access_token text NOT NULL,
  refresh_token text,
  expires_at timestamptz,
  workspace_slug text,
  app_bot text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_workspaces_plane_ws ON workspaces(plane_workspace_id);

-- CNB installations
CREATE TABLE IF NOT EXISTS cnb_installations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cnb_org_id text,
  installation_id text,
  access_token text NOT NULL,
  refresh_token text,
  expires_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cnb_installations_org ON cnb_installations(cnb_org_id);

-- Repo <-> Project mappings
CREATE TABLE IF NOT EXISTS repo_project_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  plane_workspace_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  issue_open_state_id uuid,
  issue_closed_state_id uuid,
  sync_direction sync_direction NOT NULL DEFAULT 'cnb_to_plane',
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, cnb_repo_id)
);
CREATE INDEX IF NOT EXISTS idx_rpm_repo ON repo_project_mappings(cnb_repo_id);

-- PR state mappings
CREATE TABLE IF NOT EXISTS pr_state_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  draft_state_id uuid,
  opened_state_id uuid,
  review_requested_state_id uuid,
  approved_state_id uuid,
  merged_state_id uuid,
  closed_state_id uuid,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, cnb_repo_id)
);

-- User mappings (Plane <-> external)
CREATE TABLE IF NOT EXISTS user_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_user_id uuid NOT NULL,
  cnb_user_id text,
  lark_user_id text,
  display_name text,
  connected_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_user_id),
  UNIQUE (cnb_user_id),
  UNIQUE (lark_user_id)
);

-- Cross-system issue links
CREATE TABLE IF NOT EXISTS issue_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_issue_id uuid NOT NULL,
  cnb_issue_id text,
  cnb_repo_id text,
  linked_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_issue_id, cnb_issue_id, cnb_repo_id)
);

-- Event deliveries (idempotency & retries)
CREATE TABLE IF NOT EXISTS event_deliveries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  source text NOT NULL,
  event_type text,
  delivery_id text,
  payload_sha256 text,
  status text,
  retries int NOT NULL DEFAULT 0,
  next_retry_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (source, delivery_id)
);
CREATE INDEX IF NOT EXISTS idx_event_deliveries_status ON event_deliveries(status);

-- Branch <-> Issue links
CREATE TABLE IF NOT EXISTS branch_issue_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_issue_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  branch text NOT NULL,
  is_primary boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  active boolean NOT NULL DEFAULT true,
  UNIQUE (cnb_repo_id, branch)
);

-- Commit summaries (AI)
CREATE TABLE IF NOT EXISTS commit_summaries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_issue_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  branch text NOT NULL,
  window_start timestamptz NOT NULL,
  window_end timestamptz NOT NULL,
  commit_count int NOT NULL DEFAULT 0,
  files_changed int NOT NULL DEFAULT 0,
  additions int NOT NULL DEFAULT 0,
  deletions int NOT NULL DEFAULT 0,
  summary_md text,
  model text,
  status summary_status NOT NULL DEFAULT 'pending',
  dedup_key text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (dedup_key)
);

-- Lark (Feishu) channel <-> project
CREATE TABLE IF NOT EXISTS channel_project_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  lark_chat_id text NOT NULL,
  notify_on_create boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, lark_chat_id)
);

-- Lark thread links
CREATE TABLE IF NOT EXISTS thread_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lark_thread_id text NOT NULL,
  plane_issue_id uuid NOT NULL,
  sync_enabled boolean NOT NULL DEFAULT false,
  linked_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (lark_thread_id)
);

-- Lark accounts
CREATE TABLE IF NOT EXISTS lark_accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_key text,
  app_access_token_expires_at timestamptz,
  lark_user_id text,
  plane_user_id uuid,
  connected_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_lark_accounts_tenant ON lark_accounts(tenant_key);
