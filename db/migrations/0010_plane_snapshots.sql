-- Migration 0010: Plane Snapshots (webhook-only refactor)
-- Create tables to cache Plane issue/project metadata from webhook events
-- Reduces dependency on Plane API queries for preview/notification/admin UI

-- Issue snapshots: minimal fields for preview & notification
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

-- Project snapshots: for admin UI & dropdown (optional)
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
