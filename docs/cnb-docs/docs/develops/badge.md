---
title: 徽章
permalink: https://docs.cnb.cool/zh/develops/badge.html
summary: 徽章用于展示指标数据，可通过颜色区分达标情况。云原生构建相关的徽章展示构建耗时和状态等指标，根据不同 Git 事件有细分类型，还有准备区、流水线配置等相关徽章；仓库方面有展示 fork 数和 star 数的徽章及相关访问路径和参数 。
---

徽章用于展示某个指标的数据，可通过颜色区分指标是否达标。

## 云原生构建相关徽章

### 徽章访问路径

云原生构建时产生的徽章，访问路径：

- precise: `https://cnb.cool/{group}/{repository}/-/badge/git/{sha}/{metrics}`

- latest: `https://cnb.cool/{group}/{repository}/-/badge/git/latest/{metrics}`

- 参数含义
  - `group`：仓库所在 `group`
  - `repository`：仓库名
  - `sha`: 表示 CommitId 前8位
  - `latest`：最近一次的数据
  - `metrics`: 指标名，如 `ci/status/push`，其对应的徽章表示：云原生构建时，push 事件触发的构建耗时
 
  
### 徽章类型

云原生构建时，会自动上传相关构建指标的徽章数据：

#### 1. Git 事件

- ci/status/push

push 事件构建耗时和构建状态

![push 事件](https://cnb.cool/svg/badge/push?message=pending%2C%2010s&color=pending)
![push 事件](https://cnb.cool/svg/badge/push?message=success%2C%2010s&color=success)
![push 事件](https://cnb.cool/svg/badge/push?message=failure%2C%2010s&color=failure)

- ci/status/commit.add

commit.add 事件构建耗时和构建状态

![commit.add 事件](https://cnb.cool/svg/badge/commit.add?message=pending%2C%2010s&color=pending)
![commit.add 事件](https://cnb.cool/svg/badge/commit.add?message=success%2C%2010s&color=success)
![commit.add 事件](https://cnb.cool/svg/badge/commit.add?message=failure%2C%2010s&color=failure)

- ci/status/branch.create

branch.create 事件构建耗时和构建状态

![branch.create 事件](https://cnb.cool/svg/badge/branch.create?message=pending%2C%2010s&color=pending)
![branch.create 事件](https://cnb.cool/svg/badge/branch.create?message=success%2C%2010s&color=success)
![branch.create 事件](https://cnb.cool/svg/badge/branch.create?message=failure%2C%2010s&color=failure)

- ci/status/pull_request

pull_request 事件构建耗时和构建状态

![pull_request 事件](https://cnb.cool/svg/badge/pull_request?message=pending%2C%2010s&color=pending)
![pull_request 事件](https://cnb.cool/svg/badge/pull_request?message=success%2C%2010s&color=success)
![pull_request 事件](https://cnb.cool/svg/badge/pull_request?message=failure%2C%2010s&color=failure)

- ci/status/pull_request.update

pull_request.update 事件构建耗时和构建状态

![pull_request.update 事件](https://cnb.cool/svg/badge/pull_request.update?message=pending%2C%2010s&color=pending)
![pull_request.update 事件](https://cnb.cool/svg/badge/pull_request.update?message=success%2C%2010s&color=success)
![pull_request.update 事件](https://cnb.cool/svg/badge/pull_request.update?message=failure%2C%2010s&color=failure)

- ci/status/pull_request.target

pull_request.target 事件构建耗时和构建状态

![pull_request.target 事件](https://cnb.cool/svg/badge/pull_request.target?message=pending%2C%2010s&color=pending)
![pull_request.target 事件](https://cnb.cool/svg/badge/pull_request.target?message=success%2C%2010s&color=success)
![pull_request.target 事件](https://cnb.cool/svg/badge/pull_request.target?message=failure%2C%2010s&color=failure)

- ci/status/pull_request.merged

pull_request.merged 事件构建耗时和构建状态

![pull_request.merged 事件](https://cnb.cool/svg/badge/pull_request.merged?message=pending%2C%2010s&color=pending)
![pull_request.merged 事件](https://cnb.cool/svg/badge/pull_request.merged?message=success%2C%2010s&color=success)
![pull_request.merged 事件](https://cnb.cool/svg/badge/pull_request.merged?message=failure%2C%2010s&color=failure)

- ci/status/tag_push

tag_push 事件构建耗时和构建状态

![tag_push 事件](https://cnb.cool/svg/badge/tag_push?message=pending%2C%2010s&color=pending)
![tag_push 事件](https://cnb.cool/svg/badge/tag_push?message=success%2C%2010s&color=success)
![tag_push 事件](https://cnb.cool/svg/badge/tag_push?message=failure%2C%2010s&color=failure)

#### 2. 准备工作区

- ci/git-clone-yyds

工作区大小和准备工作区产生的耗时和工作区大小：![git-clone-yyds 事件](https://cnb.cool/svg/badge/git-clone-yyds?message=2.7s%2C%20163.77%20GB&color=success)

#### 3. 流水线配置

- ci/pipeline-as-code

云原生构建的配置文件：![pipeline-as-code 事件](https://cnb.cool/svg/badge/pipeline-as-code?message=.cnb.yml&color=orange)

#### 4. 云原生开发

- code/vscode-started

准备开发环境耗时：![vscode-started](https://cnb.cool/svg/badge/准备开发环境?message=22s&color=success)

#### 5. 单元测试

使用内置任务 [testing:coverage](../build/internal-steps/README.md#coverage) 可上报单元测试徽章数据

- testing/unit/coverage

单元测试全量覆盖率

![coverage](https://cnb.cool/svg/badge/coverage?message=18.67%25&color=l2)
![coverage](https://cnb.cool/svg/badge/coverage?message=38.67%25&color=l3)
![coverage](https://cnb.cool/svg/badge/coverage?message=58.67%25&color=l4)
![coverage](https://cnb.cool/svg/badge/coverage?message=78.67%25&color=l5)

- testing/unit/coverage-pr

本次 `pull_request` 的单元测试增量覆盖率

![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=18.67%25&color=l2)
![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=38.67%25&color=l3)
![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=58.67%25&color=l4)
![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=78.67%25&color=l5)

## 仓库相关徽章

### fork 徽章

用徽章形式展示仓库的 fork 数量：![fork徽章](https://cnb.cool/svg/badge/fork?message=18&color=orange)

访问路径：`https://cnb.cool/{group}/{repository}/-/badge/fork`

参数含义：

- `group`：仓库所在组织路径
- `repository`：仓库名

### star 徽章

用徽章形式展示仓库的 star 数量：![star徽章](https://cnb.cool/svg/badge/star?message=20&color=orange)

访问路径：`https://cnb.cool/{group}/{repository}/-/badge/star`

参数含义：

- `group`：仓库所在组织路径
- `repository`：仓库名
