-- 0009_drop_connected_at.sql
-- Remove unused connected_at column from user_mappings table

ALTER TABLE user_mappings DROP COLUMN IF EXISTS connected_at;
