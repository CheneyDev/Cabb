-- 更新 repo_project_mappings 表中的 workspace_slug
-- 替换 'your-workspace-slug' 为实际的 Plane workspace slug

-- 查看当前 mapping 配置
SELECT plane_project_id, cnb_repo_id, workspace_slug 
FROM repo_project_mappings 
WHERE cnb_repo_id = '1024hub/plane-test';

-- 更新 workspace_slug（需要替换为实际值）
-- Workspace slug 可以从 Plane URL 中获取：https://app.plane.so/{workspace-slug}/projects
UPDATE repo_project_mappings 
SET workspace_slug = 'your-workspace-slug'  -- 替换为实际值
WHERE cnb_repo_id = '1024hub/plane-test';

-- 确认更新
SELECT plane_project_id, cnb_repo_id, workspace_slug 
FROM repo_project_mappings 
WHERE cnb_repo_id = '1024hub/plane-test';
