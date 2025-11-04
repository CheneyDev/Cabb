# Issue 标签同步 - 开发者指南

## 📋 概述

本指南面向 Plane-Integration 开发者，介绍架构设计、业务逻辑与实现细节。

---

## 💡 架构特性

### 分层架构

**Handler 层**（`internal/handlers/issue_label_notify.go`）：
- HTTP 请求处理、参数校验、鉴权、响应封装
- 依赖：`config.Config`、`store.DB`、`Deduper`
- 独立性：不依赖具体业务逻辑（Plane/Lark 客户端）

**业务处理层**（`processIssueLabelNotify`）：
- 异步执行，不阻塞 HTTP 响应
- 预留扩展点：同步 Plane、发送飞书通知、触发工作流

**数据层**（`internal/store`）：
- 统一数据访问接口
- 支持可选数据库（无 DB 时降级为内存去重）

### 统一错误处理

- 使用 `writeError` 统一错误响应格式
- 结构化日志记录（包含 `request_id`、`source`、`endpoint`）
- 错误码机器可读（`invalid_token`、`missing_fields` 等）

### 幂等性与安全性

详见 [API-REFERENCE.md](./API-REFERENCE.md#-安全性) 的安全性章节。

---

## 📊 业务逻辑详解

### 1. 标签过滤（_CNB 后缀）

**目的：** 只处理 CNB 平台管理的标签，避免误覆盖其他系统的标签。

**实现：**
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

**示例：**
- 输入：`["🚧 处理中_CNB", "bug", "🧑🏻‍💻 进行中：前端_CNB", "feature"]`
- 输出：`["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"]`

---

### 2. 查询映射关系

**步骤 1：** 根据 `repo_slug` 查询 repo-project 映射
```sql
SELECT plane_project_id, plane_workspace_id, cnb_repo_id
FROM repo_project_mappings
WHERE cnb_repo_id = $1 AND active = true
LIMIT 1
```

**步骤 2：** 根据 `repo_slug` + `issue_number` 查询对应的 Plane Issue
```sql
SELECT plane_issue_id FROM issue_links
WHERE cnb_repo_id = $1 AND cnb_issue_id = $2
LIMIT 1
```

**失败处理：**
- 映射不存在 → 记录日志，跳过处理
- Plane Issue 不存在 → 记录日志，跳过处理

---

### 3. 标签映射

将 CNB 标签名映射为 Plane Label UUID：
```sql
SELECT plane_label_id::text FROM label_mappings
WHERE plane_project_id = $1
  AND cnb_repo_id = $2
  AND cnb_label = ANY($3)
```

**示例：**
- CNB 标签：`["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"]`
- Plane Label IDs：`["uuid1", "uuid2"]`

---

### 4. 增量更新 Plane Issue 标签

**逻辑流程：**
```
1. 从 Plane API 读取 Issue 当前所有标签
   ↓
2. 从 label_mappings 表查询哪些标签是 CNB 管理的
   ↓
3. 过滤出非 CNB 管理的标签（需要保留）
   ↓
4. 合并：保留的标签 + 新的 CNB 标签
   ↓
5. 去重后更新到 Plane
```

**示例：**
- Plane Issue 当前标签：`["priority:high", "🚧 处理中_CNB", "bug"]`
- CNB 管理的标签（在 label_mappings）：`["🚧 处理中_CNB"]`
- CNB 新标签：`["✅ 已完成_CNB"]`
- **最终结果：** `["priority:high", "bug", "✅ 已完成_CNB"]`
  - ✅ 保留 `priority:high` 和 `bug`（非 CNB 管理）
  - ✅ 替换 `🚧 处理中_CNB` 为 `✅ 已完成_CNB`（CNB 管理）

**优点：**
- ✅ 不会误删非 CNB 管理的标签
- ✅ 只更新 CNB 负责的标签
- ✅ 其他系统（如 Plane 手动添加）的标签不受影响

---

### 5. 飞书通知（可选）

**触发条件：**
1. 配置了 `LARK_APP_ID` 和 `LARK_APP_SECRET`
2. 存在 channel-project 映射（`channel_project_mappings` 表）

**查询映射：**
```sql
SELECT lark_chat_id, notify_on_create
FROM channel_project_mappings
WHERE plane_project_id = $1
```

**通知内容：**
```
📋 Issue 标签更新

仓库：1024hub/Demo
Issue：#74 - 实现用户登录功能
状态：open
标签：🚧 处理中_CNB, 🧑🏻‍💻 进行中：前端_CNB
触发标签：🚧 处理中_CNB

🔗 查看详情：https://cnb.cool/1024hub/Demo/-/issues/74
```

**发送方式：**
- 调用飞书 API `POST /open-apis/im/v1/messages`
- 向所有绑定的群组发送文本消息
- 失败时记录日志，不阻塞主流程

---

## 📂 数据库依赖

### 必需表

1. **repo_project_mappings** - CNB 仓库 ↔ Plane 项目
2. **issue_links** - CNB Issue ↔ Plane Issue
3. **label_mappings** - CNB 标签 ↔ Plane Label ID
4. **workspaces** - Plane Bot Token

### 可选表

1. **channel_project_mappings** - 飞书群组 ↔ Plane 项目（用于通知）
2. **event_deliveries** - 事件去重记录（用于幂等性）

---

## ⚠️ 注意事项

### 1. 标签更新策略（增量更新）

**核心原则：** ✅ **只更新 CNB 管理的标签，保留其他系统的标签**

#### 工作原理

```
1. 从 Plane API 获取 Issue 当前所有标签 ID
   ↓
2. 从 label_mappings 表查询哪些标签 ID 是 CNB 管理的
   ↓
3. 过滤出非 CNB 管理的标签 ID（需保留）
   ↓
4. 合并：保留的标签 + 新的 CNB 标签
   ↓
5. 去重后更新到 Plane
```

**示例：**
```
Plane 当前标签：[priority:high, 🚧处理中_CNB, bug, feature]
CNB 管理的标签（在 label_mappings）：[🚧处理中_CNB]
CNB 新标签：[✅已完成_CNB]

执行流程：
1. 识别非 CNB 管理的标签：[priority:high, bug, feature]
2. 合并：[priority:high, bug, feature] + [✅已完成_CNB]
3. 最终结果：[priority:high, bug, feature, ✅已完成_CNB]

✅ priority:high、bug、feature 完全保留
✅ 🚧处理中_CNB 被替换为 ✅已完成_CNB
```

#### 安全性保证

**不会覆盖非 CNB 管理的标签：**
- ✅ 只有在 `label_mappings` 中配置的标签 ID 才会被识别为"CNB 管理"
- ✅ 未配置的标签会被完整保留，不受 CNB 同步影响
- ✅ Plane 手动添加的标签、其他系统添加的标签都不会被删除

**代码实现（关键逻辑）：**
```go
// 获取 CNB 管理的标签 ID 集合
cnbManagedIDs := h.db.GetCNBManagedLabelIDs(ctx, projectID, repoSlug)

// 保留非 CNB 管理的标签
for _, labelID := range currentLabelIDs {
    if !cnbManagedIDs[labelID] {  // 不在 label_mappings 中
        preservedLabelIDs = append(preservedLabelIDs, labelID)
    }
}

// 合并：保留的 + 新的 CNB 标签
finalLabelIDs := append(preservedLabelIDs, planeLabelIDs...)
```

#### 配置要求


```sql
CREATE TABLE IF NOT EXISTS label_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  cnb_label text NOT NULL,
  plane_label_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, cnb_repo_id, cnb_label)
);
```

**配置示例：**
```sql
-- 配置 CNB 标签到 Plane 标签 ID 的映射
INSERT INTO label_mappings (plane_project_id, cnb_repo_id, cnb_label, plane_label_id)
VALUES 
  ('project-uuid', '1024hub/Demo', '🚧 处理中_CNB', 'plane-label-uuid-1'),
  ('project-uuid', '1024hub/Demo', '✅ 已完成_CNB', 'plane-label-uuid-2'),
  ('project-uuid', '1024hub/Demo', '🐛 Bug_CNB', 'plane-label-uuid-3');
```

#### 优点与限制

**优点：**
- ✅ 精确控制哪些标签由 CNB 管理
- ✅ 不会误删其他系统的标签
- ✅ 支持标签重命名（修改映射表即可）
- ✅ 可审计标签映射历史

**限制：**
- ⚠️ 首次使用时会自动创建标签（默认颜色 #808080）
- ⚠️ 标签名称完全匹配或去 `_CNB` 后缀匹配

---

### 2. 标签自动创建（2025-01-04 新增）

**设计改进：** 不再需要预先手动配置 `label_mappings`，系统会自动查找或创建标签。

**工作流程：**
```
1. 收到 CNB 标签 "🚧 处理中_CNB"
   ↓
2. 查询 label_mappings 表
   ↓
3. 未找到 → 调用 Plane API 获取项目所有标签
   ↓
4. 按名称匹配（支持 "🚧 处理中_CNB" 或 "🚧 处理中"）
   ↓
5. 仍未找到 → 自动创建标签（默认颜色 #808080）
   ↓
6. 记录映射到 label_mappings
   ↓
7. 使用该标签 UUID 更新 Issue
```

**实现位置：**
- `internal/plane/client.go`: `ListProjectLabels()`, `CreateLabel()`
- `internal/handlers/issue_label_notify.go`: `findOrCreatePlaneLabel()`

**日志示例：**
```json
// 匹配到已有标签
{"event": "label.matched", "cnb_label": "🚧 处理中_CNB", "plane_name": "🚧 处理中", "label_id": "uuid"}

// 自动创建新标签
{"event": "label.created", "cnb_label": "新功能_CNB", "label_id": "uuid", "color": "#808080"}
```

**优点：**
- ✅ 零配置，开箱即用
- ✅ 自动适配 Plane 已有标签
- ✅ 自动记录映射关系
- ✅ 支持名称模糊匹配（带/不带 `_CNB` 后缀）

---

### 3. 并发处理

**当前实现：**
- 每个请求启动一个 goroutine 异步处理
- 无并发控制（如无限 goroutine）

**建议改进：**
- 使用 worker pool 限制并发数（如 100 个 worker）
- 使用消息队列（如 Redis/RabbitMQ）解耦

---

### 3. 错误重试

**当前实现：**
- 失败时记录日志，不重试
- 依赖 CNB job 的重试机制

**建议改进：**
- 对 429/5xx 错误实现指数退避重试
- 将失败事件写入重试队列

---

### 4. 性能考虑

**当前实现：**
- 每个 Issue 独立处理（多次 DB 查询 + API 调用）

**批量处理建议：**
- 接收批量请求（如 100 个 Issue）
- 批量查询数据库（减少 round trip）
- 批量调用 Plane API（如支持）

---

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
│ 5. 更新 Plane Issue（增量更新）                              │
│    - 读取当前标签 → 识别 CNB 管理 → 保留其他 → 合并去重      │
│    - 调用 Plane API PATCH /issues/{id}/                      │
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

---

## 📝 常见问题

### Q: 服务启动失败，提示数据库连接错误

**A:** 数据库连接是可选的。如果暂时没有数据库：
- 将 `DATABASE_URL` 留空或删除该配置
- 服务会使用内存去重（5 分钟 TTL）
- 日志会提示 "db connect error"，但不影响启动

---

### Q: API 返回 401 错误

**A:** 检查：
1. 请求头格式：`Authorization: Bearer <token>`（注意 "Bearer " 前缀和空格）
2. Token 值与 `INTEGRATION_TOKEN` 环境变量一致
3. Token 没有包含换行符或额外空格

---

### Q: Plane OAuth 配置是否必需？

**A:** ❌ **不必需！** Plane OAuth 配置**不是**服务启动的必需条件：
- ✅ 可以在没有 OAuth 配置时启动服务
- ✅ Issue 标签通知 API 完全不依赖 OAuth
- ❌ 只有访问 `/plane/oauth/start` 或 `/plane/oauth/callback` 时才需要

---

### Q: 当前实现会做什么业务处理？

**A:** 完整实现包括：
- ✅ 接收并验证请求
- ✅ 记录事件到数据库（如有）
- ✅ 过滤 _CNB 后缀标签
- ✅ 查询映射关系
- ✅ 增量更新 Plane Issue 标签
- ✅ 发送飞书通知（如配置）
- ✅ 返回成功响应

实现位置：`internal/handlers/issue_label_notify.go` 中的 `processIssueLabelNotify` 方法。

---

## 🎯 下一步

1. **配置数据库** - 创建必需的映射表
2. **配置 Plane Bot Token** - 启用标签同步
3. **配置飞书群组映射** - 启用飞书通知
4. **修改 CNB Job** - 参考 `.vscode/API-REFERENCE.md`
5. **端到端测试** - 在测试环境验证完整流程

---

## 📚 相关资源

- **API 规范**：`.vscode/API-REFERENCE.md`
- **配置说明**：`.vscode/ConfigNote.md`
- **测试指南**：`.vscode/LOCAL-TESTING-GUIDE.md`
- **架构文档**：`docs/ARCHITECTURE.md`
- **设计文档**：`docs/design/cnb-integration.md`
