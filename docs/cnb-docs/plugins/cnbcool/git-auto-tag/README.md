# git-auto-tag

基于 [semantic-release](https://semantic-release.gitbook.io) 分析工具制作的自动打 TAG 命令。

`semantic-release` 默认在提交记录中搜索 `feat` 、 `fix` 、 `perf` 开头的提交记录
（[详细规则](https://semantic-release.gitbook.io/semantic-release/#commit-message-format)）。
并根据最后一个 TAG，自动计算下一个版本号。
如果想增加其他类型的提交记录，可配置 `releaseRules` 参数。

使用 `git-auto-tag` 需要团队约定
[提交格式](https://www.conventionalcommits.org/en/v1.0.0/),
建议配合 `commit-lint` 相关插件使用。

## 参数

### message

- type: string
- 默认值: "Release {tag 名}"
- 必填: 否

生成 tag 时通过 `-m` 添加的 `message`，例如：`git tag -a v0.0.1 -m "Release v0.0.1"`

### tagFormat

- type: string
- 默认值: v${version}
- 必填: 否

指定 `tag` 格式，其中 `${version}` 会由生成的 `tag` 替换。如 `tag` 为 `1.2.3`，那么 `v${version}` 为 `v1.2.3`。

### dryRun

- type: boolean
- 默认值: false
- 必填: 否

是否只得到下一个版本号等变量，而非真正的生成 `tag`。

### blockWhenFail

- type: boolean
- 默认值: false
- 必填: 否

### releaseRules

- type: string
- 必填: 否

从指定的文件中读取自定义发布规则，如果没有指定，则使用默认规则，仅 `feat` 、 `fix` 、 `perf` 三种格式的提交注释会生成 `tag` 。

如需扩充发布规则可自定义 `releaseRules` 。

规则如何配置详见[文档](https://github.com/semantic-release/commit-analyzer#releaserules)

示例文件 `release-rules.json`：

```json
[
  { "type": "docs", "scope": "README", "release": "patch" },
  { "type": "refactor", "scope": "core-*", "release": "minor" },
  { "type": "refactor", "release": "patch" },
  { "scope": "no-release", "release": false }
];

```

### toFile

- type: string
- 默认值: auto_tag.json
- 必填: 否

指定文件保存执行结果，有如下字段

```js
{
  //`Object` with `version`, `gitTag` and `gitHead`
  // of the last release.  
  lastRelease,
  //`Object` with `version`, `gitTag`, `gitHead` and `notes`
  // of the release being done.
  nextRelease,
  commits,
  releases,
}
```

### branch

- type: string
- 默认值: main
- 必填: 否

自动生成 tag 依据的分支，根据这个分支获取提交日志，根据提交日志生成 tag

### repoUrlHttps

- type: string
- 必填: 否
- 默认值：通过 `git config --list` 获取到的 git 配置中的 `remote.origin.url` 的值

目标仓库仓库 https 地址

## 输出

```js
{
  tag, // 生成的 tag
}
```

## 在 云原生构建 上使用

```yaml
# .cnb.yml
main:
  push:
    - stages:
        - name: auto-tag
          image: cnbcool/git-auto-tag:latest
          settings:
            tagFormat: v\${version}
            toFile: tag_info.json
            dryRun: false
            blockWhenFail: false
            branch: $CNB_BRANCH
            repoUrlHttps: $CNB_REPO_URL_HTTPS
          exports:
            tag: NEW_TAG
        - name: show tag
          script: echo $NEW_TAG
        - name: show tag res
          script: cat tag_info.json

```
