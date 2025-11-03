# 父工作项同步问题完整修复总结

## 问题现象

在 Plane 创建带父工作项的子工作项时，同步到 CNB 后标题显示：
```
[Parent: unknown ] xxx
```

## 根本原因

**三层问题**：

1. **Plane Webhook 缺失数据**：webhook payload 中不包含 `workspace_slug` 字段
2. **后端 API 未处理**：管理接口虽然前端发送了 `workspace_slug`，但后端未接收和保存
3. **数据库配置缺失**：`repo_project_mappings` 表中 `workspace_slug` 列为空

导致系统无法调用 Plane API 获取父工作项信息，最终显示 `unknown`。

## 完整解决方案

### 代码修复（已完成）

#### 1. 增强 workspace_slug 获取策略
**文件**: `internal/handlers/plane_issue.go`

实现多层回退机制：
```
webhook 数据 → mapping 配置 → 父 issue 快照 → 项目快照
```

关键修改：
- 创建 issue 时优先使用 mapping 中的 workspace_slug
- `syncParentIssueIfNeeded` 增加快照回退逻辑
- 添加详细日志追踪来源

#### 2. 新增数据库查询方法
**文件**: `internal/store/repositories.go`

```go
// 从项目快照获取
func (d *DB) GetWorkspaceSlugByProjectID(ctx, planeProjectID) (string, error)

// 从 workspaces 表获取
func (d *DB) GetWorkspaceSlugByWorkspaceID(ctx, planeWorkspaceID) (string, error)
```

#### 3. 修复后端 API
**文件**: `internal/handlers/admin.go`

- 请求结构体添加 `WorkspaceSlug` 字段
- 创建 mapping 时保存 workspace_slug 到数据库

#### 4. 修复数据库操作
**文件**: `internal/store/repositories.go`

更新 `UpsertRepoProjectMapping`：
- INSERT SQL 包含 workspace_slug 列
- UPDATE SQL 更新 workspace_slug 字段

### 前端支持（已完成）

**文件**: `web/app/(dashboard)/mappings/page.tsx`

管理界面已支持：
- ✅ Workspace Slug 输入框
- ✅ 显示 workspace_slug 信息
- ✅ 表单提交时发送 workspace_slug

## 配置步骤（用户操作）

### 方式 1：通过管理界面（推荐）

1. 访问 `http://your-domain/admin/mappings`
2. 在创建/编辑 mapping 表单中：
   - **Workspace Slug**: 输入 Plane workspace slug（从 URL 获取）
   - 其他字段正常填写
3. 保存后自动生效

### 方式 2：通过 SQL 直接更新

```bash
# 1. 获取 workspace slug（从 Plane URL）
# https://app.plane.so/[workspace-slug]/projects

# 2. 更新数据库
psql "$DATABASE_URL" -c "
UPDATE repo_project_mappings 
SET workspace_slug = 'your-workspace-slug' 
WHERE cnb_repo_id = '1024hub/plane-test';

SELECT cnb_repo_id, workspace_slug FROM repo_project_mappings;
"
```

## 验证方法

### 1. 检查数据库配置

```bash
psql "$DATABASE_URL" -c "
SELECT cnb_repo_id, workspace_slug, active 
FROM repo_project_mappings;
"
```

确保 `workspace_slug` 列不为空。

### 2. 创建测试 Issue

在 Plane 创建一个带父工作项的子工作项。

### 3. 查看日志（成功标志）

```json
✅ 成功
{
  "event": "plane.issue.workspace_slug_resolution",
  "workspace_slug": "your-workspace-slug",
  "workspace_slug_source": "mapping"  // 或 "parent_snapshot" / "project_snapshot"
}

{
  "event": "plane.parent_issue.synced",
  "parent_plane_id": "...",
  "cnb_issue_id": "123",
  "parent_name": "父工作项标题"
}
```

```json
❌ 失败
{
  "event": "plane.parent_issue.skipped",
  "reason": "workspace_slug_unavailable_after_all_fallbacks"
}
```

### 4. 检查 CNB Issue 标题

应显示：
```
[Parent: #123 ] 子工作项标题
```

而不是：
```
[Parent: unknown ] 子工作项标题
```

## 技术细节

### 多层回退策略

系统按以下顺序尝试获取 workspace_slug：

1. **Webhook 数据**（通常为空）
   ```go
   workspaceSlug, _ := dataGetString(data, "workspace_slug")
   ```

2. **Mapping 配置**（用户手动配置）✅ 主要来源
   ```go
   if effectiveWorkspaceSlug == "" && m.WorkspaceSlug.Valid {
       effectiveWorkspaceSlug = m.WorkspaceSlug.String
   }
   ```

3. **父 Issue 快照**（如果父 issue 之前同步过）
   ```go
   snapshot, _ := h.db.GetPlaneIssueSnapshot(ctx, parentPlaneID)
   workspaceSlug = snapshot["workspace_slug"]
   ```

4. **项目快照**（如果项目有其他 issue 同步过）
   ```go
   workspaceSlug, _ := h.db.GetWorkspaceSlugByProjectID(ctx, planeProjectID)
   ```

### 自动填充机制

配置 mapping 后，系统会在每次 webhook 事件中：

1. 使用 mapping 的 workspace_slug
2. 更新 `plane_issue_snapshots` 表
3. 更新 `plane_project_snapshots` 表

后续即使 webhook 缺失数据，也能从快照获取，形成**良性循环**。

## 相关文档

- **快速修复**：`docs/QUICKFIX_PARENT_ISSUE.md`
- **详细配置**：`docs/WORKSPACE_SLUG_SETUP.md`
- **SQL 脚本**：`scripts/update_workspace_slug.sql`

## 提交历史

- `e6fa56c` - fix(admin): 后端 API 支持 workspace_slug 字段的接收和保存
- `bbb814f` - docs: 添加父工作项同步问题快速修复指南
- `96b373b` - docs: 添加 workspace_slug 配置指南和 SQL 脚本
- `9c95002` - fix(plane): 增强父工作项同步的 workspace_slug 获取策略

## FAQ

### Q: 为什么 Plane webhook 中没有 workspace_slug？
A: Plane 的 webhook payload 设计问题，只包含 workspace_id（UUID），不包含 slug。

### Q: 配置后还是失败怎么办？
A: 检查：
1. 数据库中 workspace_slug 是否正确保存
2. Plane URL 中的 slug 是否准确
3. `PLANE_SERVICE_TOKEN` 是否有权限访问该 workspace

### Q: 跨 workspace 的父子关系支持吗？
A: 不支持。父子工作项必须在同一个 workspace 中。

### Q: 需要重启服务吗？
A: 不需要。配置保存到数据库后立即生效，下次 webhook 事件会自动使用。

### Q: 已有的子 issue 会自动修复吗？
A: 不会。只影响新创建的 issue。如需修复已有 issue，需要手动更新 CNB issue 标题。
