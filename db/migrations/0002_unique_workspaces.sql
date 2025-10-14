-- Ensure uniqueness to avoid duplicate tokens per workspace
DO $$ BEGIN
  ALTER TABLE workspaces ADD CONSTRAINT uq_workspaces_plane_ws UNIQUE (plane_workspace_id, token_type);
EXCEPTION WHEN duplicate_table THEN NULL; END $$;

