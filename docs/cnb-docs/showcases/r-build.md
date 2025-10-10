# CNB 配置 R 语言项目，并且构建 docker 镜像

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该文档指导完成R语言项目的Docker镜像构建及云原生开发配置：1.通过.cnb.yml定义构建流程，使用r:4.3.1-alpine环境实现Docker登录/构建/推送操作，镜像存储于CNB制品库；2.基于rocker/r-ver的Dockerfile完成R依赖安装（dplyr/gapminder）、脚本复制及执行配置；3.云开发环境基于r-base镜像集成code-server、VSCode插件组、SSH服务及中文字符集支持。关键实现包含制品库认证机制、多阶段构建脚本和开发环境功能集成。

R 语言代码构建镜像并发布到 CNB Docker 仓库。

代码仓库：[r-build](https://cnb.cool/examples/ecosystem/r-build)

## 1. 配置 .cnb.yml 文件

r:4.3.1-alpine 构建环境 -> 启用 docker 服务 -> 打包项目 -> 登录制品库 -> 构建镜像 -> 推送镜像
```yaml
main:
  push:
    # 上传 docker 镜像到 CNB 制品库
    - services:
        - docker
      stages:
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        - name: docker build
          script: docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest .
        - name: docker push
          script: docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest
```

## 2. 配置 Dockerfile 文件

```Dockerfile
# 使用 rocker/r-ver 作为基础 R 镜像
FROM rocker/r-ver

# 在容器中创建一个目录
RUN mkdir /home/r-environment

# 安装 R 依赖包
RUN R -e "install.packages(c('dplyr', 'gapminder'))"

# 将我们的 R 脚本复制到容器中
COPY script.R /home/r-environment/script.R

# 运行 R 脚本
CMD R -e "source('/home/r-environment/script.R')"
```

## 3. 配置 R 语言项目云原生开发

在 `.ide/Dockerfile` 文件中配置云原生开发环境: 

r-base 构建环境 -> 安装 code-server 和 vscode 常用插件 -> 安装 ssh 服务 -> 指定字符集支持命令行输入中文

```Dockerfile
FROM r-base

RUN apt-get update ; apt-get install -y curl
# 安装 code-server 和 vscode 常用插件
RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension redhat.vscode-yaml \
  && code-server --install-extension dbaeumer.vscode-eslint \
  && code-server --install-extension eamodio.gitlens \
  && code-server --install-extension tencent-cloud.coding-copilot \
  && echo done

# 安装 ssh 服务，用于支持 VSCode 客户端通过 Remote-SSH 访问开发环境
RUN  apt-get install -y wget unzip openssh-server git

# 指定字符集支持命令行输入中文（根据需要选择字符集）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8

```