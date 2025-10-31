-- Migration: Add workspace_slug to repo_project_mappings
-- Reason: Since plane_credentials is dropped, store workspace_slug with each mapping
-- Date: 2025-01-XX

ALTER TABLE repo_project_mappings 
ADD COLUMN IF NOT EXISTS workspace_slug text;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_rpm_workspace_slug ON repo_project_mappings(workspace_slug);

COMMENT ON COLUMN repo_project_mappings.workspace_slug IS 'Plane workspace slug - manually entered by user as Plane does not provide workspace list API';
