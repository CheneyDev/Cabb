---
title: 默认环境变量
permalink: https://docs.cnb.cool/zh/build/build-in-env.html
summary: 该文本主要介绍了云原生构建内置的默认环境变量，这些变量只读不可写。包括基础变量（如 CI、CNB 相关等众多）、提交类变量（与构建时的提交相关）、仓库类变量（和仓库信息有关）、构建类变量（构建相关的详情）、合并类变量（涉及合并类事件相关）、远程开发类变量（远程开发相关）、Issue 类变量（Issue 事件相关）、评论类变量（评论事件相关），并对变量的具体情况进行了详细说明。
redirectFrom: /zh/build-in-env.html
---

`云原生构建` 内置了一些默认环境变量，只可读不可写，如果构建过程中尝试覆盖默认环境变量不会生效。

下文 `合并类事件` 包含以下六种：

- `pull_request`
- `pull_request.update`
- `pull_request.target`
- `pull_request.approved`
- `pull_request.changes_requested`
- `pull_request.mergeable`
- `pull_request.merged`
- `pull_request.comment`

了解更多[触发事件](./trigger-rule.md#trigger-event)。

[[TOC]]

## 基础变量

### CI

true

### CNB

true

### CNB_WEB_PROTOCOL

当前 web 使用的协议，http | https

### CNB_WEB_HOST

当前 web 使用的 HOST

### CNB_WEB_ENDPOINT

当前 web 使用的地址，包含协议、HOST、路径（如有）

### CNB_API_ENDPOINT

当前 API 地址，包含协议、HOST、路径（如有）

可与 `CNB_TOKEN` 配合在 CI 中调用 API 接口

### CNB_GROUP_SLUG

仓库所属组织路径

### CNB_GROUP_SLUG_LOWERCASE

仓库所属组织路径（小写格式）

### CNB_EVENT

值为触发构建的事件名称

事件类型参见[事件](./trigger-rule.md#trigger-event)

### CNB_EVENT_URL

- 对于由 ==合并类事件== 触发的构建，值为该 `PR` 的链接

- 对于由 `push`、`commit.add`、`branch.create`、`tag_push` 触发的构建，值为最新 `Commit` 的链接

- 否则值为空字符串

### CNB_BRANCH

- 对于由 `push`、`commit.add`、`branch.create`、`branch.delete` 触发的构建，值为当前的分支名

- 对于由 `合并类事件` 触发的构建，值为目标分支的分支名

- 对于由 `tag_push` 触发的构建，值为 `tag` 名

- 对于由 `自定义事件` 触发的构建，值为对应的分支名称

- 对于由 `crontab` 触发的构建，值为对应的分支名称

### CNB_BRANCH_SHA

- 对于 `branch.delete` 触发的构建，为空字符串

- 其他情况为 `CNB_BRANCH` 最近一次提交的 `sha`

### CNB_DEFAULT_BRANCH

仓库默认分支

### CNB_TOKEN_USER_NAME

用户临时令牌对应的用户名，固定为 `cnb`

### CNB_TOKEN

用户令牌，可用于代码提交、API 调用等

对于 pull_request 事件，权限有：

- `repo-code:r`
- `repo-pr:r`
- `repo-issue:r`
- `repo-notes:rw`
- `repo-contents:r`
- `repo-registry:r`
- `repo-commit-status:rw`
- `account-profile:r`

对于非 pull_request 事件，权限有：

- `repo-code:rw`
- `repo-pr:rw`
- `repo-issue:rw`
- `repo-notes:rw`
- `repo-contents:rw`
- `repo-registry:rw`
- `repo-commit-status:rw`
- `repo-cnb-trigger:rw`
- `repo-cnb-history:r`
- `repo-cnb-detail:r`
- `repo-basic-info:r`
- `repo-manage:r`
- `account-profile:r`
- `group-resource:r`

权限含义参考页面 个人设置 页面的 访问令牌

### CNB_TOKEN_FOR_AI

用户令牌，在 合并类事件 中由 AI 使用

权限有：

- repo-notes:rw

权限含义参考页面 `个人设置` 页面的 访问令牌

### CNB_IS_CRONEVENT

是否是定时任务事件

### CNB_DOCKER_REGISTRY

制品库 Docker 源地址

### CNB_HELM_REGISTRY

制品库 Helm 源地址

## 提交类变量

### CNB_BEFORE_SHA

- 对于由 `push`、`commit.add` 触发的构建，值为分支推送前远端仓库该分支最近一次提交的 `sha`，若是新建分支，值为 `0000000000000000000000000000000000000000`
- 对于 `branch.create` 触发的构建，值为 `0000000000000000000000000000000000000000`

### CNB_COMMIT

构建时对应的代码 sha

- 对于 `push`、`commit.add`、`branch.create` 触发的构建，是最后一次提交的 `sha`
- 对于 `tag_push`、`tag_deploy.*` 触发的构建，是该 `tag` 最后一次提交的 `sha`
- 对于 `auto_tag`、`branch.delete`、`issue.*` 类事件，是主分支最后一次提交的 `sha`
- 对于 `pull_request.merged` 触发的构建，是合并后的 `sha`
- 对于 `pull_request.target`、`pull_request.mergeable` 触发的构建，是目标分支最后一次提交的 `sha`
- 对于 `pull_request`、`pull_request.approved`、`pull_request.changes_requested`、`pull_request.comment`，代码尚未真正合并，取源分支最后一次提交的 `sha`，但构建时会进行预合并，即合并后的内容会作为最终结果
- 对于 `云原生开发`、`自定义事件` 触发的构建，是指定分支的最后一次提交的 `sha`

### CNB_COMMIT_SHORT

`CNB_COMMIT` 的缩写，取其前 8 位字符

### CNB_COMMIT_MESSAGE

CNB_COMMIT 对应的提交信息

### CNB_COMMIT_MESSAGE_TITLE

`CNB_COMMIT_MESSAGE` 的 `title` 部分，即首行

### CNB_COMMITTER

`CNB_COMMIT` 对应的提交者

### CNB_COMMITTER_EMAIL

`CNB_COMMITTER` 对应的邮箱

### CNB_NEW_COMMITS_COUNT

对于由 `commit.add` 触发的构建，值为新增的 `Commits` 的数量，最大为 99。

可结合 `git log -n` 查看新增的 `Commits`。

### CNB_IS_TAG

对于分支为 `Tag` 的构建，值为 ture

### CNB_TAG_MESSAGE

- `Tag message`：对于分支为 `Tag` 的构建, 会有该环境变量
- 否则值为空字符串

### CNB_TAG_RELEASE_TITLE

- `Release 标题`：对于分支为 `Tag` 的构建, 如果 `Release` 标题不为空，才会有值
- 否则值为空字符串

### CNB_TAG_RELEASE_DESC

- `Release 描述`：对于分支为 `Tag` 的构建, 且 `Release` 描述不为空，才会有值
- 否则值为空字符串

### CNB_TAG_IS_RELEASE

- `Tag 是否存在对应的 Release`：对于分支为 Tag 的构建，如果 `Tag` 存在对应的 `Release`，则为 true
- 否则为 `false`

### CNB_TAG_IS_PRE_RELEASE

- 对于分支为 Tag 的构建，若存在对应的 `Release`，且 `Release` 为 预发布，则值为 `true`
- 否则为 false

### CNB_IS_NEW_BRANCH

当前分支是否属于一个新创建的分支，默认为 `false`

### CNB_IS_NEW_BRANCH_WITH_UPDATE

当前分支是否属于一个新创建的分支，且带有新 commit，默认为 `false`

## 仓库类变量

### CNB_REPO_SLUG

目标仓库路径，格式为 `group_slug/repo_name`，`group_slug/sub_gourp_slog/.../repo_name`

### CNB_REPO_SLUG_LOWERCASE

目标仓库路径小写格式

### CNB_REPO_NAME

目标仓库名称

### CNB_REPO_NAME_LOWERCASE

目标仓库名称小写格式

### CNB_REPO_ID

目标仓库的 `id`

### CNB_REPO_URL_HTTPS

目标仓库仓库 https 地址

## 构建类变量

### CNB_BUILD_ID

当前构建的流水号，全局唯一

### CNB_BUILD_WEB_URL

当前构建的日志地址

### CNB_BUILD_START_TIME

当前构建的开始时间，UTC 格式，示例 `Tue, 24 Dec 2024 07:42:58 GMT`

### CNB_BUILD_USER

当前构建的触发者名称

### CNB_BUILD_USER_ID

当前构建的触发者 `id`

### CNB_BUILD_STAGE_NAME

当前构建的 `stage` 名称

### CNB_BUILD_JOB_NAME

当前构建的 `job` 名称

### CNB_BUILD_JOB_KEY

当前构建的 `job` key，同 `stage` 下唯一

### CNB_BUILD_WORKSPACE

自定义 `shell` 脚本执行的工作空间根目录

### CNB_BUILD_FAILED_MSG

流水线构建失败的错误信息，可在 `failStages` 中使用

### CNB_BUILD_FAILED_STAGE_NAME

流水线构建失败的 `stage` 的名称，可在 `failStages` 中使用

### CNB_PIPELINE_NAME

当前 `pipeline` 的 `name`，没声明时为空

### CNB_PIPELINE_KEY

当前 `pipeline` 的索引 `key`，例如 `pipeline-0`

### CNB_PIPELINE_ID

当前 `pipeline` 的 `id`，全局唯一字符串

### CNB_PIPELINE_DOCKER_IMAGE

当前 `pipeline` 所使用的 `docker image`，如：`alpine:latest`

### CNB_PIPELINE_STATUS

当前流水线的构建状态，可在 `endStages` 中查看，其可能的值包括：

- `success`：表示流水线构建成功完成。
- `error`：表示流水线构建过程中发生了错误。
- `cancel`：表示流水线构建被取消。

### CNB_RUNNER_IP

当前 `pipeline` 所在 `Runner` 的 `ip`

### CNB_CPUS

当前构建流水线可以使用的最大 `CPU` 核数

### CNB_MEMORY

当前构建流水线可以使用的最大 `内存` 大小，单位为 `GiB`

### CNB_IS_RETRY

当前构建是否由 `rebuild` 触发

### HUSKY_SKIP_INSTALL

兼容 ci 环境下 husky

## 合并类变量

### CNB_PULL_REQUEST

- 对于由 `pull_request`、`pull_request.update`、`pull_request.target` 触发的构建，值为 `true`
- 否则值为 `false`

### CNB_PULL_REQUEST_LIKE

- 对于由 `合并类事件` 触发的构建，值为 `true`
- 否则值为 `false`

### CNB_PULL_REQUEST_PROPOSER

- 对于由 `合并类事件` 触发的构建，值为提出 `PR` 者名称
- 否则值为空字符串

### CNB_PULL_REQUEST_TITLE

- 对于由 `合并类事件` 触发的构建，值为提 `PR` 时候填写的标题
- 否则值为空字符串

### CNB_PULL_REQUEST_BRANCH

- 对于由 `合并类事件` 触发的构建，值为发起 `PR` 的源分支名称
- 否则值为空字符串

### CNB_PULL_REQUEST_SHA

- 对于由 `合并类事件` 触发的构建，值为当前 `PR` 源分支最新的提交 `sha`
- 否则值为空字符串

### CNB_PULL_REQUEST_TARGET_SHA

- 对于由 `合并类事件` 触发的构建，值为当前 `PR` 目标分支最新的提交 `sha`
- 否则值为空字符串

### CNB_PULL_REQUEST_MERGE_SHA

- 对于由 `pull_request.merged` 事件触发的构建，值为当前 `PR` 合并后的 `sha`
- 对于由 `pull_request`、`pull_request.update`、`pull_request.target`、`pull_request.mergeable`、`pull_request.comment` 触发的构建，值为当前 `PR` 预合并后的 `sha`
- 否则值为空字符串

### CNB_PULL_REQUEST_SLUG

- 对于由 `合并类事件` 触发的构建，值为源仓库的仓库 slug，如 `group_slug/repo_name`，`group_slug/sub_gourp_slog/.../repo_name`
- 否则值为空字符串

### CNB_PULL_REQUEST_ACTION

对于由 `合并类事件` 触发的构建，可能的值有：

- created: 新建 PR
- code_update: 源分支 push
- status_update: 评审通过或 CI 状态变更时 `PR` 变成可合并 否则值为空字符串

### CNB_PULL_REQUEST_ID

- 对于由 `合并类事件` 触发的构建，值为当前或者关联 `PR` 的全局唯一 `id`
- 否则值为空字符串

### CNB_PULL_REQUEST_IID

- 对于由 `合并类事件` 触发的构建，值为当前或者关联 `PR` 在仓库中的编号 `iid`
- 否则值为空字符串

### CNB_PULL_REQUEST_REVIEWERS

- 对于由 `合并类事件` 触发的构建，值为评审人列表，多个以 `,` 分隔
- 否则值为空字符串

### CNB_PULL_REQUEST_REVIEW_STATE

对于由 `合并类事件` 触发的构建

- 有评审者且有人通过评审，值为 `approve`
- 有评审者但无人通过评审，值为 `unapprove`
- 否则值为空字符串

### CNB_REVIEW_REVIEWED_BY

- 对于由 `合并类事件` 触发的构建，值为同意评审的评审人列表，多个以 `,` 分隔
- 否则值为空字符串

### CNB_REVIEW_LAST_REVIEWED_BY

- 对于由 `合并类事件` 触发的构建，值为最后一个同意评审的评审人
- 否则值为空字符串

## 远程开发类变量

### CNB_VSCODE_WEB_URL

远程开发地址，仅声明了 `services : vscode` 时存在

## Issue 类变量

### CNB_ISSUE_ID

- 对于 `issue.*` 触发的构建，值为 `Issue` 全局唯一 `ID`
- 否则值为空字符串

### CNB_ISSUE_IID

- 对于 `issue.*` 触发的构建，值为 `Issue` 在仓库中的编号 `iid`
- 否则值为空字符串

### CNB_ISSUE_TITLE

- 对于 `issue.*` 触发的构建，值为 `Issue` 的 `title`
- 否则值为空字符串

### CNB_ISSUE_DESCRIPTION

- 对于 `issue.*` 触发的构建，值为 `Issue` 的 `description`
- 否则值为空字符串

### CNB_ISSUE_OWNER

- 对于 `issue.*` 触发的构建，值为 `Issue` 作者用户名
- 否则值为空字符串

### CNB_ISSUE_STATE

- 对于 `issue.*` 触发的构建，值为 `Issue` 状态：`open`、`closed`
- 否则值为空字符串

### CNB_ISSUE_IS_RESOLVED

- 对于 `issue.*` 触发的构建，表示 `Issue` 是否被解决：`true`、`false`
- 否则值为空字符串

## 评论类变量

### CNB_COMMENT_ID

- 对于评论事件触发的构建，值为评论全局唯一 `ID`
- 否则值为空字符串

### CNB_COMMENT_BODY

- 对于评论事件触发的构建，值为评论内容
- 否则值为空字符串
