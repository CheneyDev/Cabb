-- Extend thread_links with optional metadata for Plane routing
ALTER TABLE thread_links
  ADD COLUMN IF NOT EXISTS plane_project_id uuid,
  ADD COLUMN IF NOT EXISTS workspace_slug text;

-- Optional index for lookups by plane_issue_id
CREATE INDEX IF NOT EXISTS idx_thread_links_plane_issue ON thread_links(plane_issue_id);
