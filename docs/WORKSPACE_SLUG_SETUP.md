# Workspace Slug 配置指南

## 问题说明

当 Plane 创建带父工作项的子工作项时，如果 `workspace_slug` 未配置，父工作项无法同步到 CNB，导致子工作项标题显示为 `[Parent: unknown ]`。

## 获取 Workspace Slug

Workspace slug 是 Plane 工作区的唯一标识符，可以从 Plane URL 中获取：

```
https://app.plane.so/{workspace-slug}/projects
                      ^^^^^^^^^^^^^^^^
                      这就是 workspace slug
```

例如：
- URL: `https://app.plane.so/acme-corp/projects`
- Workspace slug: `acme-corp`

## 配置方法

### 方法 1：通过 SQL 直接更新（推荐）

```bash
# 使用提供的脚本
psql "$DATABASE_URL" -f scripts/update_workspace_slug.sql
```

或手动执行：

```sql
-- 1. 查看当前配置
SELECT plane_project_id, cnb_repo_id, workspace_slug 
FROM repo_project_mappings;

-- 2. 更新 workspace_slug（替换为实际值）
UPDATE repo_project_mappings 
SET workspace_slug = 'your-actual-workspace-slug'
WHERE cnb_repo_id = '1024hub/plane-test';

-- 3. 确认更新成功
SELECT plane_project_id, cnb_repo_id, workspace_slug 
FROM repo_project_mappings;
```

### 方法 2：通过管理 API（即将支持）

```bash
curl -X PATCH "http://localhost:8080/admin/mappings/repo-project/{mapping_id}" \
  -H "Content-Type: application/json" \
  -d '{"workspace_slug": "your-workspace-slug"}'
```

## 验证配置

配置完成后，创建一个新的带父工作项的子工作项，日志应显示：

```json
{
  "event": "plane.issue.workspace_slug_resolution",
  "workspace_slug": "your-workspace-slug",
  "workspace_slug_source": "mapping"
}
```

父工作项应成功同步：

```json
{
  "event": "plane.parent_issue.synced",
  "parent_plane_id": "...",
  "cnb_issue_id": "..."
}
```

## 注意事项

1. **每个 mapping 必须配置**：如果有多个 CNB 仓库映射到不同的 Plane 项目，需要为每个 mapping 配置对应的 workspace_slug
2. **后续自动填充**：配置后，系统会在后续 webhook 事件中自动更新快照表，之后即使 webhook 中缺失 workspace_slug 也能正常工作
3. **跨工作区限制**：父子工作项必须在同一个 workspace 中，系统不支持跨 workspace 的父子关系

## 故障排查

如果配置后仍然失败，检查：

1. **确认 workspace_slug 正确**：登录 Plane 确认 URL 中的值
2. **检查数据库连接**：确保应用能访问数据库
3. **查看日志**：搜索 `plane.parent_issue.skipped` 确认失败原因
4. **验证 Plane API Token**：确保 `PLANE_SERVICE_TOKEN` 有权限访问工作区

## 技术细节

系统使用以下优先级获取 workspace_slug：

1. **Webhook 数据**（通常为空）
2. **Mapping 配置**（需要手动配置）✅
3. **父 issue 快照**（如果父 issue 之前同步过）
4. **项目快照**（如果项目有其他 issue 同步过）

只要配置了第 2 步，后续第 3、4 步会自动填充，形成良性循环。
