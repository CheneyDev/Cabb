# 标签同步映射缺失 - 修复指南

## 问题现象

```
API: POST /api/v1/issues/label-sync
响应: {"code": 0, "message": "success"}  ← 200 OK
日志: {"error": "mapping_not_found"}     ← 实际失败
```

**原因：** 数据库缺少 CNB 仓库到 Plane 项目的映射。

## 修复步骤（5 分钟）

### 1. 获取 Plane Project UUID

**方法 A：自动获取（推荐）**

```bash
# 需要先配置 PLANE_SERVICE_TOKEN（见下方说明）
./scripts/get_plane_uuids.sh
```

**方法 B：从浏览器手动查找**

1. 访问 Plane 项目页面：`https://work.1024hub.org:4430/my-test/projects/test-notify/...`
2. 打开开发者工具（F12）→ Network 标签
3. 刷新页面，找到 `/api/workspaces/my-test/projects/` 请求
4. 查看响应 JSON 中的 `id` 字段（项目 UUID）

**方法 C：使用 curl 手动查询**

```bash
# 需要 PLANE_SERVICE_TOKEN
curl -H "X-API-Key: $PLANE_SERVICE_TOKEN" \
  "https://work.1024hub.org:4430/api/workspaces/my-test/projects/"
```

### 2. 配置 PLANE_SERVICE_TOKEN（如未配置）

1. 访问 Plane：https://work.1024hub.org:4430
2. 进入 **个人设置 → API Tokens**
3. 点击 **"Create Token"** 或 **"新建令牌"**
4. 权限至少选择：`project:read`, `issue:write`, `label:read`
5. 复制生成的 Token，添加到 `.env`：
   ```bash
   PLANE_SERVICE_TOKEN=plane_api_xxxxxxxxxxxxx
   ```

### 3. 执行 SQL 插入映射

获取到 UUID 后，执行以下 SQL（替换占位符）：

```bash
psql "$DATABASE_URL" << 'EOF'
-- 替换 <PROJECT_UUID>、<WORKSPACE_UUID>、<WORKSPACE_SLUG>
INSERT INTO repo_project_mappings (
  plane_project_id, plane_workspace_id, cnb_repo_id, 
  workspace_slug, active, sync_direction, created_at, updated_at
) VALUES (
  '<PROJECT_UUID>',              -- 从步骤 1 获取
  '<WORKSPACE_UUID>',            -- 从步骤 1 获取
  '1024hub/Demo/BE-test-issue',
  'my-test',
  true, 'cnb_to_plane', now(), now()
)
ON CONFLICT (plane_project_id, cnb_repo_id) DO UPDATE 
SET active = true, updated_at = now();

-- 验证
SELECT cnb_repo_id, plane_project_id::text, workspace_slug 
FROM repo_project_mappings 
WHERE cnb_repo_id = '1024hub/Demo/BE-test-issue';
EOF
```

### 4. 创建标签映射

```bash
# 查询 Plane 项目中的标签
curl -H "X-API-Key: $PLANE_SERVICE_TOKEN" \
  "https://work.1024hub.org:4430/api/workspaces/my-test/projects/<project_id>/labels/" \
  | jq '.[] | {name, id}'
```

获取标签 UUID 后，插入映射：

```sql
INSERT INTO label_mappings (plane_project_id, cnb_repo_id, cnb_label, plane_label_id)
VALUES 
  ('<PROJECT_UUID>', '1024hub/Demo/BE-test-issue', '🚧 处理中_CNB', '<LABEL_UUID_1>'),
  ('<PROJECT_UUID>', '1024hub/Demo/BE-test-issue', '🧑🏻‍💻 进行中：后端_CNB', '<LABEL_UUID_2>');
```

### 5. 验证修复

```bash
curl -X POST "https://hub.1024hub.org:8081/api/v1/issues/label-sync" \
  -H "Authorization: Bearer $INTEGRATION_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"repo_slug": "1024hub/Demo/BE-test-issue", "issue_number": 36, "labels": ["🚧 处理中_CNB"]}'
```

日志应无 `"error"` 字段。

## 常见问题

**Q: 如何获取 Plane UUID？**  
使用 `./scripts/get_plane_uuids.sh` 自动获取，或从浏览器开发者工具查看 API 响应。

**Q: PLANE_SERVICE_TOKEN 在哪里获取？**  
Plane 个人设置 → API Tokens → 创建新 Token（需要 `project:read` 权限）。

**Q: 为什么 200 但失败？**  
异步处理设计，API 立即返回 200，实际处理在后台，失败仅记录日志。

**Q: 没有 jq 工具怎么办？**  
安装：`apt install jq` 或 `brew install jq`，或手动从 JSON 响应中提取 `id` 字段。
