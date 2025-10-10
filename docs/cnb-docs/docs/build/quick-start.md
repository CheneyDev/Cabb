---
title: 快速开始
permalink: https://docs.cnb.cool/zh/build/quick-start.html
summary: “快速开始”指南以仓库为主体，通过编写配置文件描述在特定分支下发生特定事件时执行的任务，帮助用户创建并触发流水线。介绍了从创建仓库、编写配置文件到查看构建详情的步骤，并以常见需求为例详细说明了配置文件的语法和结构，强调了通过分支、事件、流水线和任务等关键要素实现流水线构建。
redirectFrom: /zh/quick-start.html
---

`云原生构建` 以仓库为主体，由 **配置文件** 描述在哪个 **分支** 下发生什么 **事件** 时执行什么 **任务**。

下面以一个最简示例开始，介绍如何一步步操作，创建仓库并触发一条流水线。

接着以开发中常见的 `Pull Request` 流水线检测需求为目标，介绍如何一步步编写配置文件。

在此基础上，可以在[最佳实践](https://cnb.cool/examples/showcase)中的选择具体场景的仓库，`fork` 仓库或复制其中配置文件，编写出符合自身需求的流水线配置。

## 编写你的第一条流水线

### 1. 创建仓库

新建一个仓库（如有跳过），创建好后，可以点击 `云原生开发` 按钮快速创建一个开发环境。

![workspace](https://docs.cnb.cool/images/quick-start/workspace.png)

选择 `WebIDE` 进入开发界面，方便快捷。

![choose-ide](https://docs.cnb.cool/images/quick-start/choose-ide.png)

### 2. 编写 .cnb.yml 配置文件

一个简单的流水线配置如下：

```yaml
# 分支名
main:
  # 事件名
  push:
    # 要执行的任务
    - stages:
        - name: echo
          script: echo "hello world"
```

添加 `CI` 配置文件 `.cnb.yml`，将该内容复制进配置文件，提交并 `push` 到远端 `main` 分支。

![add-ci](https://docs.cnb.cool/images/quick-start/add-ci.png)

即会触发流水线构建。

### 3. 查看构建详情

在仓库页面点击 `云原生构建` 可以看到构建列表。

![pipeline-push](https://docs.cnb.cool/images/quick-start/pipeline-push.png)

最新一条即是刚刚触发的 `push` 事件流水线，点进去可以看到构建详情。

## 配置说明

接下来，我们以一个常见的流水线需求，来简单介绍解释一下配置文件的语法：

**“需求：主分支有 `Pull Request` 时，触发流水线进行 lint 和 test 检测，未通过则发出通知。 ”**

我们分析下这个需求，可以从中抽取一些要素：

1. 主分支，比如 `main`。
2. 仓库事件，即 `pull_request` 事件。
3. 流水线任务：
   1. lint
   2. test
4. 失败时的任务
   1. notify

下面我们根据这些要素一步步编写 `.cnb.yml` 配置文件。

第一层属性为分支名：

```yaml
#分支名
main:
```

分支下有 `pull_request` 事件时触发构建:

```yaml{4}
#分支名
main:
  # 事件名
  pull_request:
```

事件可以执行多条流水线（并行），流水线有多个任务（串行或并行）。

这里我们简化，事件下只有一条流水线:

```yaml{6,7}
#分支名
main:
  # 事件名
  pull_request:
    # 数组类型表示可以有多个流水线
    - name: pr-check
      stages:
```

包括两个串行的任务，分别是 lint 和 test。

```yaml{9,10,11,12,13}
#分支名
main:
  # 事件名
  pull_request:
    # 数组类型表示可以有多个流水线
    - name: pr-check
      # 流水线下的多个任务
      stages:
        # 要执行的任务
        - name: lint
          script: echo "lint job"
        - name: test
          script: echo "test job"
```

如果失败的话，要发送通知，流水线下除 `stages` 表示期望执行的任务外，还有个 `failStages` 表示 `stages` 任务执行失败时要执行的任务:

```yaml{13}
#分支名
main:
  # 事件名
  pull_request:
    # 数组类型表示可以有多个流水线
    - name: pr-check
      # 流水线下的多个任务
      stages:
        - name: lint
          script: echo "lint job"
        - name: test
          script: echo "test job"
      # stages 失败时执行的任务
      failStages:
        - name: lint
          script: echo "notify to some chat group"
```

总结下，一个流水线的执行过程是：

1. 仓库发生事件
2. 确定所属分支
3. 确定事件名
4. 执行流水线
5. 执行任务
6. 失败时的任务

想了解配置文件更多用法请移步 [配置文件](./configuration.md)。

## 语法说明

想了解配置文件详细语法请移步 [语法](./grammar.md)。

## 最佳实践

更多完整示例参考 [最佳实践](https://cnb.cool/examples/showcase)。
