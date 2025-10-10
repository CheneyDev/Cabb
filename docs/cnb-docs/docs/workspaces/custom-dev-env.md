---
title: 自定义开发环境
permalink: https://docs.cnb.cool/zh/workspaces/custom-dev-env.html
summary: 文本主要介绍了自定义开发环境的两种方式。一是通过 `.cnb.yml` 编写远程开发事件流水线并指定 `pipeline.docker.image` 指定开发环境镜像，可根据镜像中是否安装 code-server 以不同模式启动开发环境 。二是通过编写仓库根目录下 `.ide/Dockerfile` 自定义开发环境，默认流水线优先使用该文件构建镜像，若不存在或构建失败则使用默认镜像 。若要同时自定义环境和启动流程，可同时编写 `.ide/Dockerfile` 和 `.cnb.yml` 。
---

## 通过 docker 镜像指定开发环境

可以通过在 `.cnb.yml` 编写远程开发事件流水线，并指定`pipeline.docker.image` 指定开发环境镜像。

```yaml{4-10}
# .cnb.yml
$:
  vscode:
    - docker:
        # 指定开发环境镜像，可以是任意可访问的镜像。
        # 如果 image 指定的镜像中已安装 code-server 代码服务，将使用单容器模式启动开发环境
        # 如果 image 指定的镜像中未安装 code-server 代码服务，将使用双容器模式启动开发环境
        # 如下镜像为 CNB 默认开发环境镜像，已安装代码服务，将使用单容器模式启动开发环境
        # 可按需替换为其他镜像
        image: cnbcool/default-dev-env:latest
      services:
        - vscode
        - docker
      # 开发环境启动后会执行的任务
      stages:
        - name: ls
          script: ls -al
```

## 通过 Dockerfile 自定义开发环境

如果通过指定镜像无法满足需求，可以自行编写 `Dockerfile` 来自定义开发环境。

在仓库根目录下增加 `.ide/Dockerfile` 文件，在 Dockerfile 中自由定制开发环境。

如果未自定义启动流水线，启动开发环境时使用默认流水线创建开发环境。
默认流水线会优先使用 `.ide/Dockerfile` 构建一个镜像，作为开发环境基础镜像。

注意：启动开发环境的默认流水线中，同时配置了 `默认镜像` 和 `.ide/Dockerfile`，
如果 `.ide/Dockerfile` 不存在或构建失败，会使用 `默认镜像` 作为开发环境基础镜像。
如果遇到启动的环境不符合预期，可以查看构建日志 `prepare` 阶段 `.ide/Dockerfile` 是否构建成功

```dockerfile
# .ide/Dockerfile

# 可将 node 替换为需要的基础镜像
FROM node:20

# 安装 code-server 和 vscode 常用插件
RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension cnbcool.cnb-welcome \
  && code-server --install-extension redhat.vscode-yaml \
  && code-server --install-extension dbaeumer.vscode-eslint \
  && code-server --install-extension waderyan.gitblame \
  && code-server --install-extension mhutchie.git-graph \
  && code-server --install-extension donjayamanne.githistory \
  && code-server --install-extension tencent-cloud.coding-copilot \
  && echo done

# 安装 ssh 服务，用于支持 VSCode 等客户端通过 Remote-SSH 访问开发环境（也可按需安装其他软件）
RUN apt-get update && apt-get install -y git wget unzip openssh-server

# 指定字符集支持命令行输入中文（根据需要选择字符集）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8
```

## 同时自定义开发环境和启动流程

如果需要同时自定义开发环境和启动流程，可以编写 `.ide/Dockerfile` 和 `.cnb.yml`。

### 自定义启动流水线

```yaml{4-9}
# .cnb.yml
$:
  vscode:
    - docker:
        build: .ide/Dockerfile
        # 可选择是否同时定义 build 和 image
        # 此时会优先使用 .ide/Dockerfile 构建镜像
        # 如果 .ide/Dockerfile 构建失败，则使用 image 指定的镜像保证环境能启动成功
        # image: cnbcool/default-dev-env:latest
      services:
        - vscode
        - docker
      # 开发环境启动后会执行的任务
      stages:
        - name: ls
          script: ls -al
```

### 自定义开发环境

```dockerfile
# .ide/Dockerfile

# 可将 node 替换为需要的基础镜像
FROM node:20

# 安装 code-server 和 vscode 常用插件
RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension cnbcool.cnb-welcome \
  && code-server --install-extension redhat.vscode-yaml \
  && code-server --install-extension dbaeumer.vscode-eslint \
  && code-server --install-extension waderyan.gitblame \
  && code-server --install-extension mhutchie.git-graph \
  && code-server --install-extension donjayamanne.githistory \
  && code-server --install-extension tencent-cloud.coding-copilot \
  && echo done

# 安装 ssh 服务，用于支持 VSCode 等客户端通过 Remote-SSH 访问开发环境（也可按需安装其他软件）
RUN apt-get update && apt-get install -y git wget unzip openssh-server

# 指定字符集支持命令行输入中文（根据需要选择字符集）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8
```
