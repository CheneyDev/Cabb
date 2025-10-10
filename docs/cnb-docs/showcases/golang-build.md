# CNB 配 Go 项目构建并推送 docker 镜像到 CNB 制品库

本文档地址：https://cnb.cool/examples/showcase

文档摘要：CNB项目使用两阶段Dockerfile构建Go服务镜像，通过配置.cnb.yml文件集成制品库登录认证、多任务流水线构建和镜像推送功能。技术实现包含环境变量控制、构建依赖管理、容器化部署，并提供云原生开发环境预置（含代码服务器、Go工具链和语言包），版本发布流程整合了变更日志采集、制品构建及二进制文件附件上传功能。

基本逻辑：配置 Dockerfile 用于构建 docker 镜像 -> 配置 .cnb.yml 文件用于构建与推送 docker 镜像到 CNB 制品库。

仓库地址：[CNB 实现 Go 项目构建](https://cnb.cool/examples/ecosystem/golang-build)

## 1. 准备 Dockerfile 用于构建 Go 项目

在项目根目录下创建一个 `Dockerfile` 文件，用于构建 docker 镜像。

```dockerfile
# 构建 go 服务镜像
FROM golang:1.22.5-alpine3.20 as builder

ARG GOPROXY
ENV GOPROXY $GOPROXY

ARG GOSUMDB
ENV GOSUMDB $GOSUMDB

WORKDIR /data/workspace

COPY main.go go.mod go.sum ./

RUN go mod vendor

RUN go build

FROM debian:bookworm

ENV TZ="Asia/Shanghai"
WORKDIR /usr/local/app

COPY --from=builder /data/workspace/helloworld /usr/local/app/helloworld

CMD ["/usr/local/app/helloworld"]
```

## 2. 配置 .cnb.yml 文件用于构建 docker 镜像

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建任务：安装依赖 -> 编译 -> 上传到目标服务器。

```yaml
main:
  push:
    # 上传 docker 镜像到 CNB 制品库
    - services:
        - docker
      stages:
        # docker 登录
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        # docker 构建镜像
        - name: docker build
          script: docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest .
        # docker 推送镜像
        - name: docker push
          script: docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest
```

## 3. 配置 Go 项目 Dockerfile 云原生开发

```Dockerfile
# 此文件为远程开发环境配置文件
FROM debian:bookworm

ENV GO_VERSION=1.22.5

RUN apt update &&\
    apt install -y wget rsync unzip openssh-server vim lsof git git-lfs locales locales-all libgit2-1.5 libgit2-dev net-tools jq curl &&\
    rm -rf /var/lib/apt/lists/*

# install golang
RUN curl -fsSLO https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz &&\
    rm -rf /usr/local/go && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz &&\
    ln -sf /usr/local/go/bin/go /usr/bin/go &&\
    ln -sf /usr/local/go/bin/gofmt /usr/bin/gofmt &&\
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2 &&\
    rm -rf go${GO_VERSION}.linux-amd64.tar.gz

# install code-server
RUN curl -fsSL https://code-server.dev/install.sh | sh
RUN code-server --install-extension dbaeumer.vscode-eslint &&\
    code-server --install-extension pinage404.git-extension-pack &&\
    code-server --install-extension redhat.vscode-yaml &&\
    code-server --install-extension esbenp.prettier-vscode &&\
    code-server --install-extension golang.go &&\
    code-server --install-extension eamodio.gitlens &&\
    code-server --install-extension mhutchie.git-graph &&\
    code-server --install-extension ms-azuretools.vscode-docker &&\
    code-server --install-extension PKief.material-icon-theme &&\
    code-server --install-extension tencent-cloud.coding-copilot &&\
    echo done

# 安装 Go Tools
ENV GOPATH /root/go
ENV PATH="${PATH}:${GOPATH}/bin"

RUN go install -v golang.org/x/tools/gopls@latest

ENV LC_ALL zh_CN.UTF-8
ENV LANG zh_CN.UTF-8
ENV LANGUAGE zh_CN.UTF-8

```

## 4. 配置 Go 项目上传二进制文件附件到 release 附件

tag_push 推送 tag 时，上传二进制文件到 release 附件

```yaml
$:
  tag_push:
    # 上传二进制包到 release 附件
    - docker:
        build: .ide/Dockerfile
      stages:
        - name: changelog
          image: cnbcool/changelog
          exports:
            latestChangeLog: LATEST_CHANGE_LOG
        - name: create release
          type: git:release
          options:
            title: release
            description: ${LATEST_CHANGE_LOG}
        - name: go mod
          script: go mod vendor
        - name: go build
          script: go build -o helloworld
        - name: release 上传附件
          image: cnbcool/attachments:latest
          settings:
            attachments:
              - helloworld
```