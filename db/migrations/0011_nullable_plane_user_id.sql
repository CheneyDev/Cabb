-- 0011_nullable_plane_user_id.sql
-- Allow plane_user_id and cnb_user_id to be nullable for more flexible user mappings
-- Primary use case: git_username + lark_user_id mapping for report notifications

-- Drop existing unique constraints
ALTER TABLE user_mappings DROP CONSTRAINT IF EXISTS user_mappings_plane_user_id_key;
ALTER TABLE user_mappings DROP CONSTRAINT IF EXISTS user_mappings_cnb_user_id_key;

-- Make plane_user_id nullable
ALTER TABLE user_mappings ALTER COLUMN plane_user_id DROP NOT NULL;

-- Create partial unique indexes (only on non-null values)
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_mappings_plane_user_id_unique 
  ON user_mappings (plane_user_id) WHERE plane_user_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_mappings_cnb_user_id_unique 
  ON user_mappings (cnb_user_id) WHERE cnb_user_id IS NOT NULL AND cnb_user_id != '';
