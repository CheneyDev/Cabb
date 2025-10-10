---
title: 单/双容器模式
permalink: https://docs.cnb.cool/zh/workspaces/double-container.html
summary: 云原生开发支持单容器和双容器两种模式。单容器模式下，开发环境和代码服务（`code-server`）在同一容器，推荐使用；双容器模式下，两者在不同容器但工作区互通，适用于开发环境容器未安装`code-server`的情况 。
---

云原生开发可使用两种模式启动：

- 单容器模式：开发环境 和 代码服务（`code-server`）在同一个容器中
- 双容器模式：开发环境 和 代码服务（`code-server`）在两个不同的容器中，两个容器的工作区（`/workspace`）目录互通

**开发环境**：可理解为用户自定义的开发环境容器（使用 `.ide/Dockerfile` 或自定义镜像 `image` 自定义）或 CNB 默认的容器环境，
其中包含了用户自己安装的软件或默认安装的软件

## 单容器模式

推荐使用单容器模式。以下方式以单容器模式启动

### 使用默认配置启动

当用户未做任何配置（即未配置启动流水线和 `.ide/Dockerfile`）时，使用默认配置启动开发环境。
默认开发环境已安装 `code-server` 服务和 `ssh` 服务，因此会使用单容器模式启动，同时支持 WebIde 和 VSCode 远程开发。

详见[文档](./default-dev-env.md)

### 使用自定义配置启动

当用户配置了启动流水线或 `.ide/Dockerfile` 时，自定义镜像或 `.ide/Dockerfile` 中已安装 `code-server` 服务，
此时开发环境和 `code-server` 服务在同一个容器中，将使用单容器模式启动。

如下三种方式(启动开发环境的镜像均安装了 code-server 服务)均使用单容器模式启动：

- [通过 docker 镜像指定开发环境](./custom-dev-env.md#同时自定义开发环境和启动流程)
- [通过 Dockerfile 指定开发环境](./custom-dev-env.md#通过Dockerfile自定义开发环境)
- [同时自定义开发环境和启动流程](./custom-dev-env.md#同时自定义开发环境和启动流程)

## 双容器模式

当开发环境容器中未安装 `code-server` 服务时，我们将会额外启动一个 `code-server` 的容器作为代码服务所在容器，
以支持 `WebIde` 和 `Vscode` 客户端访问开发环境，即双容器模式。两个容器的工作区（`/workspace`）目录互通。

以下两种方式会以双容器模式启动：

### 指定的 docker 镜像中未安装 code-server

```yaml{14-18}
# .cnb.yml
$:
  vscode:
    # 指定 docker 镜像中未安装 code-server 服务
    - docker:
        image: node:22
      services:
        - vscode
        - docker
      # 开发环境启动后会执行的任务
      stages:
        - name: ls
          script: ls -al
```

### 自定义的 Dockerfile 中未安装 code-server

```dockerfile
# .ide/Dockerfile

# 可将 node 替换为需要的基础镜像
FROM node:20

# 按需安装软件
RUN apt-get update && apt-get install -y git wget unzip openssh-server

# 指定字符集支持命令行输入中文（根据需要选择字符集）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8
```

## 双容器模式-跨容器终端

双容器模式，访问 `WebIde` 或 `VSCode` 时，实际上访问的是 `code-server` 容器。

那么要怎么在编辑器中访问开发环境容器呢？

我们在 `code-server` 容器中额外支持了一个跨容器终端，即在 WebIDE 或 VSCode 客户端中，打开终端时，
默认打开的是一个可跨容器访问开发环境容器的终端，该终端名为 `CNB`。

如果开发环境容器中不支持 `git` 命令，可切换为非 `CNB` 命名的其他终端，使用 `code-server` 自带的终端实现 `git` 操作。

## 双容器模式-插件能力限制

双容器模式，访问 `WebIde` 或 `VSCode` 时，安装的 `VSCode` 插件在 `code-server` 容器中，
可能没法实现 `Debug` 调试之类需要访问开发环境容器的能力
