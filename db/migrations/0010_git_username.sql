-- 0010_git_username.sql
-- Replace display_name with git_username for matching Git commit authors

ALTER TABLE user_mappings ADD COLUMN IF NOT EXISTS git_username text;
ALTER TABLE user_mappings DROP COLUMN IF EXISTS display_name;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_user_mappings_git_username ON user_mappings (git_username) WHERE git_username IS NOT NULL;
