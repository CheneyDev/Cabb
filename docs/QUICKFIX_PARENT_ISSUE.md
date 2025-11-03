# 父工作项同步问题快速修复

## 现象

在 Plane 创建带父工作项的子工作项时，同步到 CNB 后标题显示：`[Parent: unknown ] xxx`

## 原因

数据库中未配置 `workspace_slug`，导致系统无法调用 Plane API 获取父工作项信息。

## 快速修复（3 步）

### 1. 获取 Workspace Slug

访问你的 Plane 工作区，从浏览器地址栏复制 workspace slug：

```
https://app.plane.so/[这里就是workspace-slug]/projects
```

例如：
- URL: `https://app.plane.so/acme-corp/projects`  
- **Workspace slug**: `acme-corp` ✅

### 2. 更新数据库

方式 A - 使用提供的脚本：

```bash
# 编辑脚本，替换 'your-workspace-slug' 为实际值
nano scripts/update_workspace_slug.sql

# 执行更新
psql "$DATABASE_URL" -f scripts/update_workspace_slug.sql
```

方式 B - 直接执行 SQL（推荐）：

```bash
# 替换 YOUR_ACTUAL_SLUG 为第 1 步获取的值
psql "$DATABASE_URL" << EOF
UPDATE repo_project_mappings 
SET workspace_slug = 'YOUR_ACTUAL_SLUG' 
WHERE cnb_repo_id = '1024hub/plane-test';

-- 验证更新
SELECT cnb_repo_id, workspace_slug FROM repo_project_mappings;
EOF
```

### 3. 验证

在 Plane 创建一个新的带父工作项的子工作项，查看日志：

**✅ 成功标志**：
```json
{
  "event": "plane.issue.workspace_slug_resolution",
  "workspace_slug": "your-actual-slug",
  "workspace_slug_source": "mapping"
}

{
  "event": "plane.parent_issue.synced",
  "cnb_issue_id": "123"
}
```

**❌ 仍然失败**：
```json
{
  "workspace_slug": "",
  "workspace_slug_source": "webhook"
}
```
→ 检查第 2 步的 SQL 是否正确执行

## 示例

假设你的 Plane URL 是 `https://app.plane.so/my-team/projects`：

```bash
# 直接执行（替换 my-team 为实际值）
psql "$DATABASE_URL" -c "
UPDATE repo_project_mappings 
SET workspace_slug = 'my-team' 
WHERE cnb_repo_id = '1024hub/plane-test';

SELECT cnb_repo_id, workspace_slug FROM repo_project_mappings;
"
```

## 多仓库配置

如果有多个 CNB 仓库：

```sql
-- 查看所有 mapping
SELECT plane_project_id, cnb_repo_id, workspace_slug 
FROM repo_project_mappings;

-- 批量更新（如果都在同一个 workspace）
UPDATE repo_project_mappings 
SET workspace_slug = 'your-workspace-slug' 
WHERE workspace_slug IS NULL OR workspace_slug = '';

-- 或分别更新不同 workspace
UPDATE repo_project_mappings 
SET workspace_slug = 'workspace-a' 
WHERE cnb_repo_id = '1024hub/repo-a';

UPDATE repo_project_mappings 
SET workspace_slug = 'workspace-b' 
WHERE cnb_repo_id = '1024hub/repo-b';
```

## 后续优化

配置后，系统会：
1. 在每次 webhook 事件中自动更新快照表
2. 后续即使 webhook 数据缺失也能从快照获取
3. 父工作项信息会被缓存，无需重复调用 API

## 需要帮助？

如果配置后仍有问题，提供以下信息：

```bash
# 检查数据库配置
psql "$DATABASE_URL" -c "
SELECT plane_project_id, cnb_repo_id, workspace_slug 
FROM repo_project_mappings 
WHERE cnb_repo_id = '1024hub/plane-test';
"

# 检查最近的日志
grep "plane.parent_issue" logs/app.log | tail -20
```
