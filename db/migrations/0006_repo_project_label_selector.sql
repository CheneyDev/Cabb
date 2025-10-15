-- Add label_selector to support label-based routing for multi-repo fanout
ALTER TABLE repo_project_mappings
  ADD COLUMN IF NOT EXISTS label_selector text; -- comma/space separated label names (case-insensitive)

CREATE INDEX IF NOT EXISTS idx_rpm_plane_project ON repo_project_mappings(plane_project_id);
