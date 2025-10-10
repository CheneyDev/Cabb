# npm

这个版本是从 `plugins/npm` fork 出来修改的。

实现了 `npm publish`，支持官方仓库和私有仓库。

调整点：

- 去掉 package.json 中的 registry 检查
- 文档调整，符合 CNB 的使用场景

## 输入

- `username`: `string` 用户名
- `password`: `string` 密码
- `token`: `string` token方式鉴权
- `email`: `string` 邮箱
- `registry`: `string` registry，默认：`https://registry.npmjs.org/`
- `folder`: `string` 要发布的目录，默认当前目录。
- `fail_on_version_conflict`: `boolean` 版本存在时报错退出
- `tag`: `boolean` NPM publish tag `--tag`
- `access`: `string` NPM scoped package access `--access`

## 在 云原生构建 上使用

### 发布到 CNB 的 npm 制品库

```yaml
$:
  # tag push 时触发
  tag_push:
  - docker:
      image: node:22-alpine
    stages:
    - name: update version
      # 修改 package.json 中的version为当前tag，但不会提交代码
      script: npm version $CNB_BRANCH --no-git-tag-version
    - name: npm publish
      image: tencentcom/npm
      settings:
        username: $CNB_TOKEN_USER_NAME
        token: $CNB_TOKEN
        email: $CNB_COMMITTER_EMAIL
        registry: https://npm.cnb.cool/xxx/xxx/-/packages/
        folder: ./
        fail_on_version_conflict: true

```

### 发布到 npm 官方仓库

推送到官方仓库时不能使用密码，要使用 token 方式鉴权。

1. 将 token 保存到密钥仓库 `npm-token.yml` 文件中，然后使用 `imports` 导入。

```yaml
# npm-token.yml
NPM_USER: your-username
NPM_TOKEN: your-token
NPM_EMAIL: your-email@email.com
```

```yaml
main:
  push:
  - stages:
    - name: npm publish
      image: tencentcom/npm
      imports: https://cnb.cool/xxx/npm-token.yml
      settings:
        username: $NPM_USER
        token: $NPM_TOKEN
        email: $NPM_EMAIL
        folder: ./
        fail_on_version_conflict: true
```

## 在 Docker 上使用

```bash
docker run --rm \
  -e PLUGIN_USERNAME=xxx \
  -e PLUGIN_TOKEN=xxx \
  -e PLUGIN_EMAIL=xxx@xxx.com \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  tencentcom/npm
```
