# Git Sync Plugin

一个用于在不同 Git 平台之间同步代码的插件。支持通过 HTTPS 或 SSH 方式同步代码到其他 Git 托管平台。

例如从 CNB 同步到 GitHub，从 GitHub 同步到 CNB。

## 功能特点

- 支持 HTTPS (推荐) 和 SSH 两种认证方式
- 支持推送指定分支或所有分支
- 支持推送标签
- 支持强制推送
- 可配置 Git 用户信息
- 支持自定义 Git 服务器
- 支持私有仓库认证

## 使用方法

### 在 CNB 中使用

以下实例是从 CNB 同步到 GitHub 的，其他平台类似

```yaml
main:
  push:
    - stages:
        - name: sync to github
          image: tencentcom/git-sync
          settings:
            target_url: https://github.com/username/repo.git
            auth_type: https
            username: ${GIT_USERNAME}
            password: ${GIT_ACCESS_TOKEN}
            branch: main
```

### 在 GitHub Actions 中使用

以下是从 GitHub 同步到 CNB 的示例， 注意 CNB 的 Git user 是 `cnb`

```yaml
name: Sync to CNB
on: [push]

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Sync to CNB Repository
        run: |
          docker run --rm \
            -v ${{ github.workspace }}:${{ github.workspace }} \
            -w ${{ github.workspace }} \
            -e PLUGIN_TARGET_URL="https://cnb.cool/username/repo.git" \
            -e PLUGIN_AUTH_TYPE="https" \
            -e PLUGIN_USERNAME="cnb" \
            -e PLUGIN_PASSWORD=${{ secrets.GIT_PASSWORD }} \
            -e PLUGIN_BRANCH="main" \
            -e PLUGIN_GIT_USER="cnb" \
            -e PLUGIN_GIT_EMAIL="cnb@cnb.cool" \
            -e PLUGIN_FORCE="true" \
            tencentcom/git-sync
```

### 使用 Docker 直接运行

```bash
docker run --rm \
  -e PLUGIN_TARGET_URL="https://github.com/username/repo.git" \
  -e PLUGIN_AUTH_TYPE="https" \
  -e PLUGIN_USERNAME="your-username" \
  -e PLUGIN_PASSWORD="your-access-token" \
  -e PLUGIN_BRANCH="main" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  tencentcom/git-sync
```

## 参数说明

| 参数名     | 必填 | 默认值                | 说明                                               |
| ---------- | ---- | --------------------- | -------------------------------------------------- |
| target_url | 是   | -                     | 目标仓库的 URL，支持 HTTPS 或 SSH 格式             |
| auth_type  | 否   | https                 | 认证类型，可选值：`https` 或 `ssh`                 |
| username   | 否*  | -                     | HTTPS 认证时的用户名（*使用 HTTPS 时必填）         |
| password   | 否*  | -                     | HTTPS 认证时的密码或访问令牌（*使用 HTTPS 时必填） |
| ssh_key    | 否*  | -                     | SSH 私钥内容（*使用 SSH 时必填）                   |
| branch     | 否   | -                     | 要推送的目标分支。不指定时推送所有分支             |
| force      | 否   | false                 | 是否强制推送（使用 `--force` 选项）                |
| push_tags  | 否   | false                 | 是否推送标签                                       |
| git_user   | 否   | Git Sync Plugin       | Git 提交时使用的用户名                             |
| git_email  | 否   | git-sync@plugin.local | Git 提交时使用的邮箱                               |
| git_host   | 否   | -                     | 自定义 Git 服务器域名                              |

## 安全建议

1. 使用 HTTPS 认证时，建议使用访问令牌（Access Token）而不是实际密码
2. 确保将敏感信息（如密码、访问令牌、SSH 密钥）保存在安全的环境变量，如 CNB 的密钥仓库中
3. 如果使用 SSH 密钥，确保密钥具有适当的权限
4. 建议在目标仓库上设置适当的访问权限控制

## 常见问题

1. HTTPS 认证失败
   - 检查用户名和密码/令牌是否正确
   - 确认令牌是否具有足够的权限
   - 验证目标仓库 URL 是否正确

2. 推送失败
   - 检查是否有目标仓库的写入权限
   - 确认分支名称是否正确
   - 如果遇到冲突，考虑使用 `force: true`

3. 自定义 Git 服务器
   - 确保 `git_host` 参数设置正确
   - 检查服务器的 SSH 指纹是否正确添加

## License

[MIT License](LICENSE)
