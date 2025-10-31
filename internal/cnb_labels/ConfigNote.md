# 配置说明

## Plane 实例信息

```bash
# 自建实例 (nexsea)
实例地址: https://network.1024hub.org:4430
API 地址: https://network.1024hub.org:4430/api
工作区: my-test
项目: test-notify
Access Key: X6YdP5O.bwCu-9cmAJK3n7L57C4YgqHfLsyQjGjs
Instance ID: 29d63944cba60c76dac042c5
```

### 已有配置 (可直接使用)

```bash
PLANE_BASE_URL=https://network.1024hub.org:4430/api
PLANE_APP_BASE_URL=https://network.1024hub.org:4430
```

### 需创建 Plane OAuth 应用

❌ **需在 Plane 创建应用获取 (无法从 Access Key 推导)**:

```bash
PLANE_CLIENT_ID=<待创建>
PLANE_CLIENT_SECRET=<待创建>
PLANE_WEBHOOK_SECRET=<待创建>
PLANE_REDIRECT_URI=<线上服务地址>/plane/oauth/callback
```

**获取步骤**:
1. 访问 `https://network.1024hub.org:4430my-test/settings/integrations` (或 Developer/Apps)
2. 点击"Create Application"或"New App"
3. 填写配置 (使用线上部署地址):
   - Setup URL: `https://<your-domain>/plane/oauth/start`
   - Redirect URI: `https://<your-domain>/plane/oauth/callback`
   - Webhook URL: `https://<your-domain>/webhooks/plane`
4. 保存后复制 `Client ID`、`Client Secret`、`Webhook Secret` 填入 `.env`

**注意**:
- Plane 在线上，必须使用公网可访问的地址 (不能用 localhost)
- 本地开发可使用 ngrok 等内网穿透工具
- Webhook 验签、Bot 授权需要完整 OAuth 流程

## 环境变量配置

### 必填变量

```bash
# 基础运行
PORT=8080
DATABASE_URL=postgres://user:pass@host:5432/plane_integration?sslmode=disable
ENCRYPTION_KEY=<64字符hex>  # 用于敏感令牌加密

# Plane 集成 (OAuth/Webhook)
PLANE_CLIENT_ID=<从 Plane 应用获取>
PLANE_CLIENT_SECRET=<从 Plane 应用获取>
PLANE_WEBHOOK_SECRET=<从 Plane 应用获取>
PLANE_REDIRECT_URI=<线上地址>/plane/oauth/callback
PLANE_BASE_URL=https://network.1024hub.org:4430/api
PLANE_APP_BASE_URL=https://network.1024hub.org:4430

# CNB 集成 (回调 + API)
INTEGRATION_TOKEN=<自行生成强随机字符串>  # 用于 .cnb.yml 回调鉴权
CNB_APP_TOKEN=<从 CNB 个人设置获取>        # 用于调用 CNB API

# 飞书集成 (事件 + IM)
LARK_APP_ID=<飞书开发者后台>
LARK_APP_SECRET=<飞书开发者后台>
LARK_ENCRYPT_KEY=<飞书事件签名校验>
```

### 可选变量 (有默认值)

```bash
TIMEZONE=Local
CNB_BASE_URL=https://api.cnb.cool
CNB_OUTBOUND_ENABLED=false
ADMIN_SESSION_COOKIE=pi_admin_session
ADMIN_SESSION_TTL_HOURS=12
ADMIN_SESSION_SECURE=false  # 生产环境建议 true (需 HTTPS)
```

### 管理后台初始账号 (开发/测试环境)

```bash
ADMIN_BOOTSTRAP_EMAIL=admin@example.com
ADMIN_BOOTSTRAP_PASSWORD=ChangeMe123!
```

首次启动会自动创建该管理员账号；**生产环境请改为强密码**。

## 本地开发环境已生成密钥

```bash
# 已生成的密钥 (2025-10-29)
ENCRYPTION_KEY=bcf5b37ca2f22d52d6fee9b6a379dc7accdbd8459c345b379da97831c88ec71f
INTEGRATION_TOKEN=2ab4e3ae01e9dc1aafe1043d43d3ff06de2d674dd54705cb8f564312761f58c0
LARK_ENCRYPT_KEY=0453523b2e3ef797098eb5c5fc009b7db96fc933c93ec120ac7f32580d37dc0c
```

**注意**: 生产环境请重新生成密钥，不要使用本地开发环境的值。

## 本地 Docker PostgreSQL 配置

```bash
# 容器信息
容器名: plane-postgres
镜像: postgres:18-alpine
端口: 5432:5432
凭据: postgres/postgres
数据库: plane_integration

# 连接串
DATABASE_URL=postgres://postgres:postgres@localhost:5432/plane_integration?sslmode=disable

# 管理命令
docker stop plane-postgres    # 停止
docker start plane-postgres   # 启动
docker rm plane-postgres      # 删除 (会丢失数据)
docker logs plane-postgres    # 查看日志
```

## CNB 凭据获取说明

### INTEGRATION_TOKEN (回调鉴权)

- 用途: `.cnb.yml` 回调本服务时使用的 Bearer Token
- 生成方式:
  ```bash
  # macOS/Linux
  openssl rand -hex 32
  
  # 或使用 Python
  python3 -c 'import secrets; print(secrets.token_hex(32))'
  ```
- 注入方式: 在 CNB 仓库的"设置 → 流水线密钥"中添加 `INTEGRATION_TOKEN` 变量

### CNB_APP_TOKEN (API 访问)

- 用途: 调用 CNB OpenAPI 的访问令牌
- 获取路径:
  1. 登录 CNB 平台
  2. 进入"个人设置 → 访问令牌"
  3. 创建新令牌，选择所需权限 (至少需要 `repo:read`, `issue:write`)
  4. 复制令牌并填入 `CNB_APP_TOKEN`

## 飞书配置说明

### 必需配置 (发送消息)

- `LARK_APP_ID`: 飞书开发者后台 → 应用 → 凭证与基础信息 → App ID
- `LARK_APP_SECRET`: 同上页面 → App Secret

### 推荐配置 (事件校验)

- `LARK_ENCRYPT_KEY`: 开发者后台 → 事件与回调 → 加密策略 → Encrypt Key
  - **重要**: 请确保"不要开启事件消息体加密"，仅用于签名校验
- `LARK_VERIFICATION_TOKEN`: 同上页面 → Verification Token (兜底校验)

### 事件订阅配置

- 订阅 URL: `https://<your-domain>/webhooks/lark/events`
- 订阅事件: `im.message.receive_v1` (接收消息)
- 所需权限: `im:message.group_at_msg:readonly` 等

## 部署检查清单

- [ ] 生产环境重新生成所有密钥 (ENCRYPTION_KEY/INTEGRATION_TOKEN/LARK_ENCRYPT_KEY)
- [ ] 在 Plane 中创建 OAuth 应用并填入线上域名
- [ ] 配置生产数据库连接串 (建议启用 SSL: `sslmode=require`)
- [ ] 设置 `ADMIN_SESSION_SECURE=true` (需 HTTPS)
- [ ] 修改 `ADMIN_BOOTSTRAP_PASSWORD` 为强密码
- [ ] 在 CNB 仓库流水线密钥中注入 `INTEGRATION_TOKEN`
- [ ] 在飞书开发者后台配置事件订阅 URL
- [ ] 验证健康检查: `curl https://<your-domain>/healthz`

## 参考文档

- 完整配置说明: `README.md`
- 架构设计: `docs/ARCHITECTURE.md`
- CNB 集成: `docs/design/cnb-integration.md`
- 飞书集成: `docs/design/feishu-integration.md`
