-- Migration: Cleanup unused tables
-- Reason: Remove cnb_installations (never used) and workspaces (replaced by global PLANE_SERVICE_TOKEN)
-- Date: 2025-10-31

-- cnb_installations: 完全未使用，项目未实现 CNB OAuth
DROP TABLE IF EXISTS cnb_installations;

-- workspaces: 已废弃，现在使用全局 PLANE_SERVICE_TOKEN
-- - workspace_slug 已迁移到 repo_project_mappings (0003)
-- - workspace 元数据通过 plane_issue_snapshots/plane_project_snapshots 维护
-- - 表结构与代码不一致（代码引用不存在的 access_token/token_type 字段）
DROP TABLE IF EXISTS workspaces;

COMMENT ON TABLE repo_project_mappings IS 'Stores workspace_slug for each mapping; workspace metadata maintained in snapshot tables';
