# Issue 标签同步业务逻辑说明

## 📋 功能概述

提供两个 API 端点用于同步 CNB Issue 标签到 Plane，并可选发送飞书通知。

## 🔄 API 对比

### 完整版 API
- **端点**：`POST /api/v1/issues/label-notify`
- **字段数**：11 个（包含完整事件信息）
- **适用场景**：需要记录完整事件上下文、审计日志
- **测试脚本**：`./scripts/test-label-notify.sh`

### 简化版 API  
- **端点**：`POST /api/v1/issues/label-sync`
- **字段数**：3 个（只需核心字段）
- **适用场景**：快速标签同步，减少传输数据
- **测试脚本**：`./scripts/test-label-sync.sh`

## ✅ 已实现的业务逻辑

### 1. 标签过滤（`_CNB` 后缀）

**目的**：只处理 CNB 平台管理的标签，避免误覆盖其他系统的标签。

**逻辑**：
```go
func filterCNBLabels(labels []string) []string {
    var cnbLabels []string
    for _, label := range labels {
        if strings.HasSuffix(label, "_CNB") {
            cnbLabels = append(cnbLabels, label)
        }
    }
    return cnbLabels
}
```

**示例**：
- 输入：`["🚧 处理中_CNB", "bug", "🧑🏻‍💻 进行中：前端_CNB", "feature"]`
- 输出：`["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"]`

### 2. 查询映射关系

**步骤 1**：根据 `repo_slug` 查询 repo-project 映射
```sql
SELECT plane_project_id, plane_workspace_id, cnb_repo_id
FROM repo_project_mappings
WHERE cnb_repo_id = $1 AND active = true
LIMIT 1
```

**步骤 2**：根据 `repo_slug` + `issue_number` 查询对应的 Plane Issue
```sql
SELECT plane_issue_id FROM issue_links
WHERE cnb_repo_id = $1 AND cnb_issue_id = $2
LIMIT 1
```

**失败处理**：
- 映射不存在 → 记录日志，跳过处理
- Plane Issue 不存在 → 记录日志，跳过处理

### 3. 标签映射

**步骤**：将 CNB 标签名映射为 Plane Label UUID
```sql
SELECT plane_label_id::text FROM label_mappings
WHERE plane_project_id = $1
  AND cnb_repo_id = $2
  AND cnb_label_name = ANY($3)
```

**示例**：
- CNB 标签：`["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"]`
- Plane Label IDs：`["uuid1", "uuid2"]`

**注意**：如果没有映射配置，会跳过该标签（不会报错）。

### 4. 更新 Plane Issue

**API 调用**：
```
PATCH /api/v1/workspaces/{workspace_slug}/projects/{project_id}/issues/{issue_id}/
Authorization: Bearer <bot_token>
Content-Type: application/json

{"labels": ["uuid1", "uuid2"]}
```

**行为**：
- ✅ **覆盖式更新**：替换 Issue 的所有标签为新标签列表
- ❌ **不是增量更新**：原有标签会被替换

**注意事项**：
1. 只更新 `_CNB` 后缀的标签
2. 如果 Plane Issue 有其他系统管理的标签，需要在映射表中维护
3. 失败时记录日志，不影响其他 Issue 的处理

### 5. 飞书通知（可选）

**触发条件**：
1. 配置了 `LARK_APP_ID` 和 `LARK_APP_SECRET`
2. 存在 channel-project 映射（`channel_project_mappings` 表）

**查询映射**：
```sql
SELECT lark_chat_id, notify_on_create
FROM channel_project_mappings
WHERE plane_project_id = $1
```

**通知内容**：
```
📋 Issue 标签更新

仓库：1024hub/Demo
Issue：#74 - 实现用户登录功能
状态：open
标签：🚧 处理中_CNB, 🧑🏻‍💻 进行中：前端_CNB
触发标签：🚧 处理中_CNB

🔗 查看详情：https://cnb.cool/1024hub/Demo/-/issues/74
```

**发送方式**：
- 调用飞书 API `POST /open-apis/im/v1/messages`
- 向所有绑定的群组发送文本消息
- 失败时记录日志，不阻塞主流程

## 🔍 完整处理流程

```
┌─────────────────────────────────────────────────────────────┐
│ 1. 接收请求并校验                                            │
│    - Bearer token 鉴权                                       │
│    - JSON 解析与字段校验                                     │
│    - 内存/数据库去重（幂等性保证）                            │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. 异步处理（goroutine）                                     │
│    - 立即返回 200 OK，不阻塞 CNB job                          │
│    - 超时设置：30 秒                                          │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. 标签过滤                                                  │
│    - 提取 _CNB 后缀的标签                                     │
│    - 如无 CNB 标签 → 跳过处理                                 │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. 查询映射关系                                              │
│    - repo_slug → Plane Project                              │
│    - issue_number → Plane Issue                             │
│    - CNB Labels → Plane Label IDs                            │
│    - 任一映射失败 → 记录日志，跳过处理                         │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. 更新 Plane Issue                                          │
│    - 调用 Plane API PATCH /issues/{id}/                      │
│    - 覆盖式更新标签列表                                       │
│    - 失败 → 记录日志，不重试                                  │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. 飞书通知（可选）                                           │
│    - 查询 channel-project 映射                               │
│    - 构建通知消息                                             │
│    - 向所有绑定群组发送                                       │
│    - 失败不影响主流程                                         │
└─────────────────────────────────────────────────────────────┘
```

## 📊 数据库依赖

### 必需表
1. `repo_project_mappings` - CNB 仓库 ↔ Plane 项目
2. `issue_links` - CNB Issue ↔ Plane Issue
3. `label_mappings` - CNB 标签 ↔ Plane Label ID
4. `workspaces` - Plane Bot Token

### 可选表
1. `channel_project_mappings` - 飞书群组 ↔ Plane 项目（用于通知）
2. `event_deliveries` - 事件去重记录（用于幂等性）

## ⚠️ 注意事项

### 1. 标签覆盖策略

**当前行为**：覆盖式更新（替换所有标签）

**潜在问题**：
- 如果 Plane Issue 有非 CNB 管理的标签，会被清除
- 例如：Plane 中手动添加的 `priority:high` 标签会丢失

**解决方案**：
1. **短期**：在 `label_mappings` 中维护所有需要保留的标签
2. **长期**：实现增量更新（先读取现有标签，合并后更新）

### 2. 并发处理

**当前实现**：
- 每个请求启动一个 goroutine 异步处理
- 无并发控制（如无限 goroutine）

**建议改进**：
- 使用 worker pool 限制并发数（如 100 个 worker）
- 使用消息队列（如 Redis/RabbitMQ）解耦

### 3. 错误重试

**当前实现**：
- 失败时记录日志，不重试
- 依赖 CNB job 的重试机制

**建议改进**：
- 对 429/5xx 错误实现指数退避重试
- 将失败事件写入重试队列

### 4. 性能优化

**当前实现**：
- 每个 Issue 独立处理（多次数据库查询 + API 调用）

**批量处理建议**：
- 接收批量请求（如 100 个 Issue）
- 批量查询数据库（减少 round trip）
- 批量调用 Plane API（如支持）

## 🧪 测试建议

### 单元测试
```go
func TestFilterCNBLabels(t *testing.T) {
    input := []string{"bug_CNB", "feature", "todo_CNB"}
    output := filterCNBLabels(input)
    assert.Equal(t, []string{"bug_CNB", "todo_CNB"}, output)
}
```

### 集成测试
1. 创建测试 repo-project 映射
2. 创建测试 issue-link
3. 创建测试 label 映射
4. 调用 API 并验证 Plane Issue 标签是否更新

### 端到端测试
1. 配置真实的 CNB/Plane/飞书环境
2. 触发 CNB job
3. 验证 Plane 标签更新
4. 验证飞书通知发送

## 📝 未来扩展

1. **增量标签更新**：先读取现有标签，只更新 CNB 管理的部分
2. **标签同步方向**：支持 Plane → CNB 的反向同步
3. **冲突检测**：检测标签冲突并报警
4. **审计日志**：记录每次标签变更的详细历史
5. **批量处理**：支持一次请求同步多个 Issue

## 🔗 相关文档

- API 文档：`docs/cnb-job-integration.md`
- 测试示例：`docs/api-examples.md`
- 快速开始：`docs/QUICKSTART-LABEL-NOTIFY.md`
