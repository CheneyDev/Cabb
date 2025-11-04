-- 修复 BE-test-issue 仓库映射缺失
-- 使用说明：替换所有 <占位符> 为实际值后执行

-- 1. 参考现有映射（获取 UUID 格式示例）
SELECT cnb_repo_id, plane_project_id::text, workspace_slug 
FROM repo_project_mappings 
LIMIT 3;

-- 2. 创建 repo-project 映射
INSERT INTO repo_project_mappings (
  plane_project_id,
  plane_workspace_id,
  cnb_repo_id,
  workspace_slug,
  active,
  sync_direction,
  created_at,
  updated_at
) VALUES (
  '<PROJECT_UUID>',              -- 替换为 Plane Project ID
  '<WORKSPACE_UUID>',            -- 替换为 Plane Workspace ID
  '1024hub/Demo/BE-test-issue',
  '<WORKSPACE_SLUG>',            -- 如 'my-test'
  true,
  'cnb_to_plane',
  now(),
  now()
)
ON CONFLICT (plane_project_id, cnb_repo_id) DO UPDATE
SET active = true, updated_at = now();

-- 3. 创建标签映射（从 Plane API 获取 Label UUID）
INSERT INTO label_mappings (
  plane_project_id,
  cnb_repo_id,
  cnb_label,
  plane_label_id,
  created_at,
  updated_at
) VALUES 
  ('<PROJECT_UUID>', '1024hub/Demo/BE-test-issue', '🚧 处理中_CNB', '<LABEL_UUID_1>', now(), now()),
  ('<PROJECT_UUID>', '1024hub/Demo/BE-test-issue', '🧑🏻‍💻 进行中：后端_CNB', '<LABEL_UUID_2>', now(), now())
ON CONFLICT (plane_project_id, cnb_repo_id, cnb_label) DO UPDATE
SET plane_label_id = EXCLUDED.plane_label_id, updated_at = now();

-- 4. （可选）创建 Issue 链接（如 Issue #36 已存在于 Plane）
-- INSERT INTO issue_links (plane_issue_id, cnb_repo_id, cnb_issue_id, plane_project_id, created_at)
-- VALUES ('<PLANE_ISSUE_UUID>', '1024hub/Demo/BE-test-issue', '36', '<PROJECT_UUID>', now())
-- ON CONFLICT (cnb_repo_id, cnb_issue_id) DO NOTHING;

-- 5. 验证
SELECT cnb_repo_id, plane_project_id::text, workspace_slug, active
FROM repo_project_mappings 
WHERE cnb_repo_id = '1024hub/Demo/BE-test-issue';
