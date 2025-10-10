---
title: 跳过流水线
permalink: https://docs.cnb.cool/zh/build/skip-pipeline.html
summary: 文本主要介绍了跳过流水线的方法和事件忽略的情形。在 `push`、`commit.add` 和 `branch.create` 事件中，当最近的一个提交信息包含 `[ci skip]` 或 `[skip ci]`，或者在 `git push` 命令中使用 `-o ci.skip` 参数时，可以主动跳过流水线；同时 `云原生构建` 会忽略密钥仓库的文件变更等特定情形下的事件，以避免无意义的流水线运行。
---

## 主动跳过流水线

有时候，我们不想触发流水线，此时，我们可以下面的方式主动跳过流水线。

`push`、`commit.add` 和 `branch.create` 事件里，以下两种情况会跳过流水线：

1. 最近一个 commit message 里带 `[ci skip]` 或 `[skip ci]`
2. git push -o ci.skip

如：

```bash
git commit -m "feat: some feature [ci skip]"
git push origin main
```

或

```bash
git commit -m "feat: some feature"
git push origin main -o ci.skip
```

## 事件忽略

为避免无意义的流水线，`云原生构建` 会忽略如下情形的事件：

1. 密钥仓库的文件变更。
1. 仓库迁移时 head commit 的时间戳早于仓库创建时间十分钟前的事件。
1. 仓库设置里未勾选 `允许自动触发` 的 git 操作引起的事件。
1. 仓库设置中未勾选 `Fork 的仓库默认允许自动触发`，fork 后的仓库 `允许自动触发` 自动为未勾选状态。
