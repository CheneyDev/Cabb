# semantic-release

semantic-release 插件, 根据历史提交信息，生成 `git tag` 和 `changelog`。

`semantic-release`: v20.1.1

## 工作原理

在工作目录下执行自动运行
[semantic-release](https://github.com/semantic-release/semantic-release)工具。

`semantoc-release` 的工作时序，参考 [semantic-release#release-steps](https://github.com/semantic-release/semantic-release#release-steps)

当前工具会自动自行以下 `semantic-release` 插件：

- `@semantic-release/commit-analyzer`
- `@semantic-release/release-notes-generator`
- `@semantic-release/changelog`
- `@semantic-release/npm` (当npm.npmPublish=true时才会加载此插件)
- `@semantic-release/git`

## 注意

1. 使用此插件前，需要工作区先注入`git凭证`。
  可使用[tencentcom/git-set-credential](/docs/plugins/public/tencentcom/git-set-credential)插件实现。
1. 如需激活`npm发布`功能，需要流水线`环境变量`中注入 `NPM_TOKEN` 等 `NPM` 账户信息，
  会自动修改 `package.json` 中的版本信息，并发布 `NPM` 包。

## 参数

### tagFormat

- type: String
- 必填: 否
- 默认值: v${version}

指定 `tag` 格式，其中 `${version}` 将会由生成的目标 `tag` 号替换，如 `1.2.3`。那么 `v${version}` 即为 `v1.2.3`。

### dryRun

- type: Boolean
- 必填: 否
- 默认值: false

仅试运行，生成新的版本号，但不会执行`git tag`和发布操作。

### branch

- type: string
- 必填: 否

自动生成 tag 依据的分支，根据这个分支获取提交日志，根据提交日志生成 tag

### changelog.changelogFile

- type: String
- 必填: 否
- 默认值: `CHANGELOG.md`

指定生成的Changelog文件的文件名

### changelog.changelogTitle

- type: String
- 必填: 否
- 默认值: `-`

指定Changelog文件的标题（第一行内容）

### npm.npmPublish

- type: Boolean
- 必填: 否
- 默认值: 当传入环境变量`NPM_TOKEN`时为`true`, 否则为`false`

是否发布npm包

### npm.pkgRoot

- type: String
- 必填: 否
- 默认值: .

Directory path to publish

### npm.tarballDir

- type: Boolean
- 必填: 否
- 默认值: false

Directory path in which to write the package tarball.
If `false` the tarball is not be kept on the file system.

### git.message

- type: String
- 必填: 否
- 默认值: `chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}`

The message for the release commit

### git.assets

- type: Array[String] | Boolean
- 必填: 否
- 默认值: `['CHANGELOG.md', 'package.json', 'package-lock.json', 'npm-shrinkwrap.json']`

Files to include in the release commit.
Set to `false` to disable adding files to the release commit.

- type: Object
- 必填: 否

## 输出结果

```js
{
    // https://github.com/semantic-release/semantic-release/blob/v20.1.1/docs/developer-guide/js-api.md#lastrelease
    lastRelease.version, // 上次发布的版本号
    lastRelease.gitHead, // 上次发布的 commit
    lastRelease.gitTag, // 上次发布的 tag
    lastRelease.channel, // 上次发布的 channel
    // https://github.com/semantic-release/semantic-release/blob/v20.1.1/docs/developer-guide/js-api.md#commits
    commits,
    // https://github.com/semantic-release/semantic-release/blob/v20.1.1/docs/developer-guide/js-api.md#nextrelease
    nextRelease.type, // 本次发布的 semver 类型
    nextRelease.version, // 本次发布的版本号
    nextRelease.gitHead, // 本次发布的 commit
    nextRelease.gitTag, // 本次发布的 tag
    nextRelease.notes, // 本次发布的变更记录
    // https://github.com/semantic-release/semantic-release/blob/v20.1.1/docs/developer-guide/js-api.md#releases
    releases,
}
```

## 示例

```Dockerfile
docker run --rm \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    \
    -e PLUGIN_DRYRUN="true" \
    \
    tencentcom/semantic-release:dev
```

```yaml
main:
  push:
    - stages:
        - name: set token
          image: tencentcom/git-set-credential
          settings:
            userName: $CNB_COMMITTER
            userEmail: $CNB_COMMITTER_EMAIL
            loginUserName: $CNB_TOKEN_USER_NAME
            loginPassword: $CNB_TOKEN
        - name: semantic-release
          image: tencentcom/semantic-release
          settings:
            tagFormat: v\${version} #由于$符合默认会认为是环境变量，所有加反斜杠转义
            dryRun: true
            changelog.changelogFile: CHANGELOG.md
            changelog.changelogTitle: '-'
            npm.npmPublish: false
            npm.pkgRoot: '.'
            npm.tarballDir: false
            git.message: 'chore(release): \${nextRelease.version} [skip ci]\n\n${nextRelease.notes}'
            git.assets:
              - CHANGELOG.md
              - package.json
              - package-lock.json
              - npm-shrinkwrap.json
```
