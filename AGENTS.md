# AGENTS.md

本文件为在本仓库内工作的智能 Agent/贡献者提供统一约定。其作用域为本目录树（项目根）下的所有文件与子目录。

## 0. 重要提醒（必须）
- 在进行 Plane/CNB 相关的 API 调用、Webhook 处理、`.cnb.yml` 回调/示例、字段映射、签名校验等代码编写或调整之前，务必先查阅 `docs/` 目录中已下载的最新版官方文档与本仓库的设计文档。
- 若官方文档与现有代码/本文档存在冲突，协议细节（端点/字段/签名/状态码/限流）以 `docs/` 为准；实现需据此更新，并同步修改 `README.md`、`docs/ARCHITECTURE.md` 与相关设计文档。
- 如发现文档不明确或缺失，请在变更说明中标注“待确认”，并以特性开关/占位实现降低风险，避免推出错误行为。
- 提交前自检：是否逐项对照 `docs/` 校验了端点、字段、签名、示例与错误码？是否更新了示例与文档？

## 1. 项目背景与目标
- 将 Plane 作为“工作项中枢”，对接：
  - CNB 代码托管（Issue/PR/分支/提交摘要）。
  - 飞书（创建/链接/预览、线程↔评论双向同步、项目通知）。
- 参考设计见：`docs/design/cnb-integration.md`、`docs/design/feishu-integration.md`、`docs/design/integration-feature-spec.md`。

## 2. 代码组织（必须遵循）
- Go 版本：1.24+。
- 目录结构：
  - `cmd/server/`：入口与路由引导。
  - `internal/handlers/`：HTTP 端点（Plane/CNB/Lark/Admin/Jobs）。
  - `internal/store/`：数据访问与数据库连接。
  - `internal/version/`：版本信息。
  - `pkg/config/`：环境变量与配置。
  - `db/migrations/`：Postgres 迁移脚本。
  - `docs/`：架构与设计文档。
- 新增端点：放在 `internal/handlers/` 下，路由集中由 `internal/handlers/router.go` 注册。
- 不要随意重命名/移动现有文件与公共路由，除非用户明确要求。

## 3. 编码规范（Go/Echo）
- 依赖：优先使用标准库与 Echo；新增三方依赖需有明确收益，并在 `go.mod` 中声明，随后执行 `go mod tidy`。
- 处理器：使用 `echo.Context`，统一返回 JSON；占位实现可返回 501/Not Implemented 或 202/Accepted。
- 上下文与超时：发起外部调用前应设置 `context.Context` 与合理超时；避免阻塞。
- 命名：导出符号使用驼峰命名；SQL 列名使用蛇形（snake_case）。
- 日志：避免打印敏感信息（token/secret）；错误日志包含 request id/来源/事件类型等关键上下文。
- 不在代码中加入版权/许可证头；保持改动最小且聚焦。

## 4. 安全与校验（必须）
- Plane Webhook：校验 `X-Plane-Signature`（HMAC-SHA256(secret, raw_body)），常量时间比较，兼容 `sha256=...` 前缀。
- CNB 回调：校验 `Authorization: Bearer $INTEGRATION_TOKEN`。
- 飞书事件：支持 challenge；接入签名/时间戳校验时，避免时钟漂移导致误拒绝。
- 令牌与机密：未来实现透明加密存储；日志中绝不输出明文。

## 5. 数据库与迁移（Postgres 16）
- 迁移文件放于 `db/migrations/`，按递增编号命名：`0002_xxx.sql`、`0003_xxx.sql`...
- 迁移内容尽量幂等（`IF NOT EXISTS`）；不要修改已发布的旧迁移，新增文件进行演进。
- 本地初始化：`psql "$DATABASE_URL" -f db/migrations/0001_init.sql`。

## 6. 幂等与重试（设计遵循）
- 入站事件写入 `event_deliveries` 以去重与重试；可重试错误（429/5xx）指数退避，4xx 语义错误不重试。
- 业务去重依赖链接/映射表（如 `issue_links`/`thread_links`）。

## 7. 文档与可观测
- 新增/调整端点或数据表：务必更新 `README.md` 与必要的架构文档（`docs/ARCHITECTURE.md`）。
- 保持示例与实际实现同步（curl、环境变量名、表名/列名）。

## 9. 运行与本地验证（建议）
- 环境变量：复制 `.env.example` 并按需修改。
- 迁移：`psql "$DATABASE_URL" -f db/migrations/0001_init.sql`。
- 启动：`go run ./cmd/server`，健康检查 `GET /healthz`。

## 10. 非目标与限制
- 不擅自引入全局状态、侵入式框架或大规模目录重构。
- 不修复与任务无关的其它问题（可在总结中提示给用户）。
- 不在未授权的情况下访问外部网络或删除用户数据。

## 11. 参考
- 设计文档：`docs/design/*`
- 架构说明：`docs/ARCHITECTURE.md`
- 项目总览与使用：`README.md`

## 12. Git 提交信息规范（必须）
- 每次任务结束后，编写一条规范的 Git 提交信息并用于提交变更；推荐遵循 Conventional Commits：
  - 格式：`<type>(<scope>): <subject>`
  - 常用 type：`feat`、`fix`、`docs`、`refactor`、`chore`、`test`、`perf`、`build`。
  - 主题行不超过 50 个字符，使用祈使句；中文或英文均可，但保持一致。
  - 可选正文（换行后 72 列换行）：说明动机、关键变更点与影响范围。
  - 如关联任务/Issue，请在正文中使用引用（例如：`Refs: CNB-123` 或 `Closes: #42`）。
  - 示例：
    - `feat(cnb): 接入 issue 回调入口，验证 bearer token`
    - `docs(readme): 增加快速开始与 API 列表`

## 13. 文档编写风格（必须）
- 目标：易于理解且信息密度高（高信噪比），服务内部工程师快速决策与落地。
- 受众：内部技术读者，默认具备基础上下文（Plane/CNB/飞书），避免赘述常识与营销式语言。
- 语言：中文为主；专业术语保留英文原文（括注）以便搜索与一致性，例如“线程（thread）/ 状态（state）”。
- 结构：
  - 先结论后细节：开头给出“结果/要做什么/验收点”，随后提供背景与推导。
  - 短句与短段：每段不超过 3 行；列表每组 4–6 条，按重要性排序。
  - 明确层次：标题 1–3 个词，贴合内容；相关点合并，避免碎片化。
- 呈现规范：
  - 使用代码块展示命令/配置/SQL；使用内联代码标注路径、接口、字段名与环境变量。
  - 提供可复制的示例：请求/响应样例、默认值、边界条件、失败示例与常见错误。
  - 必要时使用 Mermaid 时序/流程/类图；图注简洁，避免大段描述。
- 一致性：
  - 术语与命名统一（CNB/Feishu/Lark/Plane、Issue/Work Item 等）；时间与时区写法一致。
  - 大小写与单位规范：ID/HTTP 方法/状态码大写；尺寸/时间单位明确。
- 更新策略：
  - 任何接口/表结构/环境变量变更，需同步更新 `README.md`、`docs/ARCHITECTURE.md` 与相关设计文档；确保示例与实现一致。
  - 在 PR/变更说明中列出“变更摘要 / 影响范围 / 迁移指引 / 验收要点（DoD）”。
  - 未决问题用“待确认”或 TODO 标注，并给出下一步行动或负责人。
- 反例（避免）：大段背景铺垫、低价值截图、含糊表述（如“可能、应该”）、无示例的描述、过长的段落。
- 自检清单（提交前快速检查）：
  - 标题是否直接表达结论与范围？
  - 是否包含可复制的命令/请求示例与默认值？
  - 是否标注限制、边界条件与失败案例？
- 是否给出文件路径/接口名便于检索？
  - 是否与代码实现保持一致（端点、字段、表名）？

## 14. API 设计与实现规范（必须）
- 协议与数据格式
  - 仅支持 `HTTP/1.1+` 与 `application/json; charset=utf-8`；所有时间使用 `RFC3339` UTC（例如：`2025-01-02T15:04:05Z`）。
  - JSON 命名使用 `snake_case`；布尔/数值使用原生类型；ID 使用 `uuid` 字符串；金额类如需，统一以最小货币单位整数表示。
- 路径与命名
  - 资源路径使用名词、短横线分隔（kebab-case），示例：`/admin/mappings/repo-project`、`/webhooks/plane`。
  - 新增通用资源建议使用复数名词集合 + 单个资源：`/projects`、`/projects/{id}`；避免动词式路径。
  - API 版本：外部稳定面向使用者时采用前缀版本 `/v1/...`；内部/回调端点可不显式版本，但保持向后兼容。
- HTTP 语义
  - `GET`（安全/幂等）、`PUT`（幂等全量）、`PATCH`（幂等部分更新）、`POST`（非幂等创建/执行）、`DELETE`（幂等删除）。
  - 成功状态码：`200`（查询/更新）、`201`（创建，附 `Location`）、`202`（异步受理）、`204`（无响应体）。
  - 失败状态码：`400`（参数错误）、`401`（未授权）、`403`（禁止）、`404`（不存在）、`409`（冲突/去重命中）、`422`（语义校验失败）、`429`（限流）、`5xx`（服务端）。
- 错误响应格式（统一）：
  - 结构：`{ "error": { "code": "invalid_signature", "message": "签名校验失败", "details": { ... }, "request_id": "..." } }`
  - `code` 为稳定、可机器识别的标识；`message` 面向人类；`details` 放字段级错误或上下文。
- 幂等与去重
  - Webhook/回调端：以 `delivery_id + payload_sha256` 在 `event_deliveries` 去重；重复请求返回 `200` 并标注 `status=duplicate`（可选）。
  - 客户端发起的 `POST`：支持 `Idempotency-Key` 请求头；服务端应以键 + 请求摘要去重（必要时返回 `409` 或已创建结果）。
- 分页/过滤/排序
  - 推荐游标分页：`limit`（<=100）、`cursor`（`next_cursor` 在响应中返回）；小表可使用 `limit`+`offset`（最大 1000）。
  - 过滤：使用命名空间前缀，示例：`filter[project_id]=...&filter[active]=true`；避免与保留参数冲突。
  - 排序：`sort=created_at,-id`（前缀 `-` 表示降序）。
  - 字段选择/展开：`fields=...`、`expand=...`（如需减少往返或展开关联）。
- 请求/响应示例与契约
  - 新增/修改端点需在文档提供最小可复制示例（请求/响应/错误），并说明边界条件与默认值。
  - 优先维护 OpenAPI 规范（如 `docs/openapi.yaml`），用于生成校验或客户端；示例与实现保持一致。
- 安全与鉴权
  - Webhook：Plane 使用 HMAC-SHA256（`X-Plane-Signature`）；飞书使用签名/时间戳；CNB 回调使用 Bearer `INTEGRATION_TOKEN`。
  - 管理端点需鉴权与审计（预留 Header 或网关层处理）；限制请求体大小（如 `1MB`）防滥用。
  - 传递与返回 `X-Request-ID`，便于追踪；在日志中保留来源、事件类型与目标资源。
- 兼容性与演进
  - 避免破坏性变更：优先“向后兼容的新增”（新增字段、可选参数、默认值）；字段重命名以新增字段 + 过渡期方式处理。
  - 若必须破坏性变更，提升主版本 `/v2`；在文档中提供迁移指引与双栈运行策略。
- 性能与限流
  - 为对端速率限制设计退避策略：`429` 返回 `Retry-After`；客户端/作业侧指数退避（10s/30s/…）。
  - 使用连接池与合理超时；长列表接口返回 `total_estimate` 而非强一致 `total`（如成本过高）。
- 并发控制
  - 支持 `ETag`/`If-Match` 以处理并发更新；无条件更新风险较高时建议启用。
- 可观测性
  - 日志结构化：至少包含 `request_id`、`source`、`endpoint`、`latency_ms`、`result`；错误日志含 `error.code`。
  - 指标：QPS、p95、错误率、429 次数与重试队列长度。
