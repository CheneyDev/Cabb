# 创建 pr 插件

支持创建 pr

## 镜像

`cnbcool/create-pr:latest`

## 支持的事件

- push
- branch.create
- branch.delete
- pull_request
- pull_request.target
- pull_request.approved
- pull_request.changes_requested
- pull_request.mergeable
- pull_request.merged
- tag_push
- vscode
- auto_tag
- tag_deploy.*

## 参数说明

- `target_branch`: 目标分支，类型为`string`，必填。这里目标分支指的是当前仓库的分支。
- `head_branch`: 源分支，类型为`string`，非必填，不填的话默认选择当前分支名`CNB_BRANCH`。这里的源分支指的是当前仓库的分支，或者跨仓库时的跨仓分支。
- `head_repo_slug`: 源分支仓库 slug，类型为`string`，跨仓库时必填。
- `title`: pr 标题，类型为`string`，非必填。
- `body`: pr 内容，类型为`string`，非必填。
- `reviewers`: pr 评审者，类型为`string`，多个参数间用`,`隔开，非必填。
- `assinees`: pr 处理人，类型为`string`，多个参数间用`,`隔开，非必填。
- `labels`: pr 标签，类型为`string`，多个参数间用`,`隔开，非必填。

## 在 云原生构建 中使用

```yml
$:
  push:
    - stages:
      - name: push 事件创建 mr
        image: cnbcool/create-pr:latest
        settings:
          title: "test"
          target_branch: "main"
          reviewers: "xxx"
          assignees: "xxx"

```

上面的代码执行后，会创建一个名为`test`的 pr，源分支为当前分支，目标分支为`main`，评审者为`xxx`，处理人为`xxx`。
