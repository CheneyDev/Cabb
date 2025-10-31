-- Schema initialization for Plane integration service (Webhook-only)
-- Requires Postgres 16+

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
DO $$ BEGIN
  CREATE TYPE sync_direction AS ENUM ('cnb_to_plane', 'bidirectional');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE summary_status AS ENUM ('pending', 'succeeded', 'failed');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Workspaces: Plane webhook installations (no OAuth tokens)
CREATE TABLE IF NOT EXISTS workspaces (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_workspace_id uuid NOT NULL UNIQUE,
  workspace_slug text,
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

-- Repo <-> Plane project mappings
CREATE TABLE IF NOT EXISTS repo_project_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  plane_workspace_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  issue_open_state_id uuid,
  issue_closed_state_id uuid,
  sync_direction sync_direction NOT NULL DEFAULT 'cnb_to_plane',
  active boolean NOT NULL DEFAULT true,
  label_selector text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, cnb_repo_id)
);
CREATE INDEX IF NOT EXISTS idx_rpm_repo ON repo_project_mappings(cnb_repo_id);
CREATE INDEX IF NOT EXISTS idx_rpm_plane_project ON repo_project_mappings(plane_project_id);

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

-- User identity mappings
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

-- PR <-> Issue links
CREATE TABLE IF NOT EXISTS pr_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_issue_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  pr_iid text NOT NULL,
  linked_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (cnb_repo_id, pr_iid)
);

-- Label mappings between CNB and Plane
CREATE TABLE IF NOT EXISTS label_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  cnb_label text NOT NULL,
  plane_label_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, cnb_repo_id, cnb_label)
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

-- Commit summaries for AI digest
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

-- Lark (Feishu) channel <-> Plane project mappings
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
  plane_project_id uuid,
  workspace_slug text,
  sync_enabled boolean NOT NULL DEFAULT false,
  linked_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (lark_thread_id)
);
CREATE INDEX IF NOT EXISTS idx_thread_links_plane_issue ON thread_links(plane_issue_id);

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

-- Chat level issue bindings (non-thread)
CREATE TABLE IF NOT EXISTS chat_issue_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lark_chat_id text NOT NULL,
  lark_thread_id text,
  plane_issue_id uuid NOT NULL,
  plane_project_id uuid,
  workspace_slug text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (lark_chat_id)
);
CREATE INDEX IF NOT EXISTS idx_chat_issue_links_chat ON chat_issue_links(lark_chat_id);

-- Unified integration mappings across connectors
CREATE TABLE IF NOT EXISTS integration_mappings (
  id bigserial PRIMARY KEY,
  scope_kind text NOT NULL CHECK (scope_kind IN ('global','plane_workspace','plane_project','cnb_repo','lark_tenant')),
  scope_id text,
  mapping_type text NOT NULL CHECK (mapping_type IN ('user','priority','label','state','pr_state','custom')),
  left_system text NOT NULL CHECK (left_system IN ('plane','cnb','lark')),
  left_type text NOT NULL,
  left_key text NOT NULL,
  right_system text NOT NULL CHECK (right_system IN ('plane','cnb','lark')),
  right_type text NOT NULL,
  right_key text NOT NULL,
  bidirectional boolean NOT NULL DEFAULT true,
  extras jsonb NOT NULL DEFAULT '{}'::jsonb,
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (scope_kind, scope_id, mapping_type, left_system, left_type, left_key, right_system, right_type, right_key)
);
CREATE INDEX IF NOT EXISTS idx_int_map_left ON integration_mappings(scope_kind, scope_id, mapping_type, left_system, left_type, left_key) WHERE active = true;
CREATE INDEX IF NOT EXISTS idx_int_map_right ON integration_mappings(scope_kind, scope_id, mapping_type, right_system, right_type, right_key) WHERE active = true;

-- Admin authentication
CREATE TABLE IF NOT EXISTS admin_users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text NOT NULL UNIQUE,
  display_name text NOT NULL,
  password_hash text NOT NULL,
  role text NOT NULL DEFAULT 'admin',
  active boolean NOT NULL DEFAULT true,
  last_login_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS admin_sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  admin_user_id uuid NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
  session_token text NOT NULL UNIQUE,
  user_agent text,
  ip_address text,
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_admin_sessions_token ON admin_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_admin_users_active ON admin_users(active);

-- Plane webhook snapshots (issues & projects)
CREATE TABLE IF NOT EXISTS plane_issue_snapshots (
  plane_issue_id uuid PRIMARY KEY,
  plane_project_id uuid NOT NULL,
  plane_workspace_id uuid NOT NULL,
  workspace_slug text NOT NULL,
  project_slug text NOT NULL,
  name text NOT NULL,
  description_html text,
  state_name text,
  priority text,
  labels text[],
  assignee_ids uuid[],
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_plane_issue_snapshots_project ON plane_issue_snapshots(plane_project_id);
CREATE INDEX IF NOT EXISTS idx_plane_issue_snapshots_workspace ON plane_issue_snapshots(plane_workspace_id);
CREATE INDEX IF NOT EXISTS idx_plane_issue_snapshots_updated ON plane_issue_snapshots(updated_at DESC);

CREATE TABLE IF NOT EXISTS plane_project_snapshots (
  plane_project_id uuid PRIMARY KEY,
  plane_workspace_id uuid NOT NULL,
  workspace_slug text NOT NULL,
  project_slug text NOT NULL,
  name text NOT NULL,
  identifier text NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_plane_project_snapshots_workspace ON plane_project_snapshots(plane_workspace_id);
CREATE INDEX IF NOT EXISTS idx_plane_project_snapshots_slug ON plane_project_snapshots(workspace_slug, project_slug);

-- Plane Service Token credentials
CREATE TABLE IF NOT EXISTS plane_credentials (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_workspace_id uuid NOT NULL,
  workspace_slug text NOT NULL,
  kind text NOT NULL DEFAULT 'service',
  token_enc text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_workspace_id, kind)
);
CREATE INDEX IF NOT EXISTS idx_plane_credentials_workspace ON plane_credentials(plane_workspace_id);
CREATE INDEX IF NOT EXISTS idx_plane_credentials_slug ON plane_credentials(workspace_slug);
COMMENT ON TABLE plane_credentials IS 'Stores Plane Service Token for outbound API calls; webhook-only refactor';
COMMENT ON COLUMN plane_credentials.token_enc IS 'Encrypted token; use transparent encryption when available';

-- Seed default global priority mappings (Plane -> CNB)
INSERT INTO integration_mappings (
  scope_kind, scope_id, mapping_type,
  left_system, left_type, left_key,
  right_system, right_type, right_key,
  bidirectional, extras, active, created_at, updated_at
)
VALUES
  ('global', NULL, 'priority', 'plane', 'priority', 'urgent', 'cnb', 'priority', 'P0', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'high',   'cnb', 'priority', 'P1', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'medium', 'cnb', 'priority', 'P2', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'low',    'cnb', 'priority', 'P3', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'none',   'cnb', 'priority', ''  , true, '{}'::jsonb, true, now(), now())
ON CONFLICT (scope_kind, scope_id, mapping_type, left_system, left_type, left_key, right_system, right_type, right_key)
DO NOTHING;
