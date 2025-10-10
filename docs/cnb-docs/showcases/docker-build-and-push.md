# docker build 并发布到 CNB Docker 仓库

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该文本详细说明了Docker镜像构建发布的流程：通过配置`.cnb.yml`实现三阶段任务（登录CNB仓库→构建镜像→推送镜像），具体包含使用环境变量`CNB_TOKEN_USER_NAME`/`CNB_TOKEN`完成身份验证，通过`docker login`命令登录仓库，借助`docker build`和`docker push`命令生成并推送标签为`latest`的镜像（路径基于`CNB_DOCKER_REGISTRY`和`CNB_REPO_SLUG_LOWERCASE`）。关键依赖Dockerfile配置镜像构建基础设置。

基本思路：
1. 配置 .cnb.yml 文件，用于登录 -> 构建 -> 推送到 CNB Docker 仓库。
2. 配置 Dockerfile 文件，用于配置 Docker 镜像的构建配置。

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建任务。
```yaml
main:
  push:
    - services:
        - docker
      stages:
        # 登录到 CNB Docker 仓库
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        # 构建 Docker 镜像
        - name: docker build
          script: docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest .
        # 发布推送到 CNB Docker 镜像仓库
        - name: docker push
          script: docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest

```