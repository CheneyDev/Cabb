# 标签同步映射缺失 - 修复指南

## 问题现象

```
API: POST /api/v1/issues/label-sync
响应: {"code": 0, "message": "success"}  ← 200 OK
日志: {"error": "mapping_not_found"}     ← 实际失败
```

**原因：** 数据库缺少 CNB 仓库到 Plane 项目的映射（`repo_project_mappings` 与 `label_mappings`）。

## 修复步骤（5 分钟）

### 步骤 1：打开项目页面

访问：https://work.1024hub.org:4430/test/projects

### 步骤 2：打开开发者工具

- **Chrome/Edge**: `F12` 或 `Ctrl+Shift+I`
- **Firefox**: `F12`
- **Safari**: `Cmd+Option+I`

### 步骤 3：查看 Network 请求获取项目 UUID

1. 切换到 **Network** 标签（网络）
2. 刷新页面（`F5` 或 `Ctrl+R`）
3. 在过滤器中输入：`projects`
4. 找到 `/api/users/me/workspaces/` 或 `/api/workspaces/test/projects/` 请求
5. 点击请求 → **Response** 标签，复制 JSON 中的 UUID：

```json
{
  "results": [
    {
      "id": "44848399-cae8-4ce6-b325-5bd913e7e1cb",      ← PROJECT_UUID
      "workspace": "4ada216e-373d-4029-ad4a-dbdadaf8f1fe", ← WORKSPACE_UUID
      "name": "test",
      "identifier": "TEST"
    }
  ]
}
```

### 步骤 4：获取标签 UUID

1. 访问项目设置 → 标签（Labels）页面
2. Network 标签中找到 `/labels/` 请求
3. 查看响应 JSON，复制标签 `id` 字段

### 步骤 5：执行 SQL 插入映射

```bash
psql "$DATABASE_URL" << 'EOF'
-- 1. 插入项目映射（替换 UUID 为步骤 3 获取的值）
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
  '44848399-cae8-4ce6-b325-5bd913e7e1cb',  -- PROJECT_UUID
  '4ada216e-373d-4029-ad4a-dbdadaf8f1fe',  -- WORKSPACE_UUID
  '1024hub/Demo/BE-test-issue',
  'my-test',
  true,
  'cnb_to_plane',
  now(),
  now()
);

-- 2. 插入标签映射（替换 UUID 为步骤 4 获取的值）
INSERT INTO label_mappings (plane_project_id, cnb_repo_id, cnb_label, plane_label_id)
VALUES 
  ('44848399-cae8-4ce6-b325-5bd913e7e1cb', '1024hub/Demo/BE-test-issue', '🚧 处理中_CNB', '<LABEL_UUID_1>'),
  ('44848399-cae8-4ce6-b325-5bd913e7e1cb', '1024hub/Demo/BE-test-issue', '🧑🏻‍💻 进行中：后端_CNB', '<LABEL_UUID_2>');

-- 3. 验证映射
SELECT cnb_repo_id, plane_project_id::text, workspace_slug FROM repo_project_mappings;
SELECT cnb_label, plane_label_id::text FROM label_mappings WHERE cnb_repo_id = '1024hub/Demo/BE-test-issue';
EOF
```

### 步骤 6：验证标签同步

```bash
curl -X POST "https://hub.1024hub.org:8081/api/v1/issues/label-sync" \
  -H "Authorization: Bearer $INTEGRATION_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"repo_slug": "1024hub/Demo/BE-test-issue", "issue_number": 36, "labels": ["🚧 处理中_CNB"]}'
```

**预期结果：** 日志无 `"error"` 字段，Plane Issue 标签已更新。

## 常见问题

**Q: 找不到 `/projects/` 请求？**  
清空 Network 标签，刷新页面。搜索包含 `workspace` 或 `project` 的请求。

**Q: 响应是 HTML 而非 JSON？**  
确认请求 URL 以 `/api/` 开头，查看 **Response** 标签（非 Preview）。

**Q: UUID 格式是什么？**  
36 字符：`xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

**Q: 为什么 API 返回 200 但失败？**  
异步处理设计，失败仅记录日志不影响 HTTP 响应。需检查服务端日志确认实际结果。

**Q: 标签映射失败怎么办？**  
检查 `plane_label_id` 是否正确，确认标签在 Plane 项目中存在。使用浏览器 Network 工具获取准确 UUID。
